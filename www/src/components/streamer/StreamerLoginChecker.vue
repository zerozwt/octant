<template>
    <div v-if="!loading">
        <n-dropdown :options="menu" @select="onMenu">
            <n-button quaternary>{{ name }}</n-button>
        </n-dropdown>
        <PasswordChanger v-model:show="showModal" realm="streamer"/>
    </div>
</template>

<script setup>
import {ref, onMounted, inject, computed} from 'vue'
import {useMessage, NDropdown, NButton} from 'naive-ui'
import {APICaller} from '@/api'
import router from "@/router"
import PasswordChanger from '../PasswordChanger.vue'

const emit = defineEmits(['login'])

const message = useMessage()
const i18n = inject("octant_locale")

const API = APICaller(router)

let loading = ref(true)
let name = ref("")

onMounted(() => {
    API.get("/api/streamer/login", {}).then((rsp) => {
        loading.value = false
        name.value = rsp.data.data.name
        emit("login")
    }).catch(err => message.error(JSON.stringify(err)))
})

let menu = computed(() => {
    return [
        {label: i18n.text.General.ChangePass, key: "change_pass"},
        {label: i18n.text.General.Logout, key: "logout"},
    ]
})

let onMenu = (key) => {
    if (key == "change_pass") {
        showModal.value = true
        return
    }
    API.get("/api/streamer/logout", {}).then((rsp) => {
        router.push("/")
    }).catch(err => message.error(JSON.stringify(err)))
}

let showModal = ref(false)
</script>

<style scoped>
</style>