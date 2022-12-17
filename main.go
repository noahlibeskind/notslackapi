package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type todo struct {
	ID        string `json:"id"`
	Item      string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos = []todo{
	{ID: "1", Item: "Clean Room", Completed: false},
	{ID: "2", Item: "Read book", Completed: false},
	{ID: "3", Item: "Record video", Completed: false},
}

func getTodos(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, todos)
}

func getTodoByID(id string) (*todo, error) {
	for i, t := range todos {
		if t.ID == id {
			return &todos[i], nil
		}
	}

	return nil, errors.New("not found")
}

func getTodo(context *gin.Context) {
	id := context.Param("id")
	todo, err := getTodoByID(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not found"})
		return
	}

	context.IndentedJSON(http.StatusOK, todo)
	return
}

func addTodo(context *gin.Context) {
	var newTodo todo

	if err := context.BindJSON(&newTodo); err != nil {
		return
	}

	todos = append(todos, newTodo)

	context.IndentedJSON(http.StatusCreated, newTodo)
}

func toggleTodoStatus(context *gin.Context) {
	id := context.Param("id")
	todo, err := getTodoByID(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not found"})
		return
	}
	// toggle status!
	todo.Completed = !todo.Completed
	context.IndentedJSON(http.StatusOK, todo)
	return
}

func deleteTodo(context *gin.Context) {
	id := context.Param("id")
	deleteIndex := -1
	for i, t := range todos {
		if t.ID == id {
			deleteIndex = i
		}
	}
	if deleteIndex == -1 {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not found"})
		return
	} else {
		// effectively removing, taking last element and putting it in deleteIndex's place
		todos[deleteIndex] = todos[len(todos)-1]
		todos = todos[0 : len(todos)-1]
		context.IndentedJSON(http.StatusOK, todos)

	}

	return
}

// func main() {
// 	router := gin.Default()
// 	router.GET("/todos", getTodos)
// 	router.GET("/todos/:id", getTodo)
// 	// Patch allows editing existing entries
// 	router.PATCH("/todos/:id", toggleTodoStatus)

// 	router.DELETE("/todos/:id", deleteTodo)
// 	router.POST("/todos", addTodo)
// 	router.Run("localhost:9090")
// }
