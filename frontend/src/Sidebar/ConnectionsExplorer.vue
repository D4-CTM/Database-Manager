<script setup lang="ts">
import CreateConnectionModal from './CreateConnectionModal.vue';
import { DbCredential } from '@/Types/Credential';
import DatabaseItem from './DatabaseItem.vue';
import { ref } from 'vue';
import axios, { AxiosResponse, HttpStatusCode } from 'axios';
import { ErrorRequest } from '@/Types/ErrorRequest';

let connections = ref([] as string[])
let credential = ref({} as DbCredential)
let showModal = ref(false)

const emits = defineEmits(['update:connections:push'])
function closeModal() {
    credential.value = ({} as DbCredential)
    showModal.value = false
}

async function patchCredential(isNew: boolean) {
    try {
        let response: AxiosResponse<string> = isNew
            ? await axios.post<string>('/api/Connection', credential.value)
            : await axios.put<string>('/api/Connection', credential.value)

        if (response.status === HttpStatusCode.Ok) {
            connections.value.push(response.data)
            closeModal()
        }
    } catch (ex) {
        if (axios.isAxiosError(ex) && ex.response) {
            const err = ex.response.data as string 
            console.error(`Server failed: ${err}`)
            alert(err)
        } else {
            console.error('Network or unknown error', ex)
        }
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
