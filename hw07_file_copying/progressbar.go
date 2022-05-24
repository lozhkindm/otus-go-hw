package main

import (
	"fmt"
	"strings"
)

type ProgressBar struct {
	size    int64
	written int64
	width   int
}

func NewProgressBar(size int64, width int) *ProgressBar {
	return &ProgressBar{size: size, width: width}
}

func (pb *ProgressBar) Write(b []byte) (n int, err error) {
	var (
		percentage    int
		progressSigns int
	)

	pb.written++
	percentage = int(pb.written * 100 / pb.size)
	if percentage > 0 {
		progressSigns = int(float64(pb.width) / float64(100) * float64(percentage))
	}

	signs := strings.Repeat("#", progressSigns)
	empties := strings.Repeat("_", pb.width-progressSigns)
	fmt.Printf("\r[%s%s] %d%%", signs, empties, percentage)

	if pb.written == pb.size {
		fmt.Println()
	}
	return len(b), nil
}
