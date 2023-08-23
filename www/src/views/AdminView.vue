<template>
    <I18nFrame>
        <template #header>
            <AdminLoginChecker @login="() => {init = true}"/>
        </template>
        <div v-if="init">
        </div>
    </I18nFrame>
</template>

<script setup>
import {ref} from 'vue'
import I18nFrame from '../components/I18nFrame.vue'
import AdminLoginChecker from '../components/admin/AdminLoginChecker.vue'

let init = ref(false)
</script>

<style scoped>
</style>