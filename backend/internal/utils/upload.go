package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"time"
)

func UploadFile(c *gin.Context, field, directory string) (string, error) {
	file, err := c.FormFile(field)
	if err != nil {
		return "", err
	}

	// 生成文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// 确保目录存在
	if err := ensureDir(directory); err != nil {
		return "", err
	}

	// 保存文件
	dst := filepath.Join(directory, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}

	return filename, nil
}

func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}
