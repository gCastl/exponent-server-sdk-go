package main

import (
	"fmt"

	expo "github.com/Terminux/exponent-server-sdk-go"
)

const token = "EXPO_TOKEN"

func main() {
	if expo.IsExpoPushToken(token) {
		message := expo.PushMessage{
			To:    token,
			Title: "Notification title",
			Body:  "Notification content",
			Data:  struct{ Value string }{"mydata"}}

		// is equivalent to expo.SendPushNotification(message)
		api, err := message.Send()
		if err != nil {
			panic(err)
		}

		fmt.Println("api result:", api)
	}
}
