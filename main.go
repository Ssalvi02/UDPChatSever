package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
)

var (
	clients   = make(map[string]*net.UDPAddr)
	mu        sync.Mutex
	broadcast = make(chan string)
)

func main() {
	addr := net.UDPAddr{
		Port: 7575,
		IP:   net.IPv4(0, 0, 0, 0),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("ERROR: Could not start server! - ", err)
		return
	}
	defer conn.Close()

	fmt.Println("UDP Server created at port " + strconv.Itoa(addr.Port))

	go handleBroadcast(conn)

	buff := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buff)
		if err != nil {
			fmt.Println("ERROR: Could not read message! - ", err)
			continue
		}

		mu.Lock()
		clients[clientAddr.String()] = clientAddr
		mu.Unlock()

		message := string(buff[:n])
		fmt.Printf("| IP [%s] | Message: %s|\n", clientAddr, message)
		broadcast <- fmt.Sprintf("[%s]: %s", clientAddr.String(), message)
	}
}

func handleBroadcast(conn *net.UDPConn) {
	for {
		msg := <-broadcast

		mu.Lock()
		for _, clientAddr := range clients {
			_, err := conn.WriteToUDP([]byte(msg), clientAddr)
			if err != nil {
				fmt.Println("ERROR: Could not send message! - ", err)
			}
		}
		mu.Unlock()
	}
}
