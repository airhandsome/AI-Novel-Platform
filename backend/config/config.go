type Config struct {
    JWT struct {
        Secret string
        Expire time.Duration
    }
    // ... 其他配置
} 