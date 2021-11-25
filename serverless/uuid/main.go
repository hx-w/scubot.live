package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

//http请求
func httpHandler(method, urlVal string) (bool, string) {
	client := &http.Client{}
	var req *http.Request
 
	req, _ = http.NewRequest(method, urlVal, nil)
	cookies_uukey := &http.Cookie{
		Name: "UUkey",
		Value: "510ff917237279fe13fbf2558309dfbd",
		HttpOnly: true,
	}
	cookies_eai := &http.Cookie{
		Name: "eai-sess",
		Value: "jq3ssf35cnhlfbodaa4ppcrsd2",
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
	fmt.Println(result)
	if len(result) == 1 {
		return true, result[0][1]
	} else {
		return false, ""
	}
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// request.QueryStringParameters
	success, retBody := httpHandler("GET", "https://wfw.scu.edu.cn/ncov/wap/default/index")
	if success {
		parse_success, uid := regHandler(retBody)
		if parse_success {
			return &events.APIGatewayProxyResponse{
				StatusCode:        200,
				Headers:           map[string]string{"Content-Type": "application/json"},
				Body:              "{\"uid\": " + uid + "}",
			}, nil
		}
	}
	return &events.APIGatewayProxyResponse{
		StatusCode:        403,
		Headers:           map[string]string{"Content-Type": "application/json"},
		Body:              "error",
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call
	lambda.Start(handler)
}
