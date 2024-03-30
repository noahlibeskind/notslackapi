package utils

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/noahlibeskind/NotSlackAPI/data"
)

// creates a channel in the workspace id
// returns all channels in that workspace
func CreateChannel(context *gin.Context) {
	// workspace id
	id := context.Param("id")

	var newChannel data.Channel
	// JSON should include channel name only
	err := context.BindJSON(&newChannel)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tokenStatus, err := ExtractTokenID(context)
	// verify auth
	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		for i, t := range data.Workspaces {
			if t.ID == id {
				// found a workspace with this id, add channel to it
				newUUID, err := exec.Command("uuidgen").Output()
				if err != nil {
					return
				}
				newChannel.ID = string(newUUID)[0 : len(newUUID)-1]

				// add id to map
				_, ok := data.Workspace_channels[id]
				if ok {
					data.Workspace_channels[id] = append(data.Workspace_channels[id], newChannel.ID)
				} else {
					data.Workspace_channels[id] = []string{newChannel.ID}
				}
				newChannel.Messages = 0

				// increment channel count of parent workspace
				count := data.Workspaces[i].Channels
				data.Workspaces[i].Channels = count + 1

				data.Channels = append(data.Channels, newChannel)

				var wsChannels = []data.Channel{}
				for _, t := range data.Channels {
					if contains(data.Workspace_channels[id], t.ID) {
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
func DeleteChannel(context *gin.Context) {
	tokenStatus, _ := ExtractTokenID(context)
	token := ExtractToken(context)

	if tokenStatus == 0 {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{"message": data.Unauthorized_message})
		return
	}
	id := context.Param("id") // channel id

	// find the workspace
	var channel_workspace data.Workspace
	for workspaceIndex, workspace := range data.Workspaces {
		for _, channelId := range data.Workspace_channels[workspace.ID] {
			if channelId == id {
				channel_workspace = workspace
				// decerement channels count of parent workspace
				count := data.Workspaces[workspaceIndex].Channels
				data.Workspaces[workspaceIndex].Channels = count - 1
			}
		}
	}
	// find index in channel list
	deleteIndex := -1
	for i, c := range data.Channels {
		if c.ID == id {
			deleteIndex = i
			// check owner is logged in user
			for _, u := range data.Users {
				if u.AccessToken == token {
					// get ID from AccessToken, if not owner, return err
					if u.ID != channel_workspace.Owner {
						context.IndentedJSON(http.StatusUnauthorized, gin.H{"message": data.Unauthorized_message})
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
		for _, mID := range data.Channel_messages[id] {
			for messageIndex, message := range data.Messages {
				if message.ID == mID {
					// delete message
					data.Messages[messageIndex] = data.Messages[len(data.Messages)-1]
					data.Messages = data.Messages[0 : len(data.Messages)-1]
				}
			}
		}
		// delete channel
		data.Channels[deleteIndex] = data.Channels[len(data.Channels)-1]
		data.Channels = data.Channels[0 : len(data.Channels)-1]
	}

	delete(data.Channel_messages, id)

	// return channels remaining in workspace
	var wsChannels = []data.Channel{}
	for _, c := range data.Channels {
		if contains(data.Workspace_channels[channel_workspace.ID], c.ID) {
			wsChannels = append(wsChannels, c)
		}
	}
	context.IndentedJSON(http.StatusOK, wsChannels)
}

// returns all channels in the specified workspace
func GetChannels(context *gin.Context) {
	tokenStatus, err := ExtractTokenID(context)

	// workspace id
	id := context.Param("id")

	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var wsChannels = []data.Channel{}
	for _, t := range data.Channels {
		if contains(data.Workspace_channels[id], t.ID) {
			wsChannels = append(wsChannels, t)
		}
	}
	context.IndentedJSON(http.StatusOK, wsChannels)
}
