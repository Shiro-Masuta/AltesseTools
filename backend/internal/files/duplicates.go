package files

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// FileHash contient le chemin et le hash
type FileHash struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

// FindDuplicates parcourt un dossier et retourne les fichiers en doublons par hash (parallélisé)
func FindDuplicates(root string) (map[string][]string, error) {
	numWorkers := runtime.NumCPU()
	fileCh := make(chan string, 100)
	resultCh := make(chan FileHash, 100)
	errCh := make(chan error, 1)

	var wg sync.WaitGroup

	// Worker pool pour calculer les hash
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range fileCh {
				hash, err := hashFileMD5(path)
				if err != nil {
					select {
					case errCh <- err:
					default:
					}
					return
				}
				resultCh <- FileHash{Path: path, Hash: hash}
			}
		}()
	}

	// Parcours des fichiers
	go func() {
		defer close(fileCh)
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return err
			}
			if !info.IsDir() {
				fileCh <- path
			}
			return nil
		})
	}()

	// Fermeture du resultCh après la fin des workers
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Collecte des résultats
	hashes := make(map[string][]string)
	for {
		select {
		case fh, ok := <-resultCh:
			if !ok {
				duplicates := make(map[string][]string)
				for h, files := range hashes {
					if len(files) > 1 {
						duplicates[h] = files
					}
				}
				return duplicates, nil
			}
			hashes[fh.Hash] = append(hashes[fh.Hash], fh.Path)
		case err := <-errCh:
			return nil, err
		}
	}

}

// DeleteDuplicates supprime les fichiers doublons en gardant le premier fichier de chaque groupe
func DeleteDuplicates(duplicates map[string][]string) ([]string, error) {
	numWorkers := runtime.NumCPU()
	taskCh := make(chan string, 100)
	deletedCh := make(chan string, 100)
	errCh := make(chan error, 1)

	var wg sync.WaitGroup

	// Workers de suppression
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range taskCh {
				if err := os.Remove(path); err != nil {
					select {
					case errCh <- err:
					default:
					}
					return
				}
				deletedCh <- path
			}
		}()
	}

	// On envoie les doublons à supprimer
	go func() {
		defer close(taskCh)
		for _, paths := range duplicates {
			if len(paths) > 1 {
				for _, dupPath := range paths[1:] {
					taskCh <- dupPath
				}
			}
		}
	}()

	// Fermeture du canal deleted après workers
	go func() {
		wg.Wait()
		close(deletedCh)
	}()

	// Collecte des chemins supprimés
	var deleted []string
	for {
		select {
		case path, ok := <-deletedCh:
			if !ok {
				return deleted, nil
			}
			deleted = append(deleted, path)
		case err := <-errCh:
			return deleted, err
		}
	}

}

// hashFileMD5 calcule le hash MD5 d'un fichier
func hashFileMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil

}
