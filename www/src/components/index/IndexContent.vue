<template>
    <h1>{{ i18n.text.Index.Title }}</h1>
    <p>{{ i18n.text.Index.SubTitle }}</p>
    <n-timeline>
        <n-timeline-item v-for="item in i18n.text.Index.Intro" type="info">
            <template #header>
                <h2>{{ item }}</h2>
            </template>
        </n-timeline-item>
    </n-timeline>
    <n-button type="primary" block strong secondary :disabled="showModal" @click="() => {showModal = true}">{{ i18n.text.Index.LoginButtonText }}</n-button>
    <n-modal v-model:show="showModal">
        <n-card :title="i18n.text.Index.LoginButtonText" style="width: 480px">
            <n-tabs size="large" animated default-value="streamer">
                <n-tab-pane name="streamer" :tab="i18n.text.Index.LoginTypes[0]">
                    <LoginPaneStreamer />
                </n-tab-pane>
                <n-tab-pane name="audience" :tab="i18n.text.Index.LoginTypes[1]">
                    <LoginPaneDD />
                </n-tab-pane>
                <n-tab-pane name="admin" :tab="i18n.text.Index.LoginTypes[2]">
                    <LoginPaneAdmin />
                </n-tab-pane>
            </n-tabs>
        </n-card>
    </n-modal>
</template>

<script setup>
import { inject, ref } from 'vue'
import { NTimeline, NTimelineItem, NButton, NModal, NCard, NTabs, NTabPane } from 'naive-ui'
import LoginPaneStreamer from './LoginPaneStreamer.vue';
import LoginPaneDD from './LoginPaneDD.vue';
import LoginPaneAdmin from './LoginPaneAdmin.vue';

const i18n = inject("octant_locale")

let showModal = ref(false)
</script>

<style scoped>
.n-timeline {
    margin: 32px 8px;
}
h2 {
    margin-top: -6px;
}
</style>