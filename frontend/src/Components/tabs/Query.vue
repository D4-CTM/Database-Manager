<script setup lang="ts">
import { TableData } from '@/Types/TableData';
import { ref } from 'vue';
import TableRenderer from '../TableRenderer.vue';
const NumberType = 'number'
const enum ResultType {
    Null = 0,
    Number = 1,
    Table = 2
}

let query = ref('')
let resultType = ref(ResultType.Table)
let result = ref<Number | TableData | null>({} as TableData)

function exec() {
    console.log(query.value)
    if (result.value === null) { 
        resultType.value = ResultType.Null
    } else if ((typeof result.value) === NumberType) {
        resultType.value = ResultType.Number
    } else {
        resultType.value = ResultType.Table
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
            <textarea v-model="query" class="form-control overflow-scroll" id="queryArea" style="max-height: 40vh;" rows="5"></textarea>
        </div>
        <div v-show='resultType === ResultType.Number'>
            <input type="text" disabled="true" class="form-control" :value="`Affected rows: ${result}`" />
        </div>
        <div style="flex: 1; overflow: scroll; min-height: 0;">
            <TableRenderer
                v-show="resultType === ResultType.Table"
                :data="result as TableData" />
        </div>
    </div>
</template>
