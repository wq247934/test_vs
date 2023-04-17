package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
)

func main() {
	menu := []string{"Menu 1", "Menu 2"}
	submenu1 := []string{"Submenu 1-1", "Submenu 1-2"}
	submenu2 := []string{"Submenu 2-1", "Submenu 2-2"}

	fmt.Println("Menu:")
	for i, m := range menu {
		fmt.Printf("%d. %s\n", i+1, m)
	}
	prompt := promptui.Select{
		Label: "Enter your choice",
		Items: menu,
	}
	_, choice, err := prompt.Run()
	if err != nil {
		fmt.Println("Invalid choice")
		return
	}
	switch choice {
	case "Menu 1":
		prompt := promptui.Select{
			Label: "Enter your choice",
			Items: submenu1,
		}
		_, subChoice, err := prompt.Run()
		if err != nil {
			fmt.Println("Invalid choice")
			return
		}
		fmt.Println(subChoice)
	case "Menu 2":
		prompt := promptui.Select{
			Label: "Enter your choice",
			Items: submenu2,
		}
		_, subChoice, err := prompt.Run()
		if err != nil {
			fmt.Println("Invalid choice")
			return
		}
		fmt.Println(subChoice)
	default:
		fmt.Println("Invalid choice")
	}
}
