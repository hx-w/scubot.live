package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Resp(StatusCode int, Body string) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: int(StatusCode),
		Headers:	map[string]string{"Content-Type": "application/json"},
		Body: 		Body,
		IsBase64Encoded: false,
	}, nil
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return Resp(200, "{\"message\": \"" + request.HTTPMethod + "\"}")
	// var reqJson map[string]interface{}
    // if err := json.Unmarshal([]byte(request.Body), &reqJson); err != nil {
	// 	return Resp(403, "{\"detail\": \"请求体错误\"")
    // }
}

func main() {
	lambda.Start(handler)
}