<script setup lang="ts">
import { TableData } from '@/Types/TableData';

defineProps<{
    data: TableData
}>()
</script>

<template>
    <table class="table">
        <thead>
            <tr>
                <th scope="col">#</th>
                <th scope="col" v-for="(col) in data.ColumnNames">{{ col }}</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="(rows, idx) in data.Rows">
                <th scope="row">{{ idx }}</th>
                <td v-for="row in rows">{{ row }}</td>
            </tr>
        </tbody>
    </table>
</template>
