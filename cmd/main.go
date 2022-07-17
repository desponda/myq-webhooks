package main

import (
	"fmt"
	"net/http"
	"os"
        "encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/desponda/myq-golang"
	"github.com/desponda/myq-webhooks/pkg/services"
)

type MyQDesiredStateRequest struct {
        SerialNumber string `json:"serial_number"`
        DesiredState string `json:"desired_state"`
}

func MyQHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
        username := os.Getenv("MYQ_USERNAME")
        password := os.Getenv("MYQ_PASSWORD")
        session := myq.Session{Username: username, Password: password}
        service := services.NewDeviceManager(services.DeviceManagerOptions{MaxRetries: 3, RetryInterval: 30}, &session)
        var desiredState MyQDesiredStateRequest

        err := json.Unmarshal([]byte(req.Body), &desiredState)
        if err != nil {
                return events.APIGatewayProxyResponse{StatusCode: 400}, err
        }

        err = service.SetDesiredState(services.DeviceDesiredState{SerialNumber: desiredState.SerialNumber, DesiredState: desiredState.DesiredState})

        if err!= nil {
                fmt.Printf("Error: %s\n", err)
                return events.APIGatewayProxyResponse{StatusCode: 500}, err
        }
  

        return events.APIGatewayProxyResponse{
                StatusCode: http.StatusOK,
            }, nil
}

func main() {
        lambda.Start(MyQHandler)
}