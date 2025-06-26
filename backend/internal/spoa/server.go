package spoa

import (
    "bufio"
    "encoding/binary"
    "fmt"
    "io"
    "log"
    "net"

    "github.com/corazawaf/coraza/v3"
    "coraza-waf/backend/internal/agent"
    "coraza-waf/backend/pkg/database"
)

// SPOAServer 服务器
type SPOAServer struct {
    addr    string
    handler *agent.Agent
}

func NewServer(addr string, waf coraza.WAF, mongo *database.MongoService) *SPOAServer {
    return &SPOAServer{
        addr:    addr,
        handler: agent.NewAgent(waf, mongo),
    }
}

// Run 启动监听
func (s *SPOAServer) Run() error {
    ln, err := net.Listen("tcp", s.addr)
    if err != nil {
        return fmt.Errorf("listen error: %w", err)
    }
    log.Printf("SPOE server listening on %s", s.addr)

    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Printf("accept error: %v", err)
            continue
        }
        go s.handleConn(conn)
    }
}

// handleConn 处理连接
func (s *SPOAServer) handleConn(conn net.Conn) {
    defer conn.Close()
    reader := bufio.NewReader(conn)

    for {
        // 读取4字节长度字段（大端）
        lenBuf := make([]byte, 4)
        if _, err := io.ReadFull(reader, lenBuf); err != nil {
            log.Printf("read length error: %v", err)
            return
        }
        frameLen := binary.BigEndian.Uint32(lenBuf)
        if frameLen == 0 {
            log.Printf("zero-length frame, closing")
            return
        }

        frameData := make([]byte, frameLen)
        if _, err := io.ReadFull(reader, frameData); err != nil {
            log.Printf("read frame data error: %v", err)
            return
        }

        msgName, headers, body, err := parseSPOEFrame(frameData)
        if err != nil {
            log.Printf("parse frame error: %v", err)
            return
        }

        clientIP := headers["client-ip"]
        var respFrame []byte
        // 官方示例联动逻辑
        switch msgName {
        case "coraza-req":
            // 这里可根据 headers/body 进行 WAF 检测，示例：
            // 1. 拦截特定路径
            if headers["path"] == "/blockme" {
                respFrame, err = buildSPOEResponse("deny", "", 0)
            } else if headers["path"] == "/redirectme" {
                respFrame, err = buildSPOEResponse("redirect", "http://example.com", 0)
            } else {
                respFrame, err = buildSPOEResponse("allow", "", 0)
            }
        case "coraza-res":
            // 可根据响应内容做二次检测
            respFrame, err = buildSPOEResponse("allow", "", 0)
        default:
            // 未知消息，返回错误
            respFrame, err = buildSPOEResponse("deny", "", 1)
        }
        if err != nil {
            log.Printf("build response error: %v", err)
            return
        }
        if _, err := conn.Write(respFrame); err != nil {
            log.Printf("write response error: %v", err)
            return
        }
    }
}

// parseSPOEFrame 示例，实际需要解析SPOE二进制协议
func parseSPOEFrame(data []byte) (msgName string, headers map[string]string, body []byte, err error) {
	headers = make(map[string]string)
	pos := 0

	// 1. 读取消息名长度和内容
	if len(data) < pos+1 {
		err = fmt.Errorf("invalid frame: no msgName length")
		return
	}
	msgNameLen := int(data[pos])
	pos++

	if len(data) < pos+msgNameLen {
		err = fmt.Errorf("invalid frame: msgName too short")
		return
	}
	msgName = string(data[pos : pos+msgNameLen])
	pos += msgNameLen

	// 2. 读取 KV 数量
	if len(data) < pos+1 {
		err = fmt.Errorf("invalid frame: no kv count")
		return
	}
	kvCount := int(data[pos])
	pos++

	// 3. 循环读取每个 key-value
	for i := 0; i < kvCount; i++ {
		if len(data) < pos+1 {
			err = fmt.Errorf("invalid frame: missing key length")
			return
		}
		keyLen := int(data[pos])
		pos++
		if len(data) < pos+keyLen {
			err = fmt.Errorf("invalid frame: key too short")
			return
		}
		key := string(data[pos : pos+keyLen])
		pos += keyLen

		if len(data) < pos+1 {
			err = fmt.Errorf("invalid frame: missing value length")
			return
		}
		valLen := int(data[pos])
		pos++
		if len(data) < pos+valLen {
			err = fmt.Errorf("invalid frame: value too short")
			return
		}
		val := string(data[pos : pos+valLen])
		pos += valLen

		headers[key] = val
	}

	// 4. 读取剩余作为 Body
	if pos < len(data) {
		body = data[pos:]
	}

	return
}


// buildSPOEResponse 构造允许动作响应帧，支持 coraza.action/coraza.data/coraza.error
func buildSPOEResponse(action string, data string, errorCode int) ([]byte, error) {
	kvs := [][]byte{}

	// coraza.action
	key := []byte("coraza.action")
	val := []byte(action)
	kvs = append(kvs, []byte{byte(len(key))}, key, []byte{byte(len(val))}, val)

	// coraza.data（可选）
	if data != "" {
		key = []byte("coraza.data")
		val = []byte(data)
		kvs = append(kvs, []byte{byte(len(key))}, key, []byte{byte(len(val))}, val)
	}

	// coraza.error（可选）
	if errorCode != 0 {
		key = []byte("coraza.error")
		val = []byte(fmt.Sprintf("%d", errorCode))
		kvs = append(kvs, []byte{byte(len(key))}, key, []byte{byte(len(val))}, val)
	}

	framePayload := []byte{}
	for _, b := range kvs {
		framePayload = append(framePayload, b...)
	}

	totalLen := uint32(len(framePayload))
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, totalLen)
	resp := append(buf, framePayload...)
	return resp, nil
}


