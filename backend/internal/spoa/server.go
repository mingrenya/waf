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
	log.Printf("[DEBUG] New connection from %s", conn.RemoteAddr())

	for {
		// 读取4字节长度字段（大端）
		lenBuf := make([]byte, 4)
		if _, err := io.ReadFull(reader, lenBuf); err != nil {
			log.Printf("[DEBUG] read length error: %v", err)
			return
		}
		frameLen := binary.BigEndian.Uint32(lenBuf)
		log.Printf("[DEBUG] frameLen: %d", frameLen)
		if frameLen == 0 {
			log.Printf("[DEBUG] zero-length frame, closing")
			return
		}

		frameData := make([]byte, frameLen)
		if _, err := io.ReadFull(reader, frameData); err != nil {
			log.Printf("[DEBUG] read frame data error: %v", err)
			return
		}
		log.Printf("[DEBUG] raw frameData: %x", frameData)

		if len(frameData) < 1 {
			log.Printf("[DEBUG] frame too short for type")
			return
		}
		frameType := frameData[0]
		log.Printf("[DEBUG] frameType: %02x", frameType)
		if frameType == 0x01 { // HELLO/ACK/控制帧，简单回 ACK
			log.Printf("[DEBUG] handle HELLO/ACK frame, reply with ACK")
			// 构造 ACK 帧（最简实现，frameType=0x01, 4字节帧ID, 1字节flags, 4字节stream-id, 4字节frame-id, 0字节payload）
			ack := make([]byte, 18) // 1+4+1+4+4+4=18
			ack[0] = 0x01 // type
			copy(ack[1:5], frameData[1:5]) // 帧ID原样返回
			ack[5] = 0x00 // flags
			copy(ack[6:10], frameData[6:10]) // stream-id
			copy(ack[10:14], frameData[10:14]) // frame-id
			// 剩余4字节payload长度为0
			ackLen := make([]byte, 4)
			binary.BigEndian.PutUint32(ackLen, uint32(len(ack)-4))
			ackFrame := append(ackLen, ack...)
			if _, err := conn.Write(ackFrame); err != nil {
				log.Printf("[DEBUG] write ACK error: %v", err)
				return
			}
			continue
		}
		if frameType != 0x02 { // 只处理 NOTIFY 帧，其它类型直接跳过
			log.Printf("[DEBUG] skip non-NOTIFY frame (type=%02x)", frameType)
			continue
		}

		msgName, headers, _, err := parseSPOEFrame(frameData)
		if err != nil || msgName == "" {
			log.Printf("[DEBUG] skip frame: parse error or empty msgName, err=%v", err)
			continue
		}
		log.Printf("[DEBUG] msgName: %s, headers: %+v", msgName, headers)

		var respFrame []byte
		// 官方示例联动逻辑
		switch msgName {
		case "coraza-req":
			if headers["path"] == "/blockme" {
				respFrame, err = buildSPOEResponse("deny", "", 0)
			} else if headers["path"] == "/redirectme" {
				respFrame, err = buildSPOEResponse("redirect", "http://example.com", 0)
			} else {
				respFrame, err = buildSPOEResponse("allow", "", 0)
			}
		case "coraza-res":
			respFrame, err = buildSPOEResponse("allow", "", 0)
		default:
			respFrame, err = buildSPOEResponse("deny", "", 1)
		}
		if err != nil {
			log.Printf("[DEBUG] build response error: %v", err)
			return
		}
		log.Printf("[DEBUG] respFrame: %x", respFrame)
		if _, err := conn.Write(respFrame); err != nil {
			log.Printf("[DEBUG] write response error: %v", err)
			return
		}
	}
}

// parseSPOEFrame 兼容 HAProxy SPOE 协议（简化版，适配常见 HAProxy 3.x 消息帧）
func parseSPOEFrame(data []byte) (msgName string, headers map[string]string, body []byte, err error) {
	headers = make(map[string]string)
	pos := 0
	log.Printf("[DEBUG][parseSPOEFrame] total len=%d", len(data))
	// 1. 读取帧类型（1字节）
	if len(data) < pos+1 {
		err = fmt.Errorf("invalid frame: no frame type")
		return
	}
	frameType := data[pos]
	log.Printf("[DEBUG][parseSPOEFrame] pos=%d frameType=%02x", pos, frameType)
	pos++
	// 2. 读取4字节帧ID
	if len(data) < pos+4 {
		err = fmt.Errorf("invalid frame: no frame id")
		return
	}
	log.Printf("[DEBUG][parseSPOEFrame] pos=%d frameID=%x", pos, data[pos:pos+4])
	pos += 4
	// 3. 读取1字节flags
	if len(data) < pos+1 {
		err = fmt.Errorf("invalid frame: no flags")
		return
	}
	log.Printf("[DEBUG][parseSPOEFrame] pos=%d flags=%02x", pos, data[pos])
	pos++
	// 4. 读取1字节消息名长度
	if len(data) < pos+1 {
		err = fmt.Errorf("invalid frame: no msgName length")
		return
	}
	msgNameLen := int(data[pos])
	log.Printf("[DEBUG][parseSPOEFrame] pos=%d msgNameLen=%d", pos, msgNameLen)
	pos++
	if msgNameLen == 0 {
		log.Printf("[DEBUG][parseSPOEFrame] msgNameLen=0，跳过该帧")
		return
	}
	if len(data) < pos+msgNameLen {
		err = fmt.Errorf("invalid frame: msgName too short, need %d, left %d", msgNameLen, len(data)-pos)
		return
	}
	msgName = string(data[pos : pos+msgNameLen])
	log.Printf("[DEBUG][parseSPOEFrame] pos=%d msgName=%s", pos, msgName)
	pos += msgNameLen
	// 5. 读取1字节 KV 数量
	if len(data) < pos+1 {
		err = fmt.Errorf("invalid frame: no kv count")
		return
	}
	kvCount := int(data[pos])
	log.Printf("[DEBUG][parseSPOEFrame] pos=%d kvCount=%d", pos, kvCount)
	pos++
	// 6. 循环读取每个 key-value
	for i := 0; i < kvCount; i++ {
		if len(data) < pos+1 {
			err = fmt.Errorf("invalid frame: missing key length")
			return
		}
		keyLen := int(data[pos])
		log.Printf("[DEBUG][parseSPOEFrame] pos=%d keyLen=%d", pos, keyLen)
		pos++
		if len(data) < pos+keyLen {
			err = fmt.Errorf("invalid frame: key too short, need %d, left %d", keyLen, len(data)-pos)
			return
		}
		key := string(data[pos : pos+keyLen])
		log.Printf("[DEBUG][parseSPOEFrame] pos=%d key=%s", pos, key)
		pos += keyLen
		if len(data) < pos+1 {
			err = fmt.Errorf("invalid frame: missing value length")
			return
		}
		valLen := int(data[pos])
		log.Printf("[DEBUG][parseSPOEFrame] pos=%d valLen=%d", pos, valLen)
		pos++
		if len(data) < pos+valLen {
			err = fmt.Errorf("invalid frame: value too short, need %d, left %d", valLen, len(data)-pos)
			return
		}
		val := string(data[pos : pos+valLen])
		log.Printf("[DEBUG][parseSPOEFrame] pos=%d val=%s", pos, val)
		pos += valLen
		headers[key] = val
	}
	if pos < len(data) {
		body = data[pos:]
		log.Printf("[DEBUG][parseSPOEFrame] pos=%d bodyLen=%d", pos, len(body))
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



