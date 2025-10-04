<script lang="ts">
  import { Rename } from "../../../../wailsjs/go/services/RenameService";
  import { writable } from "svelte/store";
  import FileSelector from "../../organismes/FileDropzone.svelte";
  import FilePreview from "../../organismes/FilePreview.svelte";
  
  let files: string[] = [];
  let newName = "";
  let prefix = "";
  let suffix = "";
  let replace = "";
  let withText = "";
  let startNumber = 1;
  let padding = 3;
  const message = writable("");
  const messageType = writable<"success" | "error" | "">("");

  function handleFiles(e: CustomEvent<{ files: string[] }>) {
    files = [...files, ...e.detail.files];
    message.set("");
  }

  function removeFile(e: CustomEvent<{ index: number }>) {
    const index = e.detail.index;
    if (index >= 0 && index < files.length) {
      files = [...files.slice(0, index), ...files.slice(index + 1)];
    }
  }

  function clearFiles() {
    files = [];
    message.set("");
  }

  async function rename() {
    if (files.length === 0) {
      message.set("Veuillez sélectionner des fichiers !");
      messageType.set("error");
      return;
    }
    const options = {
      NewName: newName,
      Prefix: prefix,
      Suffix: suffix,
      Replace: replace,
      With: withText,
      StartNumber: startNumber,
      Padding: padding,
    };
    try {
      await Rename(files, options);
      message.set(`${files.length} fichier${files.length > 1 ? 's' : ''} renommé${files.length > 1 ? 's' : ''} avec succès !`);
      messageType.set("success");
      files = [];
      newName = "";
      prefix = "";
      suffix = "";
      replace = "";
      withText = "";
    } catch (err) {
      console.error(err);
      message.set("Erreur: " + err);
      messageType.set("error");
    }
  }
</script>

<div class="min-h-screen p-4 max-w-6xl mx-auto">
  <!-- Header -->
  <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-8 pb-6 border-b border-slate-600">
    <div class="flex items-center gap-4 mb-4 sm:mb-0">
      <div class="w-14 h-14 bg-gradient-to-br from-indigo-400 to-indigo-600 rounded-2xl flex items-center justify-center shadow-lg shadow-indigo-500/30">
        <span class="material-icons-outlined text-white text-3xl">
          edit
        </span>
      </div>
      <div>
        <h1 class="text-3xl font-bold bg-gradient-to-r from-indigo-400 to-indigo-600 bg-clip-text text-transparent">
          Renommage de fichiers
        </h1>
        <p class="text-sm text-slate-400 mt-1">Renommez vos fichiers en masse avec précision</p>
      </div>
    </div>
    
    {#if files.length > 0}
      <div class="flex items-center gap-2 px-4 py-2 bg-slate-800 border border-slate-600 rounded-xl">
        <span class="px-2.5 py-1 bg-gradient-to-r from-indigo-400 to-indigo-600 text-white rounded-lg text-sm font-bold">
          {files.length}
        </span>
        <span class="text-sm text-slate-400">fichier{files.length > 1 ? 's' : ''}</span>
      </div>
    {/if}
  </div>

  <div class="space-y-6">
    <!-- File Selector -->
    <div class="bg-slate-800 border border-slate-700 rounded-2xl p-6 transition-all duration-300 hover:border-indigo-500 hover:shadow-lg hover:shadow-indigo-500/10">
      <FileSelector
        accept="*.*"
        multiple={true}
        placeholder="Glissez-déposez ou cliquez pour sélectionner"
        on:files={handleFiles}
      />
    </div>

    <!-- File Preview -->
    {#if files.length > 0}
      <div class="bg-slate-800 border border-slate-700 rounded-2xl p-6 transition-all duration-300 hover:border-indigo-500 hover:shadow-lg hover:shadow-indigo-500/10">
        <FilePreview
          {files}
          viewMode="list"
          on:remove={removeFile}
          on:clearAll={clearFiles}
        />
      </div>
    {/if}

    <!-- Options -->
    <div class="bg-slate-800 border border-slate-700 rounded-2xl p-6 transition-all duration-300 hover:border-indigo-500 hover:shadow-lg hover:shadow-indigo-500/10">
      <h3 class="text-lg font-semibold text-slate-100 mb-5 flex items-center gap-2">
        Options de renommage
      </h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
        <!-- New Name (Prioritaire) -->
        <div class="flex flex-col gap-2 md:col-span-2">
          <label for="newName" class="text-sm font-medium text-slate-400 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-indigo-400">
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
            </svg>
            Nouveau nom complet
            <span class="text-xs text-slate-500">(prioritaire - remplace toutes les autres options)</span>
          </label>
          <input
            id="newName"
            type="text"
            bind:value={newName}
            placeholder="Nom complet sans extension (ex: document_final)..."
            class="w-full px-4 py-3 bg-slate-900 border border-slate-600 rounded-xl text-slate-100 text-sm transition-all duration-200 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 placeholder:text-slate-500"
          />
        </div>

        <div class="md:col-span-2 border-t border-slate-700 pt-5 mt-2">
          <p class="text-xs text-slate-500 mb-5 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="16" x2="12" y2="12"/>
              <line x1="12" y1="8" x2="12.01" y2="8"/>
            </svg>
            Options avancées (ignorées si "Nouveau nom complet" est renseigné)
          </p>
        </div>

        <!-- Prefix -->
        <div class="flex flex-col gap-2">
          <label for="prefix" class="text-sm font-medium text-slate-400 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-indigo-400">
              <polyline points="9 18 15 12 9 6"/>
            </svg>
            Préfixe
          </label>
          <input
            id="prefix"
            type="text"
            bind:value={prefix}
            placeholder="Ajouter avant le nom..."
            disabled={newName !== ""}
            class="w-full px-4 py-3 bg-slate-900 border border-slate-600 rounded-xl text-slate-100 text-sm transition-all duration-200 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 placeholder:text-slate-500 disabled:opacity-50 disabled:cursor-not-allowed"
          />
        </div>

        <!-- Suffix -->
        <div class="flex flex-col gap-2">
          <label for="suffix" class="text-sm font-medium text-slate-400 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-indigo-400">
              <polyline points="15 18 9 12 15 6"/>
            </svg>
            Suffixe
          </label>
          <input
            id="suffix"
            type="text"
            bind:value={suffix}
            placeholder="Ajouter après le nom..."
            class="w-full px-4 py-3 bg-slate-900 border border-slate-600 rounded-xl text-slate-100 text-sm transition-all duration-200 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 placeholder:text-slate-500"
          />
        </div>

        <!-- Replace -->
        <div class="flex flex-col gap-2">
          <label for="replace" class="text-sm font-medium text-slate-400 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-indigo-400">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
            Remplacer
          </label>
          <input
            id="replace"
            type="text"
            bind:value={replace}
            placeholder="Texte à remplacer..."
            class="w-full px-4 py-3 bg-slate-900 border border-slate-600 rounded-xl text-slate-100 text-sm transition-all duration-200 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 placeholder:text-slate-500"
          />
        </div>

        <!-- With -->
        <div class="flex flex-col gap-2">
          <label for="withText" class="text-sm font-medium text-slate-400 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-indigo-400">
              <polyline points="20 6 9 17 4 12"/>
            </svg>
            Par
          </label>
          <input
            id="withText"
            type="text"
            bind:value={withText}
            placeholder="Nouveau texte..."
            class="w-full px-4 py-3 bg-slate-900 border border-slate-600 rounded-xl text-slate-100 text-sm transition-all duration-200 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 placeholder:text-slate-500"
          />
        </div>

        <!-- Start Number -->
        <div class="flex flex-col gap-2">
          <label for="startNumber" class="text-sm font-medium text-slate-400 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-indigo-400">
              <line x1="12" y1="1" x2="12" y2="23"/>
              <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/>
            </svg>
            Numéro de départ
          </label>
          <input
            id="startNumber"
            type="number"
            bind:value={startNumber}
            min="0"
            class="w-full px-4 py-3 bg-slate-900 border border-slate-600 rounded-xl text-slate-100 text-sm transition-all duration-200 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20"
          />
        </div>

        <!-- Padding -->
        <div class="flex flex-col gap-2">
          <label for="padding" class="text-sm font-medium text-slate-400 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-indigo-400">
              <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
              <line x1="9" y1="9" x2="15" y2="9"/>
              <line x1="9" y1="15" x2="15" y2="15"/>
            </svg>
            Zéros de padding
          </label>
          <input
            id="padding"
            type="number"
            bind:value={padding}
            min="1"
            max="10"
            class="w-full px-4 py-3 bg-slate-900 border border-slate-600 rounded-xl text-slate-100 text-sm transition-all duration-200 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20"
          />
        </div>
      </div>
    </div>

    <!-- Action Button -->
    <div class="flex flex-col gap-4">
      <button 
        on:click={rename} 
        disabled={files.length === 0}
        class="w-full px-8 py-4 bg-gradient-to-r from-indigo-400 to-indigo-600 text-white rounded-xl font-semibold text-base transition-all duration-300 shadow-lg shadow-indigo-500/30 hover:shadow-xl hover:shadow-indigo-500/40 hover:-translate-y-0.5 active:translate-y-0 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:translate-y-0 disabled:hover:shadow-lg flex items-center justify-center gap-3"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
          <polyline points="7 10 12 15 17 10"/>
          <line x1="12" y1="15" x2="12" y2="3"/>
        </svg>
        Renommer les fichiers
      </button>

      {#if $message}
        <div class="px-5 py-4 rounded-xl flex items-center gap-3 text-sm font-medium animate-in slide-in-from-top-2 fade-in duration-300 {$messageType === 'success' ? 'bg-emerald-500/10 border border-emerald-500 text-emerald-400' : 'bg-red-500/10 border border-red-500 text-red-400'}">
          {#if $messageType === "success"}
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="20 6 9 17 4 12"/>
            </svg>
          {:else}
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
          {/if}
          <span>{$message}</span>
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  /* Reset number input appearance */
  input[type="number"]::-webkit-outer-spin-button,
  input[type="number"]::-webkit-inner-spin-button {
    -webkit-appearance: none;
    margin: 0;
  }

  input[type="number"] {
    -moz-appearance: textfield;
  }

  /* Animation for message */
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

  .animate-in {
    animation: slide-in-from-top 0.3s ease-out;
  }
</style>