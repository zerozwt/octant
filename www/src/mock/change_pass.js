const Succ = {
    code: 0,
    msg: "",
    data: {},
}

export default {
    'post|^/api/admin/password$': opt => Succ,
}