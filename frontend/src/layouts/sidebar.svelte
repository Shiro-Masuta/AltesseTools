<script lang="ts">
    export let menuItems;
    export let selected;
    export let navigate;
</script>

<div
    class="flex h-screen"
    style="background-color: var(--background); color: var(--text-primary);"
>
    <!-- Sidebar -->
    <div
        class="w-16 flex flex-col border-r border-gray-700"
        style="background-color: var(--surface);"
    >

        <!-- Menu Items -->
        <nav
            class="flex-1 py-4 space-y-2"
            style="background-color: var(--background);"
        >
            {#each menuItems as item}
                <div class="relative group px-2">
                    <button
                        class="w-12 h-12 flex items-center justify-center rounded-xl transition-all duration-200"
                        style="
                            background-color: {selected.id === item.id
                            ? 'var(--accent)'
                            : 'transparent'};
                            color: {selected.id === item.id
                            ? 'var(--text-tiers)'
                            : 'var(--text-secondary)'};
                        "
                        on:click={() => navigate(item)}
                        on:mouseenter={(e) =>
                            (e.currentTarget.style.backgroundColor =
                                selected.id === item.id
                                    ? "var(--accent)"
                                    : "var(--surface-alt)")}
                        on:mouseleave={(e) =>
                            (e.currentTarget.style.backgroundColor =
                                selected.id === item.id
                                    ? "var(--accent)"
                                    : "transparent")}
                    >
                        <span class="material-icons text-xl">{item.icon}</span>
                    </button>

                    <!-- Tooltip -->
                    <div
                        class="absolute left-16 top-1/2 -translate-y-1/2 px-3 py-2 text-sm rounded-lg opacity-0 pointer-events-none group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10 shadow-lg"
                        style="background-color: var(--surface-alt); color: var(--text-primary); border: 1px solid var(--border); backdrop-filter: var(--backdrop-blur);"
                    >
                        {item.label}
                        <div
                            class="absolute left-0 top-1/2 -translate-y-1/2 -translate-x-1 w-2 h-2 rotate-45"
                            style="background-color: var(--surface-alt); border-left: 1px solid var(--border); border-bottom: 1px solid var(--border);"
                        ></div>
                    </div>
                </div>
            {/each}
        </nav>
    </div>

    <!-- Contenu principal -->
    <div
        class="flex-1 overflow-auto"
        style="background-color: var(--background); color: var(--text-primary);"
    >
        <slot />
    </div>
</div>
