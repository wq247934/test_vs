package otp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"strings"
	"time"
)

func GenerateOTP(secret string, timeStep int64) string {
	// 将密钥转换为字节数组
	key, _ := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	// 将时间步长转换为字节数组
	timeStepBytes := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		timeStepBytes[i] = byte(timeStep & 0xff)
		timeStep >>= 8
	}
	// 计算HMAC-SHA1哈希值
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(timeStepBytes)
	hmacHash := hmacSha1.Sum(nil)
	// 获取动态截断偏移量
	offset := int(hmacHash[19] & 0xf)
	// 计算动态截断值
	truncatedHash := int32(hmacHash[offset]&0x7f)<<24 |
		int32(hmacHash[offset+1]&0xff)<<16 |
		int32(hmacHash[offset+2]&0xff)<<8 |
		int32(hmacHash[offset+3]&0xff)
	// 计算OTP值
	otp := truncatedHash % 1000000
	// 返回OTP值
	return fmt.Sprintf("%06d", otp)
}

func VerifyOTP(secret string, otp string, timeWindow int64) bool {
	// 获取当前时间戳
	epochSeconds := time.Now().Unix()
	// 遍历时间窗口内的所有时间步长
	for i := -timeWindow; i <= timeWindow; i++ {
		// 计算当前时间步长
		timeStep := (epochSeconds + i) / 30
		// 生成当前时间步长的OTP值
		currentOTP := GenerateOTP(secret, timeStep)
		// 如果当前时间步长的OTP值与用户输入的OTP值相同，则验证成功
		if currentOTP == otp {
			return true
		}
	}
	// 验证失败
	return false
}
