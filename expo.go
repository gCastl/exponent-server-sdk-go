package expo

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
)

var version = "2.3.1"
var baseURL = "https://exp.host"
var baseAPIURL = baseURL + "/--/api/v2"

// MaxBodySizeWithoutGzip allows to set the max length of body allowed to be send without gzip.
// The MaxBodySizeWithoutGzip can be increase or decrease but it is not recommanded to set higher than 1024
var MaxBodySizeWithoutGzip = 1024

// ChunkLimit allows to set the max message in each chunk. This variable is used on ChunkPushNotifications function.
// The ChunkLimit can be increase or decrease but it is not recommanded to set higher than 100
var ChunkLimit = 100

// PushNotificationResult is the result returned by the Expo api
type PushNotificationResult struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Details struct {
		Error string `json:"error"`
	} `json:"details"`
}

// PushNotificationError is the result error returned by the Expo api
type PushNotificationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
	Stack   string `json:"stack"`
}

type pushNotificationResponse struct {
	Errors []*PushNotificationError  `json:"errors"`
	Data   []*PushNotificationResult `json:"data"`
}

// PushMessage is the message sended to the Expo api
type PushMessage struct {
	// To is an Expo push token specifying the recipient of this message.
	To string `json:"to"`

	// Title is the title to display in the notification. On iOS this is displayed only on Apple Watch.
	Title string `json:"title,omitempty"`

	// Body is the push notification content
	Body string `json:"body,omitempty"`

	// Data is a JSON object delivered to your app. It may be up to about 4KiB; the total
	// notification payload sent to Apple and Google must be at most 4KiB or else you will get a "Message Too Big" error.
	Data interface{} `json:"data,omitempty"`

	// Sound to play when the recipient receives this notification.
	// Specify "default" to play the device's default notification sound, or omit this field to play no sound.
	Sound string `json:"sound,omitempty"`

	// TTL (Time to Live) is the number of seconds for which the message may be kept around for redelivery if it hasn't been delivered yet.
	TTL int `json:"ttl"`

	// Expiration is a timestamp since the UNIX epoch specifying when the message expires.
	Expiration int `json:"expiration"`

	// Priority is the delivery priority of the message.
	// Possible values: normal | hight | default or omit field to use the default priority
	Priority string `json:"priority,omitempty"`

	// Badge is the number to display in the badge on the app icon
	Badge int `json:"badge"`
}

// Send allows to send the current message
func (p *PushMessage) Send() (r *PushNotificationResult, e *PushNotificationError, err error) {
	res, resErr, err := SendPushNotifications([]*PushMessage{p})
	if isError(err) {
		return nil, nil, err
	}

	if len(res) == 1 {
		r = res[0]
	}

	if len(resErr) == 1 {
		e = resErr[0]
	}

	return r, e, nil
}

// IsExpoPushToken determines if the token is a Expo push token
func IsExpoPushToken(token string) bool {
	return strings.HasPrefix(token, "ExponentPushToken[") && strings.HasSuffix(token, "]")
}

// SendPushNotification allows to send the message
func SendPushNotification(message *PushMessage) (*PushNotificationResult, *PushNotificationError, error) {
	return message.Send()
}

func isError(err error) bool {
	return err != nil
}

func gZipBody(body []byte) ([]byte, bool, error) {
	if len(body) < MaxBodySizeWithoutGzip {
		return body, false, nil
	}

	var err error
	var b bytes.Buffer

	w := zlib.NewWriter(&b)

	if _, err = w.Write(body); isError(err) {
		return nil, false, err
	}

	w.Close()
	return b.Bytes(), true, nil
}

func processResponse(body io.ReadCloser) (r []*PushNotificationResult, e []*PushNotificationError, err error) {
	var response pushNotificationResponse

	result, _ := ioutil.ReadAll(body)
	fmt.Println(string(result))
	if err = json.Unmarshal(result, &response); isError(err) {
		return nil, nil, err
	}

	return response.Data, response.Errors, nil
}

// SendPushNotifications allows to send several messages at the same times
// Is highly recommanded to not send more than 100 messages at once
func SendPushNotifications(messages []*PushMessage) (r []*PushNotificationResult, e []*PushNotificationError, err error) {
	var body []byte
	var gzipped bool

	if body, err = json.Marshal(messages); isError(err) {
		return nil, nil, err
	}

	if body, gzipped, err = gZipBody(body); isError(err) {
		return nil, nil, err
	}

	var req *http.Request

	if req, err = http.NewRequest("POST", baseAPIURL+"/push/send", bytes.NewBuffer(body)); isError(err) {
		return nil, nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("User-Agent", "exponent-server-sdk-node/"+version)
	req.Header.Set("Content-Type", "application/json")

	if gzipped {
		req.Header.Set("Content-Encoding", "gzip")
	}

	var resp *http.Response

	client := &http.Client{}
	if resp, err = client.Do(req); isError(err) {
		return nil, nil, err
	}
	defer resp.Body.Close()

	return processResponse(resp.Body)
}

// ChunkPushNotifications returns an array of chunks
// The chunks size is determined with the ChunkLimit variable
func ChunkPushNotifications(messages []*PushMessage) [][]*PushMessage {
	size := 1
	if len(messages) >= ChunkLimit {
		size = int(math.Ceil(float64(len(messages)) / float64(ChunkLimit)))
	}

	Chunks := make([][]*PushMessage, size)

	var Chunk int
	for _, message := range messages {
		Chunks[Chunk] = append(Chunks[Chunk], message)

		if len(Chunks[Chunk]) >= ChunkLimit {
			Chunk++
		}
	}

	return Chunks
}
