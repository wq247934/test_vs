package main

import (
	"fmt"
	"test_vs/otp/otp"
	"time"
)

func main() {
	secret := "JBSWY3DPEHPK3PXP"
	fmt.Println("Secret:", secret)
	// 获取当前时间戳
	epochSeconds := time.Now().Unix()
	// 计算当前时间步长
	timeStep := epochSeconds / 30
	// 生成当前时间步长的OTP值
	passcode := otp.GenerateOTP(secret, timeStep)
	fmt.Println("Passcode:", passcode)
}
