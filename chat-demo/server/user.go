package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// Create a new user instance
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
		server: server,
	}

	// Start a goroutine to listen for messages on the user's channel
	go user.ListenMessage()

	return user
}

// Listen for messages on the user's channel and send them to the user's connection
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}

// User goes online
func (u *User) Online() { 
	// Implementation for user going online can be added here
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	// broadcast user online message
	u.server.BroadCast(u, "is online")
	
}

// User goes offline
func (u *User) Offline() {
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	// broadcast user offline message
	u.server.BroadCast(u, "is offline")
	
}	

// send message to user
func (u *User) SendPrivateMessage(msg string) {
	u.conn.Write([]byte(msg))
}

// User sends a message
func (u *User) SendMessage(msg string) {
	if msg == "who" {
		// query online users
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := fmt.Sprintf("[%s]: online\n", user.Name)
			u.SendPrivateMessage(onlineMsg)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// rename user
		newName := msg[7:]
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendPrivateMessage("The name is already in use.\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.Name = newName
			u.server.OnlineMap[u.Name] = u
			u.server.mapLock.Unlock()
			u.SendPrivateMessage("You have successfully changed your name to " + u.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// private message format: to|username|message
		
		// get target username
		remoteUserName := strings.Split(msg, "|")[1]
		if remoteUserName == "" {
			u.SendPrivateMessage("The format of private message is incorrect. Use to|username|message.\n")
			return
		}

		remoteUser, ok := u.server.OnlineMap[remoteUserName]
		if !ok {
			u.SendPrivateMessage("The user does not exist.\n")
			return
		}

		// get message content
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.SendPrivateMessage("The format of private message is incorrect. Use to|username|message.\n")
			return
		}

		remoteUser.SendPrivateMessage("[" + u.Name + "] to you: " + content + "\n")
		
	} else {
		u.server.BroadCast(u, msg)
	}
}