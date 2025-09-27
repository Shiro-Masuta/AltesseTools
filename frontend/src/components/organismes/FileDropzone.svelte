<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { SelectImagePaths, SaveDroppedFiles } from '../../../wailsjs/go/main/App';

  export let accept = '*/*';
  export let multiple = true;
  export let placeholder = 'Glissez-déposez ou cliquez pour sélectionner';
  export let disabled = false;

  const dispatch = createEventDispatcher();
  let inputEl: HTMLInputElement;
  let isDragging = false;

  async function openFilePicker() {
    if (disabled) return;

    try {
      const result = await SelectImagePaths();
      if (!result || result.length === 0) return;

      dispatch('files', { files: result });
    } catch (err) {
      console.error('Erreur lors de la sélection d\'images:', err);
    }
  }

  async function handleFiles(files: FileList | null) {
    if (!files || disabled) return;
    
    try {
      // Convertir les fichiers en données à envoyer au backend
      const filesData = await Promise.all(
        Array.from(files).map(async (file) => {
          const arrayBuffer = await file.arrayBuffer();
          return {
            name: file.name,
            content: Array.from(new Uint8Array(arrayBuffer)),
            size: file.size
          };
        })
      );

      // Utiliser le service Go pour sauvegarder et valider
      const validPaths = await SaveDroppedFiles(filesData);
      
      if (validPaths.length > 0) {
        dispatch('files', { files: validPaths });
      } else {
        alert('Aucun fichier image valide trouvé dans la sélection');
      }
    } catch (error) {
      console.error('Erreur lors du traitement des fichiers:', error);
      alert('Erreur lors du traitement des fichiers droppés');
    }
  }

  function onDrop(e: DragEvent) {
    e.preventDefault();
    isDragging = false;
    if (disabled) return;
    
    const files = e.dataTransfer ? e.dataTransfer.files : null;
    handleFiles(files);
  }

  function onInputChange(e: Event) {
    if (disabled) return;
    const target = e.target as HTMLInputElement;
    handleFiles(target.files);
  }

  function onDragOver(e: DragEvent) {
    e.preventDefault();
    if (!disabled) {
      isDragging = true;
    }
  }

  function onDragLeave(e: DragEvent) {
    e.preventDefault();
    isDragging = false;
  }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div
  class="flex flex-col items-center justify-center gap-3 p-6 border-2 border-dashed rounded-xl text-center transition-colors duration-150"
  class:cursor-pointer={!disabled}
  class:cursor-not-allowed={disabled}
  class:bg-gray-100={isDragging && !disabled}
  class:border-blue-400={isDragging && !disabled}
  class:hover:border-blue-400={!disabled}
  class:opacity-50={disabled}
  class:border-gray-300={!isDragging}
  on:click={openFilePicker}
  on:drop={onDrop}
  on:dragover={onDragOver}
  on:dragleave={onDragLeave}
>
  <input
    bind:this={inputEl}
    type="file"
    class="hidden"
    {accept}
    {multiple}
    {disabled}
    on:change={onInputChange}
  />
  <span class="material-icons text-4xl" class:text-gray-400={!disabled} class:text-gray-300={disabled}>
    upload_file
  </span>
  <span class="text-sm" class:text-gray-500={!disabled} class:text-gray-400={disabled}>
    {placeholder}
  </span>
</div>