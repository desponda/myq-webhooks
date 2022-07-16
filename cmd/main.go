package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/desponda/myq-golang"
)

type MyEvent struct {
        Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
        username := os.Getenv("MYQ_USERNAME")
        password := os.Getenv("MYQ_PASSWORD")
        session := myq.Session{Username: username, Password: password}
        if err := session.Login(); err != nil {
                fmt.Printf("Error logging in to myq: %s\n", err)
                return "", err
        }
        retries := 10
        retry := 0
        state, err := session.DeviceState("SERIAL_NUMBER")
        for  !strings.Contains(state,"closed") && retry < retries {
                retry++
                session.SetDoorState("SERIAL_NUMBER", "close")
                time.Sleep(30 * time.Second)
                state, err = session.DeviceState("SERIAL_NUMBER")
                fmt.Printf("State: %s, Error: %s\n", state, err)
        }

        if err!= nil {
                fmt.Printf("Error: %s\n", err)
        }

        if !strings.Contains( state, "closed") { 
                return "", fmt.Errorf("device not closed")
        }

  

        return "completed", nil
}
func main() {
        lambda.Start(HandleRequest)
}