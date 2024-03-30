// NotSlack API
// created by Noah Libeskind

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/noahlibeskind/NotSlackAPI/data"
	"github.com/noahlibeskind/NotSlackAPI/utils"
)

// CORSMiddleware is a custom middleware function that wraps the cors handler in a gin.HandlerFunc
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	// initializing World Wide Workspace:
	data.Workspace_users["00000000-0000-0000-0000-000000000000"] = []string{}
	data.Workspace_channels["00000000-0000-0000-0000-000000000000"] = []string{"00000000-0000-0000-0000-000000000000"}
	data.Channel_messages["00000000-0000-0000-0000-000000000000"] = []string{"00000000-0000-0000-0000-000000000000"}

	router := gin.Default()

	router.Use(CORSMiddleware())

	router.POST("/login", utils.Login)
	router.POST("/newuser", utils.CreateUser)
	router.DELETE("/users/:uid", utils.DeleteUser)
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

	router.Run("0.0.0.0:8080")
}
