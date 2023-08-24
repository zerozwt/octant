<template>
    <n-modal :show="show" @update-show="onShow">
        <n-card :title="i18n.text.General.ChangePass" style="width: 480px">
            <n-space vertical>
                <p>{{ i18n.text.General.OldPass }}</p>
                <n-input type="password" show-password-on="mousedown" v-model:value="oldPass" />
                <p>{{ i18n.text.General.NewPass }}</p>
                <n-input type="password" show-password-on="mousedown" v-model:value="newPass" />
                <p>{{ i18n.text.General.NewPass2 }}</p>
                <n-input type="password" show-password-on="mousedown" v-model:value="newPass2" />
                <n-button type="primary" block strong :loading="loading" :disabled="disabled" @click="onChangePass">{{ i18n.text.General.ChangePass }}</n-button>
            </n-space>
        </n-card>
    </n-modal>
</template>

<script setup>
import {inject, ref, computed} from 'vue'
import {NModal, useMessage, NCard, NButton, NInput, NSpace} from 'naive-ui'
import {APICaller} from '@/api'
import router from "@/router"

const message = useMessage()
const i18n = inject("octant_locale")

const API = APICaller(router)

const props = defineProps(['show', 'realm'])
const emit = defineEmits(['update:show'])

let oldPass = ref("")
let newPass = ref("")
let newPass2 = ref("")
let loading = ref(false)

let onShow = (value) => {
    if (!value) {
        oldPass.value = ""
        newPass.value = ""
        newPass2.value = ""
    }
    emit("update:show", value)
}

let disabled = computed(() => {
    return loading.value || oldPass.value == "" || newPass.value == "" || newPass.value != newPass2.value
})

let onChangePass = () => {
    loading.value = true
    let url = `/api/${props.realm}/password`
    API.post(url, {old_password: oldPass.value, new_password: newPass.value}).then(rsp => {
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
p {
    font-size: 16px;
}
button {
    margin-top: 8px;
}
</style>