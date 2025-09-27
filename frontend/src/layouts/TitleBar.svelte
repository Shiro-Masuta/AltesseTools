<script lang="ts">
    import { WindowMinimise, Quit } from "../../wailsjs/runtime";
    import logo from "../assets/images/altesseToolsLogo.png"; // ton logo

    function closeWindow() {
        Quit();
    }

    function minimizeWindow() {
        WindowMinimise();
    }
</script>

<!-- Import Google Material Icons -->
<svelte:head>
    <link
        href="https://fonts.googleapis.com/icon?family=Material+Icons"
        rel="stylesheet"
    />
</svelte:head>

<div
    class="titlebar flex items-center justify-between bg-[var(--background)] h-12 px-4 pl-0 select-none border-b border-[var(--border)]"
    style="--wails-draggable: drag;"
>
    <div class="title flex items-center space-x-2 text-[var(--text-primary)] font-semibold tracking-wide">
        <!-- Logo à gauche -->
        <img src={logo} alt="Logo" class="w-16 h-16 rounded-md object-cover" />
        <span>AltesseTools</span>
    </div>

    <div class="buttons flex items-center space-x-1" style="--wails-draggable: no-drag;">
        <!-- Bouton Minimiser -->
        <button
            class="window-btn w-8 h-8 flex items-center justify-center rounded-lg bg-[var(--surface-alt)] hover:bg-[var(--warning)] text-[var(--text-secondary)] hover:text-[var(--background)] transition-all duration-200"
            on:click={minimizeWindow}
            aria-label="Minimiser"
        >
            <span class="material-icons text-base">minimize</span>
        </button>

        <!-- Bouton Fermer -->
        <button
            class="window-btn w-8 h-8 flex items-center justify-center rounded-lg bg-[var(--surface-alt)] hover:bg-[var(--error)] text-[var(--text-secondary)] hover:text-[var(--text-tiers)] transition-all duration-200"
            on:click={closeWindow}
            aria-label="Fermer"
        >
            <span class="material-icons text-base">close</span>
        </button>
    </div>
</div>

<style>
    .window-btn:hover {
        transform: translateY(-1px);
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.25);
    }

    .window-btn:active {
        transform: translateY(0);
    }

    .titlebar:active {
        cursor: grab;
    }

    .window-btn {
        cursor: pointer;
        box-shadow:
            0 1px 3px rgba(0, 0, 0, 0.12),
            0 1px 2px rgba(0, 0, 0, 0.24);
    }
</style>
