<script setup lang="ts">
import Expanible from '@/Components/Expanible.vue';
import { Options } from '@/Types/Options';

const props = defineProps<{
    conName: string
}>()
</script>

<template>
    <Expanible class="border-bottom" :btn-txt="conName">
        <template #CONTENT>
            <div v-for="(icon, name) in Options">
                <Expanible :btn-txt="name" :idle-icon="icon">
                    <template #CONTENT>
                        <p>Test</p>
                    </template>
                </Expanible>
            </div>
        </template>
    </Expanible>
</template>
