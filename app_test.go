package main

import (
	"testing"
)

func TestGetSystemStats(t *testing.T) {
	app := NewApp()
	stats, err := app.GetSystemStats()
	if err != nil {
		t.Fatal(err)
	}
	if stats.CPUPercent < 0 || stats.CPUPercent > 100 {
		t.Errorf("invalid CPU percent: %f", stats.CPUPercent)
	}
}

func TestRankDirectory(t *testing.T) {
	app := NewApp()
	res, err := app.RankDirectory(".")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) == 0 {
		t.Error("expected non-empty directory ranking")
	}
}
