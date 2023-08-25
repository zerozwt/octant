<template>
    <div class="el-bar">
        <div><h1>{{ i18n.text.Streamer.Menu[4] }}</h1></div>
        <div style="flex-grow: 1;"></div>
        <div><router-link to="/streamer/events/new"><n-button type="primary" secondary>{{ i18n.text.Streamer.Event.New }}</n-button></router-link></div>
    </div>
    <n-data-table
        remote
        ref="table"
        :columns="tableCols"
        :data="tableData"
        :loading="loading"
        :pagination="tablePage"
        :row-key="(row) => row.id"
        @update:page="onPageChange"
    />
    <n-modal
        v-model:show="showModal"
        preset="card"
        style="max-width: 600px;"
        :mask-closable="false"
        :title="i18n.text.Streamer.Event.ListOps[1]"
        :segmented="{content: 'soft'}">
        <n-space vertical>
            <p>{{ i18n.text.Streamer.Event.ListCols[0] }}</p>
            <n-input type="text" v-model:value="currEvent.name" />
            <p>{{ i18n.text.Streamer.Event.ListCols[1] }}</p>
            <n-input type="text" v-model:value="currEvent.content" />
        </n-space>
        <template #action>
            <n-button type="primary" block strong @click="onUpdateEvent" :loading="editing" :disabled="editBtnDisable">{{ i18n.text.Admin.Streamer.AddRoom.Confirm }}</n-button>
        </template>
    </n-modal>
</template>

<script setup>
import {ref, reactive, inject, computed, h, onMounted} from 'vue'
import {RouterLink} from 'vue-router'
import {NButton, NSpace, NPopconfirm, NDataTable, useMessage, NModal, NInput} from 'naive-ui'
import {APICaller} from '@/api'
import router from "@/router"

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

let currEvent = reactive({
    id: 0,
    name: "",
    content: "",
})

let showModal = ref(false)
let editing = ref(false)
let editBtnDisable = computed(() => {
    return editing.value || currEvent.name.length == 0 || currEvent.content.length == 0
})

let onEditEvent = (row) => {
    currEvent.id = row.id
    currEvent.name = row.name
    currEvent.content = row.content
    showModal.value = true
}

let onUpdateEvent = () => {
    editing.value = true
    API.post("/api/event/modify", {id: currEvent.id, name: currEvent.name, reward: currEvent.content}).then(rps => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}] ERROR: ${data.msg}`)
            return
        }
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {
        editing.value = false
        showModal.value = false
        onPageChange(tablePage.page)
    })
}

let onDelEvent = (row) => {
    API.post("/api/event/delete", {id: Number(row.id)}).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}] ERROR: ${data.msg}`)
            return
        }
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {onPageChange(tablePage.page)})
}

let tableCols = computed(() => {
    return [
        {key: "name", title: i18n.text.Streamer.Event.ListCols[0]},
        {key: "content", title: i18n.text.Streamer.Event.ListCols[1], ellipsis: {tooltip: true}},
        {
            key: "status",
            title: i18n.text.Streamer.Event.ListCols[2],
            width: 300,
            render(row) {
                return h("div", {class: `estatus_${row.status}`}, i18n.text.Streamer.Event.EvtStatus[row.status-1])
            }
        },
        {
            key: "ops",
            title: i18n.text.Streamer.Event.ListCols[3],
            width: 300,
            render(row) {
                return h(NSpace, {}, () => {
                    let ret = []
                    if (row.status == 4) {
                        ret.push(h(
                            RouterLink,
                            {to: `/streamer/event/${row.id}`},
                            () => h(NButton, {size: "tiny", type: "primary"}, () => i18n.text.Streamer.Event.ListOps[0])
                        ))
                    }
                    ret.push(h(NButton, {size: "tiny", type: "info", onClick: () => {onEditEvent(row)}}, () => i18n.text.Streamer.Event.ListOps[1]))
                    if (row.status != 2) {
                        ret.push(h(
                            NPopconfirm,
                            { onPositiveClick: () => {onDelEvent(row)} },
                            {
                                default: () => i18n.text.Streamer.Event.ListDelConfirm,
                                trigger: () => h(NButton, {size: "tiny", type: "error"}, () => i18n.text.Streamer.Event.ListOps[2])
                            },
                        ))
                    }
                    return ret
                })
            }
        },
    ]
})

let loadData = (page, size) => {
    loading.value = true
    API.get("/api/event/list", {params: {page: page, size: 10}}).then(rsp => {
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
.el-bar {
    margin-bottom: 16px;
    display: flex;
    place-items: center;
}
</style>

<style>
.estatus_2 {
    color: rgb(16, 96, 201)
}
.estatus_3 {
    color: rgb(208, 48, 80);
}
.estatus_4 {
    color: rgb(13, 124, 68);
}
</style>