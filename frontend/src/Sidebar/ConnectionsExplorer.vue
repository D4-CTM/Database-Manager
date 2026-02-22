<script setup lang="ts">
import CreateConnectionModal from './CreateConnectionModal.vue';
import { DbCredential } from '@/Types/Credential';
import DatabaseItem from './DatabaseItem.vue';
import { ref } from 'vue';

let credential = ref({} as DbCredential)
let showModal = ref(false)
</script>

<template>
    <div class="flex-shrink-0 py-2 bg-white bg-transparent" style="width: 280px;">
        <a class="fs-5 fw-semibold d-flex justify-content-center pb-3 mb-3 text-decoration-none border-bottom">
            <CreateConnectionModal
                v-model="showModal" :credential="credential" />
        </a>
        <ul class="list-unstyled ps-0">
            <DatabaseItem conName="test" />
        </ul>
    </div>
</template>
