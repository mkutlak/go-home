package main

import "fmt"

func main() {
	colors := map[string]string{
		"red":   "#ff0000",
		"green": "#4b7f6c",
		"white": "#ffffff",
	}

	fmt.Println(colors)

	weird := make(map[string]string)

	weird["name"] = "typist"
	weird["surname"] = "normist"
	delete(weird, "name")

	fmt.Println(weird)

	printMap(colors)
}

func printMap(m map[string]string) {
	for color, hex := range m {
		fmt.Printf("color:%s - hex:%s\n", color, hex)
	}
}
