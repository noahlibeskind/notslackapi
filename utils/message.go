package utils

import (
	"net/http"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/noahlibeskind/NotSlackAPI/data"
)

// adds a member with id memId to workspace with id wsId
// returns all messages in the channel
func CreateMessage(context *gin.Context) {
	// channel id
	id := context.Param("id")

	var newMessage data.Message
	// JSON should include message content only
	err := context.BindJSON(&newMessage)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tokenStatus, err := ExtractTokenID(context)
	token := ExtractToken(context)

	// verify auth
	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		for _, t := range data.Users {
			if t.AccessToken == token {
				// get ID from AccessToken
				newMessage.Member = t.ID
			}
		}
		for i, t := range data.Channels {
			if t.ID == id {
				// found a channel with this id, add message to it
				newUUID, err := exec.Command("uuidgen").Output()
				if err != nil {
					return
				}
				newMessage.ID = string(newUUID)[0 : len(newUUID)-1]

				// add id to map
				_, ok := data.Channel_messages[id]
				if ok {
					data.Channel_messages[id] = append(data.Channel_messages[id], newMessage.ID)
				} else {
					data.Channel_messages[id] = []string{newMessage.ID}
				}
				// get current time
				newMessage.Posted = time.Now().UTC().Format(time.RFC3339) //time.Now()
				// increment message count of parent channel
				count := data.Channels[i].Messages
				data.Channels[i].Messages = count + 1

				data.Messages = append(data.Messages, newMessage)

				var chMessages = []data.Message{}
				for _, m := range data.Messages {
					if contains(data.Channel_messages[id], m.ID) {
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
func DeleteMessage(context *gin.Context) {
	tokenStatus, _ := ExtractTokenID(context)
	token := ExtractToken(context)

	if tokenStatus == 0 {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{"message": data.Unauthorized_message})
		return
	}
	id := context.Param("id") // message id

	// find the channel
	var message_channel data.Channel
	for channelIndex, channel := range data.Channels {
		for messageIndex, messageId := range data.Channel_messages[channel.ID] {
			if messageId == id {
				message_channel = channel
				data.Channel_messages[message_channel.ID][messageIndex] = data.Channel_messages[message_channel.ID][len(data.Channel_messages[message_channel.ID])-1]
				data.Channel_messages[message_channel.ID] = data.Channel_messages[message_channel.ID][0 : len(data.Channel_messages[message_channel.ID])-1]
				// decement messages count in parent channel
				count := data.Channels[channelIndex].Messages
				data.Channels[channelIndex].Messages = count - 1
			}
		}
	}

	// find the workspace
	var channel_workspace data.Workspace
	for _, workspace := range data.Workspaces {
		for _, channelId := range data.Workspace_channels[workspace.ID] {
			if channelId == message_channel.ID {
				channel_workspace = workspace
			}
		}
	}
	// find index in message list
	deleteIndex := -1
	for i, m := range data.Messages {
		if m.ID == id {
			deleteIndex = i
			// check owner is logged in user
			for _, u := range data.Users {
				if u.AccessToken == token {
					// get ID from AccessToken, if not owner, return err
					if u.ID != channel_workspace.Owner && u.ID != m.Member {
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
		for messageIndex, message := range data.Messages {
			if message.ID == id {
				// delete message
				data.Messages[messageIndex] = data.Messages[len(data.Messages)-1]
				data.Messages = data.Messages[0 : len(data.Messages)-1]
			}
		}
	}

	// return channels remaining in workspace
	var chMessages = []data.Message{}
	for _, m := range data.Messages {
		if contains(data.Channel_messages[message_channel.ID], m.ID) {
			chMessages = append(chMessages, m)
		}
	}
	context.IndentedJSON(http.StatusOK, chMessages)
}

// returns all messages in the specified channel
func GetMessages(context *gin.Context) {
	// channel id
	id := context.Param("id")
	tokenStatus, err := ExtractTokenID(context)

	// verify auth
	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		for _, c := range data.Channels {
			if c.ID == id {
				// found a channel with this id, add message to it
				var chMessages = []data.Message{}
				for _, m := range data.Messages {
					if contains(data.Channel_messages[id], m.ID) {
						chMessages = append(chMessages, m)
					}
				}
				context.IndentedJSON(http.StatusOK, chMessages)
				return
			}
		}
	}
	context.JSON(http.StatusBadRequest, gin.H{"message": "channel not found"})
}
