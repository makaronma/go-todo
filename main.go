package main

// API reference: https://developer.todoist.com/rest/v1/#getting-started

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Todo means a single todo item
type Todo struct {
	ID          string `form:"id" json:"id" xml:"id"  binding:"required"`
	Title       string `form:"title" json:"title" xml:"title" binding:"required"`
	IsCompleted bool   `json:"isCompleted"`
}

var todos = []Todo{
	{ID: "0", Title: "title 0", IsCompleted: true},
	{ID: "1", Title: "title 1", IsCompleted: false},
	{ID: "2", Title: "title 2", IsCompleted: false},
	{ID: "3", Title: "title 3", IsCompleted: false},
	{ID: "4", Title: "title 4", IsCompleted: false},
	{ID: "5", Title: "title 5", IsCompleted: false},
}

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("/todos", getTodos)
	router.GET("/todo/:id", getTodoByID)

	router.POST("/todo", addTodo)
	router.DELETE("/todo/:id", deleteTodo)

	// complete/uncomplete
	router.POST("/todo/:id/:updatedStatus", updateTodo)

	router.Run("localhost:8080")
}

func getTodos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todos)
}

func getTodoByID(c *gin.Context) {
	id := c.Param("id")

	for _, todo := range todos {
		if todo.ID == id {
			c.IndentedJSON(http.StatusOK, todo)
			return
		}
	}
}
func addTodo(c *gin.Context) {
	var newTodo Todo

	if err := c.ShouldBind(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todos = append(todos, newTodo)

	c.JSON(http.StatusCreated, newTodo)
}

func deleteTodo(c *gin.Context) {
	id := c.Param("id")
	var index int
	for i, todo := range todos {
		if todo.ID == id {
			index = i
			break
		}
	}
	todos = append(todos[:index], todos[index+1:]...)

	c.JSON(http.StatusNoContent, gin.H{
		"message": "deleted",
	})
}

func updateTodo(c *gin.Context) {
	id := c.Param("id")
	updatedStatus := c.Param("updatedStatus")

	var newBool bool
	if updatedStatus == "complete" {
		newBool = true
	} else if updatedStatus == "uncomplete" {
		newBool = false
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "param not valild",
		})
	}

	index := -1
	for i, todo := range todos {
		if todo.ID == id {
			index = i
			break
		}
	}
	if index >= 0 {
		todos[index].IsCompleted = newBool
		c.JSON(http.StatusNoContent, gin.H{
			"message": "updated to complete",
		})
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "not found",
	})
}
