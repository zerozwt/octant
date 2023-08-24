let genSearchResult = () => {
    let ret = []

    for (let i = 0; i < 10; i++) {
        ret.push({
            uid: i,
            name: `DD_${i+1}`,
            time: "2023-08-19 20:00:00",
            gift: {
                name: "小花花",
                count: 114
            },
            sc: {
                price: 30,
                content: `哈哈 ${i+1}`,
            },
            guard: {
                level: (i%3)+1,
                count: 1,
            }
        })
    }

    return ret
}

const SearchResult = {
    code: 0,
    msg: "",
    data: {
        count: 100,
        list: genSearchResult(),
    }
}

const GiftList = {
    code: 0,
    msg: "",
    data: {
        count: 0,
        list: [
            {id: 1, name: "礼物1"},
            {id: 2, name: "礼物2"},
            {id: 3, name: "礼物3"},
            {id: 4, name: "礼物4"},
            {id: 5, name: "礼物5"},
        ],
    },
}

export default {
    'post|^/api/simple_search$': opt => SearchResult,
    'post|^/api/gifts$': opt => GiftList
}