package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logMu      sync.RWMutex
	logBuffer  []string
	mainLogger *slog.Logger
)

type HubWriter struct {
	file io.Writer
	app  *App
}

func (hw *HubWriter) Write(p []byte) (n int, err error) {
	msg := string(p)
	logMu.Lock()
	logBuffer = append(logBuffer, msg)
	if len(logBuffer) > 100 {
		logBuffer = logBuffer[1:]
	}
	logMu.Unlock()
	if hw.file != nil {
		hw.file.Write(p)
	}
	os.Stdout.Write(p)
	if hw.app != nil && hw.app.ctx != nil {
		hw.app.BroadcastLog(msg)
	}
	return len(p), nil
}

func initLogging(app *App) {
	dir, err := os.UserConfigDir()
	var file io.Writer
	if err == nil {
		logDir := filepath.Join(dir, "wat")
		os.MkdirAll(logDir, 0700)
		if f, err := os.OpenFile(filepath.Join(logDir, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600); err == nil {
			file = f
		}
	}
	hw := &HubWriter{file: file, app: app}
	h := slog.NewJSONHandler(hw, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
			}
			return a
		},
	})
	mainLogger = slog.New(h)
	slog.SetDefault(mainLogger)
}

func getRecentLogs() []string {
	logMu.RLock()
	defer logMu.RUnlock()
	res := make([]string, len(logBuffer))
	copy(res, logBuffer)
	return res
}

func logInfo(msg string, args ...any) {
	if mainLogger != nil {
		mainLogger.Info(msg, args...)
	} else {
		fmt.Printf("[INFO] "+msg+"\n", args...)
	}
}

func logError(msg string, args ...any) {
	if mainLogger != nil {
		mainLogger.Error(msg, args...)
	} else {
		fmt.Fprintf(os.Stderr, "[ERROR] "+msg+"\n", args...)
	}
}

func (a *App) GetLogs() ([]string, error) {
	return getRecentLogs(), nil
}
