package run

import (
	"log"
	"net"
)

func Start(socketPath string,run func()) {
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", socketPath, err)
	}
	defer listener.Close()

	log.Printf("Server started. Listening on %s", socketPath)

	for {
		// 接受新的连接请求
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		// 启动一个新的goroutine处理请求
		go func(conn net.Conn) {

		}(conn)
	}
}