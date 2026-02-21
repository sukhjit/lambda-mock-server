package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sukhjit/util/pkg/stringz"
)

var errorLogger = log.New(os.Stderr, "[ERROR] ", log.Llongfile)

// New will create and return handler
func New() *gin.Engine {
	router := gin.Default()

	router.GET("/status", statusHandle)
	router.GET("/delay", responseHandler(delayHandle))
	router.POST("/delay", responseHandler(delayHandle))
	router.GET("/sort", sortItemsHandle)

	return router
}

func statusHandle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format("2006-01-02 15:04:05"),
	})
}

func delayHandle(c *gin.Context) (any, int, error) {
	waitTimeStr := c.DefaultQuery("time", "2")

	waitTime, err := strconv.Atoi(waitTimeStr)
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("invalid time, should be between 1-9")
	}
	if waitTime < 1 {
		waitTime = 1
	}
	if waitTime > 9 {
		waitTime = 9
	}

	time.Sleep(time.Duration(waitTime) * time.Second)

	result := fmt.Sprintf("wait done for %d seconds... Usage: ?time=n", waitTime)

	return result, http.StatusOK, nil
}

func sortItemsHandle(ctx *gin.Context) {
	items := struct {
		Items []string `json:"items"`
	}{}
	if err := ctx.BindJSON(&items); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	slices.Sort(items.Items)

	ctx.JSON(200, items.Items)
}

func responseHandler(h func(*gin.Context) (any, int, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, code, err := h(c)
		if err != nil {
			payload := errorResponse(code, err)

			c.JSON(code, payload)

			return
		}

		c.JSON(code, data)
	}
}

func errorResponse(code int, err error) map[string]string {
	// not 5xx error
	if code < http.StatusInternalServerError {
		return map[string]string{
			"error": err.Error(),
		}
	}

	// 5xx error
	errID := stringz.RandomString(8)

	errorLogger.Printf("ErrorID: %s, %v", errID, err)

	return map[string]string{
		"error": "internal server error",
		"code":  errID,
	}
}
