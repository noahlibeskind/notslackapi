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

var mySigningKey = []byte("mysecretphrase")

// passed when login is called
// func generateJWT() (string, error) {
// 	token := jwt.New(jwt.SigningMethodHS256)
// 	claims := token.Claims.(jwt.MapClaims)

// 	claims["authorized"] = true
// 	claims["user"] = "Noah Libeskind"
// 	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

// 	tokenString, err := token.SignedString(mySigningKey)

// 	if err != nil {
// 		fmt.Errorf("Something went wrong: %s", err.Error())
// 		return "", err
// 	}

// 	return tokenString, nil
// }

// check this anytime information is provided to client
// func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.Header["Token"] != nil {
// 			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
// 				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 					return nil, fmt.Errorf("There was an error")
// 				}
// 				return mySigningKey, nil
// 			})

// 			if err != nil {
// 				fmt.Fprintf(w, err.Error())
// 			}
// 			if token.Valid {
// 				endpoint(w, r)
// 			}
// 		} else {
// 			fmt.Fprintf(w, "Not Authorized")
// 		}
// 	})
// }

func getUsers(context *gin.Context) {
	tokenStatus, err := utils.ExtractTokenID(context)

	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, users)
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

		// query := `INSERT INTO users VALUES ($1, $2, $3, $4, $5)`
		// _, err2 := q.Exec(query, newUser.ID, newUser.Name, newUser.Email, newUser.Password, newUser.AccessToken)
		// if err2 != nil {
		// 	context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		// 	return
		// }
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

// func addTodo(context *gin.Context) {
// 	var newTodo todo

// 	if err := context.BindJSON(&newTodo); err != nil {
// 		return
// 	}

// 	todos = append(todos, newTodo)

// 	context.IndentedJSON(http.StatusCreated, newTodo)
// }

// func toggleTodoStatus(context *gin.Context) {
// 	id := context.Param("id")
// 	todo, err := getTodoByID(id)
// 	if err != nil {
// 		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not found"})
// 		return
// 	}
// 	// toggle status!
// 	todo.Completed = !todo.Completed
// 	context.IndentedJSON(http.StatusOK, todo)
// 	return
// }

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

	protected := router.Group("/api/admin")
	protected.Use(utils.JwtAuthMiddleware())

	protected.GET("/member", getUsers)
	// router.GET("/todos/:id", getTodo)
	// // Patch allows editing existing entries
	// router.PATCH("/todos/:id", toggleTodoStatus)

	// router.DELETE("/todos/:id", deleteTodo)
	// router.POST("/todos", addTodo)
	router.Run("localhost:9090")
}
