package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// OptionRename définit les paramètres du batch rename
type OptionRename struct {
	NewName     string
	Prefix      string
	Suffix      string
	Replace     string
	With        string
	StartNumber int
	Padding     int
}

// RenameBatch renomme un tableau de fichiers/dossiers selon les options
func RenameBatch(paths []string, opts OptionRename) error {
	for i, oldPath := range paths {
		dir := filepath.Dir(oldPath)
		base := filepath.Base(oldPath)
		ext := filepath.Ext(base)
		name := strings.TrimSuffix(base, ext)

		// Renommer complètement si NewName est défini
		if opts.NewName != "" {
			// ajouter une numérotation automatique
			if len(paths) > 1 {
				padding := opts.Padding
				if padding == 0 {
					padding = 3 // padding par défaut
				}
				num := opts.StartNumber + i
				name = fmt.Sprintf("%s_%0*d", opts.NewName, padding, num)
			} else {
				name = opts.NewName
			}
		} else {
			// Remplacement texte
			if opts.Replace != "" {
				name = strings.ReplaceAll(name, opts.Replace, opts.With)
			}
			// Numérotation
			if opts.StartNumber >= 0 {
				num := opts.StartNumber + i
				name = fmt.Sprintf("%s%0*d", opts.Prefix, opts.Padding, num)
			} else {
				// Ajouter préfixe/suffixe si pas de numérotation
				name = opts.Prefix + name + opts.Suffix
			}
		}

		// Toujours garder l'extension
		name += ext
		newPath := filepath.Join(dir, name)

		// Vérifier si le fichier de destination existe déjà
		if newPath != oldPath {
			if _, err := os.Stat(newPath); err == nil {
				return fmt.Errorf("le fichier %s existe déjà", newPath)
			}
		}

		if err := os.Rename(oldPath, newPath); err != nil {
			return fmt.Errorf("erreur sur %s: %v", oldPath, err)
		}
	}
	return nil
}
