package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/shomali11/slacker"
)

func main() {
	os.Setenv("SLACK_BOT_TOKEN", "")
	os.Setenv("SLACK_APP_TOKEN", "")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	var wg sync.WaitGroup
	wg.Add(1)

	bot.Command("ping", &slacker.CommandDefinition{
		Handler: func(ctx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			// Get the event text from the context
			eventText := ctx.Event().Text

			// Extract the arguments after the command
			args := strings.Split(eventText, " ")[1:]

			// Check if there are any arguments
			if len(args) > 0 {
				// Process the arguments here
				for _, arg := range args {
					result := stocks(arg)
					for _, i := range result {
						response.Reply(i)
					}
				}
			} else {
				response.Reply("No arguments specified.")
			}
		},
	})

	go func() {
		defer wg.Done()
		err := bot.Listen(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()
}

