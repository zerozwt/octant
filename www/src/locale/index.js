const zhCN = {
    General: {
        ChangePass: "修改密码",
        Logout: "退出登录",
        OldPass: "当前密码",
        NewPass: "新密码",
        NewPass2: "重复新密码",
    },
    Index: {
        Title: "欢迎使用Octant",
        SubTitle: "一个兴趣使然的B站直播运营活动管理系统",
        Intro: [
            "灵活查询礼物、SC和大航海数据",
            "自由设定规则，自动计算观众名单",
            "对名单内观众群发私信",
            "收集观众收货信息",
            "一键导出下载名单",
        ],
        LoginButtonText: "登陆",
        LoginTypes: ["主播", "观众", "管理员"],
        LoginPane: {
            Pass: "密码：",
            Admin: {
                AlreadyLogin: "当前管理员已登录",
                Jump: "进入管理页面",
            },
        },
    },
    Admin: {
        Name: "管理员",
    },
};

const enUS = {
    General: {
        ChangePass: "Change password",
        Logout: "Logout",
        OldPass: "Current password",
        NewPass: "New password",
        NewPass2: "Repeat new password",
    },
    Index: {
        Title: "Welcome to Octant",
        SubTitle: "An operation management system for live streamers in Bilibili",
        Intro: [
            "灵活查询礼物、SC和大航海数据",
            "自由设定规则，自动计算观众名单",
            "对名单内观众群发私信",
            "收集观众收货信息",
            "一键导出下载名单",
        ],
        LoginButtonText: "SIGN IN",
        LoginTypes: ["Streamer", "Audience", "Administrator"],
        LoginPane: {
            Pass: "Password:",
            Admin: {
                AlreadyLogin: "Administrator has already signed in",
                Jump: "Go to admin panel",
            },
        },
    },
    Admin: {
        Name: "Administrator",
    },
};

const jaJP = {
    General: {
        ChangePass: "パスワード変更",
        Logout: "ログアウト",
        OldPass: "今のパスワード",
        NewPass: "新しいパスワード",
        NewPass2: "新しいパスワード再度入力",
    },
    Index: {
        Title: "Octantへようこそ",
        SubTitle: "ビリビリでの配信者の為の運営活動管理システム",
        Intro: [
            "灵活查询礼物、SC和大航海数据",
            "自由设定规则，自动计算观众名单",
            "对名单内观众群发私信",
            "收集观众收货信息",
            "一键导出下载名单",
        ],
        LoginButtonText: "ログイン",
        LoginTypes: ["配信者", "リスナー", "管理人"],
        LoginPane: {
            Pass: "パスワード：",
            Admin: {
                AlreadyLogin: "管理人は既に",
                Jump: "管理画面へ",
            },
        },
    },
    Admin: {
        Name: "管理人",
    },
};

const localeObj = (locale) => {
    if (locale == "en-US") return enUS;
    if (locale == "ja-JP") return jaJP;
    return zhCN;
}

const loadCache = (key, def_value) => {
    let value = window.localStorage.getItem(key)
    return value ? value : def_value
}

const defaultLang = () => {
    const supportedLangs = ["zh-CN", "en-US", "ja-JP"]
    if (supportedLangs.includes(window.navigator.language)) {
        return window.navigator.language
    }
    return "zh-CN"
}

const defaultLocaleObj = () => {
    return localeObj(loadCache("octant_lang", defaultLang()))
}

export {localeObj, loadCache, defaultLang, defaultLocaleObj}