package expo

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

const token = "ExponentPushToken[xxxxxxxxxxxxxxxxxxxxxx]"
const status200 = `{"data": [{"status": "ok"}]}`

func TestIsExpoPushToken(t *testing.T) {
	if IsExpoPushToken("badToken") {
		t.Errorf("IsExpoPushToken returned unexpected value: got true want false")
	}

	if !IsExpoPushToken(token) {
		t.Errorf("IsExpoPushToken returned unexpected value: got false want true")
	}
}

func TestChunkPushNotifications(t *testing.T) {
	messages := []*PushMessage{
		{To: "token"},
		{To: "token"},
		{To: "token"},
	}

	chunks := ChunkPushNotifications(messages)
	if len(chunks) > 1 {
		t.Errorf("ChunkPushNotifications returned unexpected chunks: chunks length got %v want 1", len(chunks))
	}

	ChunkLimit = 2
	chunks = ChunkPushNotifications(messages)
	if len(chunks) != 2 {
		t.Errorf("ChunkPushNotifications returned unexpected chunks: chunks length got %v want 2", len(chunks))
	}

	ChunkLimit = 1
	chunks = ChunkPushNotifications(messages)
	if len(chunks) != 3 {
		t.Errorf("ChunkPushNotifications returned unexpected chunks: chunks length got %v want 3", len(chunks))
	}
}

func TestSendPushNotification(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	message := PushMessage{
		To:    token,
		Title: "Notification title",
		Body:  "Notification content"}

	httpmock.RegisterResponder("POST", baseAPIURL+"/push/send",
		httpmock.NewStringResponder(200, status200))

	api, _ := message.Send()
	if api.Data[0].Status != "ok" {
		t.Errorf("SendPushNotification returned unexpected response: status got %s want ok", api.Data[0].Status)
	}
}

func TestSendPushNotifications(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	m1 := PushMessage{
		To:    token,
		Title: "Notification title",
		Body:  "Notification content"}

	m2 := PushMessage{
		To:    token,
		Title: "Notification title",
		Body:  "Notification content"}

	httpmock.RegisterResponder("POST", baseAPIURL+"/push/send",
		httpmock.NewStringResponder(200, status200))

	api, _ := SendPushNotifications([]*PushMessage{&m1, &m2})
	if api.Data[0].Status != "ok" {
		t.Errorf("SendPushNotifications returned unexpected response: status got %s want ok", api.Data[0].Status)
	}
}

func TestBodyGzip(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	message := PushMessage{
		To:    token,
		Title: "Notification title",
		Body:  "Notification content"}

	MaxBodySizeWithoutGzip = 1

	httpmock.RegisterResponder("POST", baseAPIURL+"/push/send",
		func(req *http.Request) (*http.Response, error) {
			if req.ContentLength != 110 {
				t.Errorf("SendPushNotification send unexpected message: ContentLength got %v want 110", req.ContentLength)
			}

			resp, _ := httpmock.NewJsonResponse(200, status200)
			return resp, nil
		})

	SendPushNotification(&message)
}
