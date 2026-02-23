<script setup lang="ts">
import { TableData } from '@/Types/TableData';
import { ref } from 'vue';
import TableRenderer from './TableRenderer.vue';
import { Post, Put } from '@/Helpers/HttpCaller';
const enum ResultType {
    Null = 0,
    Number = 1,
    Table = 2
}

const props = defineProps({
    conName: {
        type: String,
        required: true
    }
})

let query = ref('')
let resultType = ref(ResultType.Null)
let result = ref<Number | TableData | null>({} as TableData)

async function exec() {
    try {
        const upperQuery = query.value.trim().toUpperCase()
        const conn = props.conName.replace(' ', '%20')

        if (upperQuery.startsWith('SELECT')) {
            resultType.value = ResultType.Table
            result.value = await Post<TableData, string>(`/api/Select/${conn}`, upperQuery)
        } else {
            resultType.value = ResultType.Number
            result.value = await Put<Number, string>(`/api/Exec/${conn}`, upperQuery)
        }
    } catch (ex) {
        alert(ex)
    }
}

function clear() {
    query.value = '';
}
</script>

<template>
    <div class="d-flex flex-column">
        <nav class="navbar navbar-expand-lg navbar-light bg-light">
            <button class="btn btn-primary m-1" @click="exec">
                execute
            </button>
            <button class="btn btn-secondary m-1" @click="clear">
                clear
            </button>
        </nav>
        <div class="mb-3">
            <label for="queryArea" class="form-label">Query to execute</label>
            <textarea v-model="query" class="form-control" wrap="hard" id="queryArea" style="max-height: 40vh;" rows="5"></textarea>
        </div>
        <div v-show='resultType === ResultType.Number'>
            <input type="text" disabled="true" class="form-control" :value="`Affected rows: ${result}`" />
        </div>
        <div v-show="resultType === ResultType.Table" 
             style="flex: auto; overflow: auto; min-height: 0;">
            <TableRenderer :data="result as TableData" />
        </div>
    </div>
</template>
