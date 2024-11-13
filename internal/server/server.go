package server

import (
	"BolshoiGolangProject/internal/storage/storage"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	engine.Static("/static", "/app/web/static/")
	engine.NoRoute(func(c *gin.Context) {
		c.File("/app/web/static/index1.html")
	})
	arrayPoint := engine.Group("/array")
	scalarPoint := engine.Group("/scalar")

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	scalarPoint.POST("/set/:key", r.HandlerSet)
	scalarPoint.GET("/get/:key", r.HandlerGet)

	arrayPoint.POST("/lpush/:key", r.handlerArrLpush)
	arrayPoint.POST("/rpush/:key", r.handlerArrRpush)
	arrayPoint.POST("/raddtoset/:key", r.handlerRaddtoset)
	arrayPoint.POST("/lpop/:key", r.handlerLpopArr)
	arrayPoint.POST("/rpop/:key", r.handlerRpopArr)
	arrayPoint.POST("/lset/:key", r.handlerArrLSet)
	arrayPoint.POST("/expire/:key/:expireSeconds", r.handlerExpireSet)
	arrayPoint.GET("/lget/:key", r.handlerArrLGet)
	arrayPoint.GET("/getArr/:key", r.handlerArrGet)

	return engine
}

// @Summary		Время жизни значения
// @Description	Установить время жизни значения по ключу
// @Param			time	query	int	false	"Время жизни значения"
// @Router			/array/expire/:key/:expireSeconds [post]
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

// @Summary		Получить значение по индексу по ключу
// @Router			/array/lGet/:key [get]
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

// @Summary		Установить значение по индексу по ключу
// @Router			/array/lSet/:key [get]
func (r *Server) handlerArrLSet(ctx *gin.Context) {
	key := ctx.Param(KeyParam)

	var v Entry
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r.storage.LSet(key, uint64(v.ListInt[0]), v.Value)
	ctx.Status(http.StatusOK)
}

// @Summary		Удалить значение по индексу по ключу справа
// @Router			/array/rpop/:key [post]
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

// @Summary		Удалить значение по индексу по ключу слева
// @Router			/array/lpop/:key [post]
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

// @Summary		Вставляет по ключу только уникальные значения в массив
// @Router			/array/raddtoset/:key [post]
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

// @Summary		Вставляет слева все значения по ключу
// @Router			/array/lpush/:key [post]
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

// @Summary		Вставляет справа все значения по ключу
// @Router			/array/rpush/:key [post]
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

// @Summary		Получить списко по ключу
// @Router			/array/getArr/:key [get]
func (r *Server) handlerArrGet(ctx *gin.Context) {
	key := ctx.Param(KeyParam)
	v, expireTime, err := r.storage.CheckArr(key)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, ArrGet{List: v, ExpireAt: expireTime})
}

// @Summary		Установить скаляр по ключу
// @Description	key in path, value in json
// @Param			name	query		string	false	"Name of the user"
// @Success		200		{string}	string	"Hello, {name}"
// @Router			/scalar/set/:key [post]
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

// @Summary		Получить скаляр по ключу
// @Router			/scalar/get/:key [get]
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
