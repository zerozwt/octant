<template>
    <div class="donate">
        <div>{{ price ? `￥${price} ` : "" }}{{ time }}</div>
        <div>{{ content }}</div>
    </div>
</template>

<script setup>
const props = defineProps(['time', 'price', 'content'])
</script>

<style scoped>
.donate {
    max-width: 300px;
    background-color: white;
    border-color: rgb(239, 239, 245);
    border: solid 1px rgb(239, 239, 245);
    border-radius: 3px;
}
.donate div:first-child {
    border-color: rgb(239, 239, 245);
    border-bottom: solid 1px rgb(239, 239, 245);
    font-weight: bold;
}
.donate div {
    padding: 4px;
}
</style>