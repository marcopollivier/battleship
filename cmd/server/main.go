package main

import "fmt"


func main() {
	var board [10][10] string

	for i := 0; i<10; i++ {
		for j := 0; j<10; j++ {
			board[i][j] = "~"
		}
	}

	fmt.Println("   A B C D E F G H I J")
	for i := 0; i<10; i++ {
		label := fmt.Sprintf("%2d ", i+1)
		fmt.Print(label)
		for j := 0; j<10; j++ {
			fmt.Print(board[i][j] + " ")
		}
		fmt.Println()
	}

}
