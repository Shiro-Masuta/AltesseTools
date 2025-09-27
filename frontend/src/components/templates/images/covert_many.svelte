<script lang="ts">
    import FilePreview from "../../organismes/FilePreview.svelte";
    import FileDropzone from "../../organismes/FileDropzone.svelte";
    import ImageConversionSettings from "../../organismes/ImageConversionSettings.svelte";
    import { ConvertManyToFolder } from "../../../../wailsjs/go/services/ConverterService";
    import { SelectOutputFolder } from "../../../../wailsjs/go/main/App";
    import { EventsOn } from "../../../../wailsjs/runtime/runtime";

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

    const tiffCompressionMap: Record<TIFFCompressionType, number> = {
        None: 1,
        Deflate: 32946,
        LZW: 5,
    };

    let files: string[] = [];
    let outputFolder: string = "";
    let isConverting: boolean = false;

    let options: ConversionOptions = {
        format: "JPEG",
        quality: 85,
        lossless: false,
        pngLevel: 6,
        tiffCompress: "Deflate",
    };

    // Progression
    let convertedCount = 0;
    let totalCount = 0;

    const handleFiles = (e: CustomEvent<{ files: string[] }>) => {
        files = [...files, ...e.detail.files];
    };

    const removeFile = (e: CustomEvent<{ index: number }>) =>
        (files = files.filter((_, i) => i !== e.detail.index));

    const clearFiles = () => (files = []);

    const updateOptions = (e: CustomEvent<ConversionOptions>) =>
        (options = e.detail);

    const selectOutputFolder = async () => {
        try {
            const folder = await SelectOutputFolder();
            if (folder) outputFolder = folder;
        } catch (error) {
            console.error("Erreur lors de la sélection du dossier:", error);
        }
    };

    const convertToFolder = async () => {
        if (files.length === 0) return;
        if (!outputFolder) return;

        isConverting = true;
        convertedCount = 0;
        totalCount = files.length;

        // Écoute des événements de progression
        EventsOn("conversion-progress", (data: any) => {
            convertedCount = data.current;
        });

        try {
            await ConvertManyToFolder(files, outputFolder, {
                Format: options.format.toLowerCase(),
                Quality: options.quality,
                Lossless: options.lossless,
                PNGLevel: options.pngLevel,
                TIFFCompress: tiffCompressionMap[options.tiffCompress],
            });

            alert(
                `Conversion terminée ! ${convertedCount} / ${totalCount} fichier(s) converti(s)`,
            );
        } catch (error) {
            console.error(error);
            alert(`Erreur lors de la conversion: ${error.message || error}`);
        } finally {
            isConverting = false;
        }
    };

    $: placeholder = `Glissez-déposez ou cliquez (${files.length} fichier${files.length > 1 ? "s" : ""})`;
    $: canConvert = files.length > 0 && outputFolder && !isConverting;
</script>

<div class="space-y-6 p-4 text-[var(--text-primary)]">
    <FileDropzone on:files={handleFiles} {placeholder} />

    <FilePreview
        {files}
        viewMode="list"
        on:remove={removeFile}
        on:clearAll={clearFiles}
    />

    <ImageConversionSettings bind:options on:optionsChange={updateOptions} />

    <!-- Output folder selection -->
    <div class="space-y-3">
        <h3 class="text-lg font-semibold text-[var(--text-primary)]">
            Dossier de sortie
        </h3>
        <div class="flex gap-3">
            <input
                type="text"
                bind:value={outputFolder}
                placeholder="Sélectionnez un dossier de sortie..."
                readonly
                class="flex-1 px-3 py-2 bg-[var(--surface)] border border-[var(--border)] rounded-md text-sm text-[var(--text-primary)] placeholder-[var(--text-secondary)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] focus:border-transparent transition-all duration-300"
            />
            <button
                type="button"
                on:click={selectOutputFolder}
                class="px-4 py-2 bg-[var(--accent)] hover:bg-[var(--accent-dark)] text-[var(--text-tiers)] text-sm font-medium rounded-md transition-all duration-300 hover:shadow-[var(--shadow-accent)] hover:scale-105"
            >
                Parcourir
            </button>
        </div>
    </div>

    {#if isConverting}
        <div class="my-4">
            <p class="text-sm">
                Conversion : {convertedCount} / {totalCount} fichier(s)
            </p>
            <div
                class="w-full h-3 bg-[var(--surface)] rounded-md overflow-hidden"
            >
                <div
                    class="h-full bg-[var(--accent)] transition-all duration-300"
                    style="width: {(convertedCount / totalCount) * 100}%"
                ></div>
            </div>
        </div>
    {/if}

    <!-- Convert button -->
    {#if files.length > 0}
        <div class="flex justify-center mt-6">
            <button
                type="button"
                on:click={convertToFolder}
                disabled={!canConvert}
                class="flex items-center gap-3 px-6 py-3 font-medium rounded-lg transition-all duration-300
                       {isConverting
                    ? 'bg-[var(--warning)] text-[var(--background)] cursor-wait'
                    : canConvert
                      ? 'bg-gradient-to-r from-[var(--accent)] to-[var(--accent-dark)] text-[var(--text-tiers)] hover:shadow-[var(--shadow-accent)] hover:scale-105 hover:bg-gradient-to-r hover:from-[var(--accent-light)] hover:to-[var(--accent)]'
                      : 'bg-[var(--muted)] text-[var(--text-secondary)] cursor-not-allowed opacity-50'}"
            >
                {#if isConverting}
                    <div
                        class="animate-spin w-5 h-5 border-2 border-[var(--background)] border-t-transparent rounded-full"
                    ></div>
                    <span>Conversion en cours...</span>
                {:else}
                    <svg
                        class="w-5 h-5"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                    >
                        <path
                            fill-rule="evenodd"
                            d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z"
                            clip-rule="evenodd"
                        />
                    </svg>
                    <span
                        >Convertir vers dossier ({files.length} fichier{files.length >
                        1
                            ? "s"
                            : ""})</span
                    >
                {/if}
            </button>
        </div>
    {/if}
</div>
