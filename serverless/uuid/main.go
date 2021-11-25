package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

//http请求
func httpHandle(method, urlVal string) string {
	client := &http.Client{}
	var req *http.Request
 
	//添加cookie，key为X-Xsrftoken，value为df41ba54db5011e89861002324e63af81
        //可以添加多个cookie
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
 
	// req.Header.Add("X-Xsrftoken","b6d695bbdcd111e8b681002324e63af81")
 
	resp, err := client.Do(req)
 
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	return string(b)
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	ret := httpHandle("GET", "https://wfw.scu.edu.cn/ncov/wap/default/index")
	return &events.APIGatewayProxyResponse{
		StatusCode:        200,
		Headers:           map[string]string{"Content-Type": "text/plain"},
		Body:              ret,
		IsBase64Encoded:   false,
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call
	// fmt.Println(httpHandle("GET", "https://wfw.scu.edu.cn/ncov/wap/default/index"))
	lambda.Start(handler)
}
