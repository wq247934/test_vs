package main

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"os"
	"strings"
	_21point "test_vs/wechat/21point"
	"test_vs/wechat/chat_gpt"
	img2 "test_vs/wechat/img"
)

func main() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
	}
	//// 注册登陆二维码回调
	//bot.UUIDCallback = openwechat.PrintlnQrcodeUrl
	//
	//// 登陆
	//if err := bot.Login(); err != nil {
	//	fmt.Println(err)
	//	return
	//}

	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")

	defer reloadStorage.Close()

	// 执行热登录

	if err := bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
		fmt.Println(err)
		return
	}
	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	bot.MessageHandler = func(msg *openwechat.Message) {
		go handleMsg(msg)
	}

	// 获取所有的好友
	friends, err := self.Friends()
	fmt.Println(friends, err)

	// 获取所有的群组
	groups, err := self.Groups()
	fmt.Println(groups, err)

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}

func handleMsg(msg *openwechat.Message) {
	sender, err := msg.Sender()
	if err != nil {
		fmt.Println(err)
	}
	group, ok := sender.AsGroup()
	if !ok {
		return
	}
	user, err := msg.SenderInGroup()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(group.NickName)
	//if group.NickName != "测试群" && sender.UserName != "@@5934148189f368e1e6b8a9998aa151cfd998fd3329542cbe63f752798a8f1553" {
	//	return
	//}

	fmt.Println("msg.ToUserName => ", msg.ToUserName)
	fmt.Println(msg.Content)

	if msg.IsText() && msg.IsAt() {
		handleCli(user.NickName, msg)
	}
}

var gameMode bool

func handleCli(username string, msg *openwechat.Message) {

	if gameMode {
		playGame(username, msg)
		return
	}
	switch {
	case strings.HasSuffix(msg.Content, "开启游戏模式"):
		gameMode = true
		msg.ReplyText("已开启")
	case strings.HasSuffix(msg.Content, "reset"), strings.HasSuffix(msg.Content, "重置"):
		msg.ReplyText(chat_gpt.Reset(username))
	default:
		msg.ReplyText(chat_gpt.Chat(username, strings.ReplaceAll(msg.Content, "@bot", "")))
	}

}

func playGame(username string, msg *openwechat.Message) {
	switch {
	case strings.Contains(msg.Content, "开启游戏模式"):
		gameMode = true
		msg.ReplyText("已开启")
	case strings.Contains(msg.Content, "关闭游戏模式"):
		gameMode = false
		msg.ReplyText("已关闭")
	case strings.Contains(msg.Content, "初始化"):
		if !gameMode {
			return
		}
		_21point.StartGame()
		_, err := msg.ReplyText("初始化成功")
		if err != nil {
			fmt.Println(err)
			return
		}
	case strings.Contains(msg.Content, "开始游戏"):
		if !gameMode {
			return
		}
		result, images := _21point.StartGetCards(username)
		img, err := img2.GetCardImg(images...)
		if err != nil {
			fmt.Println(err)
		} else {
			_, err := msg.ReplyImage(img)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		defer os.Remove(img.Name())
		msg.ReplyText(result)
	case strings.Contains(msg.Content, "摸牌"):
		result, image, err := _21point.GetCard(username)
		if err != nil {
			msg.ReplyText(err.Error())
			return
		}
		img, err := img2.GetCardImg(image)
		if err != nil {
			fmt.Println(err)
		} else {
			msg.ReplyImage(img)
		}

		msg.ReplyText(result)
		result = _21point.SettleGame()
		if result != "" {
			_, err := msg.ReplyText(result)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	case strings.Contains(msg.Content, "停牌"):
		msg.ReplyText(_21point.Stop(username))
		result := _21point.SettleGame()
		if result != "" {
			_, err := msg.ReplyText(result)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	case strings.Contains(msg.Content, "游戏规则"):
		printRule(msg)

	case strings.Contains(msg.Content, "重新开始"):
		_21point.Reset()
		msg.ReplyText("游戏数据已重置,请@我 开始游戏")
	default:
		printHelp(msg)
	}
}

func printHelp(msg *openwechat.Message) {
	_, err := msg.ReplyText(`
请输入指令:
	开启游戏模式
	关闭游戏模式
	游戏规则
 	初始化
	开始游戏
	摸牌
	停牌
	重新开始
`)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func printRule(msg *openwechat.Message) {
	msg.ReplyText(`
21点，又名黑杰克（Blackjack），是一种扑克牌游戏。
游戏者的目标是使手中的牌的点数之和不超过21点且尽量大。

2点至10点的牌以牌面的点数计算，J、Q、K 每张为10点。A可记为1点或11点，若玩家会因A而爆牌则A可算为1点。

开始游戏后,玩家会获得两张初始牌.

接下来，玩家可以选择要牌或停牌。如果玩家选择要牌，则会继续给玩家发牌，直到玩家选择停牌或者爆掉（手中的牌点数之和超过21点）。

当所有玩家都停牌后或爆掉后，开始游戏结算，点数大的一方获胜。
`)
}
