package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

var CRLF = []byte("\r\n")

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		err = handle(conn)
		if err != nil {
			fmt.Println("Error handle connection: ", err.Error())
			continue
		}

	}
}

func parse(buf []byte) []string {
	str := string(buf)
	strs := strings.Split(str, string(CRLF))
	return strs
}

func getPath(str string) string {
	paths := strings.Split(str, " ")
	if len(paths) < 2 {
		return ""
	}
	return paths[1]
}

func handle(conn net.Conn) error {
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		return err
	}

	strs := parse(buf)
	path := getPath(strs[0])
	if path == "/" {
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		if err != nil {
			fmt.Println("Error Write connection: ", err.Error())
			return err
		}
		return nil
	}

	_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	if err != nil {
		fmt.Println("Error Write connection: ", err.Error())
		return err
	}
	return nil
}
