// NotSlack API
// created by Noah Libeskind

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/noahlibeskind/NotSlackAPI/data"
	"github.com/noahlibeskind/NotSlackAPI/utils"
)

var mySigningKey = []byte("notslackisnotsecure")

func main() {
	// initializing World Wide Workspace... feel free to delete this :)
	data.Workspace_users["00000000-0000-0000-0000-000000000000"] = []string{}
	data.Workspace_channels["00000000-0000-0000-0000-000000000000"] = []string{"00000000-0000-0000-0000-000000000000"}
	data.Channel_messages["00000000-0000-0000-0000-000000000000"] = []string{"00000000-0000-0000-0000-000000000000"}
	// for _, u := range data.Users {
	// 	workspace_users["00000000-0000-0000-0000-000000000000"] = append(workspace_users["00000000-0000-0000-0000-000000000000"], u.ID)
	// }

	router := gin.Default()
	router.POST("/login", utils.Login)
	router.POST("/newuser", utils.CreateUser)
	router.GET("/member", utils.GetUsers)

	router.GET("/workspace", utils.GetWorkSpaces)
	router.POST("/workspace", utils.CreateWorkSpace)
	router.DELETE("/workspace/:id", utils.DeleteWorkSpace)

	router.GET("/workspace/:id/member", utils.WorkSpaceMembers)
	router.POST("/workspace/:wsId", utils.AddWorkSpaceMember)
	router.DELETE("/workspace/remove/:wsId", utils.DeleteWorkSpaceMember)

	router.GET("/workspace/channel/:id", utils.GetChannels)
	router.POST("/workspace/channel/:id", utils.CreateChannel)
	router.DELETE("/channel/:id", utils.DeleteChannel)

	router.GET("/channel/message/:id", utils.GetMessages)
	router.POST("/channel/message/:id", utils.CreateMessage)
	router.DELETE("/message/:id", utils.DeleteMessage)

	router.Run("localhost:9090")
}
