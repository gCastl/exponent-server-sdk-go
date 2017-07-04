/*
Package expo is used to send push notifications to Expo Experiences from a Go server

A simplest example:
	package main
	import "github.com/Terminux/exponent-server-sdk-go"

	const token = "EXPO_TOKEN"

	func main() {
		if expo.IsExpoPushToken(token) {
			message := expo.PushMessage{
				To:    token,
				Title: "Notification title",
				Body:  "Notification content",
				Sound: "default",
				Badge: 1
				Data:  struct{ Value string }{"mydata"}}}

			message.Send()
		}
	}
*/
package expo
