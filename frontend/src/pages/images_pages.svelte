<script lang="ts">
  import Convert from "../components/templates/images/convert.svelte";
  import ConvertMany from "../components/templates/images/covert_many.svelte"
  import NavTab, { type Tab } from "../layouts/navTabs.svelte";

  const tabs: Tab[] = [
    {
      label: "Convert",
      content: Convert,
    },
    {
      label: "Convert Many",
      content: ConvertMany,
    },
  ];

  let currentTab: number = 0;
</script>

<div class="p-6">
  <NavTab {tabs} bind:activeTab={currentTab} />
</div>
