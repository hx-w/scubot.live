package main

import (
	"fmt"
	"os"
	"time"

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

func WorkFlowNum() int {
	result, err := rdb.Keys("scubot-*").Result()
	if err != nil {
		panic(err)
	}
	return len(result)
}

func LogNum() int {
	result, err := rdb.Keys("checkinLog-*").Result()
	if err != nil {
		panic(err)
	}
	return len(result)
}

func DayDiff() int {
	startTime, _ := time.Parse("2006-01-02", os.Getenv("START_DAY"))
	nowTime := time.Now()
	return int(nowTime.Sub(startTime).Hours() / 24)
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if initClient() != nil {
		return Resp(403, "{\"detail\": \"数据库连接错误\"}")
	}
	if typeType, ok := request.QueryStringParameters["type_"]; ok {
		if typeType == "0" {
			return Resp(200, fmt.Sprintf("{\"schemaVersion\": 1, \"label\": \"有效工作流\", \"message\": \"%d项\", \"color\": \"ad453f\"}", WorkFlowNum()))
		} else if typeType == "1" {
			return Resp(200, fmt.Sprintf("{\"schemaVersion\": 1, \"label\": \"缓存打卡日志\", \"message\": \"%d条\", \"color\": \"2f90b9\"}", LogNum()))
		} else if typeType == "2" {
			return Resp(200, fmt.Sprintf("{\"schemaVersion\": 1, \"label\": \"运行时间\", \"message\": \"%d天\", \"color\": \"3c9566\"}", DayDiff()))
		}
	}
	return Resp(403, "{\"detail\": \"参数错误\"}")
}

func main() {
	lambda.Start(handler)
}