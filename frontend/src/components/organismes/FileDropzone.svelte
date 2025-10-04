<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { SelectFiles, SaveDroppedFiles } from '../../../wailsjs/go/main/App';

  export let accept = '*/*'; // ex: "*.png;*.jpg;*.jpeg" ou "*.*" ou "image/*"
  export let multiple = true;
  export let placeholder = 'Glissez-déposez ou cliquez pour sélectionner';
  export let disabled = false;

  const dispatch = createEventDispatcher();
  let inputEl: HTMLInputElement;
  let isDragging = false;

  async function openFilePicker() {
    if (disabled) return;

    try {
      const filters = [
        {
          DisplayName: "Fichiers",
          Pattern: accept
        }
      ];

      const result = await SelectFiles("Sélectionnez vos fichiers", filters, multiple);
      if (!result || result.length === 0) return;

      dispatch('files', { files: result });
    } catch (err) {
      console.error("Erreur lors de la sélection de fichiers:", err);
    }
  }

  async function handleFiles(files: FileList | null) {
    if (!files || disabled) return;
    
    try {
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

      const savedPaths = await SaveDroppedFiles(filesData);
      
      if (savedPaths.length > 0) {
        dispatch('files', { files: savedPaths });
      }
    } catch (error) {
      console.error("Erreur lors du traitement des fichiers:", error);
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
