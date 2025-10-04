<script context="module" lang="ts">
  import type { ComponentType } from "svelte";
  export interface Tab {
    label: string;
    content: string | ComponentType;
  }
</script>

<script lang="ts">
  export let tabs: Tab[] = [];
  export let activeTab: number = 0;
  export let maxTabWidth: string = "200px"; // Largeur maximale configurable
  
  function selectTab(index: number): void {
    activeTab = index;
  }
</script>

<div class="w-full " style="color: var(--text-primary); max-width:1200px;">
  <!-- Navigation -->
  <div class="flex border-b overflow-x-auto backdrop-blur-[var(--backdrop-blur)] border border-[var(--border)] rounded-xl shadow-[var(--shadow-accent)]" style="border-color: var(--border); background-color: var(--background);">
    {#each tabs as tab, index}
      <button
        class="px-4 py-2 font-medium text-sm transition-colors duration-200 border-b-2 whitespace-nowrap flex-shrink-0"
        style="
          color: {activeTab === index ? 'var(--accent)' : 'var(--text-secondary)'};
          border-color: {activeTab === index ? 'var(--accent)' : 'transparent'};
          max-width: {maxTabWidth};
          min-width: 80px;
        "
        on:click={() => selectTab(index)}
        on:mouseenter={(e) => {
          if (activeTab !== index) e.currentTarget.style.color = 'var(--text-primary)';
        }}
        on:mouseleave={(e) => {
          if (activeTab !== index) e.currentTarget.style.color = 'var(--text-secondary)';
        }}
      >
        <span class="truncate block" title={tab.label}>
          {tab.label}
        </span>
      </button>
    {/each}
  </div>
  
  <!-- Contenu -->
  <div class="mt-4">
    {#each tabs as tab, index}
      {#if activeTab === index}
        <div class="tab-content">
          {#if typeof tab.content === "string"}
            {@html tab.content}
          {:else}
            <svelte:component this={tab.content} />
          {/if}
        </div>
      {/if}
    {/each}
  </div>
</div>

<style>
  .tab-content {
    animation: fadeIn 0.25s ease-in-out;
    background-color: var(--background);
    color: var(--text-primary);
    padding: 1rem;
    border-radius: 0.5rem;
    border: 1px solid var(--border);
    box-shadow: var(--shadow-accent);
  }
  
  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(4px); }
    to { opacity: 1; transform: translateY(0); }
  }
</style>