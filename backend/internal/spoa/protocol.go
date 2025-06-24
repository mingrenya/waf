package spoa

// SPOA 协议定义
type Protocol struct {
    Version string
    Commands []Command
}

type Command struct {
    Name string
    Args []string
}
