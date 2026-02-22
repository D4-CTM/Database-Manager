<script setup lang="ts">
import { ref } from 'vue';

const props = defineProps({
    btnTxt: {
        type: String,
        required: true
    },
    idleIcon: {
         type: String,
         default: 'bi-caret-right-fill'
    }
})

let emit = defineEmits(['beforeExpand'])

function beforeExpand() {
    if (!expand.value) {
        emit('beforeExpand', 
            () => expand.value = true)
    }
    else expand.value = false
}

let expand = ref(false)
</script>

<template>
    <li class="my-1 px-2 d-flex flex-column">
        <button class="btn p-0 my-1 text-left btn-light rounded fs-5"
            :aria-expanded="expand" @click="beforeExpand">
            <i :class="['px-1', 'bi', expand ? 'bi-caret-down-fill' : idleIcon]"/>{{ btnTxt }}
        </button>
        <div v-if="expand" class="d-flex flex-column">
            <slot name="CONTENT"/>
        </div>
    </li>
</template>
