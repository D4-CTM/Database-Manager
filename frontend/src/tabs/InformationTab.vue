<script setup lang="ts">
import { ref } from 'vue';
import TableRenderer from './TableRenderer.vue';
import { ColumnsMetadata, ConstraintMetadata, TableData } from '@/Types/TableData';
import { Get } from '@/Helpers/HttpCaller';
import { QueryPayload } from '@/Types/TabData';
const prop = defineProps<{
    conName: String,
    query: QueryPayload
    options: string[]
    conType: string
}>()

let ddl = ref<string>('')
let table = ref({
    Name: prop.options[0]
} as TableData)

console.log(prop.conName)
async function fillTable(opt: string) {
    const qd = prop.query
    let fetchedDdl = ''
    try {
        switch (opt) {
            case 'Data': {
                const result =
                    await Get<TableData>(`/api/Query/Table/${prop.conName}?Table=${qd.table}&Owner=${qd.owner}`)

                table.value = result
            } break
            case 'Columns': {
                const result =
                    await Get<ColumnsMetadata[]>(`/api/Query/Columns/${prop.conName}?Table=${qd.table}&Owner=${qd.owner}`)
                
                table.value = colToTableData(result)
            } break
            case 'Constraints': {
                const result =
                    await Get<ConstraintMetadata[]>(`/api/Query/Constraints/${prop.conName}?Table=${qd.table}&Owner=${qd.owner}`)

                table.value = conToTableData(result)
            } break
            case 'DDL': {
                const result =
                    await Get<string>(`/api/Query/DDL/${prop.conName}?Type=${prop.conType}&Name=${qd.table}&Owner=${qd.owner}`)

                fetchedDdl = result
            }
            break
        }
    } catch (ex) {
        alert(ex)
        console.log(ex)
    }

    ddl.value = fetchedDdl
    table.value.Name = opt
}
fillTable(prop.options[0])

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
            <li class="nav-item fs-6 inline" v-for="tab in options">
                <div class="nav-link" :class="[table.Name == tab ? 'active' : '']">
                    <a @click="fillTable(tab)">
                        {{ tab }}
                    </a>
                </div>
            </li>
        </ul>
    </div>
    <div class="p-2 border-top w-100 flex-grow-1 overflow-auto">
        <TableRenderer v-if="ddl === ''" :data="table" />
        <div v-if="ddl !== ''" class="mb-3 flex">
            <label for="ddlArea" class="form-label">{{ query.table }} DDL</label>
            <textarea class="form-control flex h-100" wrap="hard" id="ddlArea" rows="25" disabled="true">{{ ddl }}</textarea>
        </div>
    </div>
</template>
