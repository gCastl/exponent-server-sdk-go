# exponent-server-sdk-go
[![Build Status](https://travis-ci.org/Terminux/exponent-server-sdk-go.svg?branch=master)](https://travis-ci.org/Terminux/exponent-server-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Terminux/exponent-server-sdk-go)](https://goreportcard.com/report/github.com/Terminux/exponent-server-sdk-go)
[![GoDoc](https://godoc.org/github.com/Terminux/exponent-server-sdk-go?status.svg)](https://godoc.org/github.com/Terminux/exponent-server-sdk-go)

Server side library for working with Exponent using Go

## Installing
To install this library, simply run:
```bash
go get github.com/Terminux/exponent-server-sdk-go
```
## Usage
### Token Checker
Check if the token is a valid Expo push token
```go
  expo.IsExpoPushToken(token)
```

### Send Single Message
```go
  message := expo.PushMessage{To: token, Body: "content"}

  message.Send()
  // or
  expo.SendPushNotification(&message)
```

### Send Several Messages
```go
  expo.SendPushNotifications([]*expo.PushMessage{
    &expo.PushMessage{To: token1, Body: "first message"},
    &expo.PushMessage{To: token2, Body: "another message"},
    &expo.PushMessage{To: token3, Body: "last message"},
  })
```

### Chunks Messages
Split the chunk messages into several chunks messages
```go
  expo.ChunkPushNotifications(messages)
```

### Example
Here's a sample showcasing many features of expo.
```go
import "github.com/Terminux/exponent-server-sdk-go"

const token = "EXPO_TOKEN"

func main() {
	if expo.IsExpoPushToken(token) {
		message := expo.PushMessage{
			To:    token,
			Title: "Notification title",
			Body:  "Notification content",
			Data:  struct{ Value string }{"mydata"}}

		api, err := message.Send()
		if err != nil {
			panic(err)
		}

		fmt.Println("api result:", api)
	}
}
```

### Based on

  * https://github.com/expo/exponent-server-sdk-node
  * https://docs.expo.io/versions/v18.0.0/guides/push-notifications.html