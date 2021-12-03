package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

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

func httpHandler(method, urlVal, UUkey, eai_sess string) (bool, string) {
	client := &http.Client{}
	var req *http.Request
 
	req, _ = http.NewRequest(method, urlVal, nil)
	cookies_uukey := &http.Cookie{
		Name: "UUkey",
		Value: UUkey,
		HttpOnly: true,
	}
	cookies_eai := &http.Cookie{
		Name: "eai-sess",
		Value: eai_sess,
		HttpOnly: true,
	}
	req.AddCookie(cookies_uukey)
	req.AddCookie(cookies_eai)
 
	resp, err := client.Do(req)
 
	if err != nil {
		return false, ""
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	return true, string(b)
}

func regHandler(respBody string) (bool, string) {
	pattern := regexp.MustCompile(`var def =.*"uid":"(\d*)".*;`)
	result := pattern.FindAllStringSubmatch(respBody, -1)
	if len(result) == 1 {
		return true, result[0][1]
	} else {
		return false, ""
	}
}

func SCUBotKey(uid string) string {
	const PREFIX = "scubot-"
	return PREFIX + uid
}

func CheckUidExist(uid string) bool {
	_, err := rdb.Get(SCUBotKey(uid)).Result()
	return err != redis.Nil
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if initClient() != nil {
		return Resp(403, "{\"detail\": \"数据库连接错误\"}")
	}
	UUkey, okU := request.QueryStringParameters["UUkey"]
	eai_sess, okE := request.QueryStringParameters["eai-sess"]
	if !okE || !okU {
		return Resp(403, "{\"detail\": \"参数错误\"}")
	}

	success, retBody := httpHandler("GET", "https://wfw.scu.edu.cn/ncov/wap/default/index", UUkey, eai_sess)
	if success {
		parse_success, uid := regHandler(retBody)
		exist := "false"
		if CheckUidExist(uid) {
			exist = "true"
		}
		if parse_success {
			return &events.APIGatewayProxyResponse{
				StatusCode:        200,
				Headers:           map[string]string{"Content-Type": "application/json"},
				Body:              "{\"uid\": " + uid + ", \"exist\": " + exist + "}",
			}, nil
		}
	}
	return Resp(403, "{\"detail\": \"请求微服务服务器错误\"}")
}

func main() {
	lambda.Start(handler)
}
