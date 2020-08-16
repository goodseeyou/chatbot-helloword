package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// not good to do this.
const debugDoesCheckXSigniture = true

const requestHeaderXLineSignature = "X-Line-Signature"

type eventSource struct {
	UserID     string `json:"userId"`
	SourceType string `json:"type"`
	GroupID    string `json:"groupId"`
	RoomID     string `json:"roomId"`
}

type emojiMessage struct {
	Index     int    `json:"index"`
	Length    int    `json:"lenght"`
	ProductID string `json:"productId"`
	EmojiID   string `json:"emojiId"`
}

type messageContentProvider struct {
	ContentType        string `json:"type"`
	OriginalContentURL string `json:"originalContentUrl"`
	PreviewImageURL    string `json:"previewImageUrl"`
}

type messageUnsent struct {
	MessageID string `json:"messageId"`
}

type eventMessage struct {
	MessageType         string                 `json:"type"`
	ID                  string                 `json:"id"`
	Text                string                 `json:"text"`
	Emojis              []emojiMessage         `json:"emojis"`
	PacakgeID           string                 `json:"packageId"`
	StickerID           string                 `json:"stickerId"`
	StickerResourceType string                 `json:"stickerResourceType"`
	Duration            int                    `json:"duration"`
	ContentProvider     messageContentProvider `json:"contentProvider"`
	FileName            string                 `json:"filename"`
	FileSize            int                    `json:"fileSize"`
	Title               string                 `json:"title"`
	Address             string                 `json:"address"`
	Latitude            float64                `json:"latitude"`
	Longtitude          float64                `json:"longtitude"`
}

type event struct {
	MessageType   string        `json:"type"`
	ReplyToken    string        `json:"replyToken"`
	Source        eventSource   `json:"source"`
	Timestamp     int           `json:"timestamp"`
	Mode          string        `json:"mode"`
	Message       eventMessage  `json:"message"`
	MessageUnsent messageUnsent `json:"unsend"`
	// checkpoint member join
}

type webhookEventOjbect struct {
	Events      []event `json:"events"`
	Destination string  `json:"destination"`
}

func verifySignature(body, xLineSignature string) error {
	if xLineSignature == "" {
		return fmt.Errorf("No value of %s from request header", requestHeaderXLineSignature)
	}

	// work around
	lineBotScret := os.Getenv("LineBotScret")

	h := hmac.New(sha256.New, []byte(lineBotScret))
	io.WriteString(h, body)
	requestSignature, err := base64.StdEncoding.DecodeString(xLineSignature)
	if err != nil {
		return err
	}
	expectSignature := h.Sum(nil)

	if !hmac.Equal(expectSignature, requestSignature) {
		return fmt.Errorf("Signatures are NOT equal. request: %s | should be %s",
			xLineSignature,
			base64.StdEncoding.EncodeToString(expectSignature))
	}

	return nil
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	chatbotV1G := r.Group("/chatbot/v1")
	chatbotV1G.GET("/healthyCheck", func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})

	chatbotV1G.POST("/message", func(c *gin.Context) {
		xLineSignature := c.Request.Header.Get(requestHeaderXLineSignature)
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if err := verifySignature(string(body), xLineSignature); debugDoesCheckXSigniture && err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		tmp := webhookEventOjbect{}
		if err := json.Unmarshal(body, &tmp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		tmps, err := json.Marshal(tmp)
		fmt.Print(string(tmps))

		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8088")
}
