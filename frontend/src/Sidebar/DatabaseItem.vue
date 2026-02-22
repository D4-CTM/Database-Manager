<script setup lang="ts">
import Expansible from '@/Components/Expansible.vue';
import { Get } from '@/Helpers/HttpCaller';
import { PingResult } from '@/Types/PingResult';
import { ref } from 'vue';
import DatabaseContent from './DatabaseContent.vue';

const props = defineProps<{
    conName: string
}>()
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

</script>

<template>
    <Expansible class="border-bottom" @beforeExpand="pingConnection" :btn-txt="conName">
        <template #CONTENT>
            <Expansible v-if="dbUsers.accepted" v-for="schema in dbUsers.Schemas" 
                @beforeExpand="(expand: Function) => expand()" :btnTxt="schema">
                <template #CONTENT>
                    <DatabaseContent :conName="conName" :schema="schema"/>
                </template>
            </Expansible>
        </template>
    </Expansible>
</template>
