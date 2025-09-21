package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp string
	ServerPort int
	Name string
	Conn net.Conn
	flag int //1-boardcast 2-private chat 3-change name 0-exit
}

func (c *Client) menu() bool {
	var flag int
	fmt.Println("Welcome to go-chat, Pick a function you want to use")
	fmt.Println("1. Public Chat")
	fmt.Println("2. Private Chat")
	fmt.Println("3. Change Name")
	fmt.Println("0. Exit")
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		c.flag = flag
		return true
	} else {
		fmt.Println("Invalid choice, please try again.")
		return false
	}
}

func (c *Client) Run() {
	for c.flag != 0 {
		for c.menu() != true {
			//loop until get valid choice
		}
		switch c.flag {
		case 1:
			//public chat
		case 2:

		case 3:
			c.updateName()
		
		}

	}

}

func (c *Client) DealResponse() {
	//receive server message and print to stdout
	//should be a goroutine
	io.Copy(os.Stdout, c.Conn)
	
	// for{
	// 	buf := make([]byte, 4096)
	// 	_, err := c.Conn.Read(buf)
	// 	if err != nil {
	// 		if err != io.EOF {
	// 			fmt.Println("Error reading from server:", err)
	// 		}
	// 		break
	// 	}

	// }
}


func (c *Client) updateName() bool {
	fmt.Println("Please enter your new name:")
	var newName string
	fmt.Scanln(&newName)
	if len(newName) == 0 {
		fmt.Println("Name cannot be empty.")
		return false
	}
	c.Name = newName
	sendMsg := "rename|" + newName + "\n"
	_, err := c.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("Error sending rename request:", err)
		return false
	}
	return true
}

// func (c *Client) PublicChat() {
// 	msg := ""
// 	for msg != exit{
		
// 	}
// }

func NewClient(serverIp string, serverPort int) *Client {
	//create client object
	cli := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag : 999,
	}
	//link to server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return nil
	}
	cli.Conn = conn
	//return client
	return cli
}

var serverIp string
var serverPort int

func init(){
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set server ip(default localhost)")
	flag.IntVar(&serverPort, "port", 8888, "set server port(default 8888)")
}

func main(){
	//parse command line
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("Failed to create client")
		return
	}

	fmt.Println("Client created successfully:", client)

	go client.DealResponse()
	client.Run()
	defer client.Conn.Close()

}




