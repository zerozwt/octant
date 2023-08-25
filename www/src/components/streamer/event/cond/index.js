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
        count: node.count ? 0 : Number(node.count),
        gift_id: node.giftID ? 0 : Number(node.giftID),
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
        node.subs.forEach((value) => {
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
        gift_id: node.giftID > 0 ? String(node.giftID) : "",
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

export {CondTreeToReq, CondTreeFromReq}