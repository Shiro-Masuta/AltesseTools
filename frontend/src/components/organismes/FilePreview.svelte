<script lang="ts">
  import { createEventDispatcher, onMount } from "svelte";
  import { GetImagePreview } from "../../../wailsjs/go/services/ConverterService";
  
  export let files: string[] = [];
  export let viewMode: "grid" | "list" = "grid";
  
  const dispatch = createEventDispatcher();

  // Cache pour les aperçus d'images (seulement pour le mode grid)
  let imageCache: Record<string, string | null> = {};
  let loadingStates: Record<string, boolean> = {};

  const removeFile = (index: number) => dispatch("remove", { index });
  const clearAllFiles = () => dispatch("clearAll");

  const getName = (path: string) => 
    path.split('/').pop()?.split('\\').pop() || 'unknown';

  const getExtension = (path: string) =>
    getName(path).split('.').pop()?.toLowerCase() || '';

  const isImage = (path: string) => 
    ['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg', 'avif', 'tiff'].includes(getExtension(path));

  const getFileIcon = (path: string) => 
    isImage(path) ? 'image' : 'insert_drive_file';

  const loadPreview = async (filePath: string) => {
    // Si déjà en cache, retourner la valeur
    if (filePath in imageCache) {
      return imageCache[filePath];
    }

    // Si déjà en cours de chargement, attendre
    if (loadingStates[filePath]) {
      return null;
    }

    try {
      loadingStates[filePath] = true;
      const preview = await GetImagePreview(filePath);
      imageCache[filePath] = preview;
      return preview;
    } catch (error) {
      console.error(`Erreur chargement aperçu pour ${filePath}:`, error);
      imageCache[filePath] = null;
      return null;
    } finally {
      loadingStates[filePath] = false;
    }
  };

  // Précharger les images seulement en mode grid
  $: if (files.length > 0 && viewMode === "grid") {
    files.forEach(filePath => {
      if (isImage(filePath) && !(filePath in imageCache)) {
        loadPreview(filePath);
      }
    });
  }

  // Nettoyer le cache quand les fichiers changent
  $: {
    // Supprimer du cache les fichiers qui ne sont plus dans la liste
    const currentFiles = new Set(files);
    Object.keys(imageCache).forEach(path => {
      if (!currentFiles.has(path)) {
        delete imageCache[path];
        delete loadingStates[path];
      }
    });
  }
</script>

{#if files.length > 0}
  <!-- Header -->
  <div class="flex items-center justify-between mb-4 mt-6">
    <div class="flex items-center gap-2">
      <span class="material-icons text-lg" style="color: var(--accent);">collections</span>
      <span class="text-sm font-medium" style="color: var(--text-primary);">
        {files.length} fichier{files.length > 1 ? "s" : ""} sélectionné{files.length > 1 ? "s" : ""}
      </span>
    </div>

    <button
      on:click={clearAllFiles}
      class="flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm font-medium transition-all hover:scale-105"
      style="background: var(--error); color: var(--text-tiers); box-shadow: var(--shadow-accent);"
      title="Supprimer tous les fichiers"
    >
      <span class="material-icons text-sm">delete_sweep</span>
      Tout supprimer
    </button>
  </div>

  {#if viewMode === "grid"}
    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
      {#each files as filePath, index}
        <div class="relative group">
          {#if isImage(filePath)}
            <!-- Image avec aperçu ou placeholder -->
            <div class="w-full h-32 rounded-lg overflow-hidden" style="box-shadow: var(--shadow-accent);">
              {#if imageCache[filePath]}
                <!-- Image chargée -->
                <img
                  src={imageCache[filePath]}
                  alt={getName(filePath)}
                  class="w-full h-full object-cover"
                  on:error={() => {
                    console.error(`Erreur affichage image: ${filePath}`);
                    imageCache[filePath] = null; // Forcer le fallback
                  }}
                />
              {:else if loadingStates[filePath]}
                <!-- Chargement en cours -->
                <div 
                  class="w-full h-full flex flex-col items-center justify-center"
                  style="background: var(--surface-alt); border: 1px solid var(--border);"
                >
                  <div class="animate-spin w-6 h-6 border-2 border-accent border-t-transparent rounded-full mb-2"></div>
                  <span class="text-xs" style="color: var(--text-secondary);">Chargement...</span>
                </div>
              {:else}
                <!-- Placeholder si erreur ou pas encore chargé -->
                <div 
                  class="w-full h-full flex flex-col items-center justify-center border-2 border-dashed"
                  style="background: var(--surface-alt); border-color: var(--border);"
                >
                  <span class="material-icons text-2xl mb-1" style="color: var(--accent);">image</span>
                  <span class="text-xs" style="color: var(--text-secondary);">Aperçu</span>
                </div>
              {/if}
            </div>
          {:else}
            <!-- Fichier non-image -->
            <div 
              class="w-full h-32 border-2 border-dashed rounded-lg flex flex-col items-center justify-center"
              style="background: var(--surface-alt); border-color: var(--border);"
            >
              <span class="material-icons text-3xl mb-2" style="color: var(--text-secondary);">insert_drive_file</span>
              <span class="text-xs text-center px-2" style="color: var(--text-secondary);">
                {getExtension(filePath).toUpperCase()}
              </span>
            </div>
          {/if}

          <!-- Bouton suppression -->
          <button
            on:click={() => removeFile(index)}
            class="absolute top-2 right-2 w-6 h-6 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
            style="background: var(--error); color: var(--text-tiers);"
            title="Supprimer"
          >
            <span class="material-icons text-sm">close</span>
          </button>

          <!-- Nom du fichier -->
          <div class="absolute bottom-0 left-0 right-0 bg-black bg-opacity-70 text-white text-xs p-2 rounded-b-lg">
            <div class="truncate" title={getName(filePath)}>
              {getName(filePath)}
            </div>
          </div>
        </div>
      {/each}
    </div>
  {:else}
    <!-- Mode liste optimisé sans images -->
    <div
      class="space-y-2 max-h-96 overflow-y-auto p-2 rounded-lg"
      style="background: var(--surface); border: 1px solid var(--border);"
    >
      {#each files as filePath, index}
        <div
          class="flex items-center gap-3 p-3 rounded-lg transition-colors group cursor-pointer"
          style="background: var(--surface-alt);"
          on:mouseenter={(e) => {
            e.currentTarget.style.background = "var(--accent)";
            e.currentTarget.style.color = "var(--text-tiers)";
          }}
          on:mouseleave={(e) => {
            e.currentTarget.style.background = "var(--surface-alt)";
            e.currentTarget.style.color = "var(--text-primary)";
          }}
        >
          <!-- Index -->
          <span class="text-sm font-medium w-8 text-center" style="color: var(--text-secondary);">
            {index + 1}
          </span>

          <!-- Icône de fichier uniquement -->
          <div class="w-10 h-10 rounded flex items-center justify-center" style="background: var(--surface);">
            <span class="material-icons text-lg" style="color: var(--accent);">
              {getFileIcon(filePath)}
            </span>
          </div>

          <!-- Informations fichier -->
          <div class="flex-1 min-w-0">
            <div class="text-sm font-medium truncate" style="color: var(--text-primary);" title={filePath}>
              {getName(filePath)}
            </div>
            <div class="text-xs truncate" style="color: var(--text-secondary);" title={filePath}>
              {filePath}
            </div>
          </div>

          <!-- Extension -->
          <span class="text-xs px-2 py-1 rounded-full" 
                style="background: var(--accent); color: var(--text-tiers);">
            {getExtension(filePath).toUpperCase() || 'UNKNOWN'}
          </span>

          <!-- Bouton suppression -->
          <button
            on:click={() => removeFile(index)}
            class="w-8 h-8 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
            style="background: var(--error); color: var(--text-tiers);"
            title="Supprimer"
          >
            <span class="material-icons text-sm">close</span>
          </button>
        </div>
      {/each}
    </div>
  {/if}
{/if}