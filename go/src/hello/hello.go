package main

import "fmt"
import "example.com/greetings"

func main() {
	message := greetings.Hello("Gladys")
	fmt.Println(message)
}

