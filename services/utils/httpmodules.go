package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		msgSize, _ := strconv.ParseInt(req.URL.Path[1:], 0, 32)
		payload := make([]byte, msgSize)
		for i := 0; i < int(msgSize); i++ {
			payload[i] = 'a'
		}
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

func HttpRequest(c *http.Client, msgSize int) (*http.Response, error) {
	return c.Get(fmt.Sprintf("https://quic.localhost/%d", msgSize))
}
