package alerts

import (
	"fmt"
    "time"
	"context"
)

type fn func(context.Context)

type Poller interface {
	Poll(ctx context.Context, tickRate time.Duration, pollingFunc fn)
}

type DefaultPoller struct{}

func (p DefaultPoller) Poll(
	ctx context.Context,
	tickRate time.Duration,
	pollingFunc fn,
) {
	ticker := time.NewTicker(tickRate).C

	for {
		select {
			case <-ticker:
				fmt.Println("************************ Polled after 5 seconds***********************")
				pollingFunc(ctx)
			default:
				break
		}
	}
}

