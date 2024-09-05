package route

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	Router *gin.Engine
}

func NewServer() *Server {
	r := gin.New()

	r.Use(cors.Default())
	r.Use(logMiddleware)
	r.Use(errorHandler)

	r.NoRoute(func(c *gin.Context) {
		_ = c.Error(errors.New("not found"))
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	SetDefaultRoutes(r)

	s := &Server{Router: r}

	return s
}

func SetDefaultRoutes(r *gin.Engine) {
	// '/' -> 404(ingress 문제), '/status' -> 200 // fixme:: 계속 404내릴건지?
	r.GET("/status", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	r.HEAD("/status", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
}

type resLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w resLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func BindBodyJSON(c *gin.Context, obj any) error {
	return c.ShouldBindBodyWith(&obj, binding.JSON)
}

func logMiddleware(c *gin.Context) {
	start := time.Now()

	// record request body
	bodyMap := make(map[string]interface{})
	_ = BindBodyJSON(c, &bodyMap)
	body, _ := json.Marshal(bodyMap)

	// record response body
	blw := &resLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	// Process request
	c.Next()

	if !skipLogPath(c.Request.URL.Path) {
		//":method :url :status :response-time ms - :res[content-length] :: request: :req-body response: :res-body",
		log.WithFields(log.Fields{
			"tag":            "API",
			"method":         c.Request.Method,
			"url":            getPath(c.Request.URL.Path, c.Request.URL.RawQuery),
			"status":         c.Writer.Status(),
			"response-time":  time.Now().Sub(start).Seconds() * 1000,
			"content-length": blw.body.Len(),
			"request":        body,
			"response":       blw.body.String(),
		}).Info()
	}
}

func skipLogPath(path string) bool {
	// if '/', '/status' or contains 'swagger' -> true
	// else -> false
	return path == "/" || path == "/status" || strings.Contains(path, "swagger")
}

func getPath(path string, raw string) string {
	if raw != "" {
		return path + "?" + raw
	} else {
		return path
	}
}

func errorHandler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			switch err := r.(type) {
			case error:
				errorHandle(c, err)
			default:
				errorHandle(c, errors.New("unknown error")) // fixme:: cause 필드
			}
		}
	}()

	c.Next()

	if err := c.Errors.Last(); err != nil {
		errorHandle(c, err)
	}
}

func errorHandle(c *gin.Context, err error) {
	// fixme:: error 만들고 그 이후에 로그 상세화
	log.WithFields(log.Fields{
		"tag": "API",
		"err": err,
	}).Error("error")
	c.JSON(500, gin.H{"error": err.Error()})
}
