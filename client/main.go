package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	Host = "localhost"
	Port = "9000"
	Type = "tcp"
)

func main() {
	fmt.Printf("Connecting to %v server %v:%v\n", Type, Host, Port)
	conn, err := net.Dial(Type, Host+":"+Port)
	defer conn.Close()
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Enter your name: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	conn.Write([]byte(input)) // отправляем ник пользователя
	go receiver(conn)
	for {
		fmt.Print("Text to send: ")
		input, _ := reader.ReadString('\n')
		conn.Write([]byte(input))

	}

}

func receiver(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("\nServer Error")
			os.Exit(1)
			return
		}
		fmt.Printf("\n%v", message)
	}
}
