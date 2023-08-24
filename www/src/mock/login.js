const AdminLogin = {
    code: 0,
    msg: "",
    data: {},
};

const AdminDoLogin = {
    code: 0,
    msg: "",
    data: {},
};

const StreamerLogin = {
    code: 0,
    msg: "",
    data: {
        room_id: 123,
        name: "XXX_Channel",
        account_name: "xxx",
    }
}

export default {
    'get|^/api/admin/login$': opt => AdminLogin,
    'post|^/api/admin/login$': opt => AdminDoLogin,
    'get|^/api/admin/logout$': opt => AdminDoLogin,
    'get|^/api/streamer/login$': opt => StreamerLogin,
    'post|^/api/streamer/login$': opt => AdminDoLogin,
    'get|^/api/streamer/logout$': opt => AdminDoLogin,
}