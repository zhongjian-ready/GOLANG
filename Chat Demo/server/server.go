package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int
	// online user map
	OnlineMap map[string]*User
	// message broadcasting channel
	Message chan string
	// map lock
	mapLock sync.RWMutex

}

// create a server instance
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:   ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

// listen message channel and broadcast to all online users
func (s *Server) ListenMessager() {
	for {
		msg := <-s.Message
		// send message to all online users
		s.mapLock.Lock()
		for _, user := range s.OnlineMap {
			user.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// broadcast message to all online users
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s]: %s", user.Name, msg)
	s.Message <- sendMsg
}
 
func (s *Server) Handler(conn net.Conn) {
	// ... the real-time chat logic would go here
	// user online add user to online map
	user := NewUser(conn, s)

	// user goes online
	user.Online()

	// lister user liveness
	isLive := make(chan bool)

	// receive user messages
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil {
				fmt.Println("conn.Read error:", err)
				return
			}
			// broadcast user message
			msg := string(buf[:n-1])
			
			user.SendMessage(msg)

			// user is alive
			isLive <- true
		}
	}()

	for {
		select {
		case <-isLive:
			// do nothing, just continue and reset the timer
		case <-time.After(30 * 60 * time.Second): // 30 minutes
			// timeout kick user offline
			user.SendPrivateMessage("You have been kicked offline due to inactivity.\n")
			// destroy user resources
			close(user.C)
			// close the connection
			conn.Close()
			// user goes offline
			user.Offline()
			// exit handler
			return // runtime.Goexit()
	}
	}

	// fmt.Println("New connection established:", conn.RemoteAddr().String())
	// defer conn.Close()
}

// Start the server
func (s *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}
	// close listener socket
	defer listener.Close()

	// start message listening goroutine
	go s.ListenMessager()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept error:", err)
			continue
		}
		// do handler
		go s.Handler(conn)
	}
}