package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const RequestHeaderXLineSignature = "X-Line-Signature"

func VerifySignature(body, xLineSignature string) error {
	if xLineSignature == "" {
		return fmt.Errorf("No value of %s from request header", RequestHeaderXLineSignature)
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

	r.GET("/healthyCheck", func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})

	r.POST("/pushMsg", func(c *gin.Context) {
		xLineSignature := c.Request.Header.Get(RequestHeaderXLineSignature)
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		err = VerifySignature(string(body), xLineSignature)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8088")
}
