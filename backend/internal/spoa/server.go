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
        if msgName == "coraza-req" {
            s.handler.HandleRequest(body, headers, clientIP)
        } else if msgName == "coraza-res" {
            s.handler.HandleResponse(body, headers, clientIP)
        }

        // 这里根据WAF逻辑可以调整是否阻断，先统一允许
        respFrame, err := buildSPOEResponse(true)
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


// buildSPOEResponse 构造允许动作响应帧
func buildSPOEResponse(allow bool) ([]byte, error) {
	var action string
	if allow {
		action = "allow"
	} else {
		action = "deny"
	}

	// 构建 SPOE kv 结构：
	// [keyLen][key][valLen][val]
	key := "action"
	val := action
	kb := []byte(key)
	vb := []byte(val)

	framePayload := []byte{
		byte(len(kb)),
	}
	framePayload = append(framePayload, kb...)
	framePayload = append(framePayload, byte(len(vb)))
	framePayload = append(framePayload, vb...)

	// 前缀加上 4 字节长度字段（大端）
	totalLen := uint32(len(framePayload))
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, totalLen)

	// 拼接最终响应帧
	resp := append(buf, framePayload...)
	return resp, nil
}


