package main

import (
	"fmt"
	// "log"
	"net/http"
	"os/exec"
	"strconv"

	"time"

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
	ID       string `json:"id"`
	Name     string `json:"name"`
	Messages string `json:"messages"`
}

type message struct {
	ID      string    `json:"id"`
	Content string    `json:"content"`
	Poster  string    `json:"poster"`
	Posted  time.Time `json:"posted"`
}

// mapping of workspace ids to the users inside them
// maybe need reverse mapping of user ids to workspaces they are a part of

var users = []user{
	{ID: "00000000-0000-0000-0000-000000000000", Name: "Noah Libeskind", Email: "nlibeski@ucsc.edu", Password: "1651623", AccessToken: ""},
}

var workspaces = []workspace{
	{ID: "00000000-0000-0000-0000-000000000000", Name: "(WWW) World Wide Workspace", Channels: "0", Owner: "00000000-0000-0000-0000-000000000000"},
}

var channels = []channel{
	{ID: "00000000-0000-0000-0000-000000000000", Name: "World Chat Channel", Messages: "0"},
}

var messages = []message{}

// maps workspace IDs to IDs of users in that workspace
var workspace_users = map[string][]string{}

// maps workspace IDs to IDs of channels in that workspace
var workspace_channels = map[string][]string{}

// maps channel IDs to IDs of messages in that channel
var channel_messages = map[string][]string{}

var mySigningKey = []byte("notslackisnotsecure")

// https://play.golang.org/p/Qg_uv_inCek
// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// get all users
func getUsers(context *gin.Context) {
	tokenStatus, err := utils.ExtractTokenID(context)

	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, users)
}

// get all workspaces
// todo: only get ws associated with a user from accessToken
func getWorkSpaces(context *gin.Context) {
	tokenStatus, err := utils.ExtractTokenID(context)

	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, workspaces)
}

// creates a new workspace with logged in user as the owner
// returns the created workspace
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
		for _, u := range users {
			if u.AccessToken == token {
				// get ID from AccessToken
				newWS.Owner = u.ID
			}
		}
	}
	// not current logged in accessToken
	if newWS.Owner == "" {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Bad Credentials"})
		return
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

// creates a new workspace with logged in user as the owner
// returns the created workspace
func deleteWorkSpace(context *gin.Context) {

	tokenStatus, _ := utils.ExtractTokenID(context)
	token := utils.ExtractToken(context)

	if tokenStatus == 0 {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Bad Credentials"})
		return
	}

	id := context.Param("id")
	deleteIndex := -1
	for i, w := range workspaces {
		if w.ID == id {
			deleteIndex = i
			// check owner is logged in user
			for _, u := range users {
				if u.AccessToken == token {
					// get ID from AccessToken, if not owner, return err
					if u.ID != w.Owner {
						context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Bad Credentials"})
						return
					}
				}
			}
		}
	}
	if deleteIndex == -1 {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not found"})
		return
	} else {

		// remove all channels in workspace
		for _, chID := range workspace_channels[id] {
			for channelIndex, channel := range channels {
				if channel.ID == chID {
					// delete messages in channel
					for _, mID := range channel_messages[chID] {
						for messageIndex, message := range messages {
							if message.ID == mID {
								// delete message
								messages[messageIndex] = messages[len(messages)-1]
								messages = messages[0 : len(messages)-1]
							}
						}
					}
					// delete channel
					channels[channelIndex] = channels[len(channels)-1]
					channels = channels[0 : len(channels)-1]
				}
			}
		}
		// delete workspace
		workspaces[deleteIndex] = workspaces[len(workspaces)-1]
		workspaces = workspaces[0 : len(workspaces)-1]

		context.IndentedJSON(http.StatusOK, workspaces)
	}
	return
}

// adds a member with id memId to workspace with id wsId
// returns all members in that workspace
// todo: make sure owner of wsId == loggedInUser
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

				var wsMembers = []user{}
				for _, t := range users {
					if contains(workspace_users[wsId], t.ID) {
						wsMembers = append(wsMembers, t)
					}
				}
				context.IndentedJSON(http.StatusOK, wsMembers)
				// context.IndentedJSON(http.StatusOK, workspace_users[wsId])
				return
			}
		}
	}
}

// deletes the member with id memId to workspace with id wsId
// returns all members in that workspace
// todo: make sure owner of wsId == loggedInUser and EVERYTHING TO MAKE THIS WORK
func deleteWorkSpaceMember(context *gin.Context) {
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

				var wsMembers = []user{}
				for _, t := range users {
					if contains(workspace_users[wsId], t.ID) {
						wsMembers = append(wsMembers, t)
					}
				}
				context.IndentedJSON(http.StatusOK, wsMembers)
				// context.IndentedJSON(http.StatusOK, workspace_users[wsId])
				return
			}
		}
	}
}

// returns all users in the workspace with specified id
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
				var wsMembers = []user{}
				for _, t := range users {
					if contains(workspace_users[id], t.ID) {
						wsMembers = append(wsMembers, t)
					}
				}
				context.IndentedJSON(http.StatusOK, wsMembers)
				//context.IndentedJSON(http.StatusOK, workspace_users[id])
				return
			}
		}
	}

	context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no workspace found with this id"})
	return
}

// adds a member with id memId to workspace with id wsId
// returns all channels in that workspace
// todo: make sure owner of wsId == loggedInUser
func createChannel(context *gin.Context) {
	// workspace id
	id := context.Param("id")

	var newChannel channel
	// JSON should include channel name only
	err := context.BindJSON(&newChannel)

	tokenStatus, err := utils.ExtractTokenID(context)
	// verify auth
	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		for i, t := range workspaces {
			if t.ID == id {
				// found a workspace with this id, add channel to it
				newUUID, err := exec.Command("uuidgen").Output()
				if err != nil {
					return
				}
				newChannel.ID = string(newUUID)[0 : len(newUUID)-1]

				// add id to map
				_, ok := workspace_channels[id]
				if ok {
					workspace_channels[id] = append(workspace_channels[id], newChannel.ID)
				} else {
					workspace_channels[id] = []string{newChannel.ID}
				}
				newChannel.Messages = "0"
				count, _ := strconv.Atoi(workspaces[i].Channels)
				workspaces[i].Channels = strconv.Itoa(count + 1)
				channels = append(channels, newChannel)

				var wsChannels = []channel{}
				for _, t := range channels {
					if contains(workspace_channels[id], t.ID) {
						wsChannels = append(wsChannels, t)
					}
				}
				context.IndentedJSON(http.StatusOK, wsChannels)
				//context.IndentedJSON(http.StatusOK, workspace_channels[id])
				return
			}
		}
	}
}

// returns all channels in the specified workspace
func getChannels(context *gin.Context) {
	tokenStatus, err := utils.ExtractTokenID(context)

	// workspace id
	id := context.Param("id")

	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var wsChannels = []channel{}
	for _, t := range channels {
		if contains(workspace_channels[id], t.ID) {
			wsChannels = append(wsChannels, t)
		}
	}
	context.IndentedJSON(http.StatusOK, wsChannels)
}

// adds a member with id memId to workspace with id wsId
// returns all messages in the channel
// todo: make sure owner of wsId == loggedInUser
func createMessage(context *gin.Context) {
	// channel id
	id := context.Param("id")

	var newMessage message
	// JSON should include message content only
	err := context.BindJSON(&newMessage)

	tokenStatus, err := utils.ExtractTokenID(context)
	token := utils.ExtractToken(context)

	// verify auth
	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		for _, t := range users {
			if t.AccessToken == token {
				// get ID from AccessToken
				newMessage.Poster = t.ID
			}
		}
		for i, t := range channels {
			if t.ID == id {
				// found a channel with this id, add message to it
				newUUID, err := exec.Command("uuidgen").Output()
				if err != nil {
					return
				}
				newMessage.ID = string(newUUID)[0 : len(newUUID)-1]

				// add id to map
				_, ok := channel_messages[id]
				if ok {
					channel_messages[id] = append(channel_messages[id], newMessage.ID)
				} else {
					channel_messages[id] = []string{newMessage.ID}
				}
				// get current time
				newMessage.Posted = time.Now()
				count, _ := strconv.Atoi(channels[i].Messages)
				channels[i].Messages = strconv.Itoa(count + 1)
				messages = append(messages, newMessage)

				var chMessages = []message{}
				for _, m := range messages {
					if contains(channel_messages[id], m.ID) {
						chMessages = append(chMessages, m)
					}
				}
				context.IndentedJSON(http.StatusOK, chMessages)
				//context.IndentedJSON(http.StatusOK, channel_messages[id])
				return
			}
		}
	}
}

// returns all messages in the specified channel
func getMessages(context *gin.Context) {
	// channel id
	id := context.Param("id")
	tokenStatus, err := utils.ExtractTokenID(context)

	// verify auth
	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		for _, t := range channels {
			if t.ID == id {
				// found a channel with this id, add message to it
				var chMessages = []message{}
				for _, m := range messages {
					if contains(channel_messages[id], m.ID) {
						chMessages = append(chMessages, m)
					}
				}
				context.IndentedJSON(http.StatusOK, chMessages)
				//context.IndentedJSON(http.StatusOK, channel_messages[id])
				return
			}
		}
	}
	context.JSON(http.StatusBadRequest, gin.H{"message": "channel not found"})
	return
}

// creates a new user
// returns the user object
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

// logs in user with specified email and password
// returns the user object with a JWT if credentials are valid
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
	router := gin.Default()
	router.POST("/login", login)
	router.POST("/newuser", createUser)
	router.GET("/member", getUsers)

	router.GET("/workspace", getWorkSpaces)
	router.POST("/workspace", createWorkSpace)
	router.DELETE("/workspace/:id", deleteWorkSpace)

	router.GET("/workspace/:id/member", workSpaceMembers)
	router.POST("/workspace/:wsId/member/:memId", addWorkSpaceMember)

	router.POST("/workspace/channel/:id", createChannel)
	router.GET("/workspace/channel/:id", getChannels)

	router.POST("/channel/:id/message", createMessage)
	router.GET("/channel/:id/message", getMessages)

	router.Run("localhost:9090")
}
