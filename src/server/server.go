package server

import (
	"fmt"
	"net"
	"sync"
	"github.com/RunningKuma/It-My-First-GO/src/user"
)

type Server struct {
	Ip   string
	Port int
	//TODO:add user map&lock&channel}
	UserMap map[string]*user.User
	mapLock sync.RWMutex

	Message chan string

}


// GO routine of Listen Message
// if user has message,send to all client
func (this *Server) ListenMessager() {

	//TODO: for linsten
	for{
		msg := <-this.Message

		this.mapLock.Lock()
		for _, cli := range this.UserMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// method to boardcast message
func (this *Server) Broadcast(user *user.User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}



func (this *Server) Handler(conn net.Conn) {
	//work
	user :=  user.CreateUser(conn)

	//insert into user map
	//REMEMBER LOCK
	this.mapLock.Lock()
	this.UserMap[user.Name] = user
	this.mapLock.Unlock()
	//boardcast user online message
	//Receive user Message
	go func() {
		buf := make([]byte, SIZE_MEDIUM)
		for{
		n ,err := conn.Read(buf) 
			
		if n == 0 {
			this.Broadcast(user, "is Offline")
			return 	
		}

		if err != nil {
			fmt.Println("read error:", err)
			return 
		}
		msg := string(buf[:n-1])

		this.Broadcast(user, msg)

		}
	}()


	this.Broadcast(user, "is Login!")
}



// Create Server socket
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
		//TODO:init user map
		UserMap: make(map[string]*user.User),
		Message: make(chan string),	
	}

	return server
}

// Start Server socket
func (this *Server) Start() {

	listener, err := net.Listen("tcp4", "127.0.0.1:8080")
	if err != nil{
		fmt.Println("Listen error:", err)
		return 
	}

	defer listener.Close()
	//start listen message

	go this.ListenMessager()

	for{
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accepted error:",err)
			continue		
		}

		go this.Handler(conn)
		
	}
}

