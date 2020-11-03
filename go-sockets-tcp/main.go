package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strings"
)

type User struct {
	id int
	username string
}

type Message struct {
	User
	content string
}

func (m *Message) send(conn net.Conn) {
	_, _ = io.WriteString(conn, m.content)
}

func main() {
	l, err := net.Listen("tcp", "localhost:8080")

	if err != nil {
		log.Fatalln("error listening:", err.Error())
	}

	defer l.Close()

	users := make(map[net.Conn]User)
	newConnection := make(chan net.Conn)
	messages := make(chan Message)
	activeUsers := make(chan net.Conn)
	disconnectedUser := make(chan net.Conn)

	go handleNewConnection(l, newConnection)

	handleActions(newConnection, users, messages, activeUsers, disconnectedUser)

}

func handleNewConnection(l net.Listener, newConnection chan<- net.Conn) {
	for {
		c, err := l.Accept()

		if err != nil {
			log.Println("Error connecting:", err.Error())
			return
		}

		newConnection <- c
	}
}

func handleActions(newConnection <-chan net.Conn, users map[net.Conn]User, messages chan Message, activeUsers chan net.Conn, disconnectedUser chan net.Conn) {

	for {
		select {
		case conn := <-newConnection:
			go func(conn net.Conn) {
				reader := bufio.NewReader(conn)
				_, _ = io.WriteString(conn, "Enter your username: ")
				username, _ := reader.ReadString('\n')
				// removing the \n at the end of the username
				username = username[:len(username)-1]
				user := User{
					id: rand.Int(),
					username: username,
				}

				users[conn] = user

				message := Message{
					User:    users[conn],
					content: fmt.Sprintf("%v has joined the server\n", username),
				}

				messages <- message
				activeUsers <- conn
			}(conn)
		case conn := <- activeUsers:
			go func(conn net.Conn) {
				reader := bufio.NewReader(conn)

				for {
					message, err := reader.ReadString('\n')

					if err != nil {
						break
					}

					message = message[:len(message)-1]
					username := users[conn].username

					if strings.HasPrefix(message, "/name ") {
						s := strings.Split(message, " ")
						newUsername := s[1]

						newUser := users[conn]
						newUser.username = newUsername

						users[conn] = newUser


						messages <- Message{
							User:    users[conn],
							content: fmt.Sprintf("%v changed his username to %v\n", username, newUsername),
						}
					} else {
						messages <- Message{
							User:    users[conn],
							content: fmt.Sprintf("%v: %v\n", username, message),
						}
					}
				}

				// if user exits
				disconnectedUser <- conn
			}(conn)
		case message := <- messages:
			for conn := range users {
				// don't send the message to the sender
				if users[conn].id != message.id {
					message.send(conn)
				}
			}
		case conn := <- disconnectedUser:
			delete(users, conn)
			conn.Close()
		}
	}
}