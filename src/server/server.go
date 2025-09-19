package server

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int
	//TODO:add user map&lock&channel}
	UserMap map[string]*User
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

// method to boardcast message.
// user is sender, msg is content, will send message to Server.Message channel 
func (this *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ": " + msg

	this.Message <- sendMsg
}


//Handler API
func (this *Server) Handler(conn net.Conn) {
	//work
	user :=  CreateUser(conn, this)

	user.Online()
	//insert into user map
	//REMEMBER LOCK
	//boardcast user online message
	//Receive user Message

	is_live := make(chan bool) 
	go func() {
		buf := make([]byte, SIZE_MEDIUM)
		for{
		n ,err := conn.Read(buf) 
			
		if n == 0 {
			user.Offline()
			return 	
		}

		if err != nil {
			fmt.Println("read error:", err)
			return 
		}
		msg := string(buf[:n-1])

		user.SendMsg(msg)

		is_live <- true
		}
	}()
	
	for {
			select {
				case <-is_live:
				//do nothing,just for refresh
				case <-time.After(time.Minute*30):
					user.sendSelfMessage("You are kicked off for no activity in 30 minutes")
					kickOffline(user)
					close(user.C)
					return 

			}
		}

}



// Create Server socket
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
		//TODO:init user map
		UserMap: make(map[string]*User),
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

