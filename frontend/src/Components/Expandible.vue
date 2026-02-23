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
let slots = defineSlots<{
    OPTIONS?: {},
    CONTENT: {}
}>()

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
    <li class="my-1 px-2 d-flex flex-column w-100">
        <div class="d-inline-flex bg-white mb-1 rounded w-100">
            <button
                class="p-0 my-1 flex-grow-1 d-flex align-items-center justify-content-start btn-light rounded fs-5 border-0"
                :aria-expanded="expand" @click="beforeExpand" style="background: transparent;">
                <i :class="['px-1', 'bi', expand ? 'bi-caret-down-fill' : idleIcon]" />{{ btnTxt }}
            </button>
            <div class="align-content-center dropdown" v-if="slots.OPTIONS != null">
                <button class="border-0 btn btn-light p-0 d-flex px-2 align-items-center justify-content-center"
                    style="background: transparent;" data-bs-toggle="dropdown" aria-expanded="false">
                    <i class="bi bi-three-dots-vertical" />
                </button>

                <ul class="dropdown-menu dropdown-menu-end">
                    <slot name="OPTIONS" />
                </ul>
            </div>
        </div>
        <div v-if="expand" class="d-flex flex-column">
            <slot name="CONTENT" />
        </div>
    </li>
</template>
