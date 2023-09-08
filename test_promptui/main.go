package main

import (
	"fmt"
)

const (
	MenuDecrypt               = "解密"
	MenuKeyManage             = "密钥管理"
	MenuChangPwd              = "修改密码"
	MenuGenerateChallengeCode = "生成挑战码"
)

func main() {
	menu := []string{MenuDecrypt, MenuKeyManage}
	fmt.Println("Menu:")
	for i, m := range menu {
		fmt.Printf("%d. %s\n", i+1, m)
	}
	fmt.Print("Enter your choice: ")
	var choice int
	fmt.Scan(&choice)
	switch choice {
	case 1:
		handleDecryptMenu()
	case 2:
		handleKeyManageMenu()
	default:
		fmt.Println("Invalid choice")
	}
}

func handleDecryptMenu() {
	fmt.Println("解密.....")
}

func handleKeyManageMenu() {
	subMenu := []string{MenuChangPwd, MenuGenerateChallengeCode}
	fmt.Println("Submenu:")
	for i, m := range subMenu {
		fmt.Printf("%d. %s\n", i+1, m)
	}
	fmt.Print("Enter your choice: ")
	var subChoice int
	fmt.Scan(&subChoice)
	switch subChoice {
	case 1:
		handleChangePwdMenu()
	case 2:
		handleGenerateChallengeMenu()
	default:
		fmt.Println("Invalid choice")
	}
}

func handleGenerateChallengeMenu() {
	fmt.Println("生成挑战码....")
}

func handleChangePwdMenu() {
	fmt.Println("修改密码....")
}
