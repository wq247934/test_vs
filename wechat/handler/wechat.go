package handler

import (
	"bytes"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"os"
	"strings"
	_21point "test_vs/wechat/21point"
	"test_vs/wechat/chat_gpt"
	img2 "test_vs/wechat/img"
)

const (
	chatMode  = "WeChatMode"
	gameMode  = "GameMode"
	imageMode = "ImageMode"
)

var (
	userModMap = make(map[string]string)
)

func HandleMsg(msg *openwechat.Message) {
	sender, err := msg.Sender()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, ok := sender.AsGroup()
	if !ok {
		go handleUser(msg)
	}
	go handleGroup(msg)
}

func handleUser(msg *openwechat.Message) {

}
func handleGroup(msg *openwechat.Message) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			sender, _ := msg.Sender()
			group, _ := sender.AsGroup()
			fmt.Printf("%+v\n", group)
		}
	}()
	sender, _ := msg.Sender()
	group, _ := sender.AsGroup()
	user, err := msg.SenderInGroup()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(group.NickName)

	fmt.Println("msg.ToUserName => ", msg.ToUserName)
	fmt.Println(msg.Content)

	if msg.IsText() && msg.IsAt() {
		go handleCli(user.NickName, msg)
	}
}

func handleCli(username string, msg *openwechat.Message) {
	if strings.HasSuffix(msg.Content, "开启游戏模式") {
		userModMap[username] = gameMode
		msg.ReplyText("已开启")
	} else if strings.Contains(msg.Content, "画一个") {
		prompt := strings.ReplaceAll(msg.Content, "画一个", "")
		imgBytes, err := chat_gpt.CreateImage(prompt)
		if err != nil {
			fmt.Println(err)
			imgBytes, _ = chat_gpt.CreateImage(prompt)
			return
		}
		_, err = msg.ReplyImage(bytes.NewReader(imgBytes))
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}
	mode, ok := userModMap[username]
	if !ok {
		mode = chatMode
		userModMap[username] = mode
	}
	switch mode {
	case chatMode:
		msg.ReplyText(chat_gpt.Chat(username, strings.ReplaceAll(msg.Content, "@bot", "")))
	case gameMode:
		playGame(username, msg)
	}

}

func playGame(username string, msg *openwechat.Message) {
	switch {
	case strings.Contains(msg.Content, "关闭游戏模式"):
		userModMap[username] = chatMode
		msg.ReplyText("已关闭")
	case strings.Contains(msg.Content, "初始化"):
		_21point.StartGame()
		_, err := msg.ReplyText("初始化成功")
		if err != nil {
			fmt.Println(err)
			return
		}
	case strings.Contains(msg.Content, "开始游戏"):
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
