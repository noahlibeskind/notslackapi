package data

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	AccessToken string `json:"accessToken"`
}

type Workspace struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Channels int    `json:"channels"`
	Owner    string `json:"owner"`
}

type Channel struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Messages int    `json:"messages"`
}

type Message struct {
	ID      string `json:"id"`
	Member  string `json:"member"`
	Posted  string `json:"posted"`
	Content string `json:"content"`
}

var Users = []User{
	{ID: "00000000-0000-0000-0000-000000000000", Name: "Noah Libeskind", Email: "noah@ucsc.edu", Password: "pass", AccessToken: ""},
}

var Workspaces = []Workspace{
	{ID: "00000000-0000-0000-0000-000000000000", Name: "(WWW) World Wide Workspace", Channels: 1, Owner: "00000000-0000-0000-0000-000000000000"},
}

var Channels = []Channel{
	{ID: "00000000-0000-0000-0000-000000000000", Name: "World Chat Channel", Messages: 1},
}

var Messages = []Message{
	{ID: "00000000-0000-0000-0000-000000000000", Member: "00000000-0000-0000-0000-000000000000", Posted: "2023-01-02T00:01:01Z", Content: "Hello! Welcome to the world chat channel!"},
}

// maps workspace IDs to IDs of users in that workspace
var Workspace_users = map[string][]string{}

// maps workspace IDs to IDs of channels in that workspace
var Workspace_channels = map[string][]string{}

// maps channel IDs to IDs of messages in that channel
var Channel_messages = map[string][]string{}

var Unauthorized_message = "Invalid Credentials"
