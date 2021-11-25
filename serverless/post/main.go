package main

import (
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

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if initClient() != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode:        403,
			Headers:           map[string]string{"Content-Type": "application/json"},
			Body:              "{\"detail\": \"redis connect refused!\"}",
		}, nil
	}
	err := rdb.Set("score", request.Body, 0).Err()
    if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode:        403,
			Headers:           map[string]string{"Content-Type": "application/json"},
			Body:              "{\"detail\": \"redis refused!\"}",
		}, nil
    }

	return &events.APIGatewayProxyResponse{
		StatusCode:        200,
		Headers:           map[string]string{"Content-Type": "application/json"},
		Body:              "set!",
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call
	lambda.Start(handler)
}