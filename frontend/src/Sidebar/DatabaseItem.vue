<script setup lang="ts">
import Expansible from '@/Components/Expansible.vue';
import { Options } from '@/Types/Options';

const props = defineProps<{
    conName: string
}>()
</script>

<template>
    <Expansible class="border-bottom" :btn-txt="conName">
        <template #CONTENT>
            <div v-for="(icon, name) in Options">
                <Expansible :btn-txt="name" :idle-icon="icon">
                    <template #CONTENT>
                        <p>Test</p>
                    </template>
                </Expansible>
            </div>
        </template>
    </Expansible>
</template>
