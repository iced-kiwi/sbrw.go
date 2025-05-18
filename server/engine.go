package server

import (
	"log"

	"gosbrw/server/middleware"

	"github.com/gin-gonic/gin"
)

type Engine struct {
	router *gin.Engine
}

func NewEngine() *Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.CORS())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	engine := &Engine{
		router: r,
	}

	return engine
}

func (e *Engine) RegisterRoutes(registerFn func(router *gin.Engine)) {
	registerFn(e.router)
}

func (e *Engine) RegisterGroupRoutes(path string, registerFn func(group *gin.RouterGroup)) {
	group := e.router.Group(path)
	registerFn(group)
}

func (e *Engine) Start(address string) error {
	log.Printf("Starting server on %s", address)
	return e.router.Run(address)
}

func (e *Engine) GetRouter() *gin.Engine {
	return e.router
}
