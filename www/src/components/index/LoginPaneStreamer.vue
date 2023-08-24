<template>
    <div v-if="loading"><n-skeleton text :repeat="5" /></div>
    <div v-else>
        <div>
            <div v-if="needLogin">
                <n-space vertical>
                    <p>{{ i18n.text.Index.LoginPane.Name }}</p>
                    <n-input type="text" v-model:value="name"/>
                    <p>{{ i18n.text.Index.LoginPane.Pass }}</p>
                    <n-input type="password" show-password-on="mousedown" v-model:value="pass"/>
                    <n-button type="primary" block strong @click="onLogin" :disabled="btnDisable" :loading="logging">{{ i18n.text.Index.LoginButtonText }}</n-button>
                </n-space>
            </div>
            <div v-else>
                <n-space vertical>
                    <p>{{ i18n.text.Index.LoginPane.Streamer.Already(name) }}</p>
                    <n-button type="primary" block strong @click="() => {router.push('/streamer/data/member')}">{{ i18n.text.Index.LoginPane.Streamer.Jump }}</n-button>
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
let name = ref("")
let pass = ref("")

onMounted(() => {
    loading.value = true
    axios.get("/api/streamer/login").then(rsp => {
        let data = rsp.data
        needLogin.value = data.code != 0
        if (data.code == 0) {
            name.value = data.data.name
        }
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {loading.value = false})
})

let logging = ref(false)
let btnDisable = computed(() => {
    return logging.value || name.value == "" || pass.value == ""
})

let onLogin = () => {
    logging.value = true
    axios.post("/api/streamer/login", {name: name.value, password: pass.value}).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}] ERROR: ${data.msg}`)
            return
        }
        router.push('/streamer/data/member')
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