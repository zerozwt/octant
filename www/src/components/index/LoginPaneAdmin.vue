<template>
    <div v-if="loading"><n-skeleton text :repeat="3" /></div>
    <div v-else>
        <div>
            <div v-if="needLogin">
                <n-space vertical>
                    <p>{{ i18n.text.Index.LoginPane.Pass }}</p>
                    <n-input type="password" show-password-on="mousedown" v-model:value="pass"/>
                    <n-button type="primary" block strong @click="onLogin" :disabled="btnDisable" :loading="logging">{{ i18n.text.Index.LoginButtonText }}</n-button>
                </n-space>
            </div>
            <div v-else>
                <n-space vertical>
                    <p>{{ i18n.text.Index.LoginPane.Admin.AlreadyLogin }}</p>
                    <n-button type="primary" block strong @click="() => {router.push('/admin')}">{{ i18n.text.Index.LoginPane.Admin.Jump }}</n-button>
                </n-space>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, inject, computed } from 'vue'
import { NSkeleton, useMessage, NSpace, NButton, NInput } from 'naive-ui'
import axios from 'axios'
import router from "@/router"

const message = useMessage()
const i18n = inject("octant_locale")

let loading = ref(true)
let needLogin = ref(true)
let pass = ref("")

let loadStatus = () => {
    loading.value = true
    axios.get("/api/admin/login").then(rsp => {
        let data = rsp.data
        needLogin.value = data.code != 0
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {loading.value = false})
}

onMounted(loadStatus)

let logging = ref(false)
let btnDisable = computed(() => {
    return logging.value || pass.value == ""
})

let onLogin = () => {
    logging.value = true
    axios.post("/api/admin/login", {password: pass.value}).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}] ERROR: ${data.msg}`)
            return
        }
        router.push('/admin')
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {logging.value = false})
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