package main

import "fmt"

func if_else_example() {
	geneiric_number, lucky_nubmer := 10, 67
	var i int
	if geneiric_number > lucky_nubmer {
		i++
	} else if geneiric_number < lucky_nubmer {
		i--
	}
	fmt.Println("i is equal to", i)
}

func for_example() {
	var i int = 0
	for i < 10 {
		fmt.Println("i is equal to", i)
		i++
	}
}

func switch_case_example() {
	my_number := 5
	switch my_number {
	case 5:
		fmt.Println("my_number is 5")
	case 10:
		fmt.Println("my_number is 10")
		fallthrough
	case 15:
		fmt.Println("my_number is 15")
	case 20:
		fmt.Println("my_number is 20")
	default:
		fmt.Println("my_number is neither 5 nor 10 nor 15 nor 20")
	}
}
func main() {
	if_else_example()
	for_example()
	switch_case_example()
}
