<script lang="ts">
    import Sidebar from "./layouts/sidebar.svelte";
    import TitleBar from "./layouts/TitleBar.svelte";

    import Home from "./pages/home_pages.svelte";
    import Images from "./pages/images_pages.svelte";
    import Files from "./pages/files_pages.svelte";
    // Menu avec composant
    const menuItems = [
        {
            id: "dashboard",
            icon: "dashboard",
            label: "Dashboard",
            component: Home,
        },
        { id: "image", icon: "image", label: "Images", component: Images },
        { id: "files", icon: "folder", label: "Files", component: Files },
    ];

    let selected = menuItems[0];

    type MenuItem = {
        id: string;
        icon: string;
        label: string;
        component: typeof Home | typeof Images;
    };

    const navigate = (item: MenuItem) => {
        selected = item;
    };
</script>

<svelte:head>
    <link
        href="https://fonts.googleapis.com/icon?family=Material+Icons"
        rel="stylesheet"
    />
</svelte:head>

<!-- Layout principal -->
<div class="flex flex-col h-screen">
    <!-- TitleBar en haut, fixe -->
    <TitleBar />

    <div class="flex flex-1 overflow-hidden">
        <!-- Sidebar fixe à gauche -->
        <Sidebar {menuItems} {selected} {navigate} />

        <!-- Zone de contenu scrollable avec scrollbar décalée -->
        <div class="relative flex-1 overflow-hidden">
            <!-- Contenu principal scrollable -->
            <div class="absolute inset-0 overflow-y-auto pr-2">
                <svelte:component this={selected.component} />
            </div>

            <!-- Scrollbar custom (optionnelle) -->
            <div
                class="absolute top-0 right-0 w-1 h-full bg-transparent scrollbar-custom"
            ></div>
        </div>
    </div>
</div>
