<script context="module" lang="ts">
    export interface ConvertedFile {
        originalFile: File;
        convertedBlob: Blob;
        convertedName: string;
        format: string;
        quality?: number;
        originalSize: number;
        convertedSize: number;
        previewUrl: string;
        base64Data?: string;
    }
</script>

<script lang="ts">
    export let convertedFiles: ConvertedFile[] = [];
    export let isConverting = false;

    import { createEventDispatcher } from "svelte";
    const dispatch = createEventDispatcher();

    const download = (item: ConvertedFile) => {
        const url = URL.createObjectURL(item.convertedBlob);
        const a = document.createElement("a");
        a.href = url;
        a.download = item.convertedName;
        a.click();
        URL.revokeObjectURL(url);
    };

    const remove = (index: number) => dispatch("removeConverted", { index });
    const clearAll = () => dispatch("clearAllConverted");

    const getImageUrl = (item: ConvertedFile) => {
        if (item.base64Data) {
            return `data:image/${item.format.toLowerCase()};base64,${item.base64Data}`;
        }
        return item.previewUrl;
    };
</script>

<!-- Loading State -->
{#if isConverting}
    <div class="bg-slate-800/50 backdrop-blur-sm border border-slate-600 rounded-xl p-6 mb-6 mt-5">
        <div class="flex items-center justify-center space-x-3">
            <div class="animate-spin rounded-full h-6 w-6 border-2 border-indigo-400 border-t-transparent"></div>
            <p class="text-slate-300 font-medium">Conversion en cours...</p>
        </div>
    </div>
{/if}

<!-- Converted Files Display -->
{#if convertedFiles.length > 0}
    <div class="bg-slate-800/30 backdrop-blur-sm border border-slate-600 rounded-xl p-6 shadow-lg mt-5">
        <!-- Header -->
        <div class="flex items-center justify-between mb-6">
            <div class="flex items-center space-x-3">
                <div class="w-3 h-3 bg-emerald-500 rounded-full animate-pulse"></div>
                <h3 class="text-xl font-bold text-slate-100">
                    Images converties 
                    <span class="inline-flex items-center justify-center px-2 py-1 ml-2 text-xs font-bold text-indigo-100 bg-indigo-600 rounded-full">
                        {convertedFiles.length}
                    </span>
                </h3>
            </div>
            
            <button 
                on:click={clearAll}
                class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white font-medium rounded-lg 
                       transition-all duration-200 hover:shadow-lg hover:shadow-red-500/25 
                       focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 focus:ring-offset-slate-800"
            >
                <svg class="w-4 h-4 mr-2 inline" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                </svg>
                Tout supprimer
            </button>
        </div>

        <!-- Files Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {#each convertedFiles as item, i}
                <div class="bg-slate-700/50 backdrop-blur-sm border border-slate-600 rounded-lg p-4 
                           hover:border-indigo-500 hover:shadow-lg hover:shadow-indigo-500/10 
                           transition-all duration-300 group">
                    
                    <!-- Image Preview -->
                    <div class="relative mb-4 overflow-hidden rounded-lg">
                        <img 
                            src={getImageUrl(item)} 
                            alt={item.convertedName}
                            class="w-full h-100 object-cover transition-transform duration-300 
                                   group-hover:scale-105 border border-slate-600"
                        />
                        <div class="absolute top-2 right-2 bg-slate-900/80 backdrop-blur-sm px-2 py-1 rounded-md">
                            <span class="text-xs font-medium text-indigo-400 uppercase">{item.format}</span>
                        </div>
                    </div>

                    <!-- File Info -->
                    <div class="mb-4">
                        <h4 class="text-sm font-semibold text-slate-100 truncate mb-2" title={item.convertedName}>
                            {item.convertedName}
                        </h4>
                        
                        <div class="flex items-center justify-between text-xs text-slate-400">
                            <span class="flex items-center">
                                <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
                                    <path fill-rule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clip-rule="evenodd" />
                                </svg>
                                {Math.round(item.convertedSize/1024)}KB
                            </span>
                            
                            <span class="px-2 py-0.5 bg-slate-600 rounded-full">
                                {item.format}
                            </span>
                        </div>
                    </div>

                    <!-- Action Buttons -->
                    <div class="flex space-x-2">
                        <button 
                            on:click={() => download(item)}
                            class="flex-1 px-3 py-2 bg-emerald-600 hover:bg-emerald-700 text-white text-sm font-medium 
                                   rounded-md transition-all duration-200 hover:shadow-lg hover:shadow-emerald-500/25
                                   focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2 focus:ring-offset-slate-700"
                        >
                            <svg class="w-4 h-4 mr-1 inline" fill="currentColor" viewBox="0 0 20 20">
                                <path fill-rule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clip-rule="evenodd" />
                            </svg>
                            Télécharger
                        </button>
                        
                        <button 
                            on:click={() => remove(i)}
                            class="px-3 py-2 bg-slate-600 hover:bg-red-600 text-slate-300 hover:text-white text-sm font-medium 
                                   rounded-md transition-all duration-200 hover:shadow-lg hover:shadow-red-500/25
                                   focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 focus:ring-offset-slate-700"
                            title="Supprimer"
                        >
                            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                                <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                            </svg>
                        </button>
                    </div>
                </div>
            {/each}
        </div>
    </div>
{/if}