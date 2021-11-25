package main

import (
	"encoding/json"
	"os"
	"strconv"

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

func SCUBotKey(uid string) string {
	const PREFIX = "scubot-"
	return PREFIX + uid
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if initClient() != nil {
		return Resp(403, "{\"detail\": \"数据库连接错误\"}")
	}
	var reqJson map[string]interface{}
    if err := json.Unmarshal([]byte(request.Body), &reqJson); err != nil {
		return Resp(403, "{\"detail\": \"请求体错误\"")
    }
	content, _ := json.Marshal(reqJson["content"])
	uid := strconv.Itoa(int(reqJson["uid"].(float64)))
	err := rdb.Set(SCUBotKey(uid), string(content), 0).Err()
    if err != nil {
		return Resp(403, "{\"detail\": \"数据库操作错误\"}")
    }

	return Resp(200, "{\"message\": \"操作成功\"}")
}

func main() {
	lambda.Start(handler)
}