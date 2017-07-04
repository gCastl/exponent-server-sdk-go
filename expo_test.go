package expo

import (
	"testing"

	"github.com/jarcoal/httpmock"
)

const token = "ExponentPushToken[xxxxxxxxxxxxxxxxxxxxxx]"

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
		&PushMessage{To: "token"},
		&PushMessage{To: "token"},
		&PushMessage{To: "token"},
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
		httpmock.NewStringResponder(200, `{"data": [{"status": "ok"}]}`))

	api, _ := message.Send()
	if api.Data[0].Status != "ok" {
		t.Errorf("SendPushNotification returned unexpected response: status got %s want ok", api.Data[0].Status)
	}
}
