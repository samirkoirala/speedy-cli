package animation

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

func StartSpinner(label string) func() {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	var stop int32
	go func() {
		i := 0
		for atomic.LoadInt32(&stop) == 0 {
			fmt.Printf("\r%s %s", frames[i%len(frames)], label)
			time.Sleep(90 * time.Millisecond)
			i++
		}
	}()
	return func() {
		atomic.StoreInt32(&stop, 1)
		time.Sleep(110 * time.Millisecond)
		fmt.Print("\r\033[K")
	}
}

func Bar(value float64, max float64) string {
	if max <= 0 {
		max = 100
	}
	blocks := int((value / max) * 12)
	if blocks < 1 {
		blocks = 1
	}
	if blocks > 12 {
		blocks = 12
	}
	return strings.Repeat("█", blocks)
}

func DualGraph(a string, av float64, b string, bv float64) string {
	return fmt.Sprintf("%s: %-12s %.2f\n%s: %-12s %.2f", a, Bar(av, 200), av, b, Bar(bv, 100), bv)
}

func Progress(pct int) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	filled := pct / 10
	return "[" + strings.Repeat("=", filled) + strings.Repeat(" ", 10-filled) + "]"
}

func RocketLaunch() {
	frames := []string{
		"   🚀\n   ||\n  /__\\",
		"\n   🚀\n   ||\n  /__\\",
		"\n\n   🚀\n   ||\n  /__\\",
	}
	for i := 0; i < 6; i++ {
		fmt.Print("\033[H\033[2J")
		fmt.Println(frames[i%len(frames)])
		time.Sleep(120 * time.Millisecond)
	}
	fmt.Print("\033[H\033[2J")
}
