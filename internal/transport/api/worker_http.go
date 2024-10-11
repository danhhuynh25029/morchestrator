package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"morchestrator/internal/task"
	"net/http"
)

func (a *Api) StartTaskHandler(c *gin.Context) {
	var req task.Task
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a.Worker.AddTask(req)
	log.Printf("Add task :%v", req.ID)
	c.JSON(http.StatusCreated, gin.H{"task": req})
}

func (a *Api) GetListTaskHandler(c *gin.Context) {
	c.JSON(http.StatusOK, a.Worker.GetTasks())
}

func (a *Api) StopTaskHandler(c *gin.Context) {
	taskId := c.Param("id")
	if taskId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "taskId is required"})
		return
	}
	tID, _ := uuid.Parse(taskId)
	_, ok := a.Worker.Db[tID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	taskToStop := a.Worker.Db[tID]
	taskCopy := *taskToStop
	taskCopy.State = task.Completed
	a.Worker.AddTask(taskCopy)
	log.Printf("Added task %v to stop container %v\n", taskToStop.ID, taskToStop.ContainerID)

	c.JSON(http.StatusOK, gin.H{"task": taskToStop})
}

func (a *Api) GetStatsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, a.Worker.Stats)
}
