package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

type ClientOption func(*Client)

func WithServerIp(ip string) ClientOption {
	return func(c *Client) {
		c.ServerIp = ip
	}
}

func WithServerPort(port int) ClientOption {
	return func(c *Client) {
		c.ServerPort = port
	}
}

func NewClient(opts ...ClientOption) *Client {
	// create client object with defaults
	client := &Client{
		ServerIp:   "127.0.0.1",
		ServerPort: 8888,
		flag:       999,
	}

	for _, opt := range opts {
		opt(client)
	}

	// link server
	conn, err := net.Dial("tcp", net.JoinHostPort(client.ServerIp, fmt.Sprintf("%d", client.ServerPort)))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn

	// return client object
	return client
}

func (c *Client) DealResponse() {
	// receive server message and display it to standard output
	io.Copy(os.Stdout, c.conn)
}

func (c *Client) menu() bool {
	var flag int
	fmt.Println("1. Public Chat")
	fmt.Println("2. Private Chat")
	fmt.Println("3. Update User Name")
	fmt.Println("0. Exit")
	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		c.flag = flag
		return true
	} else {
		fmt.Println(">>> Please input valid number <<<")
		return false
	}
}

// public chat
func (c *Client) PublicChat() {
	var chatMsg string

	fmt.Println(">>> Please input message to send, exit to quit public chat <<<")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		// send to server
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := c.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write error:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>> Please input message to send, exit to quit public chat <<<")
		fmt.Scanln(&chatMsg)
	}
}

// select private chat user
func (c *Client) SelectUser() {
	sendMsg := "who\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write error:", err)
		return
	}
}
	

// private chat
func (c *Client) PrivateChat() {
	c.SelectUser()
	var remoteName string
	var chatMsg string

	fmt.Println(">>> Please input the username to chat with, exit to quit private chat <<<")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>> Please input message to send to " + remoteName + ", exit to quit private chat <<<")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			// send to server
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := c.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write error:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>> Please input message to send to " + remoteName + ", exit to quit private chat <<<")
			fmt.Scanln(&chatMsg)
		}

		c.SelectUser()
		fmt.Println(">>> Please input the username to chat with, exit to quit private chat <<<")
		fmt.Scanln(&remoteName)
	}
}

// update user name
func (c *Client) UpdateName(msg string) bool {
	fmt.Println(">>>>>> pls input name:")
	fmt.Scanln(&c.Name)
	sendMsg := "rename|" + c.Name + "\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write error:", err)
		return false
	}
	return true
}


func (c *Client) Run() {	
	for c.flag != 0 {
		for c.menu() != true {
		}
		switch c.flag {
		case 1:
			// public chat logic
			fmt.Println("Public Chat selected.")
			c.PublicChat()
		case 2:
			// private chat logic
			fmt.Println("Private Chat selected.")
			c.PrivateChat()
		case 3:
			// update user name logic
			fmt.Println("Update User Name selected.")
			c.UpdateName(c.Name)
		case 0:
			// exit
			fmt.Println("Exiting...")
			return
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Server IP address")
	flag.IntVar(&serverPort, "port", 8888, "Server port")
	flag.Parse()
}

func main() {
	client := NewClient(WithServerIp(serverIp), WithServerPort(serverPort))
	if client == nil {
		fmt.Println("Failed to create client.")
		return
	}

	go client.DealResponse()

	fmt.Println("Client created successfully:", client)
	
	client.Run()
}

