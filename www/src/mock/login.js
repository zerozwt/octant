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

export default {
    'get|^/api/admin/login$': opt => AdminLogin,
    'post|^/api/admin/login$': opt => AdminDoLogin,
    'get|^/api/admin/logout$': opt => AdminDoLogin,
}