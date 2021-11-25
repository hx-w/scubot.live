const redis = require('redis')
exports.handler = async function (event, context) {
    // parse body


    const redisClient = redis.createClient(process.env.REDIS_PORT, process.env.REDIS_HOST)
    redisClient.on('error', err => {
        console.error(err)
        return {
            statusCode: 403,
            body: '数据库连接错误'
        }
    })
    redisClient.auth(process.env.REDIS_PASSWORD)
    redisClient.get()
    redisClient.end(true);

    return {
        statusCode: 200,
        body: JSON.stringify({ message: "Hello World" + event.body })
    };
}