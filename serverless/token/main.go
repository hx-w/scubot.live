package main

import (
	"os"

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
	if input_token, ok := request.QueryStringParameters["simple_verify"]; ok {
		simple_token := os.Getenv("SIMPLE_TOKEN")
		if input_token == simple_token {
			return Resp(200, "{\"message\": \"验证成功\"}")
		} else {
			return Resp(403, "{\"detail\": \"验证失败\"}")
		}
	}
	if request.HTTPMethod != "POST" {
		return Resp(403, "{\"detail\": \"未规定的请求方法\"}")
	} else {
		token := os.Getenv("ACCESS_TOKEN")
		return Resp(200, "{\"access_token\": \"" + token + "\"}")
	}
}

func main() {
	lambda.Start(handler)
}