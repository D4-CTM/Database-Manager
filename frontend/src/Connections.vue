<script setup lang="ts">
import { provide } from 'vue';
import { TabStore } from './Types/TabData';
import ConnectionsExplorer from './Sidebar/ConnectionsExplorer.vue';
import Tabs from './tabs/Tabs.vue';

const store = new TabStore()
provide('TabStore', store)
</script>

<template>
    <main class="d-flex bg-black bg-gradient" style="height: 100vh; max-height: 100vh;">
        <ConnectionsExplorer style="overflow-y: auto; overflow-x: hidden;"/>
        <div class="w-100 bg-light d-flex flex-column overflow-hidden" style="height: 100vh;">
            <Tabs />
        </div>
    </main>
</template>
