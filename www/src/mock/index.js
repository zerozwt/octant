var Mock = await import('mockjs');

Mock.setup({timeout: "200-400"});

let confArr = [];

const files = import.meta.glob('./*.js');
for (const key in files) {
    if (key == "./index.js") continue;
    confArr = confArr.concat((await import(key)).default);
}

confArr.forEach((item) => {
    for (let [mathod_path, target] of Object.entries(item)) {
        let tmp = mathod_path.split('|');
        let method = tmp[0];
        let path = new RegExp(tmp[1]);
        console.log("MOCK", method, path);
        Mock.mock(path, method, target);
    }
});