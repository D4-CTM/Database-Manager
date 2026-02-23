<script setup lang="ts">
import { ModalType } from '@/Types/ModalType';
import Modal from '@/Components/Modal.vue';
import { DbCredential } from '@/Types/Credential'
import { computed, ref, watch } from 'vue';

const props = defineProps<{
    credential: DbCredential
    modelValue: boolean
    title: string
}>()

let showPassword = ref(false)
let conName = ref<String>('')

const emit = defineEmits(['update:modelValue', 'update:conName', 'onCancel', 'onConfirm'])

const modelProxy = computed({
get: ()=>props.modelValue,
set: (val: boolean) => emit('update:modelValue', val)
})

watch(() => props.credential.conName, () => {
    conName.value = props.credential.conName
})

const validCredentials = () =>
    props.credential.Database != '' &&
    props.credential.Port > 0 &&
    props.credential.Password != '' &&
    props.credential.User != '';

function onClose() {
    emit('onCancel')
}

function onConfirm() {
    if (!validCredentials())
        throw new Error('Invalid credentials!');

    if (!props.credential.isNew && conName.value.trim() == "")
        throw new Error('Connection name is required!')

    emit('onConfirm', conName.value)
}
</script>

<template>
    <button button class="my-2 py-2 btn btn-primary fs-5" @click="modelProxy = true">
        Add Credential
    </button>

    <Modal btnTxt="Add credential"
           :title="title" 
           v-model="modelProxy"
           :modal-type="ModalType.FORM"
           @onClose="onClose"
           @onConfirm="onConfirm">
        <template #CONTENT>
            <form>
                <div class="row">
                    <div class="col mb-3">
                        <label for="databaseInput" class="form-label">Database</label>
                        <input type="text" v-model.trim="credential.Database" class="form-control" id="databaseInput">
                    </div>
                    <div class="col mb-3">
                        <label for="portInput" class="form-label">Port</label>
                        <input type="text" v-model.number="credential.Port" class="form-control" 
                               pattern=".[\d]*" id="portInput" title="Must be a numeric value">
                    </div>
                </div>
                <div class="mb-3">
                    <label for="serverInput" class="form-label">Server</label>
                    <input type="text" v-model.trim="credential.Server" class="form-control" id="serverInput">
                </div>
                <div class="mb-3">
                    <label for="usernameInput" class="form-label">User</label>
                    <input type="text" v-model.trim="credential.User" class="form-control" id="usernameInput">
                </div>
                <div class="mb-3">
                    <label for="passwordInput" class="form-label">Password</label>
                    <div class="input-group">
                        <input :type="showPassword ? 'text' : 'password'" class="form-control" id="passwordInput"
                            v-model="credential.Password" placeholder="Enter password" />
                        <button class="btn btn-outline-secondary" type="button" @click="showPassword = !showPassword">
                            <i :class="showPassword ? 'bi bi-eye-slash' : 'bi bi-eye'"></i>
                        </button>
                    </div>
                </div>
                <div v-if="!credential.isNew" class="mb-3">
                    <label for="connectionNameInput" class="form-label">Connection Name</label>
                    <input type="text" v-model="conName" class="form-control" id="connectionNameInput">
                </div>
                <div class="mb-3 form-check">
                    <input type="checkbox" v-model="credential.ShowAll" class="form-check-input" id="showAllInput">
                    <label class="form-check-label" for="showAllInput">Show All Schemas</label>
                </div>
            </form>
        </template>
    </Modal>
</template>
