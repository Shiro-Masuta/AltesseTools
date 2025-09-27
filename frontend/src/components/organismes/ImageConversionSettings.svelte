<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  // Types basés sur votre code Go
  type ImageFormat = 'PNG' | 'JPEG' | 'WEBP' | 'AVIF' | 'BMP' | 'TIFF';
  type PNGCompressionLevel = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9;
  type TIFFCompressionType = 'Deflate' | 'LZW' | 'None';

  interface ConversionOptions {
    format: ImageFormat;
    quality: number;
    lossless: boolean;
    pngLevel: PNGCompressionLevel;
    tiffCompress: TIFFCompressionType;
  }

  // Props avec bind
  export let options: ConversionOptions = {
    format: 'JPEG',
    quality: 85,
    lossless: false,
    pngLevel: 6,
    tiffCompress: 'Deflate'
  };

  // Dispatcher pour les événements
  const dispatch = createEventDispatcher();

  // Formats supportés
  const supportedFormats: ImageFormat[] = ['PNG', 'JPEG', 'WEBP', 'AVIF', 'BMP', 'TIFF'];
  
  // Niveaux de compression PNG (0-9)
  const pngLevels: PNGCompressionLevel[] = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];
  
  // Types de compression TIFF
  const tiffCompressionTypes: TIFFCompressionType[] = ['None', 'Deflate', 'LZW'];

  // Vérifie si le format supporte la qualité
  $: supportsQuality = ['JPEG', 'WEBP', 'AVIF'].includes(options.format);
  
  // Vérifie si le format supporte le mode lossless
  $: supportsLossless = ['WEBP', 'AVIF'].includes(options.format);
  
  // Vérifie si c'est PNG pour afficher les options spécifiques
  $: isPNG = options.format === 'PNG';
  
  // Vérifie si c'est TIFF pour afficher les options spécifiques
  $: isTIFF = options.format === 'TIFF';

  // Fonction pour émettre les changements
  function handleChange() {
    dispatch('optionsChange', options);
  }

  // Réagir aux changements d'options
  $: if (options) {
    handleChange();
  }
</script>

<!-- Le reste du template reste identique -->
<div class="bg-[var(--surface)] border border-[var(--border)] rounded-lg p-6 mt-10 space-y-6">
  <div class="flex items-center gap-3 mb-6">
    <div class="w-3 h-3 rounded-full bg-[var(--accent)]"></div>
    <h3 class="text-lg font-semibold text-[var(--text-primary)]">
      Paramètres de conversion
    </h3>
  </div>

  <!-- Sélection du format -->
  <div class="space-y-3">
    <!-- svelte-ignore a11y-label-has-associated-control -->
    <label class="block text-sm font-medium text-[var(--text-secondary)]">
      Format de sortie
    </label>
    <select 
      bind:value={options.format}
      class="w-full bg-[var(--surface-alt)] border border-[var(--border)] rounded-md px-3 py-2 text-[var(--text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] focus:border-transparent"
    >
      {#each supportedFormats as format}
        <option value={format}>{format}</option>
      {/each}
    </select>
  </div>

  <!-- Qualité (pour JPEG, WebP, AVIF) -->
  {#if supportsQuality}
    <div class="space-y-3">
      <div class="flex justify-between items-center">
        <!-- svelte-ignore a11y-label-has-associated-control -->
        <label class="text-sm font-medium text-[var(--text-secondary)]">
          Qualité
        </label>
        <span class="text-sm text-[var(--accent)] font-mono">
          {options.quality}%
        </span>
      </div>
      <input 
        type="range"
        min="1"
        max="100"
        bind:value={options.quality}
        class="w-full h-2 bg-[var(--surface-alt)] rounded-lg appearance-none cursor-pointer slider"
      />
      <div class="flex justify-between text-xs text-[var(--text-secondary)]">
        <span>1%</span>
        <span>100%</span>
      </div>
    </div>
  {/if}

  <!-- Mode Lossless (pour WebP, AVIF) -->
  {#if supportsLossless}
    <div class="flex items-center justify-between">
      <div class="space-y-1">
        <!-- svelte-ignore a11y-label-has-associated-control -->
        <label class="text-sm font-medium text-[var(--text-secondary)]">
          Mode sans perte
        </label>
        <p class="text-xs text-[var(--muted)]">
          Compression sans perte de qualité
        </p>
      </div>
      <label class="relative inline-flex items-center cursor-pointer">
        <input 
          type="checkbox" 
          bind:checked={options.lossless}
          class="sr-only peer"
        />
        <div class="w-11 h-6 bg-[var(--surface-alt)] peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[var(--accent)]"></div>
      </label>
    </div>
  {/if}

  <!-- Niveau de compression PNG -->
  {#if isPNG}
    <div class="space-y-3">
      <div class="flex justify-between items-center">
        <!-- svelte-ignore a11y-label-has-associated-control -->
        <label class="text-sm font-medium text-[var(--text-secondary)]">
          Niveau de compression PNG
        </label>
        <span class="text-sm text-[var(--accent)] font-mono">
          {options.pngLevel}
        </span>
      </div>
      <select 
        bind:value={options.pngLevel}
        class="w-full bg-[var(--surface-alt)] border border-[var(--border)] rounded-md px-3 py-2 text-[var(--text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] focus:border-transparent"
      >
        {#each pngLevels as level}
          <option value={level}>
            Niveau {level} {level === 0 ? '(Aucune)' : level === 6 ? '(Défaut)' : level === 9 ? '(Maximum)' : ''}
          </option>
        {/each}
      </select>
      <p class="text-xs text-[var(--muted)]">
        0 = Aucune compression, 9 = Compression maximale
      </p>
    </div>
  {/if}

  <!-- Type de compression TIFF -->
  {#if isTIFF}
    <div class="space-y-3">
      <!-- svelte-ignore a11y-label-has-associated-control -->
      <label class="block text-sm font-medium text-[var(--text-secondary)]">
        Type de compression TIFF
      </label>
      <select 
        bind:value={options.tiffCompress}
        class="w-full bg-[var(--surface-alt)] border border-[var(--border)] rounded-md px-3 py-2 text-[var(--text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] focus:border-transparent"
      >
        {#each tiffCompressionTypes as compression}
          <option value={compression}>
            {compression} {compression === 'Deflate' ? '(Défaut)' : ''}
          </option>
        {/each}
      </select>
    </div>
  {/if}

  <!-- Résumé des paramètres -->
  <div class="bg-[var(--surface-alt)] rounded-md p-4 border-l-4 border-l-[var(--accent)]">
    <h4 class="text-sm font-medium text-[var(--text-primary)] mb-2">
      Configuration actuelle
    </h4>
    <div class="space-y-1 text-sm text-[var(--text-secondary)]">
      <div class="flex justify-between">
        <span>Format:</span>
        <span class="text-[var(--accent)]">{options.format}</span>
      </div>
      {#if supportsQuality}
        <div class="flex justify-between">
          <span>Qualité:</span>
          <span class="text-[var(--accent)]">{options.quality}%</span>
        </div>
      {/if}
      {#if supportsLossless}
        <div class="flex justify-between">
          <span>Sans perte:</span>
          <span class="text-[var(--accent)]">{options.lossless ? 'Oui' : 'Non'}</span>
        </div>
      {/if}
      {#if isPNG}
        <div class="flex justify-between">
          <span>Compression PNG:</span>
          <span class="text-[var(--accent)]">Niveau {options.pngLevel}</span>
        </div>
      {/if}
      {#if isTIFF}
        <div class="flex justify-between">
          <span>Compression TIFF:</span>
          <span class="text-[var(--accent)]">{options.tiffCompress}</span>
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  .slider::-webkit-slider-thumb {
    appearance: none;
    height: 20px;
    width: 20px;
    border-radius: 50%;
    background: var(--accent);
    cursor: pointer;
    box-shadow: 0 0 10px rgba(129, 140, 248, 0.5);
  }

  .slider::-moz-range-thumb {
    height: 20px;
    width: 20px;
    border-radius: 50%;
    background: var(--accent);
    cursor: pointer;
    border: none;
    box-shadow: 0 0 10px rgba(129, 140, 248, 0.5);
  }
</style>