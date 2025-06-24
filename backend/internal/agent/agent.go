package agent

import (
	"log"
	"strings"
	"time"

	"coraza-waf/backend/pkg/database"
	"coraza-waf/backend/pkg/logging"

	"github.com/corazawaf/coraza/v3"
	"github.com/corazawaf/coraza/v3/types"
)

type Agent struct {
	WAF   coraza.WAF
	Mongo *database.MongoService
}

func NewAgent(waf coraza.WAF, mongo *database.MongoService) *Agent {
	return &Agent{WAF: waf, Mongo: mongo}
}

func (a *Agent) HandleRequest(reqBody []byte, headers map[string]string, clientIP string) {
	start := time.Now()
	tx := a.WAF.NewTransaction()
	defer tx.ProcessLogging()

	for k, v := range headers {
		tx.AddRequestHeader(k, v)
	}
	if len(reqBody) > 0 {
		tx.WriteRequestBody(reqBody)
	}

	tx.ProcessRequestHeaders()
	tx.ProcessRequestBody()

	logEntry := &logging.WafLog{
		Timestamp:      start,
		ClientIP:       clientIP,
		RequestMethod:  headers["method"],
		RequestURI:     headers["path"],
		Host:           headers["host"],
		UserAgent:      headers["user-agent"],
		Referer:        headers["referer"],
		RequestHeaders: mapToString(headers),
		RequestBody:    string(reqBody),
		WafAction:      getAction(tx),
		ResponseStatus: 200,
		Latency:        time.Since(start).Milliseconds(),
	}

	if err := a.Mongo.InsertLog(logEntry); err != nil {
		log.Printf("InsertLog error: %v", err)
	}
}

func (a *Agent) HandleResponse(respBody []byte, headers map[string]string, clientIP string) {
	start := time.Now()
	tx := a.WAF.NewTransaction()
	defer tx.ProcessLogging()

	for k, v := range headers {
		tx.AddResponseHeader(k, v)
	}
	if len(respBody) > 0 {
		tx.WriteResponseBody(respBody)
	}

	tx.ProcessResponseHeaders(200, "OK")
	tx.ProcessResponseBody()

	logEntry := &logging.WafLog{
		Timestamp:       start,
		ClientIP:        clientIP,
		ResponseHeaders: mapToString(headers),
		ResponseBody:    string(respBody),
		WafAction:       getAction(tx),
		ResponseStatus:  200,
		Latency:         time.Since(start).Milliseconds(),
	}

	if err := a.Mongo.InsertLog(logEntry); err != nil {
		log.Printf("InsertLog error: %v", err)
	}
}

func getAction(tx types.Transaction) string {
    if it := tx.Interruption(); it != nil {
        return "BLOCK"
    }
    return "ALLOW"
}

func mapToString(m map[string]string) string {
	var sb strings.Builder
	for k, v := range m {
		sb.WriteString(k + ": " + v + "\n")
	}
	return sb.String()
}

