// internal/middleware/waf.go
package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"
	"strconv"
	
	"github.com/corazawaf/coraza/v3"
	"github.com/corazawaf/coraza/v3/types"
	"github.com/gin-gonic/gin"

	"coraza-waf/backend/pkg/database"
	"coraza-waf/backend/pkg/logging"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}

// WAFMiddleware applies Coraza WAF to each request/response
func WAFMiddleware(waf coraza.WAF, mongo *database.MongoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 1. Capture request body
		var reqBody []byte
		if c.Request.Body != nil {
			reqBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewReader(reqBody))
		}

		// 2. Wrap response writer
		blw := &bodyLogWriter{ResponseWriter: c.Writer, buf: &bytes.Buffer{}}
		c.Writer = blw

		// 3. Start transaction
		tx := waf.NewTransaction()
		defer func() {
			tx.ProcessLogging()
			tx.Close()
		}()

		// 4. Feed request into WAF
		for k, vs := range c.Request.Header {
			for _, v := range vs {
				tx.AddRequestHeader(k, v)
			}
		}
		if len(reqBody) > 0 {
			tx.WriteRequestBody(reqBody)
		}
		tx.ProcessRequestHeaders()
		tx.ProcessRequestBody()

		// 5. If blocked, abort
		if it := tx.Interruption(); it != nil {
			c.AbortWithStatus(it.Status)
			entry := createWafLog(c, start, string(reqBody), "", tx, "BLOCK")
			go mongo.InsertLog(entry)
			return
		}

		// 6. Continue handler chain
		c.Next()

		// 7. Capture response
		respBody := blw.buf.Bytes()
		for k, vs := range c.Writer.Header() {
			for _, v := range vs {
				tx.AddResponseHeader(k, v)
			}
		}
		if len(respBody) > 0 {
			tx.WriteResponseBody(respBody)
		}
		// Note: ProcessResponseHeaders wants (status, contentType)
		tx.ProcessResponseHeaders(c.Writer.Status(), c.Writer.Header().Get("Content-Type"))
		tx.ProcessResponseBody()

		// 8. Log allowed response
		entry := createWafLog(c, start, string(reqBody), string(respBody), tx, "ALLOW")
		go mongo.InsertLog(entry)
	}
}

func createWafLog(c *gin.Context, start time.Time, reqBody, respBody string, tx types.Transaction, action string) *logging.WafLog {
	// pick first matched rule
	var ruleID, ruleMsg string
	if rules := tx.MatchedRules(); len(rules) > 0 {
		r := rules[0].Rule()
		ruleID = strconv.Itoa(r.ID())
	}

	return &logging.WafLog{
		Timestamp:       start,
		ClientIP:        c.ClientIP(),
		RequestMethod:   c.Request.Method,
		RequestURI:      c.Request.RequestURI,
		ServerProtocol:  c.Request.Proto,
		Host:            c.Request.Host,
		UserAgent:       c.Request.UserAgent(),
		Referer:         c.Request.Referer(),
		RequestHeaders:  headersToString(c.Request.Header),
		RequestBody:     truncateString(reqBody, 2048),
		ResponseStatus:  c.Writer.Status(),
		ResponseHeaders: headersToString(c.Writer.Header()),
		ResponseBody:    truncateString(respBody, 4096),
		WafAction:       action,
		RuleID:          ruleID,
		RuleMessage:     ruleMsg,
		Latency:         time.Since(start).Milliseconds(),
		RequestID:       c.GetHeader("X-Request-ID"),
	}
}

func headersToString(h map[string][]string) string {
	var sb strings.Builder
	for k, vs := range h {
		sb.WriteString(k + ": " + strings.Join(vs, ",") + "\n")
	}
	return sb.String()
}

func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}

