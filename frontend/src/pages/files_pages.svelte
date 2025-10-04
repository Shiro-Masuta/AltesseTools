<script lang="ts">
    import Duplicates from "../components/templates/files/duplicates.svelte";
  import Rename from "../components/templates/files/rename.svelte";
  import NavTab, { type Tab } from "../layouts/navTabs.svelte";
const tabs: Tab[] = [
  {
    label: "Renommer",
    content: Rename,
  },
  {
    label: "Recherche de doublons",
    content: Duplicates,
  },
];
  let currentTab: number = 0;
</script>

<div class="p-6">
  <NavTab {tabs} bind:activeTab={currentTab} />
</div>
