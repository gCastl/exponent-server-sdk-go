package main

import (
	"fmt"
	"os"

	"github.com/Terminux/exponent-server-sdk-go"
)

// command to run example: EXPO_TOKEN=your_expo_token go run main.go
func main() {
	token := os.Getenv("EXPO_TOKEN")

	if expo.IsExpoPushToken(token) {
		message := expo.PushMessage{
			To:    token,
			Title: "Notification title",
			Body:  "Notification content",
			Data:  struct{ Value string }{"mydata"}}

		// is equivalent to expo.SendPushNotification(message)
		apiRes, apiErr, err := message.Send()
		if err != nil {
			panic(err)
		}

		fmt.Println("apiRes:", apiRes)
		fmt.Println("apiErr:", apiErr)
	}
}
