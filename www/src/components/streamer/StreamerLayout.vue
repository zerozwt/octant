<template>
    <n-layout has-sider>
        <n-layout-sider>
            <n-menu :value="menuValue" :options="opts" @update-value="onMenu" :default-expand-all="true"/>
        </n-layout-sider>
        <n-layout-content>
            <slot></slot>
        </n-layout-content>
    </n-layout>
</template>

<script setup>
import {ref, inject, computed} from 'vue'
import {NLayout, NLayoutSider, NLayoutContent, NMenu} from 'naive-ui'
import router from "@/router"

const props = defineProps(["menuValue"])
let menuValue = ref(props.menuValue)
const i18n = inject("octant_locale")

let opts = computed(() => {
    return [
        {
            label: i18n.text.Streamer.Menu[0],
            key: "search",
            children: [
                {label: i18n.text.Streamer.Menu[1], key: "member"},
                {label: i18n.text.Streamer.Menu[2], key: "sc"},
                {label: i18n.text.Streamer.Menu[3], key: "gift"},
            ],
        },
        {label: i18n.text.Streamer.Menu[4], key: "events"},
        {label: i18n.text.Streamer.Menu[5], key: "dm"},
    ]
})

let onMenu = (key, item) => {
    menuValue.value = key

    if (key == "member") {
        router.push("/streamer/data/member")
        return
    }
    if (key == "sc") {
        router.push("/streamer/data/sc")
        return
    }
    if (key == "gift") {
        router.push("/streamer/data/gift")
        return
    }
    if (key == "events") {
        router.push("/streamer/events")
        return
    }
    if (key == "dm") {
        router.push("/streamer/dm")
        return
    }
}
</script>

<style scoped>
.n-layout-sider {
    border-right: solid 1px;
    border-color: #cfcfcf;
    background-color: #fafafa;
}
.n-layout {
    min-height: 100%;
    height: 100%;
}
</style>