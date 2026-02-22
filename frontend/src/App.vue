<script setup lang="ts">
import { ref } from 'vue';
import ConnectionsExplorer from './Sidebar/ConnectionsExplorer.vue';
import Tabs from './Components/tabs/Tabs.vue';
import { TabData } from './Types/TabData';

let tabs = ref([] as TabData[])

function removeTab(idx: number) {
    tabs.value.splice(idx, 1)
}
</script>

<template>
    <main class="d-flex bg-black bg-gradient" style="height: 100vh; max-height: 100vh;">
        <ConnectionsExplorer />
        <div class="w-100 bg-light">
            <Tabs @removeTab="removeTab"
                  :tabs="tabs"/>
        </div>
    </main>
</template>
