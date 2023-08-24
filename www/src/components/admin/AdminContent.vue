<template>
    <div class="title">
        <div><h2>{{ i18n.text.Admin.Streamer.Title }}</h2></div>
        <div></div>
        <div><n-button type="primary" secondary @click="() => {showAdd = true}">{{ i18n.text.Admin.Streamer.Add }}</n-button></div>
    </div>
    <n-data-table
        remote
        ref="table"
        :columns="tableCols"
        :data="tableData"
        :loading="loading"
        :pagination="tablePage"
        :row-key="(row) => row.room_id"
        @update:page="onPageChange"
    />
    <AdminAddStreamer v-model:show="showAdd" @add="() => {loadData(tablePage.page, tablePage.pageSize)}"/>
    <AdminResetStreamer v-model:show="showReset" :streamer="currentStreamer"/>
</template>

<script setup>
import {ref, reactive, onMounted, inject, computed, h} from 'vue'
import {useMessage, NButton, NDataTable, NSpace, NPopconfirm} from 'naive-ui'
import {APICaller} from '@/api'
import router from "@/router"
import AdminAddStreamer from './AdminAddStreamer.vue'
import AdminResetStreamer from './AdminResetStreamer.vue'

const message = useMessage()
const i18n = inject("octant_locale")

const API = APICaller(router)

let loading = ref(false)
let tableData = ref([])
let tablePage = reactive({
    page: 1,
    pageCount: 1,
    pageSize: 10,
})

let currentStreamer = reactive({
    room_id: 0,
    name: "",
    account_name: "",
})

let showReset = ref(false)

let onResetPass = (row) => {
    currentStreamer.room_id = row.room_id
    currentStreamer.name = row.name
    currentStreamer.account_name = row.account_name
    showReset.value = true
}

let onDeleteStreamer = (row) => {
    API.post("/api/admin/streamer/delete", {id: Number(row.room_id)}).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}] ERROR: ${data.msg}`)
            return
        }
        loadData(tablePage.page, tablePage.pageSize)
    }).catch(err => message.error(JSON.stringify(err)))
}

let tableCols = computed(() => {
    return [
        {key: "room_id", title: i18n.text.Admin.Streamer.Cols[0]},
        {key: "name", title: i18n.text.Admin.Streamer.Cols[1]},
        {key: "account_name", title: i18n.text.Admin.Streamer.Cols[2]},
        {
            key: "room_id",
            title: i18n.text.Admin.Streamer.Cols[3],
            render(row) {
                return h(NSpace, {}, () => [
                    h(NButton, {size: "tiny", type: "primary", onClick: () => {onResetPass(row)}}, () => i18n.text.Admin.Streamer.Reset),
                    h(
                        NPopconfirm,
                        { onPositiveClick: () => {onDeleteStreamer(row)} },
                        {
                            default: () => i18n.text.Admin.Streamer.DelConfirm,
                            trigger: () => h(NButton, {size: "tiny", type: "error"}, () => i18n.text.Admin.Streamer.Delete)
                        },
                    ),
                ])
            },
        },
    ]
})

let loadData = (page, size) => {
    loading.value = true
    API.get("/api/admin/streamer/list", {params: {page: page, size: 10}}).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}] ERROR: ${data.msg}`)
            return
        }
        data = data.data
        tablePage.page = page
        tablePage.pageCount = Math.ceil(data.count / size)
        tablePage.itemCount = data.count
        tableData.value = data.list
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {loading.value = false})
}

let onPageChange = (page) => {loadData(page, tablePage.size)}

onMounted(() => {onPageChange(1)})

// -----------------------------------------------------------------

let showAdd = ref(false)

</script>

<style scoped>
.title {
    display: flex;
    margin-bottom: 16px;
}
.title > div {
    flex-grow: 1;
}
.title > div:first-child, div:last-child {
    flex-grow: 0;
}
</style>