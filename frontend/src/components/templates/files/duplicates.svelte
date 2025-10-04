<script lang="ts">
  import {
    Search,
    Delete,
  } from "../../../../wailsjs/go/services/DuplicateService";
  import { SelectFolder } from "../../../../wailsjs/go/main/App";

  let folder = "";
  let duplicates: Record<string, string[]> = {};
  let deleted: string[] = [];
  let loading = false;
  let error = "";
  let showConfirmModal = false;

  async function selectFolder() {
    try {
      const selected = await SelectFolder("Sélectionnez un dossier");
      if (selected) folder = selected;
    } catch (e: any) {
      error = "Erreur lors de la sélection du dossier";
    }
  }

  async function searchDuplicates() {
    if (!folder) return;
    loading = true;
    error = "";
    duplicates = {};
    deleted = [];

    try {
      duplicates = await Search(folder);
    } catch (e: any) {
      error = e.message || "Erreur inconnue";
    } finally {
      loading = false;
    }
  }

  function openConfirmModal() {
    showConfirmModal = true;
  }

  function closeConfirmModal() {
    showConfirmModal = false;
  }

  async function confirmDelete() {
    showConfirmModal = false;
    loading = true;
    error = "";
    deleted = [];

    try {
      deleted = await Delete(duplicates);
      for (let hash in duplicates) {
        duplicates[hash] = duplicates[hash].filter((f) => !deleted.includes(f));
        if (duplicates[hash].length <= 1) delete duplicates[hash];
      }
    } catch (e: any) {
      error = e.message || "Erreur lors de la suppression";
    } finally {
      loading = false;
    }
  }

  $: duplicateCount = Object.keys(duplicates).length;
  $: filesToDelete = Object.values(duplicates).reduce((acc, files) => acc + files.length - 1, 0);
</script>

<div class="min-h-screen p-4 max-w-6xl mx-auto">
  <!-- Header -->
  <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-8 pb-6 border-b border-slate-600">
    <div class="flex items-center gap-4 mb-4 sm:mb-0">
      <div class="w-14 h-14 bg-gradient-to-br from-indigo-400 to-indigo-600 rounded-2xl flex items-center justify-center shadow-lg shadow-indigo-500/30">
        <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-white">
          <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/>
          <rect x="8" y="2" width="8" height="4" rx="1" ry="1"/>
        </svg>
      </div>
      <div>
        <h1 class="text-3xl font-bold bg-gradient-to-r from-indigo-400 to-indigo-600 bg-clip-text text-transparent">
          Recherche de doublons
        </h1>
        <p class="text-sm text-slate-400 mt-1">Détectez et supprimez les fichiers en double</p>
      </div>
    </div>
    
    {#if duplicateCount > 0}
      <div class="flex items-center gap-2 px-4 py-2 bg-slate-800 border border-slate-600 rounded-xl">
        <span class="px-2.5 py-1 bg-gradient-to-r from-indigo-400 to-indigo-600 text-white rounded-lg text-sm font-bold">
          {duplicateCount}
        </span>
        <span class="text-sm text-slate-400">groupe{duplicateCount > 1 ? 's' : ''}</span>
      </div>
    {/if}
  </div>

  <div class="space-y-6">
    <!-- Folder Selection -->
    <div class="bg-slate-800 border border-slate-700 rounded-2xl p-6 transition-all duration-300 hover:border-indigo-500 hover:shadow-lg hover:shadow-indigo-500/10">
      <h3 class="text-lg font-semibold text-slate-100 mb-5 flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-indigo-400">
          <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
        </svg>
        Dossier à analyser
      </h3>
      
      <div class="flex gap-3">
        <div class="flex-1 relative">
          <input
            type="text"
            placeholder="Aucun dossier sélectionné..."
            bind:value={folder}
            class="w-full px-4 py-3 bg-slate-900 border border-slate-600 rounded-xl text-slate-100 text-sm transition-all duration-200 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 placeholder:text-slate-500"
            readonly
          />
          {#if folder}
            <div class="absolute right-3 top-1/2 -translate-y-1/2 w-2 h-2 bg-emerald-500 rounded-full animate-pulse"></div>
          {/if}
        </div>
        <button
          on:click={selectFolder}
          class="px-5 py-3 bg-slate-900 border border-slate-600 text-slate-100 rounded-xl font-medium text-sm transition-all duration-200 hover:border-indigo-500 hover:bg-slate-800 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
          disabled={loading}
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
          </svg>
          Parcourir
        </button>
      </div>
    </div>

    <!-- Action Buttons -->
    <div class="flex gap-3">
      <button
        on:click={searchDuplicates}
        disabled={!folder || loading}
        class="flex-1 px-8 py-4 bg-gradient-to-r from-indigo-400 to-indigo-600 text-white rounded-xl font-semibold text-base transition-all duration-300 shadow-lg shadow-indigo-500/30 hover:shadow-xl hover:shadow-indigo-500/40 hover:-translate-y-0.5 active:translate-y-0 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:translate-y-0 disabled:hover:shadow-lg flex items-center justify-center gap-3"
      >
        {#if loading && !deleted.length}
          <svg class="animate-spin w-5 h-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Analyse en cours...
        {:else}
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"/>
            <path d="m21 21-4.35-4.35"/>
          </svg>
          Rechercher les doublons
        {/if}
      </button>

      <button
        on:click={openConfirmModal}
        disabled={duplicateCount === 0 || loading}
        class="px-6 py-4 bg-red-600 text-white rounded-xl font-semibold text-base transition-all duration-300 shadow-lg shadow-red-500/30 hover:shadow-xl hover:shadow-red-500/40 hover:-translate-y-0.5 active:translate-y-0 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:translate-y-0 disabled:hover:shadow-lg flex items-center justify-center gap-2"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="3 6 5 6 21 6"/>
          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
        </svg>
        Supprimer
      </button>
    </div>

    <!-- Error Message -->
    {#if error}
      <div class="px-5 py-4 rounded-xl flex items-center gap-3 text-sm font-medium animate-in slide-in-from-top-2 fade-in duration-300 bg-red-500/10 border border-red-500 text-red-400">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10"/>
          <line x1="12" y1="8" x2="12" y2="12"/>
          <line x1="12" y1="16" x2="12.01" y2="16"/>
        </svg>
        <span>{error}</span>
      </div>
    {/if}

    <!-- Duplicates Results -->
    {#if duplicateCount > 0}
      <div class="bg-slate-800 border border-slate-700 rounded-2xl p-6 transition-all duration-300 hover:border-indigo-500 hover:shadow-lg hover:shadow-indigo-500/10">
        <div class="flex items-center justify-between mb-5">
          <h3 class="text-lg font-semibold text-slate-100 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-indigo-400">
              <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/>
              <rect x="8" y="2" width="8" height="4" rx="1" ry="1"/>
            </svg>
            Doublons détectés
          </h3>
          <span class="text-sm text-slate-400">
            {duplicateCount} groupe{duplicateCount > 1 ? 's' : ''} trouvé{duplicateCount > 1 ? 's' : ''}
          </span>
        </div>

        <div class="space-y-4 max-h-[500px] overflow-y-auto scrollbar-thin scrollbar-thumb-slate-600 scrollbar-track-slate-800">
          {#each Object.entries(duplicates) as [hash, files]}
            <div class="bg-slate-900 border border-slate-600 rounded-xl p-4 transition-all duration-200 hover:border-indigo-500/50">
              <div class="flex items-center gap-2 mb-3">
                <div class="h-px flex-1 bg-gradient-to-r from-transparent via-slate-600 to-transparent"></div>
                <span class="text-xs text-slate-500 font-mono px-3 py-1 bg-slate-800 rounded-full border border-slate-600">
                  {hash.substring(0, 16)}...
                </span>
                <div class="h-px flex-1 bg-gradient-to-r from-transparent via-slate-600 to-transparent"></div>
              </div>

              <div class="space-y-2">
                {#each files as file, i}
                  <div class="p-3 bg-slate-800 border border-slate-700 rounded-lg transition-all duration-200 hover:border-indigo-500/30 hover:translate-x-1">
                    <div class="flex items-center gap-3">
                      <span
                        class="inline-flex items-center px-2.5 py-1 rounded-md text-xs font-semibold uppercase tracking-wider border {i === 0
                          ? 'bg-emerald-500/15 text-emerald-400 border-emerald-500/30'
                          : 'bg-red-500/15 text-red-400 border-red-500/30'}"
                      >
                        {#if i === 0}
                          <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-1">
                            <polyline points="20 6 9 17 4 12"/>
                          </svg>
                          Original
                        {:else}
                          <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-1">
                            <line x1="18" y1="6" x2="6" y2="18"/>
                            <line x1="6" y1="6" x2="18" y2="18"/>
                          </svg>
                          Doublon
                        {/if}
                      </span>
                      <div class="flex-1 min-w-0">
                        <p class="text-sm text-slate-400 truncate">{file}</p>
                      </div>
                    </div>
                  </div>
                {/each}
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <!-- Success Message -->
    {#if deleted.length > 0}
      <div class="px-5 py-4 rounded-xl flex items-start gap-3 text-sm font-medium animate-in slide-in-from-top-2 fade-in duration-300 bg-emerald-500/10 border border-emerald-500 text-emerald-400">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="flex-shrink-0 mt-0.5">
          <polyline points="20 6 9 17 4 12"/>
        </svg>
        <div class="flex-1">
          <p class="font-semibold mb-2">{deleted.length} fichier{deleted.length > 1 ? 's' : ''} supprimé{deleted.length > 1 ? 's' : ''} avec succès</p>
          <div class="space-y-1 max-h-32 overflow-y-auto scrollbar-thin scrollbar-thumb-emerald-500/30 scrollbar-track-transparent">
            {#each deleted as file}
              <p class="text-xs text-emerald-400/80">✓ {file}</p>
            {/each}
          </div>
        </div>
      </div>
    {/if}
  </div>
</div>

<!-- Confirmation Modal -->
{#if showConfirmModal}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div 
    class="fixed inset-0 flex items-center justify-center p-4 animate-in fade-in duration-200"
    style="background-color: rgba(0, 0, 0, 0.7); z-index: 9999;"
    on:click={closeConfirmModal}
  >
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div 
      class="bg-slate-800 border border-slate-700 rounded-2xl p-6 max-w-md w-full shadow-2xl animate-in zoom-in-95 duration-200" 
      on:click|stopPropagation
    >
      <div class="flex items-start gap-4 mb-4">
        <div class="w-12 h-12 rounded-xl bg-red-500/20 flex items-center justify-center flex-shrink-0">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-red-400">
            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
            <line x1="12" y1="9" x2="12" y2="13"/>
            <line x1="12" y1="17" x2="12.01" y2="17"/>
          </svg>
        </div>
        <div>
          <h3 class="text-lg font-bold text-slate-100 mb-2">
            Confirmer la suppression
          </h3>
          <p class="text-sm text-slate-400">
            Vous êtes sur le point de supprimer <span class="font-semibold text-red-400">{filesToDelete} fichier{filesToDelete > 1 ? 's' : ''}</span> en double.
            Cette action est <span class="font-semibold">irréversible</span>.
          </p>
        </div>
      </div>
      
      <div class="flex gap-3 mt-6">
        <button
          on:click={closeConfirmModal}
          class="flex-1 px-4 py-3 bg-slate-900 border border-slate-600 text-slate-100 rounded-xl font-medium text-sm transition-all duration-200 hover:border-slate-500 hover:bg-slate-800"
        >
          Annuler
        </button>
        <button
          on:click={confirmDelete}
          class="flex-1 px-4 py-3 bg-red-600 text-white rounded-xl font-semibold text-sm transition-all duration-200 hover:bg-red-700 shadow-lg shadow-red-500/30 hover:shadow-xl hover:shadow-red-500/40"
        >
          Supprimer
        </button>
      </div>
    </div>
  </div>
{/if}

<style>

  /* Animations */
  @keyframes slide-in-from-top {
    from {
      opacity: 0;
      transform: translateY(-10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @keyframes fade-in {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  @keyframes zoom-in {
    from {
      opacity: 0;
      transform: scale(0.95);
    }
    to {
      opacity: 1;
      transform: scale(1);
    }
  }

  .animate-in {
    animation: slide-in-from-top 0.3s ease-out;
  }

  .zoom-in-95 {
    animation: zoom-in 0.2s ease-out;
  }

  .fade-in {
    animation: fade-in 0.2s ease-out;
  }
</style>