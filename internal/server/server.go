package server

import (
	"encoding/json"
	"net/http"
	"project_1/internal/storage/storage"

	"github.com/gin-gonic/gin"
)

type Server struct {
	host    string
	storage *storage.Storage
}

type Entry struct {
	Value string `json: "value"`
}

func New(host string, st *storage.Storage) *Server {
	s := &Server{
		host:    host,
		storage: st,
	}
	return s
}

func (r *Server) newApi() *gin.Engine {
	engine := gin.New()

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	engine.GET("/hello-world", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Working! ^^")
	})

	engine.POST("/scalar/set/:key", r.handlerSet)
	engine.GET("/scalar/get/:key", r.handlerGet)

	//engine.POST("")

	return engine
}

func (r *Server) handlerSet(ctx *gin.Context) {
	key := ctx.Param("key")

	var v Entry

	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	r.storage.Set(key, v.Value)
	r.storage.SaveToJSON("data.json")
	ctx.Status(http.StatusOK)
}

func (r *Server) handlerGet(ctx *gin.Context) {
	key := ctx.Param("key")

	v, err := r.storage.Get(key)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, Entry{Value: v})
}

func (r *Server) Start() {
	r.newApi().Run(r.host)
}
