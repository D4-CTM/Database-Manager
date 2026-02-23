<script setup lang="ts">
import { ref } from 'vue';
import TableRenderer from './TableRenderer.vue';
import { ColumnsMetadata, ConstraintMetadata, TableData } from '@/Types/TableData';
import { Get } from '@/Helpers/HttpCaller';
import { QueryPayload } from '@/Types/TabData';
const prop = defineProps<{
    conName: String,
    query: QueryPayload
}>()

const options = [
    'Table',
    'Columns',
    'Constraints'
]
let table = ref({
    Name: options[0]
} as TableData)

async function fillTable(idx: number) {
    const qd = prop.query
    try {
        switch (idx) {
            case 0: {
                const result =
                    await Get<TableData>(`/api/Query/Table/${prop.conName}?Table=${qd.table}&Owner=${qd.owner}`)

                table.value = result
            } break
            case 1: {
                const result =
                    await Get<ColumnsMetadata[]>(`/api/Query/Columns/${prop.conName}?Table=${qd.table}&Owner=${qd.owner}`)
                
                table.value = colToTableData(result)
            } break
            case 2: {
                const result =
                    await Get<ConstraintMetadata[]>(`/api/Query/Constraints/${prop.conName}?Table=${qd.table}&Owner=${qd.owner}`)

                table.value = conToTableData(result)
            } break
        }
    } catch (ex) {
        alert(ex)
        console.log(ex)
    }

    table.value.Name = options[idx]
}
fillTable(0)

function conToTableData(col: ConstraintMetadata[]) {
    let rows: string[][] = []
    col.forEach(x => rows.push([
        x.ConstraintName,
        x.ConstraintType,
        x.ColumnName,
        x.SearchCondition,
        x.RefConstraint,
        x.Position.toString()
    ]))
    return {
        ColumnNames: [
            'ConstraintName',
            'ConstraintType',
            'ColumnName',
            'SearchCondition',
            'RefConstraint',
            'Position',
        ],
        Rows: rows
    } as TableData
}

function colToTableData(col: ColumnsMetadata[]) {
    let rows: string[][] = []
    col.forEach(x => rows.push([
        x.Name,
        x.DataType,
        x.DefaultData,
        x.Detail,
        x.IsNullable ? '[X]' : '[ ]',
        x.IsIdentity ? '[X]' : '[ ]',
        x.OrdPosition.toString()
    ]))
    return {
        ColumnNames: [
            'Name',
            'DataType',
            'DefaultData',
            'Detail',
            'IsNullable',
            'IsIdentity',
            'OrdPosition'
        ],
        Rows: rows
    } as TableData
}
</script>

<template>
    <div class="d-flex flex-column w-100">
        <ul class="nav nav-tabs">
            <li class="nav-item fs-6 inline" v-for="(tab, idx) in options">
                <div class="nav-link" :class="[table.Name == tab ? 'active' : '']">
                    <a @click="fillTable(idx)">
                        {{ tab }}
                    </a>
                </div>
            </li>
        </ul>
    </div>
    <div class="p-2 border-top w-100 flex-grow-1 overflow-auto">
        <TableRenderer :data="table" />
    </div>
</template>
