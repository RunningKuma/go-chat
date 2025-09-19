package server

import (
	"net"
	"strings"
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



//boardcast if user is Login
func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.UserMap[this.Name] = this
	this.server.mapLock.Unlock()
	this.server.Broadcast(this, "is Login!")
}

//boardcast if user is Offline
func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.UserMap, this.Name)
	this.server.mapLock.Unlock()
	this.server.Broadcast(this, "is Offline!")
}

//send message to self, not boardcast
func (this *User) sendSelfMessage(msg string) {
	this.Conn.Write([]byte(msg + "\n"))
}
//kick user offline(can use in many cases)
func kickOffline(u *User) {
	u.Conn.Close()
}

//user send message API
func (this *User) SendMsg(msg string) {
	if msg == "whoison"{
		this.server.mapLock.Lock()
		onlineMsg := "These user(s) are online:\n"
		for _, user := range this.server.UserMap {
			 onlineMsg += user.Addr + ":" + user.Name + "\n"
		}
		this.server.mapLock.Unlock()
		this.sendSelfMessage(onlineMsg)
		return
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		_, ok := this.server.UserMap[newName]
		if ok {
			this.sendSelfMessage("This name has been used!")
			return
		} else {
			this.server.mapLock.Lock()
			delete(this.server.UserMap, this.Name)
			this.server.UserMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName
			this.sendSelfMessage("You have renamed to:" + this.Name)
			return 
		}
	} else if msg == "exit" {
		this.sendSelfMessage("Thanks for using,bye!")
		kickOffline(this)
		return
	}
	this.server.Broadcast(this, msg)
}

//listen user channel,if user has message, send to client
func (this *User) ListenMessage(){
	for{
		msg := <-this.C

		this.Conn.Write([]byte(msg + "\n"))
		
	}
}