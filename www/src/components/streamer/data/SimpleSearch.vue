<template>
    <n-card :title="i18n.text.Streamer.Menu[0]" :segmented="{content: true}">
        <n-grid :cols="5" x-gap="8" y-gap="8">
            <n-gi><div>{{ i18n.text.Streamer.Data.TimeRange }}</div></n-gi>
            <n-gi :span="4"><div class="filter-var"><n-date-picker v-model:value="timeRange" type="datetimerange"/></div></n-gi>
            <n-gi><div>{{ i18n.text.Streamer.Data.Cols[0] }}</div></n-gi>
            <n-gi :span="4"><div class="filter-var"><n-input type="text" v-model:value="uid" :allow-input="(value) => !value || /^\d+$/.test(value)"/></div></n-gi>
            <n-gi><div>{{ i18n.text.Streamer.Data.Cols[1] }}</div></n-gi>
            <n-gi :span="4"><div class="filter-var"><n-input type="text" v-model:value="name"/></div></n-gi>
            <n-gi v-if="isMember"><div>{{ i18n.text.Streamer.Data.MemberLevel }}</div></n-gi>
            <n-gi :span="4" v-if="isMember">
                <div class="filter-var">
                    <n-space>
                        <n-checkbox v-model:checked="member1">{{ i18n.text.Streamer.Data.Member[0] }}</n-checkbox>
                        <n-checkbox v-model:checked="member2">{{ i18n.text.Streamer.Data.Member[1] }}</n-checkbox>
                        <n-checkbox v-model:checked="member3">{{ i18n.text.Streamer.Data.Member[2] }}</n-checkbox>
                    </n-space>
                </div>
            </n-gi>
            <n-gi v-if="isSC"><div>{{ i18n.text.Streamer.Data.SCContent }}</div></n-gi>
            <n-gi v-if="isSC" :span="4"><div class="filter-var"><n-input type="text" v-model:value="scContent"/></div></n-gi>
            <n-gi v-if="isGift"><div>{{ i18n.text.Streamer.Menu[3] }}</div></n-gi>
            <n-gi v-if="isGift" :span="4"><div class="filter-var"><n-select v-model:value="giftID" :options="gifts"/></div></n-gi>
            <n-gi></n-gi>
            <n-gi><div class="filter-var"><n-button type="primary" @click="() => {onPageChange(1)}">{{ i18n.text.Streamer.Data.Search }}</n-button></div></n-gi>
        </n-grid>
    </n-card>
    <n-data-table
        remote
        ref="table"
        :columns="tableCols"
        :data="tableData"
        :loading="loading"
        :pagination="tablePage"
        :row-key="(row) => row.key"
        @update:page="onPageChange"
    />
</template>

<script setup>
import {ref, reactive, inject, computed, h, onMounted} from 'vue'
import {NCard, NDataTable, NInput, NCheckbox, NButton, NDatePicker, useMessage, NGrid, NGi, NSpace, NSelect} from 'naive-ui'
import router from "@/router"
import {APICaller} from '@/api'
import dayjs from 'dayjs'

const renderTime = (ts) => dayjs.unix(Math.floor(ts/1000)).format("YYYYMMDDHHmmss")

const props = defineProps(["type"])
const i18n = inject("octant_locale")
const message = useMessage()

let dummyKey = ref(0)

let timeRange = ref([Date.now()-7*24*3600*1000, Date.now()])
let uid = ref("")
let name = ref("")
let giftID = ref("")
let scContent = ref("")
let member1 = ref(true)
let member2 = ref(true)
let member3 = ref(true)

const API = APICaller(router)

let loading = ref(false)
let tableData = ref([])
let tablePage = reactive({
    page: 1,
    pageCount: 1,
    pageSize: 10,
})

let isMember = props.type == "member"
let isSC = props.type == "sc"
let isGift = props.type == "gift"

let gifts = ref([
    {label: "", value: ""},
])

let tableCols = computed(() => {
    if (isMember) {
        return [
            {key: "uid", title: i18n.text.Streamer.Data.Cols[0]},
            {key: "name", title: i18n.text.Streamer.Data.Cols[1]},
            {key: "time", title: i18n.text.Streamer.Data.Cols[2]},
            {
                key: "member",
                title: i18n.text.Streamer.Menu[1],
                render(row) {
                    return h("span", {}, `${i18n.text.Streamer.Data.Member[row.guard.level-1]} x${row.guard.count}`)
                }
            },
        ]
    }
    if (isSC) {
        return [
            {key: "uid", title: i18n.text.Streamer.Data.Cols[0]},
            {key: "name", title: i18n.text.Streamer.Data.Cols[1]},
            {key: "time", title: i18n.text.Streamer.Data.Cols[2]},
            {
                key: "sc",
                title: i18n.text.Streamer.Menu[2],
                render(row) {
                    return h("span", {}, `￥${row.sc.price} ${row.sc.content}`)
                }
            },
        ]
    }
    if (isGift) {
        return [
            {key: "uid", title: i18n.text.Streamer.Data.Cols[0]},
            {key: "name", title: i18n.text.Streamer.Data.Cols[1]},
            {key: "time", title: i18n.text.Streamer.Data.Cols[2]},
            {
                key: "gift",
                title: i18n.text.Streamer.Menu[3],
                render(row) {
                    return h("span", {}, `${row.gift.name} x${row.gift.count}`)
                }
            },
        ]
    }
    return [
        {key: "uid", title: i18n.text.Streamer.Data.Cols[0]},
        {key: "name", title: i18n.text.Streamer.Data.Cols[1]},
        {key: "time", title: i18n.text.Streamer.Data.Cols[2]},
    ]
})

let loadData = (page, size) => {
    let req = {
        datasource: props.type,
        page: page,
        size: 10,
        start_time: renderTime(timeRange.value[0]),
        end_time: renderTime(timeRange.value[1]),
        filter: {
            sender_uid: 0,
            sender_name: "",
            gift_id: 0,
            sc_content: "",
            guard_level: [],
        },
    }
    if (uid.value.length > 0) {
        req.filter.sender_uid = Number(uid.value)
    }
    if (name.value.length > 0) {
        req.filter.sender_name = name.value
    }
    if (giftID.value.length > 0) {
        req.filter.gift_id = Number(giftID.value)
    }
    if (scContent.value.length > 0) {
        req.filter.sc_content = scContent.value
    }
    if (member1.value) {
        req.filter.guard_level.push(1)
    }
    if (member2.value) {
        req.filter.guard_level.push(2)
    }
    if (member3.value) {
        req.filter.guard_level.push(3)
    }

    loading.value = true
    API.post("/api/simple_search", req).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}] ERROR: ${data.msg}`)
            return
        }
        data = data.data
        tablePage.page = page
        tablePage.pageCount = Math.ceil(data.count / 10)
        tablePage.itemCount = data.count

        let tmp = []
        data.list.forEach((value) => {
            value.key = dummyKey.value++
            tmp.push(value)
        })
        tableData.value = tmp
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {loading.value = false})
}

let onPageChange = (page) => {loadData(page, tablePage.size)}

onMounted(() => {
    onPageChange(1)
    if (isGift) {
        API.post("/api/gifts", {}).then(rsp => {
            let data = rsp.data
            if (data.code != 0) {
                message.error(`[${data.code}] ERROR: ${data.msg}`)
                return
            }
            let list = data.data.list
            list.forEach((value) => {
                gifts.value.push({label: value.name, value: String(value.id)})
            })
        }).catch(err => message.error(JSON.stringify(err)))
    }
})
</script>

<style scoped>
.n-card {
    margin-bottom: 16px;
}
.filter-var {
    display: flex;
}
.filter-var .n-input,.n-select {
    max-width: 480px;
}
</style>