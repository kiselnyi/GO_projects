package main
import "fmt"

func if_else_example(){
	geneiric_number, lucky_nubmer := 10, 67
	var i int
	if geneiric_number > lucky_nubmer{
		i++
	} else if geneiric_number < lucky_nubmer{
		i--
	}
	fmt.Println("i is equal to", i)
}

func for_example(){
	var i int = 0
	for i < 10{
		fmt.Println("i is equal to", i)
		i++
	}
}

func switch_case_example(){
	
}
func main() {
	if_else_example()
	for_example()
}
