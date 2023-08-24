<template>
    <n-modal preset="card" style="width: 480px" :title="i18n.text.Admin.Streamer.Add" :mask-closable="false" :show="show" @update-show="onShow">
        <n-space vertical>
            <p>{{ i18n.text.Admin.Streamer.AddRoom.ID }}</p>
            <n-input-number v-model:value="RoomID" :show-button="false" />
            <p>{{ i18n.text.Admin.Streamer.AddRoom.Name }}</p>
            <n-input type="text" v-model:value="Account"/>
            <p>{{ i18n.text.Admin.Streamer.AddRoom.Pass }}</p>
            <n-input type="password" v-model:value="Pass" show-password-on="mousedown"/>
            <p>{{ i18n.text.Admin.Streamer.AddRoom.Pass2 }}</p>
            <n-input type="password" v-model:value="Pass2" show-password-on="mousedown"/>
        </n-space>
        <template #action>
            <n-button type="primary" block strong :disabled="commitDisable" :loading="loading" @click="commit">
                {{ i18n.text.Admin.Streamer.AddRoom.Confirm }}
            </n-button>
        </template>
    </n-modal>
</template>

<script setup>
import {ref, inject, computed} from 'vue'
import {useMessage, NButton, NSpace, NInput, NModal, NInputNumber} from 'naive-ui'
import {APICaller} from '@/api'
import router from "@/router"

const message = useMessage()
const i18n = inject("octant_locale")
const props = defineProps(['show'])
const emit = defineEmits(['update:show', 'add'])

const API = APICaller(router)

let RoomID = ref(0)
let Account = ref("")
let Pass = ref("")
let Pass2 = ref("")

let loading = ref(false)

let commitDisable = computed(() => {
    return loading.value || RoomID.value <= 0 || Account.value == "" || Pass.value == "" || Pass.value != Pass2.value
})

let onShow = (value) => {
    if (!value) {
        RoomID.value = 0
        Account.value = ""
        Pass.value = ""
        Pass2.value = ""
    }
    emit("update:show", value)
}

let commit = () => {
    loading.value = true
    API.post("/api/admin/streamer/add", {room_id: RoomID.value, name: Account.value, password: Pass.value}).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}] ERROR: ${data.msg}`)
            return
        }
        onShow(false)
        emit("add")
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {
        loading.value = false
    })
}
</script>

<style scoped>
</style>