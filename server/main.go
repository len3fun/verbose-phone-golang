package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	Host = "127.0.0.1"
	Port = "9000"
	Type = "tcp"
)

type User struct {
	Name       string
	Connection net.Conn
}

var connectedUsers = make([]User, 0)

func main() {
	fmt.Printf("Starting %v server on %v:%v\n", Type, Host, Port)
	l, err := net.Listen(Type, Host+":"+Port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}

		go handleConnection(c)
		fmt.Printf("Client %v connected\n", c.RemoteAddr().String())

	}
}

func handleConnection(conn net.Conn) {
	user := User{conn.RemoteAddr().String(), conn}
	buffer, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		conn.Write([]byte("Some error"))
		fmt.Println("Some error")
		conn.Close()
		return
	}
	user.Name = string(buffer[:len(buffer)-1])
	connectedUsers = append(connectedUsers, user)
	for {
		buffer, err := bufio.NewReader(conn).ReadBytes('\n')

		if err != nil {
			fmt.Println("Client left")
			conn.Close()
			deleteClient(conn)
			return
		}

		log.Println(user.Name, ":", string(buffer[:len(buffer)-1]))
		sendMessageToOtherUsers(conn, buffer)
	}
}

func deleteClient(conn net.Conn) {
	tmp := connectedUsers[:0]
	for _, user := range connectedUsers {
		if user.Connection != conn {
			tmp = append(tmp, user)
		}
	}
	connectedUsers = tmp
}

func sendMessageToOtherUsers(conn net.Conn, message []byte) {
	for _, user := range connectedUsers {
		if user.Connection != conn {
			user.Connection.Write(message)
		}
	}
}
