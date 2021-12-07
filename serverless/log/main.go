package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-redis/redis"
)

var rdb *redis.Client
 
func initClient() (err error) {
    rdb = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
    })
    _, err = rdb.Ping().Result()
    if err != nil {
        return err
    }
    return nil
}

func Resp(StatusCode int, Body string) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: int(StatusCode),
		Headers:	map[string]string{"Content-Type": "application/json"},
		Body: 		Body,
		IsBase64Encoded: false,
	}, nil
}

func SCUBotLogKey(uid string) string {
	const PREFIX = "checkinLog-"
	return PREFIX + uid + "-*"
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if initClient() != nil {
		return Resp(403, "{\"detail\": \"数据库连接错误\"}")
	}
	if uid, ok := request.QueryStringParameters["uid"]; ok {
		result, err := rdb.Keys(SCUBotLogKey(uid)).Result()
		if err != nil {
			return Resp(403, "{\"detail\": \"数据库操作错误\"}")
		}
		if len(result) == 0 {
			return Resp(200, "{\"message\": \"[]\"}")
		}
		valRes, err := rdb.MGet(result...).Result()
		if err != nil {
			return Resp(403, "{\"detail\": \"数据库操作错误\"}")
		}
		resJson, err := json.Marshal(valRes)
		if err != nil {
			return Resp(403, "{\"detail\": \"数据库操作错误\"}")
		}
		return Resp(200, "{\"message\": \"" + string(resJson) + "\"}")
	}
	return Resp(403, "{\"detail\": \"参数错误\"}")
}

func main() {
	lambda.Start(handler)
}