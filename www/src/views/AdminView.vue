<template>
    <I18nFrame>
        <template #header>
            <AdminLoginChecker @login="() => {init = true}"/>
        </template>
        <div v-if="init" class="admin-container">
            <AdminContent />
        </div>
    </I18nFrame>
</template>

<script setup>
import {ref} from 'vue'
import I18nFrame from '../components/I18nFrame.vue'
import AdminLoginChecker from '../components/admin/AdminLoginChecker.vue'
import AdminContent from '../components/admin/AdminContent.vue'

let init = ref(false)
</script>

<style scoped>
.admin-container {
    max-width: 1280px;
    margin: 8px auto;
    padding: 0 8px;
}
</style>