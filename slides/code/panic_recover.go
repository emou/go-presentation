package main

import "fmt"

// START OMIT
func openPresent() {
	panic("В кутията има паяк")
}

func lookUnderTree() {
	defer func() {
		if r := recover(); r != nil { // HL
			fmt.Println(r)
			fmt.Println("Използваш чехъл")
		}
	}()
	fmt.Println("Търсиш под елхата")
	openPresent()
}

func main() {
	lookUnderTree()
	fmt.Println("Най-добрият подарък!")
}

// END OMIT
