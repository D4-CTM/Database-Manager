<script setup lang="ts">
import { Get } from '@/Helpers/HttpCaller';
import { PingResult } from '@/Types/PingResult';
import { inject, ref } from 'vue';
import DatabaseContent from './DatabaseContent.vue';
import Expandible from '@/Components/Expandible.vue';
import { TabData, TabOptions, TabStore } from '@/Types/TabData';

const props = defineProps<{
    conName: string
}>()
const emit = defineEmits(['openModal', 'delContact'])
let store = inject<TabStore>('TabStore')
let dbUsers = ref({} as PingResult)

async function pingConnection(expand: Function) {
    try {
        let result = await Get<PingResult>(`/api/Ping/${props.conName}`)
        result.accepted = true
        dbUsers.value = result
        expand()
    } catch (ex) {
        alert(ex)
        dbUsers.value.accepted = false
    }
}

async function deleteConnection() {
    emit('delContact', props.conName,() => dbUsers.value = {} as PingResult)
}

function openModal() {
    emit('openModal', props.conName)
}

function openSqlEditor() {
    const conn = props.conName.trim();
    store.add({
        Title: `${conn}.sql`,
        Type: TabOptions.Query,
        Payload: conn
    } as TabData)
}
</script>

<template>
    <Expandible class="border-bottom" @beforeExpand="pingConnection" :btn-txt="conName">
        <template #CONTENT>
            <Expandible v-if="dbUsers.accepted" v-for="schema in dbUsers.Schemas"
                @beforeExpand="(expand: Function) => expand()" :btnTxt="schema">
                <template #CONTENT>
                    <DatabaseContent :conName="conName" :schema="schema" />
                </template>
            </Expandible>
        </template>
        <template #OPTIONS>
            <li><button class="dropdown-item" @click="deleteConnection">Delete</button></li>
            <li><button class="dropdown-item" @click="openModal">Edit Connection</button></li>
            <li><button class="dropdown-item" @click="openSqlEditor">Sql Editor</button></li>
        </template>
    </Expandible>
</template>
