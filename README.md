NotSlack API 

CSE 118 Students beware: the some of the URLs for your calls to the API are changed! Here are the updates:

/newuser (POST) - allows you to create an account on the app with name, email, and password in a JSON object. A logged in instance of the new user is returned. 
*** Access Tokens are only valid for one hour, then the user must login again.

Example:
Body:
{
    "name": "John Doe",
    "email": "jdoe@gmail.com",
    "password": "fluffykitten123"
}
Response:
{
    "id": "00000000-0000-0000-0000-000000000000",
    "name": "John Doe",
    "email": "jdoe@gmail.com",
    "password": "fluffykitten123",
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NzQ0NDI4OTUsInVzZXJfaWQiOjY5fQ.DKt53365JXo2wXb6ukU8_VBN9dWTx44BOllNLIq2QXQ"
}


/login (POST) - pass with JSON object email and password and the user object with an accessToken will be returned upon successful login.
Example:
Body:
{
    "email": "jdoe@gmail.com",
    "password": "fluffykitten123"
}
Response:
{
    "id": "00000000-0000-0000-0000-000000000000",
    "name": "John Doe",
    "email": "jdoe@gmail.com",
    "password": "fluffykitten123",
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NzQ0NDI4OTUsInVzZXJfaWQiOjY5fQ.DKt53365JXo2wXb6ukU8_VBN9dWTx44BOllNLIq2QXQ"
}

/member (GET) - returns all members on the app as a JSON encoded list of the following format:
Example:
Body: None
Response:
[{
    "id": "00000000-0000-0000-0000-000000000000",
    "name": "John Doe",
    "email": "jdoe@gmail.com",
    "password": "fluffykitten123",
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NzQ0NDI4OTUsInVzZXJfaWQiOjY5fQ.DKt53365JXo2wXb6ukU8_VBN9dWTx44BOllNLIq2QXQ"
}, ...]

/workspace (GET) - returns all workspaces.
Example:
Body: None
Response:
[{
    "id": "00000000-0000-0000-0000-000000000000",
    "name": "(WWW) World Wide Workspace",
    "channels": 1,
    "owner": "00000000-0000-0000-0000-000000000000"
}, ...]


/workspace (POST) - creates a new workspace with a specified name. Returns that workspace object.
Example:
Body: 
{
    "name" : "Private Workspace"
}
Response:
{
    "id": "00000000-0000-0000-0000-000000000001",
    "name": "Private Workspace",
    "channels": 0,
    "owner": "00000000-0000-0000-0000-000000000000"
}

/workspace/{id} (DELETE) - deletes workspace with specified id. Returns list of remaining workspaces.
Example:
Body: None
Response:
[{
    "id": "00000000-0000-0000-0000-000000000000",
    "name": "(WWW) World Wide Workspace",
    "channels": 1,
    "owner": "00000000-0000-0000-0000-000000000000"
}, ...]

/workspace/channel/{id} (GET) - returns all the channels in the workspace with the specified id.
Example:
Body: None
Response: 
[{
    "id": "00000000-0000-0000-0000-000000000000",
    "name": "First channel",
    "messages": 0
}, ...]

/workspace/channel/{id} (POST) - creates a channel in the workspace with the specified id. Include the name of the channel to be created, a list of the channel objects in that workspace will be returned.
Example:
Body:
{
    "name" : "First channel"
}
Response:
[{
    "id": "00000000-0000-0000-0000-000000000000",
    "name": "First channel",
    "messages": 0
}, ...]

/channel/{id} (DELETE) - deletes the channel with the specified id. Returns the remaining channels in that workspace.
Example:
Body: None
Response: 
[{
    "id": "00000000-0000-0000-0000-000000000001",
    "name": "Other channel",
    "messages": 0
}, ...]


/channel/message/{id} (GET) - returns all the messages in the channel with the specified id

/channel/message/{id} (POST) - creates a message in the channel with the specified id. Include the message content in JSON, returns a list of the message objects in that channel.
Example:
Body:
{
    "content" : "Hello this is John!"
}
Response:
[{
    "id": "00000000-0000-0000-0000-000000000000",
    "member": "00000000-0000-0000-0000-000000000000",
    "posted": "2023-01-23T19:37:38Z",
    "content": "Hello this is John!"
}, ...]


/message/{id} (DELETE) - deletes a message with the specified id. 
Example:
Body: None
Response:
[{
    "id": "00000000-0000-0000-0000-000000000000",
    "member": "00000000-0000-0000-0000-000000000000",
    "posted": "2023-01-23T19:37:45Z",
    "content": "Another message"
}, ...]