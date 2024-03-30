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

// deletes a user
// returns 204 no content or 404 if not found
func DeleteUser(context *gin.Context) {
	tokenStatus, _ := ExtractTokenID(context)
	token := ExtractToken(context)
	if tokenStatus == 0 {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{"message": data.Unauthorized_message})
		return
	}
	id := context.Param("uid")
	deleteIndex := -1
	for i, u := range data.Users {
		// check that user to delete is logged in user
		if u.AccessToken == token {
			// get ID from AccessToken, if not owner, return err
			if u.ID != id {
				context.IndentedJSON(http.StatusUnauthorized, gin.H{"message": data.Unauthorized_message})
				return
			}
		}
		if u.ID == id {
			deleteIndex = i
		}
	}
	if deleteIndex == -1 {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not found"})
		return
	} else {
		// remove from users array
		data.Users[deleteIndex] = data.Users[len(data.Users)-1]
		data.Users = data.Users[0 : len(data.Users)-1]
		// look in each workspace
		for _, workspace := range data.Workspaces {
			// scan the users in that workspace
			for wuIndex, wuId := range data.Workspace_users[workspace.ID] {
				if wuId == id {
					data.Workspace_users[workspace.ID][wuIndex] = data.Workspace_users[workspace.ID][len(data.Workspace_users[workspace.ID])-1]
					data.Workspace_users[workspace.ID] = data.Workspace_users[workspace.ID][0 : len(data.Workspace_users[workspace.ID])-1]
				}
			}
		}

		for mIndex, m := range data.Messages {
			// scan for any messages posted by this member
			if m.Member == id {
				data.Messages[mIndex] = data.Messages[len(data.Messages)-1]
				data.Messages = data.Messages[0 : len(data.Messages)-1]
				// reduce channel's message count:
				for cIndex, c := range data.Channels {
					if contains(data.Channel_messages[c.ID], m.ID) {
						data.Channels[cIndex].Messages -= 1
					}
				}
			}
		}
	}
	context.IndentedJSON(http.StatusOK, nil)
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
					context.IndentedJSON(http.StatusUnauthorized, gin.H{"message": data.Unauthorized_message})
					return
				}
			}
		}
	}

	context.IndentedJSON(http.StatusUnauthorized, gin.H{"message": data.Unauthorized_message})
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
