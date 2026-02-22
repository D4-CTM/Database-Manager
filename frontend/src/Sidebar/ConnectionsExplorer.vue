<script setup lang="ts">
import CreateConnectionModal from './CreateConnectionModal.vue';
import { DbCredential } from '@/Types/Credential';
import DatabaseItem from './DatabaseItem.vue';
import { ref } from 'vue';
import { Get, PostOrPut } from '@/Helpers/HttpCaller';

let connections = ref([] as string[])
let credential = ref({} as DbCredential)
let showModal = ref(false)

async function loadConnections() {
    try {
        const data = await Get<string[]>('/api/Connection/list')
        connections.value = data
    } catch (ex) {
        alert(ex)
    }
}
loadConnections()

function closeModal() {
    credential.value = ({} as DbCredential)
    showModal.value = false
}

async function patchCredential(isNew: boolean) {
    try {
        const data = await PostOrPut<string, DbCredential>('/api/Connection', credential.value, isNew)
        if (isNew) {
            connections.value.push(data)
        }
    } catch (ex) {
        alert(ex)
    }
}
</script>

<template>
    <div class="flex-shrink-0 py-2 bg-white bg-transparent" style="width: 280px;">
        <a class="fs-5 fw-semibold d-flex justify-content-center pb-3 mb-3 text-decoration-none border-bottom">
            <CreateConnectionModal @onConfirm="() => patchCredential(true)" @onCancel="closeModal()"
                title="Create new contact" v-model="showModal" :credential="credential" />
        </a>
        <ul class="list-unstyled ps-0">
            <DatabaseItem v-for="con in connections" :conName="con" />
        </ul>
    </div>
</template>
