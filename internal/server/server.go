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
	Value    string   `json: "value"`
	ExpireAt int64    `json: expireAt`
	List     []string `json: "list"`
	ListInt  []int    `json: "listInt"`
	Element  string   `json: "element`
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

	engine.POST("/scalar/set/:key", r.HandlerSet)
	engine.GET("/scalar/get/:key", r.HandlerGet)

	engine.POST("/array/Lpush/:key", r.handlerArrLpush)     //+
	engine.POST("array/Rpush/:key", r.handlerArrRpush)      //+
	engine.POST("array/Raddtoset/:key", r.handlerRaddtoset) //+
	engine.POST("array/Lpop/:key", r.handlerLpopArr)        //+
	engine.POST("array/Rpop/:key", r.handlerRpopArr)        //+
	engine.POST("array/LSet/:key", r.handlerArrLSet)        //+
	engine.GET("array/LGet/:key", r.handlerArrLGet)         //+
	engine.GET("/array/getArr/:key", r.handlerArrGet)       //+

	return engine
}

func (r *Server) handlerArrLGet(ctx *gin.Context) {
	key := ctx.Param("key")

	var v Entry

	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	value, err := r.storage.LGet(key, v.ListInt[0])

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadGateway)
		return
	}

	ctx.JSON(http.StatusOK, Entry{Element: value}) //value
}

func (r *Server) handlerArrLSet(ctx *gin.Context) {
	key := ctx.Param("key")

	var v Entry

	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r.storage.LSet(key, uint64(v.ListInt[0]), v.Element)
	ctx.AbortWithStatus(http.StatusOK)
	r.storage.SaveToJSON("data.json")
}

func (r *Server) handlerRpopArr(ctx *gin.Context) {
	key := ctx.Param("key")

	var v Entry

	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r.storage.Rpop(key, v.ListInt)
}

func (r *Server) handlerLpopArr(ctx *gin.Context) {
	key := ctx.Param("key")

	var v Entry

	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r.storage.Lpop(key, v.ListInt)
	r.storage.SaveToJSON("data.json")
	ctx.AbortWithStatus(http.StatusOK)
}

func (r *Server) handlerRaddtoset(ctx *gin.Context) {
	key := ctx.Param("key")

	var v Entry

	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r.storage.Raddtoset(key, v.List)
	r.storage.SaveToJSON("data.json")
	ctx.AbortWithStatus(http.StatusOK)
}

func (r *Server) handlerArrLpush(ctx *gin.Context) {
	key := ctx.Param("key")

	var v Entry

	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r.storage.Lpush(key, v.List, v.ExpireAt)
	r.storage.SaveToJSON("data.json")
	ctx.Status(http.StatusOK)
}

func (r *Server) handlerArrRpush(ctx *gin.Context) {
	key := ctx.Param("key")

	var v Entry

	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r.storage.Rpush(key, v.List, v.ExpireAt)
	r.storage.SaveToJSON("data.json")
	ctx.Status(http.StatusOK)
}

func (r *Server) handlerArrGet(ctx *gin.Context) {
	key := ctx.Param("key")
	v, expireTime, err := r.storage.Check_arr(key)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, Entry{List: v, ExpireAt: expireTime})
}

func (r *Server) HandlerSet(ctx *gin.Context) {
	key := ctx.Param("key")

	var v Entry

	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	r.storage.Set(key, v.Value, v.ExpireAt)
	r.storage.SaveToJSON("data.json")
	ctx.Status(http.StatusOK)
}

func (r *Server) HandlerGet(ctx *gin.Context) {
	key := ctx.Param("key")

	v, expire, err := r.storage.Get(key)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, Entry{Value: v, ExpireAt: expire})
}

func (r *Server) Start() {
	r.newApi().Run(r.host)
}