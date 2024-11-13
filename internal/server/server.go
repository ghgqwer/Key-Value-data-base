package server

import (
	"encoding/json"
	"net/http"
	"project_1/internal/storage/storage"
	"strconv"

	"github.com/gin-gonic/gin"
)

const KeyParam = "key"
const DataJson = "data.json"

type Server struct {
	host    string
	storage *storage.Storage
}

type Entry struct {
	Value    string   `json: "value"`
	ExpireAt int64    `json: expireAt`
	List     []string `json: "list"`
	ListInt  []int    `json: "listInt"`
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

	arrayPoint := engine.Group("/array")
	scalarPoint := engine.Group("/scalar")

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	scalarPoint.POST("/set/:key", r.HandlerSet)
	scalarPoint.GET("/get/:key", r.HandlerGet)

	arrayPoint.POST("/Lpush/:key", r.handlerArrLpush)      //+
	arrayPoint.POST("/Rpush/:key", r.handlerArrRpush)      //+
	arrayPoint.POST("/Raddtoset/:key", r.handlerRaddtoset) //+
	arrayPoint.POST("/Lpop/:key", r.handlerLpopArr)        //+
	arrayPoint.POST("/Rpop/:key", r.handlerRpopArr)        //+
	arrayPoint.POST("/LSet/:key", r.handlerArrLSet)        //+
	arrayPoint.POST("/Expire/:key/:expireSeconds", r.handlerExpireSet)
	arrayPoint.GET("/LGet/:key", r.handlerArrLGet)  //+
	arrayPoint.GET("/getArr/:key", r.handlerArrGet) //+

	return engine
}

func (r *Server) handlerExpireSet(ctx *gin.Context) {
	key := ctx.Param(KeyParam)
	expireAt, err := strconv.Atoi(ctx.Param("expireSeconds"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	r.storage.Expire(key, int64(expireAt))
	ctx.AbortWithStatus(http.StatusOK)
}

type LGet struct {
	Value string
}

func (r *Server) handlerArrLGet(ctx *gin.Context) {
	key := ctx.Param(KeyParam)

	var v Entry
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if len(v.ListInt) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "Index didnt give")
		return
	}

	value, err := r.storage.LGet(key, v.ListInt[0])
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, LGet{Value: value})
}

func (r *Server) handlerArrLSet(ctx *gin.Context) {
	key := ctx.Param(KeyParam)

	var v Entry
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r.storage.LSet(key, uint64(v.ListInt[0]), v.Value)
	ctx.AbortWithStatus(http.StatusOK)
}

func (r *Server) handlerRpopArr(ctx *gin.Context) {
	key := ctx.Param(KeyParam)

	var v Entry
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if _, err := r.storage.Rpop(key, v.ListInt[0], v.ListInt[1]); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
}

func (r *Server) handlerLpopArr(ctx *gin.Context) {
	key := ctx.Param(KeyParam)

	var v Entry
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if _, err := r.storage.Lpop(key, v.ListInt[0], v.ListInt[1]); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	ctx.AbortWithStatus(http.StatusOK)
}

func (r *Server) handlerRaddtoset(ctx *gin.Context) {
	key := ctx.Param(KeyParam)

	var v Entry
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r.storage.Raddtoset(key, v.List)
	ctx.AbortWithStatus(http.StatusOK)
}

func (r *Server) handlerArrLpush(ctx *gin.Context) {
	key := ctx.Param(KeyParam)

	var v Entry
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if _, err := r.storage.Lpush(key, v.List, v.ExpireAt); err != nil {
		ctx.AbortWithStatus(http.StatusConflict)
	}
	ctx.Status(http.StatusOK)
}

func (r *Server) handlerArrRpush(ctx *gin.Context) {
	key := ctx.Param(KeyParam)

	var v Entry
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if _, err := r.storage.Rpush(key, v.List, v.ExpireAt); err != nil {
		ctx.AbortWithStatus(http.StatusConflict)
		return
	}
	ctx.Status(http.StatusOK)
}

type ArrGet struct {
	List     []string
	ExpireAt int64
}

func (r *Server) handlerArrGet(ctx *gin.Context) {
	key := ctx.Param(KeyParam)
	v, expireTime, err := r.storage.CheckArr(key)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, ArrGet{List: v, ExpireAt: expireTime})
}

func (r *Server) HandlerSet(ctx *gin.Context) {
	key := ctx.Param(KeyParam)

	var v Entry
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := r.storage.Set(key, v.Value, v.ExpireAt); err != nil {
		ctx.AbortWithStatus(http.StatusConflict)
		return
	}
	ctx.Status(http.StatusOK)
}

func (r *Server) HandlerScalarExpire(ctx *gin.Context) {
	key := ctx.Param(KeyParam)
	expireAt, err := strconv.Atoi(ctx.Param("expireSeconds"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var v Entry
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	r.storage.Set(key, v.Value, int64(expireAt))
	ctx.AbortWithStatus(http.StatusOK)

}

type getValues struct {
	Value    string
	ExpireAt int64
}

func (r *Server) HandlerGet(ctx *gin.Context) {
	key := ctx.Param(KeyParam)
	v, expire, err := r.storage.Get(key)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, getValues{Value: v, ExpireAt: expire})
}

func (r *Server) Start() {
	r.newApi().Run(r.host)
}
