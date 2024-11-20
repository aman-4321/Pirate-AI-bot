package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/henomis/lingoose/llm/openai"
	"github.com/henomis/lingoose/thread"
)

func main() {
	reset := flag.Bool("reset", false, "a boolean for resetting the database")
	flag.Parse()

	if *reset {
		err := deleteKey()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error deleting key from db")
			os.Exit(1)
		}
		fmt.Println("API Key Successfully deleted")
		os.Exit(0)
	}
	fmt.Println("Welcome to Pirate Bot")
	fmt.Println("---------------------")
	err := getOrSetKey()
	if err != nil {
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter Text :")
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	myThread := thread.New().AddMessages(
		thread.NewSystemMessage().AddContent(
			thread.NewTextContent("All replies must be given in a pirate style of speech"),
		),

		thread.NewUserMessage().AddContent(
			thread.NewTextContent(text),
		),
	)

	err = openai.New().Generate(context.Background(), myThread)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error: something went wrong. Please check you API Key & account")
		os.Exit(1)
	}

	fmt.Println("Pirate : ", myThread.LastMessage().Contents[0].AsString())
	os.Exit(0)
}
