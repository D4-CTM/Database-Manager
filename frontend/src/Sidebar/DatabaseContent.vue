<script setup lang="ts">
import Expandible from '@/Components/Expandible.vue';
import { Get } from '@/Helpers/HttpCaller';
import { Options } from '@/Types/Options';
import { QueryPayload, TabData, TabOptions, TabStore } from '@/Types/TabData';
import { inject, ref } from 'vue';
const props = defineProps({
    schema: {
        type: String,
        required: true
    },
    conName: {
        type: String,
        required: true
    }
})
let store = inject<TabStore>('TabStore')
let content = ref<Record<string, string[]>>({})

async function fetchContent(expand: Function, opt: string) {
    try {
        let result = await Get<string[]>(`/api/${opt}/${props.conName}/${props.schema}`)
        content.value[`${props.conName}:${opt}`] = result
        expand()
    } catch (ex) {
        alert(ex)
    }
}

async function fetchData(opt: string) {
    try {
        const table = opt.replace(' ', '%20')
        const schema = props.schema.replace(' ', '%20')

        const tab: TabData = {
            Title: opt,
            Type: TabOptions.Table,
            conName: props.conName,
            Payload: {
                table: table,
                owner: schema
            } as QueryPayload
        }

        store.add(tab)
    } catch (ex) {
        alert(ex)
    }
}
</script>

<template>
    <Expandible v-for="(icon, name) in Options" @beforeExpand="(expand: Function) => fetchContent(expand, name)"
        :btn-txt="name" :idle-icon="icon">
        <template #CONTENT>
            <button class="btn btn-light my-1 mx-2" v-for="val in content[`${conName}:${name}`]"
                    @click="fetchData(val)">
                <i class="bi bi-box mx-1" />{{ val }}
            </button>
        </template>
    </Expandible>
</template>
