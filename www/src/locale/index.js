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
            "一键导出观众名单",
        ],
        LoginButtonText: "登陆",
        LoginTypes: ["主播", "观众", "管理员"],
        LoginPane: {
            Name: "账户名：",
            Pass: "密码：",
            Admin: {
                AlreadyLogin: "当前管理员已登录",
                Jump: "进入管理页面",
            },
            Streamer: {
                Already: (name) => `${name}，欢迎再次光临Octant`,
                Jump: "进入系统",
            },
        },
    },
    Admin: {
        Name: "管理员",
        Streamer: {
            Title: "主播管理",
            Add: "添加主播",
            Cols: ["直播间ID", "主播名称", "账号名称", "操作"],
            Reset: "重置密码",
            Delete: "删除主播",
            DelConfirm: "确定要删除这名主播？",
            AddRoom: {
                ID: "直播间号码：",
                Name: "账户名：",
                Pass: "初始密码：",
                Pass2: "再次输入初始密码：",
                Confirm: "确定",
            },
            ResetPass: {
                Warn: "重置密码后，这名主播的观众需要重新填写收货地址信息。",
                Name: "主播名称：",
            },
        },
    },
    Streamer: {
        Menu: [
            "数据查询",
            "大航海",
            "醒目留言",
            "礼物",
            "活动管理",
            "私信群发",
        ],
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
            "一键导出观众名单",
        ],
        LoginButtonText: "SIGN IN",
        LoginTypes: ["Streamer", "Audience", "Administrator"],
        LoginPane: {
            Name: "Account:",
            Pass: "Password:",
            Admin: {
                AlreadyLogin: "Administrator has already signed in",
                Jump: "Go to admin panel",
            },
            Streamer: {
                Already: (name) => `Dear ${name}, welcome back to Octant`,
                Jump: "Enter dashboard",
            },
        },
    },
    Admin: {
        Name: "Administrator",
        Streamer: {
            Title: "Streamers",
            Add: "Add live streamer",
            Cols: ["Live room ID", "Streamer name", "Account", "Operation"],
            Reset: "Reset Password",
            Delete: "Delete",
            DelConfirm: "Are you sure to delete this streamer?",
            InputRoomID: "Live room ID:",
            AddRoom: {
                ID: "Live room ID:",
                Name: "Account:",
                Pass: "Initial password:",
                Pass2: "Repeat initial password:",
                Confirm: "Confirm",
            },
            ResetPass: {
                Warn: "After resetting password, audiences' delivery addresses have to be gathered again.",
                Name: "Live streamer:",
            },
        },
    },
    Streamer: {
        Menu: [
            "Data query",
            "Membership",
            "Super chat",
            "Gift",
            "Activities",
            "Direct messages",
        ],
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
            "一键导出观众名单",
        ],
        LoginButtonText: "ログイン",
        LoginTypes: ["配信者", "リスナー", "管理人"],
        LoginPane: {
            Name: "ユーザー名：",
            Pass: "パスワード：",
            Admin: {
                AlreadyLogin: "管理人は既にログインした",
                Jump: "管理画面へ",
            },
            Streamer: {
                Already: (name) => `お帰りなさいませ、${name}様`,
                Jump: "システムへ",
            },
        },
    },
    Admin: {
        Name: "管理人",
        Streamer: {
            Title: "配信者管理",
            Add: "配信者を追加する",
            Cols: ["配信部屋ID", "配信者", "アカウント", "操作"],
            Reset: "パスワードをリセット",
            Delete: "デリート",
            DelConfirm: "この配信者を消去しますか？",
            AddRoom: {
                ID: "配信部屋番号：",
                Name: "ユーザー名：",
                Pass: "初期パスワード：",
                Pass2: "初期パスワード再度入力：",
                Confirm: "OK",
            },
            ResetPass: {
                Warn: "パスワードがリセットされた後、この配信者のリスナーさん達の届け先情報を再収集する必要があります。",
                Name: "配信者：",
            },
        },
    },
    Streamer: {
        Menu: [
            "データ検索",
            "メンバーシップ",
            "スーパーチャット",
            "プレゼント",
            "イベント管理",
            "DM管理",
        ],
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