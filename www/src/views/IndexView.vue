<template>
    <I18nFrame>
        <div class="container">
            <div class="box">
                <header>
                    <img alt="Octant logo" class="logo" src="/logo.png"/>
                </header>
                <div>
                    <IndexContent />
                </div>
            </div>
        </div>
    </I18nFrame>
</template>

<script setup>
import I18nFrame from '../components/I18nFrame.vue'
import IndexContent from '../components/index/IndexContent.vue'
</script>

<style scoped>
header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
}

.logo {
    margin: 0 2rem 0 0;
}

.container {
    display: flex;
    place-items: center;
    min-height: 100%;
    height: 100%;
}

.box {
    display: flex;
    place-items: center;
    margin: 0 auto;
}
</style>