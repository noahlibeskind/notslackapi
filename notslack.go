package main

import (
	"fmt"
	// "log"
	"net/http"
	"os/exec"

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
	Channels int    `json:"channels"`
	Owner    string `json:"owner"`
}

type channel struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Messages int    `json:"messages"`
}

type message struct {
	ID      string `json:"id"`
	Member  string `json:"member"`
	Posted  string `json:"posted"`
	Content string `json:"content"`
}

// mapping of workspace ids to the users inside them
// maybe need reverse mapping of user ids to workspaces they are a part of

var users = []user{
	{ID: "00000000-0000-0000-0000-000000000000", Name: "Noah Libeskind", Email: "noah@ucsc.edu", Password: "noah", AccessToken: ""},
}

var workspaces = []workspace{
	{ID: "00000000-0000-0000-0000-000000000000", Name: "(WWW) World Wide Workspace", Channels: 1, Owner: "00000000-0000-0000-0000-000000000000"},
}

var channels = []channel{
	{ID: "00000000-0000-0000-0000-000000000000", Name: "World Chat Channel", Messages: 1},
}

var messages = []message{
	{ID: "00000000-0000-0000-0000-000000000000", Member: "00000000-0000-0000-0000-000000000000", Posted: "2023-01-02T00:01:01ZZZ", Content: "Thanks for not providing this, Dr. Harrison"},
}

var bad_rq_message = "Invalid Credentials"

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

// get all other users (but not logged in user)
func getUsers(context *gin.Context) {
	tokenStatus, err := utils.ExtractTokenID(context)
	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, users)
}

// get all workspaces
func getWorkSpaces(context *gin.Context) {
	tokenStatus, _ := utils.ExtractTokenID(context)
	token := utils.ExtractToken(context)
	loggedInUser := ""

	// encoded JSON should only include name
	if tokenStatus == 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": bad_rq_message})
		return
	} else {
		for _, u := range users {
			if u.AccessToken == token {
				// get ID from AccessToken
				loggedInUser = u.ID
			}
		}
	}
	var userWorkspaces = []workspace{}
	for _, w := range workspaces {
		if contains(workspace_users[w.ID], loggedInUser) || w.Owner == loggedInUser {
			userWorkspaces = append(userWorkspaces, w)
		}
	}
	context.IndentedJSON(http.StatusOK, userWorkspaces)
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
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": bad_rq_message})
		return
	}
	newUUID, err := exec.Command("uuidgen").Output()
	if err != nil {
		return
	}
	newWS.ID = string(newUUID)[0 : len(newUUID)-1]
	newWS.Channels = 0
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
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": bad_rq_message})
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
						context.IndentedJSON(http.StatusNotFound, gin.H{"message": bad_rq_message})
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

		delete(workspace_channels, id)
		delete(workspace_users, id)

		context.IndentedJSON(http.StatusOK, workspaces)
	}
	return
}

// adds a member with id memId to workspace with id wsId
// returns all members in that workspace
func addWorkSpaceMember(context *gin.Context) {
	// workspace id
	wsId := context.Param("wsId")
	memId := context.Request.URL.Query().Get("mid")

	tokenStatus, err := utils.ExtractTokenID(context)

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
func deleteWorkSpaceMember(context *gin.Context) {
	// workspace id
	wsId := context.Param("wsId")
	memId := context.Request.URL.Query().Get("mid")

	tokenStatus, err := utils.ExtractTokenID(context)
	token := utils.ExtractToken(context)

	// verify auth
	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deleteIndex := -1
	for i, w := range workspaces {
		if w.ID == wsId {
			deleteIndex = i
			// check owner is logged in user
			for _, u := range users {
				if u.AccessToken == token {
					// get ID from AccessToken, if user is not owner, return err
					if u.ID != w.Owner {
						context.IndentedJSON(http.StatusNotFound, gin.H{"message": bad_rq_message})
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
		for _, t := range workspaces {
			if t.ID == wsId {
				// found a workspace with this id
				_, ok := workspace_users[wsId]
				if ok {
					// found user with memID in workspace
					// workspace_users[wsId] = append(workspace_users[wsId], memId)
					for userIndex, user := range workspace_users[wsId] {
						if user == memId {
							workspace_users[wsId][userIndex] = workspace_users[wsId][len(workspace_users[wsId])-1]
							workspace_users[wsId] = workspace_users[wsId][0 : len(workspace_users[wsId])-1]
							context.IndentedJSON(http.StatusOK, workspace_users[wsId])
							return
						}
					}

				} else {
					// didn't find user with memID in workspace
					context.JSON(http.StatusBadRequest, gin.H{"error": "User not found in this workspace"})
					return
				}
			}
		}
	}
}

// returns all users in the workspace with specified id
func workSpaceMembers(context *gin.Context) {
	// workspace id
	id := context.Param("id")

	tokenStatus, err := utils.ExtractTokenID(context)
	token := utils.ExtractToken(context)
	loggedInUser := ""

	// verify auth
	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, u := range users {
		if u.AccessToken == token {
			// get ID from AccessToken
			loggedInUser = u.ID
		}
	}

	for _, t := range workspaces {
		if t.ID == id {
			// found a workspace with this id
			var wsMembers = []user{}
			for _, u := range users {
				if contains(workspace_users[id], u.ID) && u.ID != loggedInUser {
					wsMembers = append(wsMembers, u)
				}
			}
			context.IndentedJSON(http.StatusOK, wsMembers)
			return
		}
	}

	context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no workspace found with this id"})
	return
}

// creates a channel in the workspace id
// returns all channels in that workspace
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
				newChannel.Messages = 0

				// increment channel count of parent workspace
				count := workspaces[i].Channels
				workspaces[i].Channels = count + 1

				channels = append(channels, newChannel)

				var wsChannels = []channel{}
				for _, t := range channels {
					if contains(workspace_channels[id], t.ID) {
						wsChannels = append(wsChannels, t)
					}
				}
				context.IndentedJSON(http.StatusOK, wsChannels)
				return
			}
		}
	}
}

// deletes channel with id
// returns all channels in that workspace
func deleteChannel(context *gin.Context) {
	tokenStatus, _ := utils.ExtractTokenID(context)
	token := utils.ExtractToken(context)

	if tokenStatus == 0 {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": bad_rq_message})
		return
	}
	id := context.Param("id") // channel id

	// find the workspace
	var channel_workspace workspace
	for workspaceIndex, workspace := range workspaces {
		for _, channelId := range workspace_channels[workspace.ID] {
			if channelId == id {
				channel_workspace = workspace
				// decerement channels count of parent workspace
				count := workspaces[workspaceIndex].Channels
				workspaces[workspaceIndex].Channels = count - 1
			}
		}
	}
	// find index in channel list
	deleteIndex := -1
	for i, c := range channels {
		if c.ID == id {
			deleteIndex = i
			// check owner is logged in user
			for _, u := range users {
				if u.AccessToken == token {
					// get ID from AccessToken, if not owner, return err
					if u.ID != channel_workspace.Owner {
						context.IndentedJSON(http.StatusNotFound, gin.H{"message": bad_rq_message})
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
		for _, mID := range channel_messages[id] {
			for messageIndex, message := range messages {
				if message.ID == mID {
					// delete message
					messages[messageIndex] = messages[len(messages)-1]
					messages = messages[0 : len(messages)-1]
				}
			}
		}
		// delete channel
		channels[deleteIndex] = channels[len(channels)-1]
		channels = channels[0 : len(channels)-1]
	}

	delete(channel_messages, id)

	// return channels remaining in workspace
	var wsChannels = []channel{}
	for _, c := range channels {
		if contains(workspace_channels[channel_workspace.ID], c.ID) {
			wsChannels = append(wsChannels, c)
		}
	}
	context.IndentedJSON(http.StatusOK, wsChannels)
	return
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
				newMessage.Member = t.ID
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
				newMessage.Posted = time.Now().UTC().Format(time.RFC3339) //time.Now()
				// increment message count of parent channel
				count := channels[i].Messages
				channels[i].Messages = count + 1

				messages = append(messages, newMessage)

				var chMessages = []message{}
				for _, m := range messages {
					if contains(channel_messages[id], m.ID) {
						chMessages = append(chMessages, m)
					}
				}
				context.IndentedJSON(http.StatusOK, chMessages)
				return
			}
		}
		context.JSON(http.StatusBadRequest, gin.H{"error": "Channel does not exist"})
		return
	}
}

// deletes message with specified ID
func deleteMessage(context *gin.Context) {
	tokenStatus, _ := utils.ExtractTokenID(context)
	token := utils.ExtractToken(context)

	if tokenStatus == 0 {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": bad_rq_message})
		return
	}
	id := context.Param("id") // message id

	// find the channel
	var message_channel channel
	for channelIndex, channel := range channels {
		for messageIndex, messageId := range channel_messages[channel.ID] {
			if messageId == id {
				message_channel = channel
				channel_messages[message_channel.ID][messageIndex] = channel_messages[message_channel.ID][len(channel_messages[message_channel.ID])-1]
				channel_messages[message_channel.ID] = channel_messages[message_channel.ID][0 : len(channel_messages[message_channel.ID])-1]
				// decement messages count in parent channel
				count := channels[channelIndex].Messages
				channels[channelIndex].Messages = count - 1
			}
		}
	}

	// find the workspace
	var channel_workspace workspace
	for _, workspace := range workspaces {
		for _, channelId := range workspace_channels[workspace.ID] {
			if channelId == message_channel.ID {
				channel_workspace = workspace
			}
		}
	}
	// find index in message list
	deleteIndex := -1
	for i, m := range messages {
		if m.ID == id {
			deleteIndex = i
			// check owner is logged in user
			for _, u := range users {
				if u.AccessToken == token {
					// get ID from AccessToken, if not owner, return err
					if u.ID != channel_workspace.Owner && u.ID != m.Member {
						context.IndentedJSON(http.StatusNotFound, gin.H{"message": bad_rq_message})
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
		for messageIndex, message := range messages {
			if message.ID == id {
				// delete message
				messages[messageIndex] = messages[len(messages)-1]
				messages = messages[0 : len(messages)-1]
			}
		}
	}

	// return channels remaining in workspace
	var chMessages = []message{}
	for _, m := range messages {
		if contains(channel_messages[message_channel.ID], m.ID) {
			chMessages = append(chMessages, m)
		}
	}
	context.IndentedJSON(http.StatusOK, chMessages)
	return
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
		for _, c := range channels {
			if c.ID == id {
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
					context.IndentedJSON(http.StatusBadRequest, gin.H{"message": bad_rq_message})
					return
				}
			}
		}
	}

	context.IndentedJSON(http.StatusBadRequest, gin.H{"message": bad_rq_message})
	return
}

func main() {
	// initializing World Wide Workspace... feel free to delete this :)
	workspace_channels["00000000-0000-0000-0000-000000000000"] = []string{"00000000-0000-0000-0000-000000000000"}
	workspace_users["00000000-0000-0000-0000-000000000000"] = []string{}
	channel_messages["00000000-0000-0000-0000-000000000000"] = []string{"00000000-0000-0000-0000-000000000000"}
	for _, u := range users {
		workspace_users["00000000-0000-0000-0000-000000000000"] = append(workspace_users["00000000-0000-0000-0000-000000000000"], u.ID)
	}

	router := gin.Default()
	router.POST("/login", login)
	router.POST("/newuser", createUser)
	router.GET("/member", getUsers)

	router.GET("/workspace", getWorkSpaces)
	router.POST("/workspace", createWorkSpace)
	router.DELETE("/workspace/:id", deleteWorkSpace)

	router.GET("/workspace/:id/member", workSpaceMembers)
	router.POST("/workspace/:wsId", addWorkSpaceMember)
	router.DELETE("/workspace/remove/:wsId", deleteWorkSpaceMember)

	router.GET("/workspace/channel/:id", getChannels)
	router.POST("/workspace/channel/:id", createChannel)
	router.DELETE("/channel/:id", deleteChannel)

	router.GET("/channel/message/:id", getMessages)
	router.POST("/channel/message/:id", createMessage)
	router.DELETE("/message/:id", deleteMessage)

	router.Run("localhost:9090")
}
