package main

import (
	"fmt"
	gdocker "github.com/docker/go-docker"
	"morchestrator/internal/task"
	"os"
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
	fmt.Printf("create a test container\n")
	dockerTask, createResult := createContainer()
	if createResult.Error != nil {
		fmt.Printf("%v", createResult.Error)
		os.Exit(1)
	}
	time.Sleep(time.Second * 5)
	fmt.Printf("stopping container %s\n", createResult.ContainerId)
	_ = stopContainer(dockerTask, createResult.ContainerId)
}
