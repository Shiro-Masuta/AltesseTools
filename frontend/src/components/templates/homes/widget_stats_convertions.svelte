<script lang="ts">
  import { onMount } from 'svelte';
  import { GetWidgetStats } from "../../../../wailsjs/go/services/StatsService";

  let stats: any = null;
  let loading = true;

  async function loadStats() {
    try {
      loading = true;
      stats = await GetWidgetStats();
    } catch (e) {
      console.error('Erreur stats:', e);
    } finally {
      loading = false;
    }
  }

  // Stats calculées
  $: savedMB = stats ? Math.round((stats.total_original_size - stats.total_compressed_size) / (1024 * 1024)) : 0;
  $: compressionRatio = stats && stats.total_original_size > 0 
    ? ((stats.total_original_size - stats.total_compressed_size) / stats.total_original_size * 100) 
    : 0;
  $: topFormat = stats?.formats ? Object.entries(stats.formats)
    .reduce((max, [format, data]: [string, any]) => 
      data.count > (max.count || 0) ? { format, count: data.count } : max, 
      { format: 'N/A', count: 0 }
    ).format : 'N/A';

  onMount(() => {
    loadStats();
    const interval = setInterval(loadStats, 30000);
    return () => clearInterval(interval);
  });
</script>

<div class="bg-[var(--background)] backdrop-blur-[var(--backdrop-blur)] border border-[var(--border)] shadow-[var(--shadow-accent)] rounded-xl m-5 p-6 max-w-sm text-[var(--text-primary)]">
  <!-- Header -->
  <div class="flex items-center justify-between mb-6">
    <h2 class="text-xl font-semibold">Statistiques</h2>
    <button
      on:click={loadStats}
      disabled={loading}
      class="p-2 text-[var(--text-secondary)] hover:text-[var(--accent)] hover:bg-[var(--surface-alt)] 
             rounded-lg transition-all duration-200 disabled:opacity-50"
    >
      <svg class="w-4 h-4 fill-none stroke-current stroke-2 {loading ? 'animate-spin' : ''}" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
      </svg>
    </button>
  </div>

  {#if loading && !stats}
    <!-- Loading -->
    <div class="space-y-3">
      <div class="h-3 bg-[var(--surface-alt)] rounded animate-pulse w-3/4"></div>
      <div class="h-3 bg-[var(--surface-alt)] rounded animate-pulse w-1/2"></div>
      <div class="h-3 bg-[var(--surface-alt)] rounded animate-pulse w-2/3"></div>
    </div>
  {:else if stats}
    <!-- Stats Grid -->
    <div class="grid grid-cols-2 gap-3 mb-4">
      <div class="bg-[var(--surface-alt)] border border-[var(--border)] rounded-lg p-4 text-center 
                  hover:-translate-y-1 transition-transform duration-200">
        <div class="text-2xl font-bold text-[var(--accent)]">{stats.total_converted}</div>
        <div class="text-xs text-[var(--text-secondary)]">Convertis</div>
      </div>
      
      <div class="bg-[var(--surface-alt)] border border-[var(--border)] rounded-lg p-4 text-center 
                  hover:-translate-y-1 transition-transform duration-200">
        <div class="text-2xl font-bold text-[var(--success)]">{savedMB}MB</div>
        <div class="text-xs text-[var(--text-secondary)]">Économisés</div>
      </div>
      
      <div class="bg-[var(--surface-alt)] border border-[var(--border)] rounded-lg p-4 text-center 
                  hover:-translate-y-1 transition-transform duration-200">
        <div class="text-lg font-bold text-[var(--accent-light)]">{topFormat}</div>
        <div class="text-xs text-[var(--text-secondary)]">Format pop.</div>
      </div>
      
      <div class="bg-[var(--surface-alt)] border border-[var(--border)] rounded-lg p-4 text-center 
                  hover:-translate-y-1 transition-transform duration-200">
        <div class="text-lg font-bold text-[var(--text-primary)]">{compressionRatio.toFixed(1)}%</div>
        <div class="text-xs text-[var(--text-secondary)]">Compression</div>
      </div>
    </div>

    <!-- Équivalences rétro -->
    {#if stats.total_saved_cd !== undefined || stats.total_saved_floppy !== undefined}
      <div class="bg-[var(--surface-alt)] border border-[var(--border)] rounded-lg p-3">
        <div class="text-sm text-[var(--text-secondary)] mb-2">Équivalences rétro</div>
        <div class="flex justify-between text-sm">
          {#if stats.total_saved_cd !== undefined}
            <span class="flex items-center gap-1">💿 {stats.total_saved_cd} CD</span>
          {/if}
          {#if stats.total_saved_floppy !== undefined}
            <span class="flex items-center gap-1">💾 {stats.total_saved_floppy} Disquettes</span>
          {/if}
        </div>
      </div>
    {/if}
  {/if}
</div>