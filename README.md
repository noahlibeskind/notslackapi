
# NotSlack API

This API supports the backend of a workplace messaging app that closely resembles an already existing one used by many professionals, hence the name "Not Slack".

This API is built in Golang and utilizes helpful packages like Gin from Gin-Gonic and Dave Grijalva's jwt-go package to support basic security implemented in this API.

### Below are the various URIs of this API and descriptions of their functionalities

---

## /newuser (POST)

Allows you to create an account on the app with name, email, and password in a JSON object. A logged in instance of the new user is returned.  
*Access Tokens are only valid for one hour, then the user must login again*

**Example:**

**Body:**
```json
{
    "name": "John Doe",
    "email": "jdoe@gmail.com",
    "password": "fluffykitten123"
}
```

**Response:**
```json
{
    "id": "00000000-0000-0000-0000-000000000000",
    "name": "John Doe",
    "email": "jdoe@gmail.com",
    "password": "fluffykitten123",
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NzQ0NDI4OTUsInVzZXJfaWQiOjY5fQ.DKt53365JXo2wXb6ukU8_VBN9dWTx44BOllNLIq2QXQ"
}
```

---

## /login (POST)

Pass with JSON object email and password, and the user object with an accessToken will be returned upon successful login.

**Example:**

**Body:**
```json
{
    "email": "jdoe@gmail.com",
    "password": "fluffykitten123"
}
```

**Response:**
```json
{
    "id": "00000000-0000-0000-0000-000000000000",
    "name": "John Doe",
    "email": "jdoe@gmail.com",
    "password": "fluffykitten123",
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NzQ0NDI4OTUsInVzZXJfaWQiOjY5fQ.DKt53365JXo2wXb6ukU8_VBN9dWTx44BOllNLIq2QXQ"
}
```

---

## /member (GET)

Returns all members on the app as a JSON encoded list of the following format.

**Example:**

**Body:** None

**Response:**
```json
[{
    "id": "00000000-0000-0000-0000-000000000000",
    "name": "John Doe",
    "email": "jdoe@gmail.com",
    "password": "fluffykitten123",
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NzQ0NDI4OTUsInVzZXJfaWQiOjY5fQ.DKt53365JXo2wXb6ukU8_VBN9dWTx44BOllNLIq2QXQ"
}]
```

---

## /workspace (GET)

Returns all workspaces.

**Example:**

**Body:** None

**Response:**
```json
[{
    "id": "00000000-0000-0000-0000-000000000000",
    "name": "(WWW) World Wide Workspace",
    "channels": 1,
    "owner": "00000000-0000-0000-0000-000000000000"
}]
```

---

## /workspace (POST)

Creates a new workspace with a specified name. Returns that workspace object.

**Example:**

**Body:**
```json
{
    "name" : "Private Workspace"
}
```

**Response:**
```json
{
    "id": "00000000-0000-0000-0000-000000000001",
    "name": "Private Workspace",
    "channels": 0,
    "owner": "00000000-0000-0000-0000-000000000000"
}
```

---

## /workspace/{id} (DELETE)

Deletes workspace with the specified id. Returns a list of remaining workspaces.

**Example:**

**Body:** None

**Response:**
```json
[{
    "id": "00000000-0000-0000-0000-000000000000",
    "name": "(WWW) World Wide Workspace",
    "channels": 1,
    "owner": "00000000-0000-0000-0000-000000000000"
}]
```

---

## /workspace/channel/{id} (GET)

Returns all the channels in the workspace with the specified id.

**Example:**

**Body:** None

**Response:**
```json
[{
    "id": "00000000-0000-0000-0000-000000000000",
    "name": "First channel",
    "messages": 0
}]
```

---

## /workspace/channel/{id} (POST)

Creates a channel in the workspace with the specified id. Includes the name of the channel to be created, a list of the channel objects in that workspace will be returned.

**Example:**

**Body:**
```json
{
    "name" : "First channel"
}
```

**Response:**
```json
[{
    "id": "00000000-0000-0000-0000-000000000000",
    "name": "First channel",
    "messages": 0
}]
```

---

## /channel/{id} (DELETE)

Deletes the channel with the specified id. Returns the remaining channels in that workspace.

**Example:**

**Body:** None

**Response:**
```json
[{
    "id": "00000000-0000-0000-0000-000000000001",
    "name": "Other channel",
    "messages": 0
}]
```

---

## /channel/message/{id} (GET)

Returns all the messages in the channel with the specified id.

---

## /channel/message/{id} (POST)

Creates a message in the channel with the specified id. Include the message content in JSON, returns a list of the message objects in that channel.

**Example:**

**Body:**
```json
{
    "content" : "Hello this is John!"
}
```

**Response:**
```json
[{
    "id": "00000000-0000-0000-0000-000000000000",
    "member": "00000000-0000-0000-0000-000000000000",
    "posted": "2023-01-23T19:37:38Z",
    "content": "Hello this is John!"
}]
```

---

## /message/{id} (DELETE)

Deletes a message with the specified id.

**Example:**

**Body:** None

**Response:**
```json
[{
    "id": "00000000-0000-0000-0000-000000000000",
    "member": "00000000-0000-0000-0000-000000000000",
    "posted": "2023-01-23T19:37:45Z",
    "content": "Another message"
}]
```

---
