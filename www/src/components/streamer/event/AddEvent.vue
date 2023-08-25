<template>
    <div class="ec-container">
        <div class="el-bar">
            <div><h1>{{ i18n.text.Streamer.Event.New }}</h1></div>
            <div style="flex-grow: 1;"></div>
            <div><n-button type="primary" :disabled="createDisabled" @click="submit">{{ i18n.text.Admin.Streamer.AddRoom.Confirm }}</n-button></div>
        </div>
        <n-card :title="i18n.text.Streamer.Event.Add.Basic" :segmented="{content: true}">
            <n-space vertical>
                <p>{{ i18n.text.Streamer.Event.ListCols[0] }}</p>
                <n-input type="text" v-model:value="name" />
                <p>{{ i18n.text.Streamer.Event.ListCols[1] }}</p>
                <n-input type="textarea" rows="8" v-model:value="content" />
                <n-checkbox v-model:checked="hidden">{{ i18n.text.Streamer.Event.Add.Hidden }}</n-checkbox>
            </n-space>
        </n-card>
        <div style="height: 16px"></div>
        <n-card :title="i18n.text.Streamer.Event.Add.Cond" :segmented="{content: true}" v-if="treeReady">
            <CondTree :tree="condTree"/>
        </n-card>
    </div>
</template>

<script setup>
import {ref, reactive, inject, provide, computed, onMounted} from 'vue'
import {NButton, useMessage, NCard, NInput, NSpace, NCheckbox} from 'naive-ui'
import CondTree from './CondTree.vue'
import {APICaller} from '@/api'
import router from "@/router"
import {CondTreeToReq, CreateCondTreeHandler} from './cond'

const API = APICaller(router)
const i18n = inject("octant_locale")
const message = useMessage()

let name = ref("")
let content = ref("")
let hidden = ref(false)

let gifts = ref([
    {label: "", value: ""},
])

let condTree = reactive({
    cid: 0,
    type: "or",
    subs: [],
})
let nextCID = ref(1)

let treeReady = computed(() => gifts.value.length > 1)

let validate = (node) => {
    if (node.type == "and" || node.type == "or") {
        if (node.subs.length == 0) {
            return false
        }
        for (let i = 0; i < node.subs.length; ++i) {
            if (!validate(node.subs[i])) return false
        }
        return true
    }
    if (node.type == "member") {
        return node.count > 0 && (node.member1 || node.member2 || node.member3)
    }
    if (node.type == "sc") {
        return node.count > 0
    }
    if (node.type == "gift") {
        return node.count > 0 && node.giftID.length > 0
    }
    return false
}

let createDisabled = computed(() => {
    return name.value.length == 0 || content.value.length == 0 || !validate(condTree)
})

let treeEventHandler = CreateCondTreeHandler(condTree, nextCID)
provide("octant_cte", treeEventHandler)

onMounted(() => {
    API.post("/api/gifts", {}).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}] ERROR: ${data.msg}`)
            return
        }
        let list = data.data.list
        let tmp = [{label: "", value: ""}]
        list.forEach((value) => {
            tmp.push({label: value.name, value: String(value.id)})
        })
        gifts.value = tmp
    }).catch(err => message.error(JSON.stringify(err)))
})

provide("octant_gifts", gifts)

let submit = () => {
    API.post("/api/event/add", {name: name.value, reward: content.value, hidden: hidden.value ? 1 : 0, conditions: CondTreeToReq(condTree)}).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}] ERROR: ${data.msg}`)
            return
        }
        router.push("/streamer/events")
    }).catch(err => message.error(JSON.stringify(err)))
}
</script>

<style scoped>
.el-bar {
    margin-bottom: 16px;
    display: flex;
    place-items: center;
}
.ec-container {
    max-width: 1280px;
    margin: 0 auto;
}
</style>