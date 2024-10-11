package main

import (
	"fmt"
	gdocker "github.com/docker/go-docker"
	"github.com/gin-gonic/gin"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"log"
	"morchestrator/internal/manager"
	"morchestrator/internal/task"
	"morchestrator/internal/transport/api"
	"morchestrator/internal/worker"
	"os"
	"strconv"
	"time"
)

func createContainer() (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name:  "morchestrator-mysql",
		Image: "mysql:8.0",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=root123",
		},
	}
	dc, _ := gdocker.NewEnvClient()
	d := task.Docker{
		Client: dc,
		Config: c,
	}
	result := d.Run()
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil, nil
	}
	fmt.Printf(
		"Container %s is running with config %v\n", result.ContainerId, c)
	return &d, &result
}

func stopContainer(d *task.Docker, id string) *task.DockerResult {
	result := d.Stop(id)
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil
	}
	fmt.Printf("Container %s has been stopped and removed\n", result.ContainerId)
	return &result
}

func main() {
	port, _ := strconv.Atoi(os.Getenv("MOR_PORT"))

	router := gin.Default()
	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	apiHttp := api.Api{
		Router: router,
		Worker: &w,
	}
	apiHttp.InitRoute()
	if port == 0 {
		port = 8080
	}

	workers := []string{fmt.Sprintf("%s:%d", "localhost", port)}
	m := manager.New(workers)
	go runTasks(&w)
	go w.CollectStats()
	go func() {
		if err := router.Run(fmt.Sprintf(":%v", port)); err != nil {
			panic(err)
		}
	}()

	for i := 0; i < 3; i++ {
		t := task.Task{
			ID:    uuid.New(),
			Name:  fmt.Sprintf("test-container-%d", i),
			State: task.Scheduled,
			Image: "hello-world",
		}
		te := task.TaskEvent{
			ID:    uuid.New(),
			State: task.Running,
			Task:  t,
		}
		m.AddTask(te)
		m.SendWorker()
	}
	go func() {
		for {
			fmt.Printf("[Manager] Updating tasks from %d workers\n", len(m.Workers))
			m.UpdateTasks()
			time.Sleep(15 * time.Second)
		}
	}()
	for {
		for _, t := range m.TaskDb {
			fmt.Printf("[Manager] Task: id: %s, state: %d\n", t.ID, t.State)
			time.Sleep(15 * time.Second)
		}
	}
}

func runTasks(w *worker.Worker) {
	for {
		if w.Queue.Len() != 0 {
			result := w.RunTask()
			if result.Error != nil {
				log.Printf("Error running task: %v\n", result.Error)
			}
		} else {
			log.Printf("No tasks to process currently.\n")
		}
		log.Println("Sleeping for 10 seconds.")
		time.Sleep(10 * time.Second)
	}
}
