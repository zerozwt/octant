const EventList = {
    code: 0,
    msg: "",
    data: {
        count: 100,
        list: [
            {id: 1, name: "活动1", content:"特典内容", status: 1, hidden: true},
            {id: 2, name: "活动2", content:"特典内容", status: 2, hidden: false},
            {id: 3, name: "活动3", content:"特典内容", status: 3, hidden: true},
            {id: 4, name: "活动4", content:"特典内容", status: 4, hidden: false},
            {id: 5, name: "活动5", content:"特典内容", status: 1, hidden: true},
            {id: 6, name: "活动6", content:"特典内容", status: 2, hidden: false},
            {id: 7, name: "活动7", content:"特典内容", status: 3, hidden: true},
            {id: 8, name: "活动8", content:"特典内容", status: 4, hidden: false},
            {id: 9, name: "活动9", content:"特典内容", status: 1, hidden: true},
            {id: 10, name: "活动10", content:"特典内容", status: 2, hidden: false},
        ],
    }
}

const OK = {
    code: 0,
    msg: "",
    data:{},
}

const Detail = {
    code: 0,
    msg: "",
    data: {
        id: 1,
        name: "爆大米",
        reward: "1.特典1\n2.特典2\n3.特典3\n4.特典4\n\n还有呵呵",
        status: 4,
        hidden: 0,
        conditions: {
            type: "and",
            sub_conditions: [
                {
                    type: "member",
                    start_time: "20230819000000",
                    end_time: "20230820050000",
                    mode: "total",
                    count: 1,
                    guard_levels: [1,2],
                    gift_id: 0,
                },
                {
                    type: "sc",
                    start_time: "20230819000000",
                    end_time: "20230820050000",
                    mode: "once",
                    count: 1000,
                    guard_levels: [],
                    gift_id: 0,
                },
                {
                    type: "gift",
                    start_time: "20230819000000",
                    end_time: "20230820050000",
                    mode: "once",
                    count: 1000,
                    guard_levels: [],
                    gift_id: 5,
                },
            ],
        }
    }
}

const genUserList = () => {
    let ret = []
    for (let i = 0; i < 10; i++) {
        let item = {
            uid: i+1,
            name: `User_${i+1}`,
            block: i % 3 == 0,
            cols: {
                member: [],
                gift: [],
                sc: [],
            }
        }
        if (i % 2 == 0) {
            item.cols.member.push({
                time: "2023-08-19 20:00:00 CST",
                level: ((i/2)%3)+1,
                count: (i/2)+1,
            })
        }
        if (i % 3 == 0) {
            item.cols.sc.push({
                time: "2023-08-19 20:00:00 CST",
                price: (((i/2)%3)+1)*100,
                content: `Super chat #${i}`,
            })
        }
        if (i % 2 == 1) {
            item.cols.gift.push({
                time: "2023-08-19 20:00:00 CST",
                gift_id: (i%4)+1,
                gift_name: `Gift_${i}`,
                price: ((i%4)+1)*500,
                count: i*10+10,
            })
        }
        ret.push(item)
    }
    return ret
}

const UserList = {
    code: 0,
    msg: "",
    data: {
        count: 100,
        list: genUserList(),
    }
}

export default {
    'get|^/api/event/list': opt => EventList,
    'post|^/api/event/delete$': opt => OK,
    'post|^/api/event/modify$': opt => OK,
    'post|^/api/event/add$': opt => OK,
    'get|^/api/event/detail': opt =>Detail,
    'post|^/api/event/user/list$': opt => UserList,
    'post|^/api/event/user/block$': opt => OK,
    'post|^/api/event/user/unblock$': opt => OK,
}