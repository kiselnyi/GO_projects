package main

import "fmt"

func main() {
	var my_text, my_funny_text = "Hello World!", "ha-ha, six-seven!"
	fmt.Print(my_text, "\n")
	var my_number = 67
	fmt.Println(my_number, my_funny_text)
	var my_64_float float64 = 67.67676767
	fmt.Println(my_64_float, my_64_float+my_64_float)
	new_type_for_variable := byte(2)
	fmt.Println(new_type_for_variable)
	fmt.Println("And now I will try to print my variable with an %%T format specifier")
	fmt.Printf("var type for new_type_for_variable: %T\n", new_type_for_variable, "\n", "looks like it is a byte type", "\n")
	fmt.Println("looks like this text have some artifacts, i will try to print it in other way!")
	fmt.Printf("var type for new_type_for_variable: %T\nlooks like it is a byte type\n", new_type_for_variable)
	fmt.Println("Also, it will be nice to try to print something in different language, for example in Russian:\nПривет мир!")
}
