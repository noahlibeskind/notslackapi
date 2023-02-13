package data

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
	{ID: "00000000-0000-0000-0000-000000000000", Member: "00000000-0000-0000-0000-000000000000", Posted: "2023-01-02T00:01:01ZZZ", Content: "Hello! Welcome to the world chat channel!"},
}
