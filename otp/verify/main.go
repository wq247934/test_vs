package main

import (
	"fmt"
	"test_vs/otp/otp"
)

func main() {
	secret := "JBSWY3DPEHPK3PXP"
	fmt.Println("Secret:", secret)

	isValid := otp.VerifyOTP(secret, "813458", 60)
	fmt.Println("Is OTP valid?", isValid)

}
