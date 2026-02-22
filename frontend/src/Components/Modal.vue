<script setup lang="ts">
import { ModalType } from '@/Types/ModalType'
import { watch, onUnmounted } from 'vue'

const props = defineProps({
    modelValue: { type: Boolean, required: true },
    modalType: { type: String, required: true },
    btnTxt: { type: String, required: true },
    title: { type: String, required: true },
    confirmTxt: { type: String, default: 'Save Changes' }
})

const emit = defineEmits(['update:modelValue', 'onConfirm', 'onClose'])

function close() {
    try {
        emit('onClose')
    } catch (ex) {
        alert(ex)
    }
    emit('update:modelValue', false)
}

function confirm() {
    try {
        emit('onConfirm')

        close()
    } catch (ex) {
        alert(ex)
    }
}

watch(() => props.modelValue, (val) => {
    if (val) {
        document.body.classList.add('modal-open')
    } else {
        document.body.classList.remove('modal-open')
    }
})

onUnmounted(() => {
    document.body.classList.remove('modal-open')
})
</script>

<template>
    <div v-if="modelValue" class="modal-backdrop fade show"></div>

    <div v-if="modelValue" class="modal fade show" style="display: block;" tabindex="-1">
        <div class="modal-dialog modal-dialog-centered modal-dialog-scrollable">
            <div class="modal-content">

                <div class="modal-header">
                    <h5 class="modal-title">{{ title }}</h5>
                    <button class="btn-close" @click="close"></button>
                </div>

                <div class="modal-body">
                    <slot name="CONTENT" />
                </div>

                <div v-if="modalType === ModalType.FORM" class="modal-footer">
                    <button class="btn btn-secondary" @click="close">Close</button>
                    <button class="btn btn-primary" @click="confirm">
                        {{ confirmTxt }}
                    </button>
                </div>

            </div>
        </div>
    </div>
</template>
