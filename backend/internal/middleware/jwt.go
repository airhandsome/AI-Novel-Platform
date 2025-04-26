package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "net/http"
    "strings"
    "time"
)

type Claims struct {
    UserID uint `json:"user_id"`
    jwt.RegisteredClaims
}

const (
    JWTSecret = "your_jwt_secret" // 在实际应用中应该从配置文件读取
)

func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从 Header 中获取 token
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header is required",
            })
            c.Abort()
            return
        }

        // 检查 token 格式
        parts := strings.SplitN(authHeader, " ", 2)
        if !(len(parts) == 2 && parts[0] == "Bearer") {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid authorization format",
            })
            c.Abort()
            return
        }

        // 解析 token
        token, err := jwt.ParseWithClaims(parts[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte(JWTSecret), nil
        })

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid or expired token",
            })
            c.Abort()
            return
        }

        // 验证 token
        claims, ok := token.Claims.(*Claims)
        if !ok || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid token claims",
            })
            c.Abort()
            return
        }

        // 检查token是否过期
        if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Token has expired",
            })
            c.Abort()
            return
        }

        // 将用户ID存储在上下文中
        c.Set("userID", claims.UserID)
        c.Next()
    }
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint) (string, error) {
    claims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // token有效期24小时
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(JWTSecret))
}

// GetUserID 从上下文中获取用户ID的辅助函数
func GetUserID(c *gin.Context) uint {
    userID, exists := c.Get("userID")
    if !exists {
        return 0
    }
    return userID.(uint)
} 