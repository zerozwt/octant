<template>
    <div v-if="isGroup">
        <n-card :title="i18n.text.Streamer.Event.Add.Group.Title" :segmented="{content: true}" size="small">
            <template #header-extra>
                <n-button type="error" size="small" v-if="node.cid > 0 && !readonly" @click="delSelf">{{ i18n.text.Streamer.Event.Add.Delete }}</n-button>
            </template>
            <div class="group-type-container">
                <p>{{ i18n.text.Streamer.Event.Add.Group.Content[0] }}</p>
                <n-select
                    :value="node.type"
                    :options="groupOpts"
                    @update:value="(value) => {teh.onGroupChangeType(node, value)}"
                    style="max-width: 100px; margin: 0 8px;"
                />
                <p>{{ i18n.text.Streamer.Event.Add.Group.Content[1] }}</p>
                <div style="flex-grow: 1;"></div>
                <n-space v-if="!readonly">
                    <n-dropdown trigger="click" :options="addCondOpts" @select="(key) => {teh.addCond(node, key)}" :show-arrow="true">
                        <n-button type="primary" size="small" secondary>{{ i18n.text.Streamer.Event.Add.Group.Add }}</n-button>
                    </n-dropdown>
                    <n-button type="primary" size="small" secondary @click="() => {teh.addSubGroup(node)}">{{ i18n.text.Streamer.Event.Add.Group.AddGroup }}</n-button>
                </n-space>
            </div>
            <n-space vertical><slot></slot></n-space>
        </n-card>
    </div>
    <div v-else>
        <n-card :title="condTitle" :segmented="{content: true}" size="small">
            <template #header-extra>
                <n-button type="error" size="small" @click="delSelf" v-if="!readonly">{{ i18n.text.Streamer.Event.Add.Delete }}</n-button>
            </template>
            <n-grid x-gap="8" y-gap="8" :cols="totalCols">
                <n-gi><div>{{ i18n.text.Streamer.Data.TimeRange }}</div></n-gi>
                <n-gi :span="valueCols">
                    <div class="cond-value">
                        <n-date-picker type="datetimerange" :value="node.timeRange" @update:value="(value) => {teh.updateTimeRange(node, value)}"/>
                    </div>
                </n-gi>
                <n-gi :span="totalCols" v-if="isSC">
                    <div class="cond-value cond-sc">
                        <n-select
                            :value="node.mode"
                            :options="scModeOpts"
                            @update:value="(value) => {teh.updateMode(node, value)}"
                        />
                        <p>{{ i18n.text.Streamer.Event.Add.SC.Content[0] }}</p>
                        <n-input-number :value="node.count" :show-button="false" @update:value="(value) => {teh.updateCount(node, value)}"><template #prefix>￥</template></n-input-number>
                        <p>{{ i18n.text.Streamer.Event.Add.SC.Content[1] }}</p>
                    </div>
                </n-gi>
                <n-gi v-if="isMember"><div>{{ i18n.text.Streamer.Data.MemberLevel }}</div></n-gi>
                <n-gi :span="4" v-if="isMember">
                    <div class="cond-value">
                        <n-space>
                            <n-checkbox :checked="node.member1" @update:checked="(value) => {teh.updateMember1(node, value)}">{{ i18n.text.Streamer.Data.Member[0] }}</n-checkbox>
                            <n-checkbox :checked="node.member2" @update:checked="(value) => {teh.updateMember2(node, value)}">{{ i18n.text.Streamer.Data.Member[1] }}</n-checkbox>
                            <n-checkbox :checked="node.member3" @update:checked="(value) => {teh.updateMember3(node, value)}">{{ i18n.text.Streamer.Data.Member[2] }}</n-checkbox>
                        </n-space>
                    </div>
                </n-gi>
                <n-gi v-if="isMember"><div>{{ i18n.text.Streamer.Event.Add.Member.Count }}</div></n-gi>
                <n-gi :span="4" v-if="isMember">
                    <div class="cond-value cond-member">
                        <n-select
                            :value="node.mode"
                            :options="memberModeOpts"
                            @update:value="(value) => {teh.updateMode(node, value)}"
                        />
                        <p>{{ i18n.text.Streamer.Event.Add.Member.Content[0] }}</p>
                        <n-input-number :value="node.count" :show-button="false" @update:value="(value) => {teh.updateCount(node, value)}"/>
                        <p>{{ i18n.text.Streamer.Event.Add.Member.Content[1] }}</p>
                    </div>
                </n-gi>
                <n-gi v-if="isGift"><div>{{ i18n.text.Streamer.Menu[3] }}</div></n-gi>
                <n-gi :span="4" v-if="isGift">
                    <div class="cond-value cond-sc"><n-select :value="node.giftID" :options="gifts" @update:value="(value) => {teh.updateGift(node, value)}"/></div>
                </n-gi>
                <n-gi v-if="isGift"><div>{{ i18n.text.Streamer.Event.Add.Gift.Count }}</div></n-gi>
                <n-gi :span="4" v-if="isGift">
                    <div class="cond-value cond-member">
                        <n-select
                            :value="node.mode"
                            :options="memberModeOpts"
                            @update:value="(value) => {teh.updateMode(node, value)}"
                        />
                        <p>{{ i18n.text.Streamer.Event.Add.Gift.Content[0] }}</p>
                        <n-input-number :value="node.count" :show-button="false" @update:value="(value) => {teh.updateCount(node, value)}"/>
                        <p>{{ i18n.text.Streamer.Event.Add.Gift.Content[1] }}</p>
                    </div>
                </n-gi>
            </n-grid>
        </n-card>
    </div>
</template>

<script setup>
import {inject, computed} from 'vue'
import {NCard, NButton, NSelect, NSpace, NDropdown, NGrid, NGi, NDatePicker, NInputNumber, NCheckbox} from 'naive-ui'

const props = defineProps(['node', 'readonly'])
const i18n = inject("octant_locale")
const teh = inject("octant_cte")
const gifts = inject("octant_gifts")

const isGroup = props.node.type == "and" || props.node.type == "or"
const isMember = props.node.type == "member"
const isSC = props.node.type == "sc"
const isGift = props.node.type == "gift"

const totalCols = 5
const valueCols = totalCols - 1

let groupOpts = computed(() => {
    return [
        {label: i18n.text.Streamer.Event.Add.Group.Opts[0], value: "or"},
        {label: i18n.text.Streamer.Event.Add.Group.Opts[1], value: "and"},
    ]
})

let addCondOpts = computed(() => {
    return [
        {label: i18n.text.Streamer.Menu[1], key: "member"},
        {label: i18n.text.Streamer.Menu[2], key: "sc"},
        {label: i18n.text.Streamer.Menu[3], key: "gift"},
    ]
})

let condTitle = computed(() => {
    if (isMember) {
        return i18n.text.Streamer.Menu[1]
    }
    if (isSC) {
        return i18n.text.Streamer.Menu[2]
    }
    if (isGift) {
        return i18n.text.Streamer.Menu[3]
    }
    return ""
})

let delSelf = () => {teh.deleteNode(props.node.cid)}

let scModeOpts = computed(() => {
    return [
        {label: i18n.text.Streamer.Event.Add.SC.Opts[0], value: "total"},
        {label: i18n.text.Streamer.Event.Add.SC.Opts[1], value: "once"},
    ]
})

let memberModeOpts = computed(() => {
    return [
        {label: i18n.text.Streamer.Event.Add.Member.Opts[0], value: "total"},
        {label: i18n.text.Streamer.Event.Add.Member.Opts[1], value: "once"},
    ]
})
</script>

<style scoped>
.group-type-container {
    display: flex;
    place-items: center;
    margin-bottom: 8px;
}
.cond-value{
    display: flex;
    place-items: center;
}
.cond-sc .n-select {
    max-width: 300px;
    margin-right: 8px;
}
.cond-sc .n-input-number {
    margin: 0 8px;
}
.cond-member .n-select {
    max-width: 160px;
    margin-right: 8px;
}
.cond-member .n-input-number {
    margin: 0 8px;
}
</style>