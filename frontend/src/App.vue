<script setup lang="ts">
import ConnectionsExplorer from './Navbar/ConnectionsExplorer.vue';
import { DbCredential } from './Types/Credential';

let Cred: DbCredential = {
    Database: '',
    Password: '',
    Server: '',
    ShowAll: false,
    User: '',
    Port: 0
}

</script>

<template>
    <ConnectionsExplorer
        :credential="Cred" />
</template>
