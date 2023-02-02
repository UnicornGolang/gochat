package utils

import "time"

type TimerFunc func(interface{}) bool

func Timer(delay, tick time.Duration, fun TimerFunc, param interface{}) {
	go func() {
		if fun == nil {
			return
		}
		t := time.NewTimer(delay)
		for range t.C {
			if !fun(param) {
				return
			}
			t.Reset(tick)
		}
	}()
}
