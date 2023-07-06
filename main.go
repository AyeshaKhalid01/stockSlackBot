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


const (
	baseURL = "https://www.alphavantage.co/query"
)

func stocks(input string) []string {
	var stockInfo []string
	flag.Parse()
	if len(input) == 0 {
		log.Fatalf("Input one stock symbol", os.Args[0])
	}

	apiKey := "" // Replace with your actual API key

	// Build the API URL
	url := fmt.Sprintf("%s?function=GLOBAL_QUOTE&symbol=%s&apikey=%s", baseURL, input, apiKey)

	// Send the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the JSON response
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("failed")
		log.Fatal(err)
	}

	// Check if the API call was successful
	if _, ok := result["Global Quote"]; !ok {
		log.Printf("Failed to fetch stock quote for symbol: %s", input)
		return []string{fmt.Sprintf("Failed to fetch stock quote for symbol, check your internet connection: %s", input)}
	}

	// Extract the relevant information from the response
	quote := result["Global Quote"].(map[string]interface{})
	// Check if the API call was successful
	if len(quote) == 0 {
		log.Printf("Are you sure you entered the correct stock symbol: %s", input)
		return []string{fmt.Sprintf("Are you sure you entered the correct stock symbol, cannot find symbol: %s", input)}
	}

	currentPrice := quote["05. price"].(string)
	highPrice := quote["03. high"].(string)
	lowPrice := quote["04. low"].(string)

	stockInfo = append(stockInfo, fmt.Sprintf("Symbol: %s\n", input))
	stockInfo = append(stockInfo, fmt.Sprintf("Current Price: $%s\n", currentPrice))
	stockInfo = append(stockInfo, fmt.Sprintf("52 Week High: $%s\n", highPrice))
	stockInfo = append(stockInfo, fmt.Sprintf("52 Week Low: $%s\n", lowPrice))
	return stockInfo
}

