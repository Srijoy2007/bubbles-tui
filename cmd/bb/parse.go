package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Srijoy2007/bubbles-tui/internal/store"
)

func parseTime(s string) (int, error) {
	s = strings.TrimSpace(s)
	var h, m int
	if strings.Contains(s, ":") {
		if _, err := fmt.Sscanf(s, "%d:%d", &h, &m); err != nil {
			return 0, fmt.Errorf("bad time %q", s)
		}
	} else {
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("bad time %q", s)
		}
		h = n
	}
	if h < 0 || h > 24 || m < 0 || m > 59 {
		return 0, fmt.Errorf("time out of range: %q", s)
	}
	return h*60 + m, nil
}


func parseAdd(date, raw string) (store.Block, error) {
	raw = strings.TrimSpace(raw)
	parts := strings.SplitN(raw, " ", 2)
	if len(parts) < 2 {
		return store.Block{}, fmt.Errorf("usage: <start>-<end> <title>")
	}
	timePart, title := parts[0], strings.TrimSpace(parts[1])
	tp := strings.SplitN(timePart, "-", 2)
	if len(tp) != 2 {
		return store.Block{}, fmt.Errorf("usage: <start>-<end> <title>")
	}
	start, err := parseTime(tp[0])
	if err != nil {
		return store.Block{}, err
	}
	end, err := parseTime(tp[1])
	if err != nil {
		return store.Block{}, err
	}
	return store.Block{Date: date, Start: start, End: end, Title: title}, nil
}
