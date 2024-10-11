package api

import (
	"github.com/gin-gonic/gin"
	"morchestrator/internal/worker"
)

type Api struct {
	Host   string
	Port   string
	Router *gin.Engine
	Worker *worker.Worker
}

func (a *Api) InitRoute() {
	api := a.Router.Group("/api")
	{
		api.POST("/task", a.StartTaskHandler)
		api.DELETE("/task/:id", a.StopTaskHandler)
		api.GET("/tasks", a.GetListTaskHandler)
		api.GET("/stats", a.GetStatsHandler)
	}
}
