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
	img      string
}

func StartGame() {
	cards = []*card{
		{"黑桃A", 11, false, "ace_of_spades.jpg"}, {"黑桃2", 2, false, "2_of_spades.jpg"}, {"黑桃3", 3, false, "3_of_spades.jpg"}, {"黑桃4", 4, false, "4_of_spades.jpg"}, {"黑桃5", 5, false, "5_of_spades.jpg"}, {"黑桃6", 6, false, "6_of_spades.jpg"}, {"黑桃7", 7, false, "7_of_spades.jpg"}, {"黑桃8", 8, false, "8_of_spades.jpg"}, {"黑桃9", 9, false, "9_of_spades.jpg"}, {"黑桃10", 10, false, "10_of_spades.jpg"}, {"黑桃J", 10, false, "jack_of_spades.jpg"}, {"黑桃Q", 10, false, "queen_of_spades.jpg"}, {"黑桃K", 10, false, "king_of_spades.jpg"},
		{"红心A", 11, false, "ace_of_hearts.jpg"}, {"红心2", 2, false, "2_of_hearts.jpg"}, {"红心3", 3, false, "3_of_hearts.jpg"}, {"红心4", 4, false, "4_of_hearts.jpg"}, {"红心5", 5, false, "5_of_hearts.jpg"}, {"红心6", 6, false, "6_of_hearts.jpg"}, {"红心7", 7, false, "7_of_hearts.jpg"}, {"红心8", 8, false, "8_of_hearts.jpg"}, {"红心9", 9, false, "9_of_hearts.jpg"}, {"红心10", 10, false, "10_of_hearts.jpg"}, {"红心J", 10, false, "jack_of_hearts.jpg"}, {"红心Q", 10, false, "queen_of_hearts.jpg"}, {"红心K", 10, false, "king_of_hearts.jpg"},
		{"樱花A", 11, false, "ace_of_clubs.jpg"}, {"樱花2", 2, false, "2_of_clubs.jpg"}, {"樱花3", 3, false, "3_of_clubs.jpg"}, {"樱花4", 4, false, "4_of_clubs.jpg"}, {"樱花5", 5, false, "5_of_clubs.jpg"}, {"樱花6", 6, false, "6_of_clubs.jpg"}, {"樱花7", 7, false, "7_of_clubs.jpg"}, {"樱花8", 8, false, "8_of_clubs.jpg"}, {"樱花9", 9, false, "9_of_clubs.jpg"}, {"樱花10", 10, false, "10_of_clubs.jpg"}, {"樱花J", 10, false, "jack_of_clubs.jpg"}, {"樱花Q", 10, false, "queen_of_clubs.jpg"}, {"樱花K", 10, false, "king_of_clubs.jpg"},
		{"方块A", 11, false, "ace_of_diamonds.jpg"}, {"方块2", 2, false, "2_of_diamonds.jpg"}, {"方块3", 3, false, "3_of_diamonds.jpg"}, {"方块4", 4, false, "4_of_diamonds.jpg"}, {"方块5", 5, false, "5_of_diamonds.jpg"}, {"方块6", 6, false, "6_of_diamonds.jpg"}, {"方块7", 7, false, "7_of_diamonds.jpg"}, {"方块8", 8, false, "8_of_diamonds.jpg"}, {"方块9", 9, false, "9_of_diamonds.jpg"}, {"方块10", 10, false, "10_of_diamonds.jpg"}, {"方块J", 10, false, "jack_of_diamonds.jpg"}, {"方块Q", 10, false, "queen_of_diamonds.jpg"}, {"方块K", 10, false, "king_of_diamonds.jpg"},
	}
}

func pushCard() (*card, error) {
	rand.Seed(uint64(time.Now().UnixNano()))
	for i := 0; i < len(cards); i++ {
		index := rand.Intn(len(cards))
		card := cards[index]
		if !card.used {
			card.used = true
			return card, nil
		}
	}
	StartGame()
	return pushCard()
}

type user struct {
	username        string
	cards           []*card
	currentPointNum int64
	finished        bool
}

var users = make(map[string]*user)

func StartGetCards(username string) (msg string, images []string) {
	_, ok := users[username]
	if ok {
		msg = username + ",您已开始游戏,请勿重复操作!"
		return
	}
	card1, err := pushCard()
	if err != nil {
		msg = err.Error()
		return
	}
	card2, err := pushCard()
	if err != nil {
		msg = err.Error()
		return
	}
	user := &user{
		username:        username,
		cards:           []*card{card1, card2},
		currentPointNum: card1.pointNum + card2.pointNum,
	}
	images = []string{card1.img, card2.img}
	users[username] = user
	msg = fmt.Sprintf("%s,您获得了两张牌, [%s] 和 [%s],您当前点数为%d点", username, card1.name, card2.name, user.currentPointNum)
	if user.currentPointNum == 21 {
		msg = fmt.Sprintf("%s\n\n恭喜您,天选之子!\n\n已自动为您停牌,请等待游戏结算!", msg)
		user.finished = true
	}
	return
}

func GetCard(username string) (msg string, image string, err error) {
	user, ok := users[username]
	if !ok {
		err = errors.New(username + ",请先开始游戏")
		return
	}
	if user.finished {
		err = errors.New(username + ",您当前已无法继续摸牌,请等待游戏结算!")
		return
	}
	var card *card
	card, err = pushCard()
	if err != nil {
		return
	}
	user.cards = append(user.cards, card)
	image = card.img
	user.calPointNum()
	if user.currentPointNum > 21 {
		user.finished = true
		msg = fmt.Sprintf(`
boom!!!! 
%s您获得了一张 [%s],
总点数为%d`, username, card.name, user.currentPointNum)
		return
	}
	currentCardStr := "您当前牌为 "
	for _, c := range user.cards {
		currentCardStr += fmt.Sprintf("[%s] ", c.name)
	}
	msg = fmt.Sprintf("%s,您获得了一张 [%s], %s,您当前点数为 [%d]", username, card.name, currentCardStr, user.currentPointNum)
	return
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
	defer Reset()
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
