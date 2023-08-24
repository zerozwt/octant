<template>
    <I18nFrame>
        <template #header><StreamerLoginChecker @login="() => {ready = true}"/></template>
        <div v-if="ready" class="streamer-container">
            <StreamerLayout :menu-value="menuValue"><slot></slot></StreamerLayout>
        </div>
    </I18nFrame>
</template>

<script setup>
import {ref} from 'vue'
import I18nFrame from '../I18nFrame.vue'
import StreamerLayout from './StreamerLayout.vue'
import StreamerLoginChecker from './StreamerLoginChecker.vue'

const props = defineProps(["menuValue"])

let ready = ref(false)
</script>

<style scoped>
.streamer-container {
    height: 100%;
    min-height: 100%;
}
</style>