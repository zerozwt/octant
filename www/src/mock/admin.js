const StreamerList = {
    code: 0,
    msg: "",
    data: {
        count: 100,
        list: [
            {room_id: 1, name: "主播1", account_name: "streamer_1"},
            {room_id: 2, name: "主播2", account_name: "streamer_2"},
            {room_id: 3, name: "主播3", account_name: "streamer_3"},
            {room_id: 4, name: "主播4", account_name: "streamer_4"},
            {room_id: 5, name: "主播5", account_name: "streamer_5"},
            {room_id: 6, name: "主播6", account_name: "streamer_6"},
            {room_id: 7, name: "主播7", account_name: "streamer_7"},
            {room_id: 8, name: "主播8", account_name: "streamer_8"},
            {room_id: 9, name: "主播9", account_name: "streamer_9"},
            {room_id: 10, name: "主播10", account_name: "streamer_10"},
        ],
    },
}

const AddStreamer = {
    code: 0,
    msg: "",
    data: {}
}

export default {
    'get|^/api/admin/streamer/list': opt => StreamerList,
    'post|^/api/admin/streamer/add$': opt => AddStreamer,
    'post|^/api/admin/streamer/delete$': opt => AddStreamer,
    'post|^/api/admin/streamer/reset$': opt => AddStreamer,
}