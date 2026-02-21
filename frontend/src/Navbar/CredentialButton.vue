<script setup lang="ts">
import { ModalType } from '@/Types/ModalType';
import Modal from '@/Components/Modal.vue';
import { DbCredential } from '@/Types/Credential'
import { computed, ref } from 'vue';

const props = defineProps<{
    credential: DbCredential
    modelValue: boolean
}>()

let showPassword = ref(false)

const emit = defineEmits(['update:modelValue'])

const modelProxy = computed({
get: ()=>props.modelValue,
set: (val: boolean) => emit('update:modelValue', val)
})

const validCredentials = () =>
    props.credential.Database != '' &&
    props.credential.Port > 0 &&
    props.credential.Password != '' &&
    props.credential.User != '';

function PatchCredentials() {
    if (!validCredentials())
        throw new Error('Invalid credentials!');
}
</script>

<template>
    <button button class="btn btn-primary" @click="modelProxy = true">
        Add Credential
    </button>

    <Modal btnTxt="Add credential"
           title="This is test" 
           v-model="modelProxy"
           :modal-type="ModalType.FORM"
           @confirm="PatchCredentials">
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
                <div class="mb-3 form-check">
                    <input type="checkbox" v-model="credential.ShowAll" class="form-check-input" id="showAllInput">
                    <label class="form-check-label" for="showAllInput">Show All Schemas</label>
                </div>
            </form>
        </template>
    </Modal>
</template>
