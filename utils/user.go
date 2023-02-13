package utils

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/noahlibeskind/NotSlackAPI/data"
)

// creates a new user
// returns the user object
func CreateUser(context *gin.Context) {
	var newUser data.User

	err := context.BindJSON(&newUser)
	// encoded Json should include name, email, and password
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Server error"})
		return
	} else {
		for _, t := range data.Users {
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
		newUser.AccessToken, _ = GenerateToken()

		data.Workspace_users["00000000-0000-0000-0000-000000000000"] = append(data.Workspace_users["00000000-0000-0000-0000-000000000000"], newUser.ID)

		data.Users = append(data.Users, newUser)
		context.IndentedJSON(http.StatusOK, newUser)
		return
	}
}

// logs in user with specified email and password
// returns the user object with a JWT if credentials are valid
func Login(context *gin.Context) {
	var loginUser data.User

	err := context.BindJSON(&loginUser)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Server error"})
		return
	} else {
		for i, t := range data.Users {
			if t.Email == loginUser.Email {
				if t.Password == loginUser.Password {
					data.Users[i].AccessToken, err = GenerateToken()
					if err != nil {
						fmt.Printf("Err: %s", err)
						return
					}

					context.IndentedJSON(http.StatusOK, data.Users[i])
					return
				} else {
					context.IndentedJSON(http.StatusBadRequest, gin.H{"message": data.Bad_rq_message})
					return
				}
			}
		}
	}

	context.IndentedJSON(http.StatusBadRequest, gin.H{"message": data.Bad_rq_message})
	return
}

// get all users
func GetUsers(context *gin.Context) {
	tokenStatus, err := ExtractTokenID(context)
	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, data.Users)
}
