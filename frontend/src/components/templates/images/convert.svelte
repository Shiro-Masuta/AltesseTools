<script lang="ts">
  import FilePreview from "../../organismes/FilePreview.svelte";
  import FileDropzone from "../../organismes/FileDropzone.svelte";
  import ImageConversionSettings from "../../organismes/ImageConversionSettings.svelte";
  import FileConvertPreview, {
    type ConvertedFile,
  } from "../../organismes/FileConvertPreview.svelte";
  import { ConvertManyFromFilesToBase64Parallel } from "../../../../wailsjs/go/services/ConverterService";
  import { EventsOn, EventsOff } from "../../../../wailsjs/runtime/runtime";
  import { onMount, onDestroy } from "svelte";

  type ImageFormat = "PNG" | "JPEG" | "WEBP" | "AVIF" | "BMP" | "TIFF";
  type PNGCompressionLevel = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9;
  type TIFFCompressionType = "Deflate" | "LZW" | "None";

  interface ConversionOptions {
    format: ImageFormat;
    quality: number;
    lossless: boolean;
    pngLevel: PNGCompressionLevel;
    tiffCompress: TIFFCompressionType;
  }

  interface ProgressData {
    current: number;
    total: number;
    path: string;
  }

  let files: string[] = [];
  let convertedFiles: ConvertedFile[] = [];
  let isConverting = false;
  let convertedCount = 0;
  let totalCount = 0;

  const MAX_FILES = 10;
  const TIFF_MAP = { None: 0, Deflate: 1, LZW: 2 };
  const MIME_MAP: Record<ImageFormat, string> = {
    JPEG: "image/jpeg",
    PNG: "image/png",
    WEBP: "image/webp",
    AVIF: "image/avif",
    BMP: "image/bmp",
    TIFF: "image/tiff",
  };

  let options: ConversionOptions = {
    format: "JPEG",
    quality: 85,
    lossless: false,
    pngLevel: 6,
    tiffCompress: "Deflate",
  };

  // Gestion des événements Wails
  onMount(() => {
    EventsOn("conversion-progress", (data: ProgressData) => {
      convertedCount = data.current;
      totalCount = data.total;
    });
  });

  onDestroy(() => {
    EventsOff("conversion-progress");
  });

  const base64ToBlob = (base64: string, mimeType: string) => {
    const bytes = atob(base64);
    const array = new Uint8Array(bytes.length);
    for (let i = 0; i < bytes.length; i++) array[i] = bytes.charCodeAt(i);
    return new Blob([array], { type: mimeType });
  };

  const getFileName = (path: string) =>
    path.split("/").pop()?.split("\\").pop() || "unknown";

  const generateName = (path: string, format: string) => {
    const name = getFileName(path).split(".").slice(0, -1).join(".") || "image";
    return `${name}_converted.${format.toLowerCase()}`;
  };

  const handleFiles = (e: CustomEvent<{ files: string[] }>) => {
    const remaining = MAX_FILES - files.length;
    if (remaining <= 0)
      return alert(`Limite atteinte : maximum ${MAX_FILES} images`);

    const toAdd = e.detail.files.slice(0, remaining);
    files = [...files, ...toAdd];

    if (e.detail.files.length > remaining) {
      alert(`${remaining} fichier(s) ajouté(s). Limite ${MAX_FILES} atteinte.`);
    }
  };

  const removeFile = (e: CustomEvent<{ index: number }>) =>
    (files = files.filter((_, i) => i !== e.detail.index));

  const clearFiles = () => (files = []);

  const removeConverted = (e: CustomEvent<{ index: number }>) =>
    (convertedFiles = convertedFiles.filter((_, i) => i !== e.detail.index));

  const clearConverted = () => (convertedFiles = []);

  const convert = async () => {
    if (!files.length) return alert("Aucune image à convertir");

    isConverting = true;
    convertedCount = 0;
    totalCount = files.length;

    try {
      const goOptions = {
        Format: options.format.toLowerCase(),
        Quality: options.quality,
        Lossless: options.lossless,
        PNGLevel: options.pngLevel,
        TIFFCompress: TIFF_MAP[options.tiffCompress] || 1,
      };

      const results = await ConvertManyFromFilesToBase64Parallel(
        files,
        goOptions,
      );
      const mimeType = MIME_MAP[options.format];

      const converted = results.map((base64, i) => {
        const blob = base64ToBlob(base64, mimeType);
        return {
          originalFile: new File([], getFileName(files[i]), {
            type: "image/jpeg",
          }),
          convertedBlob: blob,
          convertedName: generateName(files[i], options.format),
          format: options.format,
          quality: ["JPEG", "WEBP", "AVIF"].includes(options.format)
            ? options.quality
            : undefined,
          originalSize: 0,
          convertedSize: blob.size,
          previewUrl: `data:${mimeType};base64,${base64}`,
          base64Data: base64,
        };
      });

      convertedFiles = [...convertedFiles, ...converted];
    } catch (error) {
      alert(`Erreur: ${error.message || error}`);
    } finally {
      isConverting = false;
      convertedCount = 0;
      totalCount = 0;
    }
  };

  $: isLimitReached = files.length >= MAX_FILES;
  $: placeholder = isLimitReached
    ? `Limite atteinte (${files.length}/${MAX_FILES})`
    : `Glissez-déposez ou cliquez (${files.length}/${MAX_FILES})`;
</script>

<div class="space-y-6 p-4 text-[var(--text-primary)]">
  <div
    class="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-8 pb-6 border-b border-slate-600"
  >
    <div class="flex items-center gap-4 mb-4 sm:mb-0">
      <div
        class="w-14 h-14 bg-gradient-to-br from-indigo-400 to-indigo-600 rounded-2xl flex items-center justify-center shadow-lg shadow-indigo-500/30"
      >
        <span class="material-icons-outlined text-white text-3xl"> sync_alt </span>
      </div>
      <div>
        <h1
          class="text-3xl font-bold bg-gradient-to-r from-indigo-400 to-indigo-600 bg-clip-text text-transparent"
        >
          Conversion d'image
        </h1>
        <p class="text-sm text-slate-400 mt-1">
          Convertisser vos images (le temps de traitement depend de
          la machine)
        </p>
      </div>
    </div>
  </div>

  <FileDropzone
    accept="*.png;*.jpg;*.jpeg;*.webp;*.avif;*.bmp;*.tiff;*.gif"
    multiple={true}
    placeholder="Sélectionnez ou déposez vos images"
    on:files={handleFiles}
  />

  <FilePreview
    {files}
    viewMode="grid"
    on:remove={removeFile}
    on:clearAll={clearFiles}
  />

  <ImageConversionSettings bind:options />

  {#if files.length > 0}
    <div class="flex justify-center mt-6">
      <button
        on:click={convert}
        disabled={isConverting}
        class="flex items-center gap-3 px-6 py-3 font-medium rounded-lg transition-all duration-300
               {isConverting
          ? 'bg-[var(--warning)] text-[var(--background)] cursor-wait'
          : 'bg-gradient-to-r from-[var(--accent)] to-[var(--accent-dark)] text-[var(--text-tiers)] hover:shadow-[var(--shadow-accent)] hover:scale-105 hover:bg-gradient-to-r hover:from-[var(--accent-light)] hover:to-[var(--accent)]'}"
      >
        {#if isConverting}
          <div
            class="animate-spin w-5 h-5 border-2 border-[var(--background)] border-t-transparent rounded-full"
          ></div>
          <span>Conversion en cours...</span>
        {:else}
          <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
            <path
              fill-rule="evenodd"
              d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z"
              clip-rule="evenodd"
            />
          </svg>
          <span
            >Convertir {files.length} image{files.length > 1 ? "s" : ""}</span
          >
        {/if}
      </button>
    </div>
  {/if}

  <!-- Barre de progression simple -->
  {#if isConverting}
    <div class="my-4">
      <p class="text-sm text-[var(--text-secondary)]">
        Conversion : {convertedCount} / {totalCount} fichier(s)
      </p>
      <div class="w-full h-3 bg-[var(--surface)] rounded-md overflow-hidden">
        <div
          class="h-full bg-[var(--accent)] transition-all duration-300"
          style="width: {(convertedCount / totalCount) * 100}%"
        ></div>
      </div>
    </div>
  {/if}

  <FileConvertPreview
    {convertedFiles}
    {isConverting}
    on:removeConverted={removeConverted}
    on:clearAllConverted={clearConverted}
  />
</div>
