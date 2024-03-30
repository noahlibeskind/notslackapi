package utils

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/noahlibeskind/NotSlackAPI/data"
)

// derived from https://play.golang.org/p/Qg_uv_inCek
// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// get all workspaces
func GetWorkSpaces(context *gin.Context) {

	tokenStatus, _ := ExtractTokenID(context)
	token := ExtractToken(context)
	loggedInUser := ""
	// encoded JSON should only include name
	if tokenStatus == 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": data.Bad_rq_message})
		return
	} else {
		for _, u := range data.Users {
			if u.AccessToken == token {
				// get ID from AccessToken
				loggedInUser = u.ID
			}
		}
	}
	var userWorkspaces = []data.Workspace{}
	for _, w := range data.Workspaces {
		if contains(data.Workspace_users[w.ID], loggedInUser) || w.Owner == loggedInUser {
			userWorkspaces = append(userWorkspaces, w)
		}
	}
	context.IndentedJSON(http.StatusOK, userWorkspaces)
}

// creates a new workspace with logged in user as the owner
// returns the created workspace
func CreateWorkSpace(context *gin.Context) {
	var newWS data.Workspace

	tokenStatus, _ := ExtractTokenID(context)
	token := ExtractToken(context)

	err := context.BindJSON(&newWS)
	// encoded JSON should only include name
	if err != nil || tokenStatus == 0 {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": err})
		return
	} else {
		for _, u := range data.Users {
			if u.AccessToken == token {
				// get ID from AccessToken
				newWS.Owner = u.ID
			}
		}
	}
	// not current logged in accessToken
	if newWS.Owner == "" {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": data.Bad_rq_message})
		return
	}
	newUUID, err := exec.Command("uuidgen").Output()
	if err != nil {
		return
	}
	newWS.ID = string(newUUID)[0 : len(newUUID)-1]
	newWS.Channels = 0
	data.Workspaces = append(data.Workspaces, newWS)
	context.IndentedJSON(http.StatusOK, newWS)
	return
}

// creates a new workspace with logged in user as the owner
// returns the created workspace
func DeleteWorkSpace(context *gin.Context) {
	tokenStatus, _ := ExtractTokenID(context)
	token := ExtractToken(context)
	if tokenStatus == 0 {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": data.Bad_rq_message})
		return
	}
	id := context.Param("id")
	deleteIndex := -1
	for i, w := range data.Workspaces {
		if w.ID == id {
			deleteIndex = i
			// check owner is logged in user
			for _, u := range data.Users {
				if u.AccessToken == token {
					// get ID from AccessToken, if not owner, return err
					if u.ID != w.Owner {
						context.IndentedJSON(http.StatusNotFound, gin.H{"message": data.Bad_rq_message})
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
		for _, chID := range data.Workspace_channels[id] {
			for channelIndex, channel := range data.Channels {
				if channel.ID == chID {
					// delete messages in channel
					for _, mID := range data.Channel_messages[chID] {
						for messageIndex, message := range data.Messages {
							if message.ID == mID {
								// delete message
								data.Messages[messageIndex] = data.Messages[len(data.Messages)-1]
								data.Messages = data.Messages[0 : len(data.Messages)-1]
							}
						}
					}
					// delete channel
					data.Channels[channelIndex] = data.Channels[len(data.Channels)-1]
					data.Channels = data.Channels[0 : len(data.Channels)-1]
				}
			}
		}
		// delete workspace
		data.Workspaces[deleteIndex] = data.Workspaces[len(data.Workspaces)-1]
		data.Workspaces = data.Workspaces[0 : len(data.Workspaces)-1]
		delete(data.Workspace_channels, id)
		delete(data.Workspace_users, id)
		context.IndentedJSON(http.StatusOK, data.Workspaces)
	}
	return
}

// adds a member with id memId to workspace with id wsId
// returns all members in that workspace
func AddWorkSpaceMember(context *gin.Context) {
	// workspace id
	wsId := context.Param("wsId")
	memId := context.Request.URL.Query().Get("mid")
	tokenStatus, err := ExtractTokenID(context)
	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		for _, t := range data.Workspaces {
			if t.ID == wsId {
				// found a workspace with this id, add member to it
				_, ok := data.Workspace_users[wsId]
				if ok {
					data.Workspace_users[wsId] = append(data.Workspace_users[wsId], memId)
				} else {
					data.Workspace_users[wsId] = []string{memId}
				}

				var wsMembers = []data.User{}
				for _, t := range data.Users {
					if contains(data.Workspace_users[wsId], t.ID) {
						wsMembers = append(wsMembers, t)
					}
				}
				context.IndentedJSON(http.StatusOK, wsMembers)
				return
			}
		}
		return
	}
}

// deletes the member with id memId to workspace with id wsId
// returns all members in that workspace
func DeleteWorkSpaceMember(context *gin.Context) {
	// workspace id
	wsId := context.Param("wsId")
	memId := context.Request.URL.Query().Get("mid")
	tokenStatus, err := ExtractTokenID(context)
	token := ExtractToken(context)
	// verify auth
	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deleteIndex := -1
	for i, w := range data.Workspaces {
		if w.ID == wsId {
			deleteIndex = i
			// check owner is logged in user
			for _, u := range data.Users {
				if u.AccessToken == token {
					// get ID from AccessToken, if user is not owner, return err
					if u.ID != w.Owner {
						context.IndentedJSON(http.StatusNotFound, gin.H{"message": data.Bad_rq_message})
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
		for _, t := range data.Workspaces {
			if t.ID == wsId {
				// found a data.Workspace with this id
				_, ok := data.Workspace_users[wsId]
				if ok {
					// found user with memID in data.Workspace
					for userIndex, user := range data.Workspace_users[wsId] {
						if user == memId {
							data.Workspace_users[wsId][userIndex] = data.Workspace_users[wsId][len(data.Workspace_users[wsId])-1]
							data.Workspace_users[wsId] = data.Workspace_users[wsId][0 : len(data.Workspace_users[wsId])-1]
							context.IndentedJSON(http.StatusOK, data.Workspace_users[wsId])
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
		return
	}
}

// returns all users in the workspace with specified id
func WorkSpaceMembers(context *gin.Context) {
	// workspace id
	id := context.Param("id")
	tokenStatus, err := ExtractTokenID(context)
	token := ExtractToken(context)
	loggedInUser := ""
	// verify auth
	if err != nil || tokenStatus == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, u := range data.Users {
		if u.AccessToken == token {
			// get ID from AccessToken
			loggedInUser = u.ID
		}
	}
	for _, t := range data.Workspaces {
		if t.ID == id {
			// found a workspace with this id
			var wsMembers = []data.User{}
			for _, u := range data.Users {
				if contains(data.Workspace_users[id], u.ID) && u.ID != loggedInUser {
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
