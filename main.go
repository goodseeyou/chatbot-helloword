package main

import (
	"github.com/goodseeyou/chatbot-helloworld/model/request"

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
const debugDoesCheckXSigniture = false

const requestHeaderXLineSignature = "X-Line-Signature"

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

		tmp := request.WebhookEventOjbect{}
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
