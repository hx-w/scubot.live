const redis = require('redis')
exports.handler = async function (event, context) {
    const redisClient = redis.createClient(27014, 'mc.csgowiki.top')
    redisClient.on('error', err => {
        console.error(err)
    })
    redisClient.auth(process.env.REDIS_PASSWORD)
    redisClient.set("test", event.body, function (error, res) {
        if (error) {
            console.log(error);
        } else {
            console.log(res);
        };
    });
    redisClient.end(true);

    return {
        statusCode: 200,
        body: JSON.stringify({ message: "Hello World" + event.body })
    };
}