package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/gin-gonic/gin"
	"github.com/sukhjit/lambda-mock-server/model"
	"github.com/sukhjit/lambda-mock-server/repo"
	"github.com/sukhjit/lambda-mock-server/repo/dynamodb"
	"github.com/sukhjit/util"
)

var (
	errorLogger         = log.New(os.Stderr, "[ERROR] ", log.Llongfile)
	documentRepo        repo.Document
	errInvalidID        = errors.New("Invalid document id")
	errDocumentNotFound = errors.New("Document not found")
)

// New will create and return handler
func New(awsRegion, documentTable string) *gin.Engine {
	router := gin.Default()

	documentRepo = dynamodb.New(awsRegion, documentTable)

	router.GET("/status", statusHandle)
	router.POST("/add", responseHandler(addDocHandle))
	router.GET("/get/:id", responseHandler(getDocHandle))
	router.GET("/delay", responseHandler(delayHandle))

	return router
}

func statusHandle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format("2006-01-02 15:04:05"),
	})
}

func addDocHandle(c *gin.Context) (interface{}, int, error) {
	document := &model.Document{
		Date: time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := c.BindJSON(&document); err != nil {
		return nil, http.StatusBadRequest, err
	}

	err := validateDocument(document)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// convert interface to json string
	document.Body = gabs.Wrap(document.Body).String()

	err = documentRepo.Add(document)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return gin.H{
		"created": document,
	}, http.StatusOK, nil
}

func getDocHandle(c *gin.Context) (interface{}, int, error) {
	docID := c.Param("id")

	document, err := documentRepo.Get(docID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if len(document.ID) == 0 {
		return nil, http.StatusNotFound, errDocumentNotFound
	}

	// convert json string to json
	jsonParsed, err := gabs.ParseJSON([]byte(fmt.Sprintf("%v", document.Body)))
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("Failed to parse json")
	}

	document.Body = jsonParsed

	return gin.H{
		"document": document,
	}, http.StatusOK, nil
}

func delayHandle(c *gin.Context) (interface{}, int, error) {
	waitTimeStr := c.DefaultQuery("time", "2")
	waitTime, err := strconv.Atoi(waitTimeStr)
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("Invalid time, should be between 1-9")
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

func validateDocument(document *model.Document) error {
	if len(document.ID) == 0 {
		return errors.New("ID cannot be empty")
	}

	if document.Body == nil {
		return errors.New("Body cannot be empty")
	}

	return nil
}

func responseHandler(h func(*gin.Context) (interface{}, int, error)) gin.HandlerFunc {
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
	errID := util.RandomString(8)

	errorLogger.Printf("ErrorID: %s, %v", errID, err)

	return map[string]string{
		"error": "internal server error",
		"code":  errID,
	}
}
