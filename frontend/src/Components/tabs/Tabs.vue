<script setup lang="ts">
import { TabData, TabOptions } from '@/Types/TabData';
import { ref } from 'vue';
import Query from './Query.vue';

const props = defineProps({
    tabs: {
        type: Array<TabData>,
        required: true
    },
})

const emit = defineEmits(['removeTab'])

let selectedIdx = ref(0)

function removeTab(idx: number) {
    const old = selectedIdx.value
    selectedIdx.value = old > 0 ? old - 1 : old
    emit('removeTab', idx)
}
</script>

<template>
    <div class="d-flex" style="max-height: 100vh;">
        <ul class="nav nav-tabs">
            <li class="nav-item fs-5 inline" v-for="(tab, idx) in tabs">
                <div class="nav-link" :class="[selectedIdx == idx ? 'active' : '']">
                    <a @click="selectedIdx = idx">
                        {{ tab.Title }}
                    </a>
                    <button class="btn"
                        @click="removeTab(idx)">
                        &times;
                    </button>
                </div>
            </li>
        </ul>
    </div>
    <div v-if="tabs.length > 0" class="p-2 border-top w-100">
        <Query v-show="tabs[selectedIdx].Type === TabOptions.Query" />
    </div>
</template>
