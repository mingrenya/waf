// pkg/logging/waf_log.go
package logging

import (
	"time"
)

type WafLog struct {
	Timestamp       time.Time `bson:"timestamp"`
	ClientIP        string    `bson:"client_ip"`
	RequestMethod   string    `bson:"request_method"`
	RequestURI      string    `bson:"request_uri"`
	ServerProtocol  string    `bson:"server_protocol"`
	Host            string    `bson:"host"`
	UserAgent       string    `bson:"user_agent"`
	Referer         string    `bson:"referer"`
	RequestHeaders  string    `bson:"request_headers,omitempty"`
	RequestBody     string    `bson:"request_body,omitempty"`
	ResponseStatus  int       `bson:"response_status"`
	ResponseHeaders string    `bson:"response_headers,omitempty"`
	ResponseBody    string    `bson:"response_body,omitempty"`
	WafAction       string    `bson:"waf_action"` // ALLOW/BLOCK
	RuleID          string    `bson:"rule_id"`
	RuleMessage     string    `bson:"rule_message"`
	Latency         int64     `bson:"latency"` // 毫秒
	RequestID       string    `bson:"request_id"`
	UpstreamAddr    string    `bson:"upstream_addr,omitempty"`
	SecurityLevel   string    `bson:"security_level,omitempty"`
}
