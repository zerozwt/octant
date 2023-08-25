<template>
    <n-data-table
        remote
        ref="table"
        :columns="tableCols"
        :data="tableData"
        :loading="loading"
        :pagination="tablePage"
        :row-key="(row) => row.uid"
        @update:page="onPageChange"
    />
</template>

<script setup>
import {ref, reactive, inject, computed, onMounted, h} from 'vue'
import {useMessage, NButton, NDataTable, NSpace} from 'naive-ui'
import {APICaller} from '@/api'
import router from "@/router"
import {CondHasType} from './cond'
import DonateItem from './DonateItem.vue'

const props = defineProps(['detail'])
const API = APICaller(router)
const i18n = inject("octant_locale")
const message = useMessage()

let loading = ref(false)
let tableData = ref([])
let tablePage = reactive({
    page: 1,
    pageCount: 1,
    pageSize: 10,
})

let blockUser = (row) => {
    let url = `/api/event/user/${row.block ? "un" : ""}block`
    API.post(url, {event_id: Number(props.detail.id), uid: Number(row.uid)}).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}] ERROR: ${data.msg}`)
            return
        }
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {onPageChange(1)})
}

const donateWidth = 320

let tableCols = computed(() => {
    let ret = [
        {key: "uid", title: i18n.text.Streamer.Data.Cols[0]},
        {key: "name", title: i18n.text.Streamer.Data.Cols[1]},
    ]
    if (CondHasType(props.detail.conditions, "member")) {
        ret.push({
            key: "member",
            title: i18n.text.Streamer.Menu[1],
            width: donateWidth,
            render(row) {
                return h(
                    NSpace,
                    {vertical: true},
                    () => row.cols.member.map((item) => h(
                        DonateItem,
                        {time: item.time, price: 0, content: `${i18n.text.Streamer.Data.Member[item.level-1]} x${item.count}`},
                        null,
                    ))
                )
            },
        })
    }
    if (CondHasType(props.detail.conditions, "sc")) {
        ret.push({
            key: "sc",
            title: i18n.text.Streamer.Menu[2],
            width: donateWidth,
            render(row) {
                return h(
                    NSpace,
                    {vertical: true},
                    () => row.cols.sc.map((item) => h(
                        DonateItem,
                        {time: item.time, price: item.price, content: item.content},
                        null,
                    ))
                )
            },
        })
    }
    if (CondHasType(props.detail.conditions, "gift")) {
        ret.push({
            key: "gift",
            title: i18n.text.Streamer.Menu[3],
            width: donateWidth,
            render(row) {
                return h(
                    NSpace,
                    {vertical: true},
                    () => row.cols.gift.map((item) => h(
                        DonateItem,
                        {time: item.time, price: item.price*item.count/1000, content: `${item.gift_name} x${item.count}`},
                        null,
                    ))
                )
            },
        })
    }
    ret.push({
        key: "op",
        title: i18n.text.Streamer.Event.ListCols[i18n.text.Streamer.Event.ListCols.length-1],
        width: 100,
        render(row) {
            return h(
                NButton,
                {size: "tiny", type: row.block ? "primary" : "error", onClick: () => {blockUser(row)}},
                () => i18n.text.Streamer.Event.Detail.Block[row.block ? 0 : 1],
            )
        }
    })
    return ret
})

let loadData = (page, size) => {
    loading.value = true
    API.post("/api/event/user/list", {page: page, size: 10, event_id: Number(props.detail.id)}).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}] ERROR: ${data.msg}`)
            return
        }
        data = data.data
        tablePage.page = page
        tablePage.pageCount = Math.ceil(data.count / 10)
        tablePage.itemCount = data.count
        tableData.value = data.list
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {loading.value = false})
}

let onPageChange = (page) => {loadData(page, tablePage.size)}

onMounted(() => {onPageChange(1)})
</script>

<style scoped>
</style>