const axios = require('axios')

async function get_uuid(cookies_dict, UA) {
    const pattern = /var def =(.*);!?/g
    const url = 'https://wfw.scu.edu.cn/ncov/wap/default/index'
    return new Promise((resolve, reject) => {
        axios
            .get(url, {
                headers: {
                    'Host': 'wfw.scu.edu.cn',
                    'User-Agent': UA,
                    'Cookie': `eai-sess=${cookies_dict.eai_sess}; UUkey=${cookies_dict.UUkey}`
                }
            })
            .then((resp) => {
                var got = pattern.exec(resp.data)
                if (got == null || got == undefined || got.length <= 1) {
                    reject('parse error')
                }
                defJson = JSON.parse(got[1])
                resolve(defJson['uid'])
            })
            .catch((error) => {
                reject(error)
            })
    })
}

module.exports = {
    get_uuid
}