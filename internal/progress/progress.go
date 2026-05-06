package progress

import (
	"fmt"
	"sync"
	"time"
)

type Progress struct {
	mu       sync.Mutex
	current  int
	total    int
	workers  int
	start    time.Time
	lastLine string
}

func New(workers int) *Progress {
	return &Progress{
		mu:      sync.Mutex{},
		current: 0,
		total:   0,
		workers: workers,
		start:   time.Now(),
	}
}

func (p *Progress) Update(current, total int) {
	p.mu.Lock()
	p.current = current
	p.total = total
	p.mu.Unlock()

	percentage := 0.0
	if total > 0 {
		percentage = float64(current) / float64(total) * 100
	}

	elapsed := time.Since(p.start).Seconds()
	var speed string
	if elapsed > 0 && current > 0 {
		rate := float64(current) / elapsed
		speed = fmt.Sprintf("%.1f file/s", rate)
	} else {
		speed = "0.0 file/s"
	}

	barWidth := 30
	filled := 0
	if total > 0 {
		filled = int(float64(barWidth) * float64(current) / float64(total))
	}
	if filled > barWidth {
		filled = barWidth
	}
	if filled < 0 {
		filled = 0
	}

	bar := make([]byte, barWidth)
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar[i] = '='
		} else if i == filled {
			bar[i] = '>'
		} else {
			bar[i] = ' '
		}
	}

	line := fmt.Sprintf("\r[%s] %.2f%% (%d/%d) - %s", string(bar), percentage, current, total, speed)
	fmt.Printf("%s", line)
}

func (p *Progress) Finish() {
	p.mu.Lock()
	current := p.current
	total := p.total
	p.mu.Unlock()

	percentage := 100.0
	if total > 0 {
		percentage = float64(current) / float64(total) * 100
	}

	barWidth := 30
	bar := ""
	for i := 0; i < barWidth; i++ {
		bar += "="
	}

	elapsed := time.Since(p.start).Seconds()
	var speed string
	if elapsed > 0 && current > 0 {
		rate := float64(current) / elapsed
		speed = fmt.Sprintf("%.1f file/s", rate)
	} else {
		speed = "0.0 file/s"
	}

	line := fmt.Sprintf("\r[%s] %.2f%% (%d/%d) - %s\n", bar, percentage, current, total, speed)
	fmt.Printf("\r%s", line)
}