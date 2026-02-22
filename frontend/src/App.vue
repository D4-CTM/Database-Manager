<script setup lang="ts">
import { ref } from 'vue';
import ConnectionsExplorer from './Sidebar/ConnectionsExplorer.vue';
import { DbCredential } from './Types/Credential';
import Tabs from './Components/tabs/Tabs.vue';
import { TabData, TabOptions } from './Types/TabData';

let Cred: DbCredential = {
    Database: '',
    Password: '',
    Server: '',
    ShowAll: false,
    User: '',
    Port: 0
}
let tabs = ref([
    {
        Title: 'Query #1',
        Type: TabOptions.Query, 
    },
] as TabData[])

function removeTab(idx: number) {
    tabs.value.splice(idx, 1)
}
</script>

<template>
    <main class="d-flex bg-black bg-gradient" style="height: 100vh; max-height: 100vh;">
        <ConnectionsExplorer
            :credential="Cred" />
        <div class="w-100 bg-light">
            <Tabs @removeTab="removeTab"
                  :tabs="tabs"/>
        </div>
    </main>
</template>
