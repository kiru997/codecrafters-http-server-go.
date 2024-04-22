package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/pkg/constants"
	"github.com/codecrafters-io/http-server-starter-go/pkg/helper"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

var rootDir string

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	flag.StringVar(&rootDir, "directory", "", "Root directory to conduct a file search on")
	flag.Parse()
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
		go func() {
			err = handle(conn)
			if err != nil {
				fmt.Println("Error handle connection: ", err.Error())
				return
			}
		}()
	}
}

func parse(buf []byte) []string {
	str := string(buf)
	strs := strings.Split(str, string(constants.CRLF))
	return strs
}

func getPath(str string) string {
	paths := strings.Split(str, " ")
	if len(paths) < 2 {
		return ""
	}
	return paths[1]
}

func getMethod(str string) string {
	paths := strings.Split(str, " ")
	if len(paths) < 1 {
		return ""
	}
	return paths[0]
}

func getHeader(strs []string, key string) string {
	for _, v := range strs {
		if strings.Contains(v, key) {
			return v
		}
	}

	return ""
}

func subPath(path string) []string {
	paths := strings.Split(path, constants.Slash)
	return paths
}

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

func handle(conn net.Conn) error {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		return err
	}
	buf = buf[:n]

	strs := parse(buf)
	path := getPath(strs[0])
	subPaths := subPath(path)
	method := getMethod(strs[0])

	switch subPaths[1] {
	case "":
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		if err != nil {
			fmt.Println("Error Write connection: ", err.Error())
			return err
		}
		return nil
	case "echo":
		_, err = conn.Write(helper.NewResponse(http.StatusOK, []byte(strings.Join(subPaths[2:], constants.Slash)), "text/plain"))
		if err != nil {
			fmt.Println("Error Write connection: ", err.Error())
			return err
		}
		return nil
	case "user-agent":
		userAgent := getHeader(strs, "User-Agent")
		userAgent = strings.Replace(userAgent, "User-Agent: ", "", 1)
		_, err = conn.Write(helper.NewResponse(http.StatusOK, []byte(userAgent), "text/plain"))
		if err != nil {
			fmt.Println("Error Write connection: ", err.Error())
			return err
		}
		return nil
	case "files":
		fileName := strings.Join(subPaths[2:], constants.Slash)
		path := filepath.Join(rootDir, fileName)
		if method == "POST" {
			file, err := create(path)
			if err != nil {
				_, err = conn.Write(helper.NewResponse(http.StatusInternalServerError, []byte{}, ""))
				if err != nil {
					fmt.Println("Error Write connection: ", err.Error())
					return err
				}
				return nil
			}
			defer file.Close()

			err = os.WriteFile(path, []byte(strs[len(strs)-1]), 0666)
			if err != nil {
				_, err = conn.Write(helper.NewResponse(http.StatusInternalServerError, []byte{}, ""))
				if err != nil {
					fmt.Println("Error Write connection: ", err.Error())
					return err
				}
				return nil
			}

			_, err = conn.Write(helper.NewResponse(http.StatusCreated, []byte{}, ""))
			if err != nil {
				fmt.Println("Error Write connection: ", err.Error())
				return err
			}
			return nil
		}

		_, err := os.Stat(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				_, err = conn.Write(helper.NewResponse(http.StatusNotFound, []byte{}, ""))
				if err != nil {
					fmt.Println("Error Write connection: ", err.Error())
					return err
				}
			}
			_, err = conn.Write(helper.NewResponse(http.StatusInternalServerError, []byte{}, ""))
			if err != nil {
				fmt.Println("Error Write connection: ", err.Error())
				return err
			}
		}

		data, err := os.ReadFile(path)
		if err != nil {
			_, err = conn.Write(helper.NewResponse(http.StatusInternalServerError, []byte{}, ""))
			if err != nil {
				fmt.Println("Error Write connection: ", err.Error())
				return err
			}
		}
		_, err = conn.Write(helper.NewResponse(http.StatusOK, data, "application/octet-stream"))
		if err != nil {
			fmt.Println("Error Write connection: ", err.Error())
			return err
		}
	}

	_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	if err != nil {
		fmt.Println("Error Write connection: ", err.Error())
		return err
	}
	return nil
}
