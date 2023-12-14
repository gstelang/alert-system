package main

import (
	"context"
	"fmt"
	"encoding/json"
	"github.com/gstelang/alert-system/alerts"
	"time"
	// "reflect"
)

type AlertState string
const (
	StatePass AlertState = "PASS"
	StateWarn AlertState = "WARN"
	StateCritical AlertState = "CRITICAL"
)

const (
	PassThreshold = 50
	WarnThreshold = 100
	CriticalThreshold = 200
)

const MaxThresholdVal = 200

type alertInfo struct {
	times int
	lastNotifyTime time.Time
	repeatIntervalSec int
	state AlertState
}

var alertMap = make(map[string]alertInfo)

func printJson(alert *alerts.Alert) {
	e, err := json.Marshal(alert)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(string(e))
}

var client alerts.Client
func getClientInstance() alerts.Client {
	if client == nil {
		client = alerts.NewClient("")
	}
	return client
}

func passBeteenStates(passCh <-chan string, warnCh <-chan string, criticalCh <-chan string) {
	for {
		select {
		case val := <- passCh:
			fmt.Println("pass channel", val)
		case val := <- warnCh:
			fmt.Println("warn channel", val)
		case val := <- criticalCh:
			fmt.Println("critical channel", val)
		default:
			// fmt.Println("do nothing")
		}
	}
}

func getAlertState(val float32) AlertState {
	if val >= CriticalThreshold {
		return StateCritical
	} else if val >= WarnThreshold {
		return StateWarn
	} else if val >= PassThreshold {
		return StatePass
	} else {
		// TODO: temporary to satisfy go compiler 
		// update as per business logic
		return StatePass
	}
}

func processAlerts(alerts []*alerts.Alert) {

	passCh := make(chan string)
	warnCh := make(chan string)
	criticalCh := make(chan string)

	fmt.Println(alertMap)
	go passBeteenStates(passCh, warnCh, criticalCh)

	for _, alert := range alerts {
		// printJson(alert)
		// fmt.Printf("warnVal = %T\n", warnVal)
		// fmt.Printf("criticalVal = %T\n", criticalVal)
		criticalVal := alert.Critical.Value
		warnVal := alert.Warn.Value
		alertName := alert.Name

		if criticalVal == MaxThresholdVal { // TODO: adjust  
			name := alert.Name
			msg := "Critical threshold exceeded"

			// Add to alertMap

			entry, ok := alertMap[alertName]

			if ok {
				entry.times = entry.times + 1
			} else {
				entry = alertInfo{
					times: 1,
					repeatIntervalSec: alert.RepeatIntervalSecs,
					lastNotifyTime: time.Time{},
				}	
			}
			// change it to warnVal
			entry.state = getAlertState(criticalVal)

			hasTimeElapsed := false
			if  entry.lastNotifyTime.IsZero() == false {
				duration :=  time.Duration(entry.repeatIntervalSec) * time.Second
				totalTime := entry.lastNotifyTime.Add(duration)
				hasTimeElapsed = totalTime.Before(time.Now())
			}
			// if state and timeElapsed but first time notification
			if (entry.state == StateCritical && (hasTimeElapsed ||  entry.lastNotifyTime.IsZero())) {
				fmt.Println("*********** after repeat interval **********")
				notify(name, msg)
				entry.lastNotifyTime = time.Now()
			}
			// reassign the copy
			alertMap[alertName] = entry
		}

		if warnVal >= CriticalThreshold {
			fmt.Println("write to critical channel")
			criticalCh <- alertName
		} else if warnVal >= WarnThreshold {
			fmt.Println("write to warn channel")
			warnCh <- alertName
		} else if warnVal >= PassThreshold {
			fmt.Println("write to pass channel")
			passCh <- alertName
		} else {
			// resolve
		}

	}
}

func queryAlerts(ctx context.Context) {
	alertClient := getClientInstance()
	alerts, err := alertClient.QueryAlerts(context.Background())
	if err != nil {
		fmt.Printf("error querying alerts: %+v\n", err)
	} else {
		fmt.Println("here - QueryAlerts")
		processAlerts(alerts)
	}
}

// Queries a particular one.
func queryByName(name string) {
	alertClient := getClientInstance()
	value, err := alertClient.Query(context.Background(), name)
	if err != nil {
		fmt.Printf("error resolving: %v\n", err)
	} else {
		fmt.Println("query - ", name)
		fmt.Printf("value queried: %v\n", value)
	}
}

func notify(name string, msg string) {
	alertClient := getClientInstance()

	err := alertClient.Notify(context.Background(), name, msg)
	if err != nil {
		fmt.Println("notify - ", name)
		fmt.Printf("error notifying: %v\n", err)
	} else {
		fmt.Println("notified - ", name)
	}
}

func resolve(name string) {
	alertClient := getClientInstance()
	err := alertClient.Resolve(context.Background(), name)
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