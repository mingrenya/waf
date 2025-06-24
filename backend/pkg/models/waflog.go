package models

import "time"

type WafLog struct {
	Timestamp       time.Time `bson:"timestamp"`
	ClientIP        string    `bson:"client_ip"`
	RequestMethod   string    `bson:"request_method,omitempty"`
	RequestURI      string    `bson:"request_uri,omitempty"`
	Host            string    `bson:"host,omitempty"`
	UserAgent       string    `bson:"user_agent,omitempty"`
	Referer         string    `bson:"referer,omitempty"`
	RequestHeaders  string    `bson:"request_headers,omitempty"`
	RequestBody     string    `bson:"request_body,omitempty"`
	ResponseHeaders string    `bson:"response_headers,omitempty"`
	ResponseBody    string    `bson:"response_body,omitempty"`
	WafAction       string    `bson:"waf_action"`
	ResponseStatus  int       `bson:"response_status"`
	Latency        int64     `bson:"latency_ms"`
}

