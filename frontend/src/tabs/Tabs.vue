<script setup lang="ts">
import { TabOptions, TabStore } from '@/Types/TabData';
import { inject } from 'vue';
import Query from './Query.vue';
import TableInformation from './TableInformation.vue';

let store = inject<TabStore>('TabStore')

function removeTab(title: string) {
    if (!store.removeAt(title))
        alert(`Unable to remove tab ${title}`)
}
</script>

<template>
    <div class="d-flex flex-column w-100">
        <ul class="nav nav-tabs">
            <li class="nav-item fs-5 inline" v-for="(tab, idx) in store.list()">
                <div class="nav-link" :class="[store.currentIdx() == idx ? 'active' : '']">
                    <a @click="store.selectedTab.value = idx">
                        {{ tab.Title }}
                    </a>
                    <button class="btn" @click="removeTab(tab.Title)">
                        &times;
                    </button>
                </div>
            </li>
        </ul>
    </div>
    <div v-if="store.length() > 0" class="p-2 border-top w-100 flex-grow-1 overflow-auto">
        <Query v-if="store.currentTabType() === TabOptions.Query" :conName="store.currentConn()" />
        <TableInformation v-else-if="store.currentTabType() === TabOptions.Table" :query="store.currentPayload()" :conName="store.currentConn()" />
    </div>
</template>
