<template>
    <div class="header">
        <div><img alt="Octant logo" class="logo" src="/logo.png" width="48" height="48"/></div>
        <div>Octant</div>
        <div class="header-bar"></div>
        <div class="slot">
            <n-config-provider :locale="locale" :date-locale="dateLocale">
                <n-message-provider><slot name="header"></slot></n-message-provider>
            </n-config-provider>
        </div>
        <div>
            <n-popselect :options="menuOpts" trigger="click" @update:value="onLangSelete" :value="lang">
                <n-button quaternary>{{ menuText }}</n-button>
            </n-popselect>
        </div>
    </div>
    <div class="content">
        <n-config-provider :locale="locale" :date-locale="dateLocale">
            <n-message-provider>
                <slot></slot>
            </n-message-provider>
        </n-config-provider>
    </div>
    <div class="footer"><div>powered by octant</div></div>
</template>

<script setup>
import { ref, computed, provide, reactive } from 'vue'
import { NConfigProvider, NButton, NMessageProvider, NPopselect } from 'naive-ui'
import {zhCN, dateZhCN, enUS, dateEnUS, jaJP, dateJaJP} from 'naive-ui'
import {loadCache, defaultLang, defaultLocaleObj, localeObj} from '@/locale'

let lang = ref(loadCache("octant_lang", defaultLang()))
let i18n = reactive({text: defaultLocaleObj()})

provide("octant_locale", i18n)

let menuText = computed(() => {
    if (lang.value == "en-US") return "Language"
    if (lang.value == "ja-JP") return "言語"
    return "语言"
})

let locale = computed(() => {
    if (lang.value == "en-US") return enUS
    if (lang.value == "ja-JP") return jaJP
    return zhCN
})

let dateLocale = computed(() => {
    if (lang.value == "en-US") return dateEnUS
    if (lang.value == "ja-JP") return dateJaJP
    return dateZhCN
})

let menuOpts = [
    {label: "中文", value: "zh-CN"},
    {label: "English", value: "en-US"},
    {label: "日本語", value: "ja-JP"},
]

let onLangSelete = (key) => {
    lang.value = key
    window.localStorage.setItem("octant_lang", key)
    i18n.text = localeObj(key)
}
</script>

<style scoped>
.header {
    flex: 0 1 48px;
    border-bottom: solid 1px;
    border-color: #cfcfcf;

    display: flex;
    place-items: center;
}

.header-bar {
    flex-grow: 1;
}

.header > div {
    margin: 0 8px;
}

.slot {
    display: flex;
    place-items: center;
    place-content: center;
}

.content {
    flex: 1 1 auto;
}

.n-config-provider {
    height: 100%;
}

.footer {
    flex: 0 1 40px;
    display: flex;
    border-top: solid 1px;
    border-color: #cfcfcf;
}
.footer > div {
    margin: 0 auto;
    padding: 16px;
}
</style>