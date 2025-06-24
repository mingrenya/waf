package agent

import (
	"fmt"
	"net"
	"io"
	"log"
)

// 简易模拟 SPOE TCP 服务，接收数据，调用 Agent 处理请求与响应
func StartServer(addr string, agent *Agent) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	fmt.Printf("SPOE server listening on %s\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		go handleConnection(conn, agent)
	}
}

func handleConnection(conn net.Conn, agent *Agent) {
	defer conn.Close()
	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %v", err)
			}
			break
		}
		// 简单示范，假设收到的就是请求体和请求头，实际SPOE需解析协议
		reqBody := buffer[:n]
		headers := map[string]string{
			"method":     "GET",
			"path":       "/test",
			"host":       "localhost",
			"user-agent": "Go-SPOE-Client",
		}
		clientIP := conn.RemoteAddr().String()

		agent.HandleRequest(reqBody, headers, clientIP)

		// 回复模拟响应
		respBody := []byte("HTTP/1.1 200 OK\r\n\r\nHello from SPOE agent")
		respHeaders := map[string]string{
			"Content-Type": "text/plain",
		}
		agent.HandleResponse(respBody, respHeaders, clientIP)
	}
}

