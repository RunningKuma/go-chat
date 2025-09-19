package server

import (
	"net"
	"strings"
	"errors"
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

//broadcast if user is Login
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


func (this *User) isCorrectFormat(s string) ([]string, error) {
	parts := strings.Split(s,"|")
	if len(parts) < 3 && parts[0] == "to"{
			this.sendSelfMessage("The format of private message is: to|name|message")
			return nil, errors.New("error")
	}
	return parts, nil
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
	} else if len(msg) >= 3 && msg[:3] == "to|" {
		//structure: to|name|message
		//1. get user name
		parts, err := this.isCorrectFormat(msg)
		if err != nil{
			return
		}
		remoteName := parts[1]
		if remoteName == "" {
			this.sendSelfMessage("The format of private message is: to|name|message")
			return
		}
		
		//2. get User object
		targetUser, ok := this.server.UserMap[remoteName]
		if !ok {
			this.sendSelfMessage(remoteName + " seems Not online now")
			return 
		}
		//3. get message and send by user.C
		message := parts[2]

		if message == "" {
			this.sendSelfMessage("Empty message is not allowed")
			return 
		}
		//fix issue that other user may kick or changename by private chat
		targetUser.sendSelfMessage(this.Name + " said to you in private: " + message)
	} else if msg == "exit" {
		this.sendSelfMessage("Thanks for using,bye!")
		kickOffline(this)
		return
	} else {
		this.server.Broadcast(this, msg)
	}
		
}


func (this *User) ListenMessage(){
	for{
		msg := <-this.C

		this.Conn.Write([]byte(msg + "\n"))
		
	}
}