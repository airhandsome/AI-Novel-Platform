package utils

import (
    "fmt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "log"
    "os"
    "time"
)

func InitDB(username, password, host, port, dbname string) (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        username,
        password,
        host,
        port,
        dbname,
    )
    
    // 配置GORM日志
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold: time.Second,
            LogLevel:      logger.Info,
            Colorful:      true,
        },
    )
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: newLogger,
    })
    if err != nil {
        return nil, err
    }
    
    // 配置连接池
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    return db, nil
} 