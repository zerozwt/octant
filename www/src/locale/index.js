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
        Data: {
            Cols: ["UID", "用户名", "赠送时间"],
            Member: ["总督", "提督", "舰长"],
            TimeRange: "时间范围（北京时间）",
            SCContent: "醒目留言内容",
            MemberLevel: "大航海等级",
            Search: "搜索",
        },
        Event: {
            New: "创建活动",
            ListCols: ["活动名称", "活动特典内容", "隐藏活动", "状态", "操作"],
            ListOps: ["观众名单", "编辑信息", "删除活动"],
            ListDelConfirm: "确定要删除这个活动？",
            EvtStatus: ["数据收集中", "名单计算中", "名单计算错误", "名单计算完成"],
            Add: {
                Basic: "基本信息",
                Cond: "参与条件",
                Delete: "删除",
                Hidden: "对观众隐藏本活动",
                Group: {
                    Title: "条件组",
                    Content: ["满足以下", "条件"],
                    Opts: ["任意", "全部"],
                    Add: "添加条件",
                    AddGroup: "添加条件组",
                },
                SC: {
                    Content: ["金额达到", ""],
                    Opts: ["累计", "单次"],
                },
                Member: {
                    Count: "数量",
                    Content: ["达到", "个月"],
                    Opts: ["总共", "单次"],
                },
                Gift: {
                    Count: "数量",
                    Content: ["达到", "个"],
                    Opts: ["总共", "单次"],
                },
            },
            Detail: {
                Block: ["解除屏蔽", "屏蔽"],
                Download: "下载名单",
            },
        },
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
        Data: {
            Cols: ["UID", "User name", "Time"],
            Member: ["Soutoku(总督)", "Teitoku(提督)", "Kanchou(舰长)"],
            TimeRange: "Time range (China Standard Time)",
            SCContent: "Super chat content",
            MemberLevel: "Membership level",
            Search: "Search",
        },
        Event: {
            New: "Create activity",
            ListCols: ["Activity name", "Special offer content", "Hidden activity", "Status", "Operation"],
            ListOps: ["Audience list", "Edit activity", "Delete"],
            ListDelConfirm: "Are you sure to delete this activity?",
            EvtStatus: ["Gathering data", "Calculating list", "Calculation error", "Ready"],
            Add: {
                Basic: "Basic information",
                Cond: "Conditions",
                Delete: "Delete",
                Hidden: "This activity will be hidden from audiences",
                Group: {
                    Title: "Condition group",
                    Content: ["Fulfil ", " of the following conditions"],
                    Opts: ["anyone", "all"],
                    Add: "Add condition",
                    AddGroup: "Add condition group",
                },
                SC: {
                    Content: ["reach", ""],
                    Opts: ["Sum of all super chats", "Single super chat"],
                },
                Member: {
                    Count: "Quantity",
                    Content: ["", "months"],
                    Opts: ["Total", "Single time"],
                },
                Gift: {
                    Count: "Quantity",
                    Content: ["", ""],
                    Opts: ["Total", "Single time"],
                },
            },
            Detail: {
                Block: ["Unblock", "Block"],
                Download: "Download list",
            },
        },
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
        Data: {
            Cols: ["UID", "ユーザー名", "贈り時間"],
            Member: ["総督", "提督", "艦長"],
            TimeRange: "時間帯（中国時間）",
            SCContent: "スーパーチャット内容",
            MemberLevel: "メンバーシップレベル",
            Search: "検索",
        },
        Event: {
            New: "イベントを開催する",
            ListCols: ["イベント名", "特典内容", "隠しイベント", "状態", "操作"],
            ListOps: ["リスナー名簿", "情報変更", "デリート"],
            ListDelConfirm: "このイベントを消去しますか？",
            EvtStatus: ["データ収集中", "リスナー名簿計算中", "計算エラー", "名簿計算完了"],
            Add: {
                Basic: "基本情報",
                Cond: "参加条件",
                Delete: "デリート",
                Hidden: "このイベントをリスナーに隠す",
                Group: {
                    Title: "条件グループ",
                    Content: ["下記の条件の", "を満たす"],
                    Opts: ["どれ一つ", "全て"],
                    Add: "条件を追加する",
                    AddGroup: "条件グループを追加する",
                },
                SC: {
                    Content: ["が", "以上になります"],
                    Opts: ["全部のスーパーチャットの総合金額", "一個のスーパーチャットの金額"],
                },
                Member: {
                    Count: "数",
                    Content: ["で", "ヶ月分"],
                    Opts: ["全部", "一回"],
                },
                Gift: {
                    Count: "数",
                    Content: ["で", "個分"],
                    Opts: ["全部", "一回"],
                },
            },
            Detail: {
                Block: ["ブロック解除", "ブロック"],
                Download: "リスナー名簿をダウンロード",
            },
        },
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