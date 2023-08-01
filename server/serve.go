package server

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	port   string
	engine *gin.Engine
}

func NewServer(port string, engine *gin.Engine) (*Server, error) {
	return &Server{
		port:   port,
		engine: engine,
	}, nil
}

func (s *Server) Start() {

	serve := s.engine

	serve.Use(CORSmanager)
	gin.SetMode(gin.ReleaseMode)

	serve.POST("/api/players", func(ctx *gin.Context) {
		createPlayer(ctx)
	})
	serve.PUT("/api/players/:id", func(ctx *gin.Context) {
		updatePlayer(ctx)
	})
	serve.DELETE("/api/players/:id", func(ctx *gin.Context) {
		deletePlayer(ctx)
	})
	serve.GET("/api/players", func(ctx *gin.Context) {
		getAllPlayers(ctx)
	})
	serve.GET("/api/players/rank/:val", func(ctx *gin.Context) {
		getPlayerByRank(ctx)
	})
	serve.GET("/api/players/random", func(ctx *gin.Context) {
		getRandomPlayer(ctx)
	})

	serve.Run(s.port)
}
