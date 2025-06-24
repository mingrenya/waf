package config

type Config struct {}

func Init(path string) error {
    // 这里什么都不做，直接返回nil
    return nil
}

func Get() *Config {
    return &Config{}
}

func NewWAFConfig() interface{} {
    // 占位，返回nil即可
    return nil
}

