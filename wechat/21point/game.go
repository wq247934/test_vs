package _21point

import (
	"errors"
	"fmt"
	"golang.org/x/exp/rand"
	"sort"
	"strings"
	"time"
)

var cards []*card

type card struct {
	name     string
	pointNum int64
	used     bool
}

func StartGame() {
	cards = []*card{
		{"黑桃A", 11, false}, {"黑桃2", 2, false}, {"黑桃3", 3, false}, {"黑桃4", 4, false}, {"黑桃5", 5, false}, {"黑桃6", 6, false}, {"黑桃7", 7, false}, {"黑桃8", 8, false}, {"黑桃9", 9, false}, {"黑桃10", 10, false}, {"黑桃J", 10, false}, {"黑桃Q", 10, false}, {"黑桃K", 10, false},
		{"红心A", 11, false}, {"红心2", 2, false}, {"红心3", 3, false}, {"红心4", 4, false}, {"红心5", 5, false}, {"红心6", 6, false}, {"红心7", 7, false}, {"红心8", 8, false}, {"红心9", 9, false}, {"红心10", 10, false}, {"红心J", 10, false}, {"红心Q", 10, false}, {"红心K", 10, false},
		{"樱花A", 11, false}, {"樱花2", 2, false}, {"樱花3", 3, false}, {"樱花4", 4, false}, {"樱花5", 5, false}, {"樱花6", 6, false}, {"樱花7", 7, false}, {"樱花8", 8, false}, {"樱花9", 9, false}, {"樱花10", 10, false}, {"樱花J", 10, false}, {"樱花Q", 10, false}, {"樱花K", 10, false},
		{"方块A", 11, false}, {"方块2", 2, false}, {"方块3", 3, false}, {"方块4", 4, false}, {"方块5", 5, false}, {"方块6", 6, false}, {"方块7", 7, false}, {"方块8", 8, false}, {"方块9", 9, false}, {"方块10", 10, false}, {"方块J", 10, false}, {"方块Q", 10, false}, {"方块K", 10, false},
	}
}

func pushCard() (*card, error) {
	if len(cards) == 0 {
		StartGame()
	}
	rand.Seed(uint64(time.Now().UnixNano()))
	for i := 0; i < len(cards); i++ {
		index := rand.Intn(len(cards))
		card := cards[index]
		if !card.used {
			card.used = true
			return card, nil
		}
	}
	return nil, errors.New("牌堆已无牌,请重新开始")
}

type user struct {
	username        string
	cards           []*card
	currentPointNum int64
	finished        bool
}

var users = make(map[string]*user)

func StartGetCards(username string) string {
	_, ok := users[username]
	if ok {
		return username + ",您已开始游戏,请勿重复操作!"
	}
	card1, err := pushCard()
	if err != nil {
		return err.Error()
	}
	card2, err := pushCard()
	if err != nil {
		return err.Error()
	}
	user := &user{
		username:        username,
		cards:           []*card{card1, card2},
		currentPointNum: card1.pointNum + card2.pointNum,
	}
	users[username] = user
	return fmt.Sprintf("%s,您获得了两张牌, [%s] 和 [%s]", username, card1.name, card2.name)
}

func GetCard(username string) (string, error) {
	user, ok := users[username]
	if !ok {
		return "", errors.New(username + ",请先开始游戏")
	}
	if user.finished {
		return "", errors.New(username + ",您当前已无法继续摸牌,请等待游戏结算!")
	}
	card, err := pushCard()
	if err != nil {
		return "", err
	}
	user.cards = append(user.cards, card)
	user.calPointNum()
	if user.currentPointNum > 21 {
		user.finished = true
		return fmt.Sprintf(`
boom!!!! 
%s您获得了一张 [%s],
总点数为%d`, username, card.name, user.currentPointNum), nil
	}
	currentCardStr := "您当前牌为 "
	for _, c := range user.cards {
		currentCardStr += fmt.Sprintf("[%s] ", c.name)
	}
	return fmt.Sprintf("%s,您获得了一张 [%s], %s,您当前点数为 [%d]", username, card.name, currentCardStr, user.currentPointNum), nil
}

func Stop(username string) string {
	user, ok := users[username]
	if !ok {
		return username + ",请先开始游戏"
	}
	user.finished = true
	return fmt.Sprintf("%s已停止摸牌,请等待游戏结算", username)
}

func (u *user) calPointNum() {
	var num int64
	var ANum int
	for _, c := range u.cards {
		num += c.pointNum
		if strings.Contains(c.name, "A") {
			ANum += 1
		}
	}
	for i := 0; i < ANum; i++ {
		if num > 21 {
			num -= 10
		}
	}
	u.currentPointNum = num
}

func SettleGame() string {
	for _, u := range users {
		if !u.finished {
			return ""
		}
	}
	var userArr []*user
	var failedUserArr []*user
	for _, u := range users {
		if u.currentPointNum > 21 {
			failedUserArr = append(failedUserArr, u)
		} else {
			userArr = append(userArr, u)
		}
	}
	if len(userArr) == 0 {
		return "很不幸,全部出局!"
	}
	sort.Slice(userArr, func(i, j int) bool {
		return userArr[i].currentPointNum > userArr[j].currentPointNum
	})
	var result = "当前排名为:\n"
	for i, u := range userArr {
		result += fmt.Sprintf("第%d名  %s  点数:%d\n", i, u.username, u.currentPointNum)
	}
	result += "出局者:"
	for _, u := range failedUserArr {
		result += u.username + ","
	}
	return result + "\n\n最终冠军🏆是:" + userArr[0].username
}

func Reset() {
	users = make(map[string]*user)
}
