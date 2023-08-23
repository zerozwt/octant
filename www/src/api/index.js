import axios from 'axios'

const loginChecker = (router, rsp, resolve) => {
    let data = rsp.data;
    if (data.code == 114514 || data.code == 1919 || data.code == 810) {
        router.push("/")
        return
    }
    resolve(rsp)
}

const APICaller = (router) => {
    return {
        get: (path, config) => new Promise((resolve, reject) => {
            axios.get(path, config).then(rsp => {
                loginChecker(router, rsp, resolve)
            }).catch(err => reject(err))
        }),
        post: (path, data) => new Promise((resolve, reject) => {
            axios.post(path, data).then(rsp => {
                loginChecker(router, rsp, resolve)
            }).catch(err => reject(err))
        }),
    }
}

export { APICaller }