package mocks

import (
    "github.com/stretchr/testify/mock"
    "coraza-waf/backend/logger" // 例如，模块路径为 "coraza-waf/backend"
)

// 定义一个 mock 接口
type LogMock struct {
    mock.Mock
}

// Log 方法实现
func (m *LogMock) Log(msg string) {
    m.Called(msg)
}

