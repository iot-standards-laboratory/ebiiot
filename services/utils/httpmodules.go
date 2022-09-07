package utils

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		dec := json.NewDecoder(req.Body)
		body := map[string]float64{}
		dec.Decode(&body)
		msgSize := int32(body["msg_size"])
		payload := make([]byte, msgSize)
		resp.Write(payload)
	})

	return mux
}

func NewGinMux() *gin.Engine {
	mux := gin.New()
	mux.Any("/*any", func(c *gin.Context) {
		body := map[string]float64{}
		c.BindJSON(&body)
		msgSize := int32(body["msg_size"])
		payload := make([]byte, msgSize)
		c.Writer.Write(payload)
	})

	return mux
}

func HttpRequest(c *http.Client, msgSize int32) (*http.Response, error) {

	payload, _ := json.Marshal(&map[string]int32{
		"msg_size": msgSize,
	})

	return c.Post("https://quic.localhost", "application/json", bytes.NewReader(payload))
}
