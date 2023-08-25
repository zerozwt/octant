const EventList = {
    code: 0,
    msg: "",
    data: {
        count: 100,
        list: [
            {id: 1, name: "活动1", content:"特典内容", status: 1},
            {id: 2, name: "活动2", content:"特典内容", status: 2},
            {id: 3, name: "活动3", content:"特典内容", status: 3},
            {id: 4, name: "活动4", content:"特典内容", status: 4},
            {id: 5, name: "活动5", content:"特典内容", status: 1},
            {id: 6, name: "活动6", content:"特典内容", status: 2},
            {id: 7, name: "活动7", content:"特典内容", status: 3},
            {id: 8, name: "活动8", content:"特典内容", status: 4},
            {id: 9, name: "活动9", content:"特典内容", status: 1},
            {id: 10, name: "活动10", content:"特典内容", status: 2},
        ],
    }
}

const OK = {
    code: 0,
    msg: "",
    data:{},
}

export default {
    'get|^/api/event/list': opt => EventList,
    'post|^/api/event/delete$': opt => OK,
    'post|^/api/event/modify$': opt => OK,
}