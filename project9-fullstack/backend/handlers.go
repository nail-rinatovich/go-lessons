package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// taskInput is the JSON body for creating and updating tasks
type taskInput struct {
	Title       string `json:"title" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
	Done        bool   `json:"done"`
}

// TaskHandler holds the HTTP handlers for the /api/tasks resource
type TaskHandler struct {
	store *TaskStore
}

// registerRoutes attaches the task endpoints to a route group
func (h *TaskHandler) registerRoutes(rg *gin.RouterGroup) {
	rg.GET("/tasks", h.list)
	rg.POST("/tasks", h.create)
	rg.GET("/tasks/:id", h.get)
	rg.PUT("/tasks/:id", h.update)
	rg.DELETE("/tasks/:id", h.delete)
}

// GET /api/tasks?status=all|active|done
func (h *TaskHandler) list(c *gin.Context) {
	status := c.DefaultQuery("status", "all")
	if status != "all" && status != "active" && status != "done" {
		jsonError(c, http.StatusBadRequest, "status must be one of: all, active, done")
		return
	}
	tasks, err := h.store.List(status)
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// GET /api/tasks/:id
func (h *TaskHandler) get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	task, err := h.store.Get(id)
	if err != nil {
		storeError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

// POST /api/tasks
func (h *TaskHandler) create(c *gin.Context) {
	input, ok := bindInput(c)
	if !ok {
		return
	}
	task, err := h.store.Create(input.Title, input.Description)
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, task)
}

// PUT /api/tasks/:id
func (h *TaskHandler) update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	input, ok := bindInput(c)
	if !ok {
		return
	}
	task, err := h.store.Update(id, input.Title, input.Description, input.Done)
	if err != nil {
		storeError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

// DELETE /api/tasks/:id
func (h *TaskHandler) delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.store.Delete(id); err != nil {
		storeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// bindInput parses and validates the JSON body
func bindInput(c *gin.Context) (taskInput, bool) {
	var input taskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		jsonError(c, http.StatusBadRequest, "invalid input: title is required (max 100 chars), description max 500 chars")
		return input, false
	}
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)
	if input.Title == "" {
		jsonError(c, http.StatusBadRequest, "title must not be blank")
		return input, false
	}
	return input, true
}

// parseID reads the :id URL parameter
func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		jsonError(c, http.StatusBadRequest, "invalid task id")
		return 0, false
	}
	return id, true
}

// storeError maps store errors to HTTP responses
func storeError(c *gin.Context, err error) {
	if errors.Is(err, ErrNotFound) {
		jsonError(c, http.StatusNotFound, err.Error())
		return
	}
	serverError(c, err)
}

func serverError(c *gin.Context, err error) {
	log.Println("internal error:", err)
	jsonError(c, http.StatusInternalServerError, "internal server error")
}

func jsonError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
