<template>
    <n-modal preset="card" style="width: 480px" :title="i18n.text.Admin.Streamer.Reset" :mask-closable="false" :show="show" @update-show="onShow">
        <n-space vertical>
            <n-alert :title="i18n.text.Admin.Streamer.ResetPass.Warn" type="warning"> </n-alert>
            <p>{{ `${i18n.text.Admin.Streamer.AddRoom.ID} ${streamer.room_id}` }}</p>
            <p>{{ `${i18n.text.Admin.Streamer.AddRoom.Name} ${streamer.account_name}` }}</p>
            <p>{{ `${i18n.text.Admin.Streamer.ResetPass.Name} ${streamer.name}` }}</p>
            <p>{{ i18n.text.General.NewPass }}</p>
            <n-input type="password" v-model:value="Pass" show-password-on="mousedown"/>
            <p>{{ i18n.text.General.NewPass2 }}</p>
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
import {useMessage, NButton, NSpace, NInput, NModal, NAlert} from 'naive-ui'
import {APICaller} from '@/api'
import router from "@/router"

const message = useMessage()
const i18n = inject("octant_locale")
const props = defineProps(['show', 'streamer'])
const emit = defineEmits(['update:show'])

const API = APICaller(router)

let Pass = ref("")
let Pass2 = ref("")
let loading = ref(false)

let commitDisable = computed(() => {
    return loading.value || Pass.value == "" || Pass.value != Pass2.value
})

let onShow = (value) => {
    if (!value) {
        Pass.value = ""
        Pass2.value = ""
    }
    emit("update:show", value)
}

let commit = () => {
    loading.value = true
    API.post("/api/admin/streamer/reset", {id: Number(props.streamer.room_id), password: Pass.value}).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}] ERROR: ${data.msg}`)
            return
        }
        onShow(false)
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {
        loading.value = false
    })
}
</script>

<style scoped>
</style>