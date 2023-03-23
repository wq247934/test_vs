package main

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"test_vs/wechat/handler"
)

func main() {
	startWechat()
}

func startWechat() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式
	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")

	defer reloadStorage.Close()

	// 执行热登录
	if err := bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
		fmt.Println(err)
		return
	}
	// 获取登陆的用户
	_, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	bot.MessageHandler = handler.HandleMsg

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
