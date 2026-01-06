package main

import (
	"fmt"
	"rahu/lexer"
	"rahu/utils"
	"time"
)

func main() {
	start := time.Now()
	input := utils.ParseFile("test.py")
	l := lexer.New(input)

	for {
		tok := l.NextToken()

		if tok.Type == lexer.EOF {
			break
		}
	}

	duration := time.Since(start)
	fmt.Printf("lexed in %v\n", duration)
}
