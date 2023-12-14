package main

import (
	"context"
	"fmt"
	"encoding/json"
	"github.com/gstelang/alert-system/alerts"
	"time"
)

func printJson(alert *alerts.Alert) {
	e, err := json.Marshal(alert)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(string(e))
}


// Build an alert execution engine that:
// - 1. Queries for alerts using the query_alerts API and execute them at the specified interval
// - 2. Alerts will not change over time, so only need to be loaded once at start
// - 3. The basic alert engine will send notifications whenever it sees a value that exceeds the critical threshold.
// - 4. Add support for repeat intervals, so that if an alert is continuously firing it will only re-notify after the repeat interval.
// - 5. Incorporate using the warn threshold in the alerting engine - now an alert can go between states PASS <-> WARN <-> CRITICAL <-> PASS.

// const (
// 	StatePass = "PASS"
// 	StateWarn = "WARN"
// 	StateCritical = "CRITICAL"
// )

// const RepeatIntervals = 10

// TODO: make it singleton
// note: client := alerts.NewClient("") does not work outside.
var client = alerts.NewClient("")

func queryAlerts(ctx context.Context) {
	alerts, err := client.QueryAlerts(context.Background())
	if err != nil {
		fmt.Printf("error querying alerts: %+v\n", err)
	} else {
		fmt.Println("here - QueryAlerts")
		for _, alert := range alerts {
			printJson(alert)
			// fmt.Printf("%v\n", alert)	
		}
	}
}

// Queries a particular one.
func queryByName(name string) {
	value, err := client.Query(context.Background(), name)
	if err != nil {
		fmt.Printf("error resolving: %v\n", err)
	} else {
		fmt.Println("query - ", name)
		fmt.Printf("value queried: %v\n", value)
	}
}

func notify(name string, msg string) {
	err := client.Notify(context.Background(), name, msg)
	if err != nil {
		fmt.Println("notify - ", name)
		fmt.Printf("error notifying: %v\n", err)
	}
}

func resolve(name string) {
	err := client.Resolve(context.Background(), name)
	if err != nil {
		fmt.Println("resolve - ", name)
		fmt.Printf("error resolving: %v\n", err)
	}
}

func main() {
	const pollerInterval = 5 * time.Second
	alertPoller := alerts.DefaultPoller{}

	// Poll for notifications
	alertPoller.Poll(context.Background(), pollerInterval, queryAlerts)

	queryName := "test-query-1"
	alertName := "alert-1"
	notifyMessage := "test-message"

	// Query by name
	queryByName(queryName)

	// Notification
	notify(alertName, notifyMessage)

	// Resolve
	resolve(alertName)
}