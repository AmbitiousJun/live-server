package ratelimits_test

import (
	"log"
	"testing"
	"time"

	"github.com/AmbitiousJun/live-server/internal/util/ratelimits"
)

func TestNewBucket(t *testing.T) {
	b := ratelimits.NewBucket(1, time.Second*10, 3)
	defer b.Destroy()
	for i := 1; i <= 10; i++ {
		b.Consume(1)
		log.Println(i)
		if i == 1 {
			time.Sleep(time.Second * 40)
		}
	}
}
