package controller

import (
	"task-recommender/internal/service"
)

type TaskController struct {
	service *service.TaskService
}

func NewTaskController(service *service.TaskService) *TaskController {
	return &TaskController{service: service}
}

func (c *TaskController) AddTask(title, description string) (int, error) {
	return c.service.AddTask(title, description)
}

func (c *TaskController) ListTasks() (interface{}, error) {
	return c.service.ListTasks()
}

func (c *TaskController) CompleteTask(id int) error {
	return c.service.CompleteTask(id)
}

func (c *TaskController) DeleteTask(id int) error {
	return c.service.DeleteTask(id)
}
