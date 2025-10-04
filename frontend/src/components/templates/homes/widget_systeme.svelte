<script lang="ts">
  import { onMount } from "svelte";
  import { EventsOn } from "../../../../wailsjs/runtime/runtime";
  import { GetNow, Start } from "../../../../wailsjs/go/services/SystemService";
  import type { system } from "../../../../wailsjs/go/models"; // <-- types générés

  let system: system.SystemMinimalInfo | null = null;

  const formatBytes = (bytes: number): string => {
    const sizes = ["B", "KB", "MB", "GB", "TB"];
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return (bytes / Math.pow(1024, i)).toFixed(1) + " " + sizes[i];
  };

  const getStatusColor = (percent: number): string => {
    if (percent > 80) return "stroke-[var(--error)]";
    if (percent > 60) return "stroke-[var(--warning)]";
    return "stroke-[var(--success)]";
  };

  const getTextColor = (percent: number): string => {
    if (percent > 80) return "text-[var(--error)]";
    if (percent > 60) return "text-[var(--warning)]";
    return "text-[var(--success)]";
  };

  const formatMhz = (mhz: number): string => {
    return mhz > 1000 ? `${(mhz / 1000).toFixed(1)}GHz` : `${mhz}MHz`;
  };

  onMount(async () => {
    EventsOn("system:update", (info: system.SystemMinimalInfo) => {
      system = info;
    });

    system = await GetNow();
    await Start();
  });
</script>

<!-- Container principal -->
<div
  class="p-6 m-5 bg-[var(--background)] backdrop-blur-[var(--backdrop-blur)] border border-[var(--border)] rounded-xl shadow-[var(--shadow-accent)]"
>
  {#if system}
    <!-- Grille responsive pour CPU et RAM -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
      <!-- Section CPU -->
      <div class="flex items-center gap-4">
        <!-- Icône CPU -->
        <div class="flex-shrink-0">
          <span
            class="material-icons-outlined text-[var(--text-secondary)] text-2xl"
            >memory</span
          >
        </div>

        <!-- Graphique circulaire CPU -->
        <div class="relative w-12 h-12 flex-shrink-0">
          <svg class="w-full h-full -rotate-90" viewBox="0 0 36 36">
            <!-- Cercle de fond -->
            <path
              class="stroke-[var(--muted)]"
              stroke-dasharray="100, 100"
              stroke-width="3"
              fill="none"
              d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
            />
            <!-- Cercle de progression -->
            <path
              class={getStatusColor(system.cpu.percent)}
              stroke-dasharray="{system.cpu.percent}, 100"
              stroke-width="3"
              fill="none"
              d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
            />
          </svg>
          <!-- Pourcentage au centre -->
          <div class="absolute inset-0 flex items-center justify-center">
            <span class="text-xs font-bold {getTextColor(system.cpu.percent)}">
              {system.cpu.percent.toFixed(0)}%
            </span>
          </div>
        </div>

        <!-- Informations CPU -->
        <div class="min-w-0 flex-1">
          <h3 class="text-sm font-medium text-[var(--text-primary)] mb-1">
            CPU
          </h3>
          <div class="space-y-0.5">
            <p class="text-xs text-[var(--text-secondary)]">
              {system.cpu.cores} Cores / {system.cpu.threads} Threads
            </p>
            <p class="text-xs text-[var(--text-secondary)]">
              {formatMhz(system.cpu.mhz)}
            </p>
          </div>
        </div>
      </div>

      <!-- Section RAM -->
      <div class="flex items-center gap-4">
        <!-- Icône RAM -->
        <div class="flex-shrink-0">
          <span
            class="material-icons-outlined text-[var(--text-secondary)] text-2xl"
            >developer_board</span
          >
        </div>

        <!-- Graphique circulaire RAM -->
        <div class="relative w-12 h-12 flex-shrink-0">
          <svg class="w-full h-full -rotate-90" viewBox="0 0 36 36">
            <!-- Cercle de fond -->
            <path
              class="stroke-[var(--muted)]"
              stroke-dasharray="100, 100"
              stroke-width="3"
              fill="none"
              d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
            />
            <!-- Cercle de progression -->
            <path
              class={getStatusColor(system.memory.percent)}
              stroke-dasharray="{system.memory.percent}, 100"
              stroke-width="3"
              fill="none"
              d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
            />
          </svg>
          <!-- Pourcentage au centre -->
          <div class="absolute inset-0 flex items-center justify-center">
            <span
              class="text-xs font-bold {getTextColor(system.memory.percent)}"
            >
              {system.memory.percent.toFixed(0)}%
            </span>
          </div>
        </div>

        <!-- Informations RAM -->
        <div class="min-w-0 flex-1">
          <h3 class="text-sm font-medium text-[var(--text-primary)] mb-1">
            RAM
          </h3>
          <div class="space-y-0.5">
            <p class="text-xs text-[var(--text-secondary)]">
              {formatBytes(system.memory.used)} / {formatBytes(
                system.memory.total,
              )}
            </p>
            <p class="text-xs text-[var(--text-secondary)]">
              {formatBytes(system.memory.free)} libres
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Section Disques -->
    {#if system.disk && system.disk.length > 0}
      <div class="border-t border-[var(--border)]/30 pt-6">
        <!-- En-tête des disques -->
        <div class="flex items-center gap-3 mb-4">
          <span
            class="material-icons-outlined text-[var(--text-secondary)] text-2xl"
            >storage</span
          >
          <h3 class="text-sm font-medium text-[var(--text-primary)]">
            Stockage
          </h3>
        </div>

        <!-- Grille des disques -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {#each system.disk as disk}
            <div
              class="flex items-center gap-3 p-3 bg-[var(--surface-alt)]/50 rounded-lg border border-[var(--border)]/20"
            >
              <!-- Graphique circulaire disque -->
              <div class="relative w-10 h-10 flex-shrink-0">
                <svg class="w-full h-full -rotate-90" viewBox="0 0 36 36">
                  <!-- Cercle de fond -->
                  <path
                    class="stroke-[var(--muted)]"
                    stroke-dasharray="100, 100"
                    stroke-width="4"
                    fill="none"
                    d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
                  />
                  <!-- Cercle de progression -->
                  <path
                    class={getStatusColor(disk.percent)}
                    stroke-dasharray="{disk.percent}, 100"
                    stroke-width="4"
                    fill="none"
                    d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
                  />
                </svg>
                <!-- Pourcentage au centre -->
                <div class="absolute inset-0 flex items-center justify-center">
                  <span class="text-xs font-bold {getTextColor(disk.percent)}">
                    {disk.percent.toFixed(0)}%
                  </span>
                </div>
              </div>

              <!-- Informations disque -->
              <div class="min-w-0 flex-1">
                <h4
                  class="text-sm font-medium text-[var(--text-primary)] truncate"
                >
                  {disk.mountpoint}
                </h4>
                <p class="text-xs text-[var(--text-secondary)] mt-1">
                  {formatBytes(disk.used)} / {formatBytes(disk.total)}
                </p>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}
  {:else}
    <!-- État de chargement -->
    <div class="flex flex-col items-center justify-center py-12 gap-4">
      <div
        class="w-8 h-8 border-3 border-[var(--accent)] border-t-transparent rounded-full animate-spin"
      ></div>
      <p class="text-[var(--text-secondary)] text-sm font-medium">
        Chargement des informations système...
      </p>
    </div>
  {/if}
</div>
