package _21point

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	StartGame()
	fmt.Println(StartGetCards("wang"))
	fmt.Println(StartGetCards("张"))
	fmt.Println(StartGetCards("李"))
	fmt.Println(GetCard("wang"))
	fmt.Println(GetCard("wang"))

	fmt.Println(Stop("wang"))
	fmt.Println(SettleGame())
	fmt.Println(GetCard("张"))

	fmt.Println(Stop("张"))
	fmt.Println(SettleGame())
	fmt.Println(GetCard("李"))

	fmt.Println(Stop("李"))
	fmt.Println(SettleGame())
}
