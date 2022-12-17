package main

import (
	"fmt"
	// "log"
	"net/http"
	"os/exec"

	// "time"

	"github.com/noahlibeskind/NotSlackAPI/utils"

	// jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type user struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	AccessToken string `json:"accessToken"`
}

type workspace struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Channels string `json:"channels"`
	Owner    string `json:"owner"`
}

type channel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Messages    string `json:"messages"`
	WorkspaceID string
}

type message struct {
	ID     string `json:"id"`
	Poster string `json:"poster"`
	Posted string `json:"posted"`
}

// mapping of workspace ids to the users inside them
// maybe need reverse mapping of user ids to workspaces they are a part of

var users = []user{
	{ID: "00000000-0000-0000-0000-000000000000", Name: "Noah Libeskind", Email: "nlibeski@ucsc.edu", Password: "1651623", AccessToken: ""},
}

var workspaces = []workspace{
	{ID: "00000000-0000-0000-0000-000000000000", Name: "(WWW) World Wide Workspace", Channels: "1"},
}

var workspace_users = map[string][]string{}

var mySigningKey = []byte("mysecretphrase")

func getUsers(context *gin.Context) {
	tokenStatus, err := utils.ExtractTokenID(context)

	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, users)
}

func getWorkSpaces(context *gin.Context) {
	tokenStatus, err := utils.ExtractTokenID(context)

	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, workspaces)
}

func createWorkSpace(context *gin.Context) {
	var newWS workspace

	tokenStatus, _ := utils.ExtractTokenID(context)
	token := utils.ExtractToken(context)

	err := context.BindJSON(&newWS)
	// encoded JSON should only include name
	if err != nil || tokenStatus == 0 {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": err})
		return
	} else {
		for _, t := range users {
			if t.AccessToken == token {
				// get ID from AccessToken
				newWS.Owner = t.ID
			}
		}
	}
	newUUID, err := exec.Command("uuidgen").Output()
	if err != nil {
		return
	}
	newWS.ID = string(newUUID)[0 : len(newUUID)-1]
	newWS.Channels = "0"
	workspaces = append(workspaces, newWS)
	context.IndentedJSON(http.StatusOK, newWS)
	return
}

func addWorkSpaceMember(context *gin.Context) {
	// workspace id
	wsId := context.Param("wsId")
	memId := context.Param("memId")

	tokenStatus, err := utils.ExtractTokenID(context)
	// verify auth
	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		for _, t := range workspaces {
			if t.ID == wsId {
				// found a workspace with this id, add member to it
				_, ok := workspace_users[wsId]
				if ok {
					workspace_users[wsId] = append(workspace_users[wsId], memId)
				} else {
					workspace_users[wsId] = []string{memId}
				}

				context.IndentedJSON(http.StatusOK, workspace_users[wsId])
				return
			}
		}
	}
}

func workSpaceMembers(context *gin.Context) {
	// workspace id
	id := context.Param("id")

	tokenStatus, err := utils.ExtractTokenID(context)
	// verify auth
	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		for _, t := range workspaces {
			if t.ID == id {
				// found a workspace with this id
				context.IndentedJSON(http.StatusOK, workspace_users[id])
				return
			}
		}
	}

	context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no workspace found with this id"})
	return
}

// func getTodoByID(id string) (*todo, error) {
// 	for i, t := range todos {
// 		if t.ID == id {
// 			return &todos[i], nil
// 		}
// 	}

//		return nil, errors.New("not found")
//	}
func createUser(context *gin.Context) {
	var newUser user

	err := context.BindJSON(&newUser)
	// encoded Json should include name, email, and password
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Server error"})
		return
	} else {
		for _, t := range users {
			if t.Email == newUser.Email {
				context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "User with this email already exists"})
				return
			}
		}
		newUUID, err := exec.Command("uuidgen").Output()
		if err != nil {
			return
		}
		newUser.ID = string(newUUID)[0 : len(newUUID)-1]
		newUser.AccessToken, _ = utils.GenerateToken()

		users = append(users, newUser)
		context.IndentedJSON(http.StatusOK, newUser)
		return
	}
}

func login(context *gin.Context) {
	var loginUser user

	err := context.BindJSON(&loginUser)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Server error"})
		return
	} else {
		for i, t := range users {
			if t.Email == loginUser.Email {
				if t.Password == loginUser.Password {
					users[i].AccessToken, err = utils.GenerateToken()
					if err != nil {
						fmt.Printf("Err: %s", err)
						return
					}
					fmt.Println("Token:")
					fmt.Printf("%s", users[i].AccessToken)
					context.IndentedJSON(http.StatusOK, users[i])
					return
				} else {
					context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad credentials"})
					return
				}
			}
		}
	}

	context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad credentials"})
	return
}

// func deleteTodo(context *gin.Context) {
// 	id := context.Param("id")
// 	deleteIndex := -1
// 	for i, t := range todos {
// 		if t.ID == id {
// 			deleteIndex = i
// 		}
// 	}
// 	if deleteIndex == -1 {
// 		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not found"})
// 		return
// 	} else {
// 		// effectively removing, taking last element and putting it in deleteIndex's place
// 		todos[deleteIndex] = todos[len(todos)-1]
// 		todos = todos[0 : len(todos)-1]
// 		context.IndentedJSON(http.StatusOK, todos)

// 	}

// 	return
// }

func main() {
	// newUUID, err := exec.Command("uuidgen").Output()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Generated UUID:")
	// fmt.Printf("%s", newUUID)

	router := gin.Default()
	router.POST("/login", login)
	router.POST("/newuser", createUser)

	router.GET("/member", getUsers)
	router.GET("/workspace", getWorkSpaces)

	router.POST("/workspace", createWorkSpace)
	router.GET("/workspace/:id/member", workSpaceMembers)
	router.POST("/workspace/:wsId/member/:memId", addWorkSpaceMember)
	// router.GET("/todos/:id", getTodo)
	// // Patch allows editing existing entries
	// router.PATCH("/todos/:id", toggleTodoStatus)

	// router.DELETE("/todos/:id", deleteTodo)
	// router.POST("/todos", addTodo)
	router.Run("localhost:9090")
}
