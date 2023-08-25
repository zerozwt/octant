<template>
    <StreamerFrame menu-value="events"><EventDetail :eid="Number($route.params.id)"/></StreamerFrame>
</template>

<script setup>
import StreamerFrame from '../components/streamer/StreamerFrame.vue'
import EventDetail from '../components/streamer/event/EventDetail.vue'
</script>

<style scoped>
</style>