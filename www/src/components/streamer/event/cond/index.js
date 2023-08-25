import dayjs from 'dayjs'

const CondTreeToReq = (node) => {
    if (node.type == "or" || node.type == "and") {
        if (node.subs.length == 1) {
            return CondTreeToReq(node.subs[0])
        }
        let ret = {type: node.type, sub_conditions: []}
        node.subs.forEach((value) => {
            ret.sub_conditions.push(CondTreeToReq(value))
        })
        return ret
    }
    let levels = []
    if (node.member1) levels.push(1)
    if (node.member2) levels.push(2)
    if (node.member3) levels.push(3)
    return {
        type: node.type,
        start_time: dayjs.unix(node.timeRange[0]/1000).format("YYYYMMDDHHmmss"),
        end_time: dayjs.unix(node.timeRange[1]/1000).format("YYYYMMDDHHmmss"),
        mode: node.mode,
        count: node.count ? Number(node.count) : 0,
        gift_id: node.giftID ? Number(node.giftID) : 0,
        guard_levels: levels,
    }
}

let cidCounter = 0

const tsFromStr = (str) => dayjs(str, "YYYYMMDDHHmmss").unix()*1000

const CondTreeFromReq = (node) => {
    if (node.type == "or" || node.type == "and") {
        if (node.sub_conditions.length == 1) {
            return CondTreeFromReq(node.sub_conditions[0])
        }
        let ret = {cid: cidCounter++, type: node.type, subs: []}
        node.sub_conditions.forEach((value) => {
            ret.subs.push(CondTreeFromReq(value))
        })
        return ret
    }
    let ret = {
        cid: cidCounter++,
        type: node.type,
        timeRange: [tsFromStr(node.start_time), tsFromStr(node.end_time)],
        mode: node.mode,
        count: node.count,
        giftID: node.gift_id > 0 ? String(node.gift_id) : "",
        member1: false,
        member2: false,
        member3: false,
    }
    if (node.guard_levels) {
        node.guard_levels.forEach((value) => {
            if (value == 1) ret.member1 = true
            if (value == 2) ret.member2 = true
            if (value == 3) ret.member3 = true
        })
    }
    return ret
}

//-----------------------------------------------------------

let findFather = (id, node) => {
    if (node.cid == id) return null
    if (node.type == "and" || node.type == "or") {
        for (let i = 0; i < node.subs.length; i++) {
            if (node.subs[i].cid == id) return node
            let ret = findFather(id, node.subs[i])
            if (ret) return ret
        }
    }
    return null
}

let removeChild = (node, cid) => {
    if (!node) return
    let tmp = []
    node.subs.forEach((value) => {
        if (value.cid != cid) {tmp.push(value)}
    })
    node.subs = tmp
}

let CreateCondTreeHandler = (condTree, nextCID) => {
    return {
        onGroupChangeType(node, type) {
            node.type = type
        },
        deleteNode(id) {
            if (id <= 0) return
            removeChild(findFather(id, condTree), id)
        },
        addSubGroup(node) {
            node.subs.push({cid: nextCID.value++, type: "or", subs:[]})
        },
        addCond(node, type) {
            node.subs.push({
                cid: nextCID.value++,
                type: type,
                timeRange: [Date.now()-7*24*3600*1000, Date.now()],
                mode: "total",
                count: type == "sc" ? 0 : 1,
                giftID: "",
                member1: true,
                member2: true,
                member3: true,
            })
        },
        updateTimeRange(node, value) {
            node.timeRange = value
        },
        updateMode(node, value) {
            node.mode = value
        },
        updateCount(node, value) {
            node.count = value
        },
        updateMember1(node, value) {
            node.member1 = value
        },
        updateMember2(node, value) {
            node.member2 = value
        },
        updateMember3(node, value) {
            node.member3 = value
        },
        updateGift(node, value) {
            node.giftID = value
        },
    }
}

const ReadOnlyHandler = {
    onGroupChangeType(node, type) {},
    deleteNode(id) {},
    addSubGroup(node) {},
    addCond(node, type) {},
    updateTimeRange(node, value) {},
    updateMode(node, value) {},
    updateCount(node, value) {},
    updateMember1(node, value) {},
    updateMember2(node, value) {},
    updateMember3(node, value) {},
    updateGift(node, value) {},
}

const CondHasType = (node, type) => {
    if (node.type == type) return true
    if (node.type == "and" || node.type == "or") {
        for (let i = 0; i < node.subs.length; i++) {
            if (CondHasType(node.subs[i], type)) return true
        }
    }
    return false
}

export {CondTreeToReq, CondTreeFromReq, CreateCondTreeHandler, ReadOnlyHandler, CondHasType}