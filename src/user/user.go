package user

import (
	"net"
)

// User type
type User struct {
	Name string
	Addr string
	C chan string
	Conn net.Conn
}


// create User API
//	should get a connection,return user object
func CreateUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		Conn: conn,
	}
	
	go user.ListenMessage()

	return user
}




//listen user channel,if user has message, send to client

func (this *User) ListenMessage(){
	for{
		msg := <-this.C

		this.Conn.Write([]byte(msg + "\n"))
		
	}
}