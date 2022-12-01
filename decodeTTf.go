package main

import (
	"fmt"
	"os"
	font "font/fontLib"
)


func main() {

	argNum := len(os.Args)

	if argNum < 2 {
		fmt.Printf("usage is:\n%s [ttf file]\n", os.Args[0])
		os.Exit(-1)
	}

	fontobj, err := font.GetFont(os.Args[1])
	if err != nil {
		fmt.Printf("error GetFont: %v\n", err)
		os.Exit(-1)
	}

	fmt.Println("**** success *****")
	fmt.Printf("font: %v\n", fontobj)
}
