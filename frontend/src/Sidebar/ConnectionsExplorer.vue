<script setup lang="ts">
import CreateConnectionModal from './CreateConnectionModal.vue';
import { DbCredential, PutCredentialResult } from '@/Types/Credential';
import DatabaseItem from './DatabaseItem.vue';
import { ref } from 'vue';
import { Delete, Get, Post, Put } from '@/Helpers/HttpCaller';

let connections = ref([] as string[])
let credential = ref({ isNew: true } as DbCredential)
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

async function openModal(conName: string) {
    try {
        credential.value = await Get<DbCredential>(`/api/Connection?conName=${conName}`)
        credential.value.conName = conName
        credential.value.isNew = false
        showModal.value = true
    } catch(ex) {
        alert(ex)
    } 
}

function closeModal() {
    credential.value = ({} as DbCredential)
    showModal.value = false
}

async function patchCredential(conName: String) {
    try {
        if (credential.value.isNew) {
            const data = await Post<string, DbCredential>('/api/Connection', credential.value)
            connections.value.push(data)
        } else {
            const newName = conName.trim().replace(' ', '_')
            const data = await Put<PutCredentialResult, DbCredential>(`/api/Connection/${credential.value.conName}?newName=${newName}`, credential.value)
            var idx = connections.value.findIndex(x => x === data.OldName)
            if (idx === -1) new Error(`Unable to find ${data.OldName}`)
            connections.value.splice(idx, 1)
            connections.value.push(data.NewName)
            console.log(`removed: ${data.OldName} for: ${data.NewName}`)
        }
    } catch (ex) {
        alert(ex)
    }
}

async function deleteConnection(conName: string, restart: Function) {
    try {
        var idx = connections.value.findIndex(x => x === conName)
        if (idx === -1) new Error(`Unable to find ${conName}`)

        await Delete<void>(`/api/Connection/${conName}`)
        connections.value.splice(idx, 1)
        restart()
    } catch (ex) {
        alert(ex)
    }
}
</script>

<template>
    <div class="flex-shrink-0 py-2 bg-white bg-transparent" style="width: 280px;">
        <a class="fs-5 fw-semibold d-flex justify-content-center pb-3 mb-3 text-decoration-none border-bottom w-100">
            <CreateConnectionModal @onConfirm="patchCredential" @onCancel="closeModal()"
                title="Create new contact" v-model="showModal" :credential="credential" />
        </a>
        <ul class="list-unstyled ps-0">
            <DatabaseItem v-for="con in connections" 
                          @delContact="deleteConnection"
                          @openModal="openModal"
                          :conName="con" />
        </ul>
    </div>
</template>
