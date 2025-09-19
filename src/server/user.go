package server

import (
	"net"

)

// User type
type User struct {
	Name string
	Addr string
	C chan string
	Conn net.Conn

	server *Server
}


// create User API
//	should get a connection,return user object
func CreateUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		Conn: conn,
		server: server,
	}
	
	go user.ListenMessage()

	return user
}


func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.UserMap[this.Name] = this
	this.server.mapLock.Unlock()
	this.server.Broadcast(this, "is Login!")
}

func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.UserMap, this.Name)
	this.server.mapLock.Unlock()
	this.server.Broadcast(this, "is Offline!")
}

func (this *User) SendMsg(msg string) {
	this.server.Broadcast(this, msg)
}

//listen user channel,if user has message, send to client

func (this *User) ListenMessage(){
	for{
		msg := <-this.C

		this.Conn.Write([]byte(msg + "\n"))
		
	}
}