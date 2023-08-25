<template>
    <div v-if="ready == 2">
        <n-card :title="detail.name" :segmented="{content: true}" style="margin-bottom: 16px;">
            <template #header-extra>
                <n-space>
                    <n-button type="primary" secondary @click="() => {showCond = true}">{{ i18n.text.Streamer.Event.Add.Cond }}</n-button>
                    <a :href="`/api/event/user/dl?id=${eid}`" target="_blank"><n-button type="primary" secondary>{{ i18n.text.Streamer.Event.Detail.Download }}</n-button></a>
                </n-space>
            </template>
            <n-space vertical><p style="margin-left: 16px;" v-for="item in reward">{{ item }}</p></n-space>
        </n-card>
        <UserList :detail="detail"/>
        <n-modal
            v-model:show="showCond"
            style="max-width: 960px;"
            preset="card"
            :title="i18n.text.Streamer.Event.Add.Cond"
            :segmented="{content: true}">
            <n-scrollbar style="max-height: 720px;"><CondTree :tree="detail.conditions" :readonly="true"/></n-scrollbar>
        </n-modal>
    </div>
</template>

<script setup>
import {ref, reactive, inject, provide, computed, onMounted} from 'vue'
import {useMessage, NButton, NCard, NModal, NScrollbar, NSpace} from 'naive-ui'
import {APICaller} from '@/api'
import router from "@/router"
import {CondTreeFromReq, ReadOnlyHandler} from './cond'
import CondTree from './CondTree.vue'
import UserList from './UserList.vue'

const API = APICaller(router)
const i18n = inject("octant_locale")
const props = defineProps(['eid'])
const message = useMessage()

let ready = ref(0)
let detail = reactive({})

let gifts = ref([
    {label: "", value: ""},
])
provide("octant_gifts", gifts)
provide("octant_cte", ReadOnlyHandler)

let reward = computed(() => {
    if (!detail.reward) return []
    return detail.reward.split("\n")
})

let showCond = ref(false)

onMounted(() => {
    API.get("/api/event/detail", {params: {id: props.eid}}).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}] ERROR: ${data.msg}`)
            return
        }
        data = data.data
        if (data.status != 4) {
            // not ready
            router.push("/streamer/events")
            return
        }
        detail.id = data.id
        detail.name = data.name
        detail.reward = data.reward
        detail.hidden = data.hidden != 0
        detail.conditions = CondTreeFromReq(data.conditions)
        ready.value++
    }).catch(err => message.error(JSON.stringify(err)))

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
        ready.value++
    }).catch(err => message.error(JSON.stringify(err)))
})
</script>

<style scoped>
</style>