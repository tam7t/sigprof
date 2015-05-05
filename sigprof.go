package sigprof

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"syscall"
	"time"
)

// Interface for writing files

type ProfileWriter interface {
	Writer() io.Writer
	Close() error
}

type StderrWriter struct{}

func (w StderrWriter) Writer() io.Writer {
	return os.Stderr
}

func (w StderrWriter) Close() error {
	return nil
}

type StdoutWriter struct{}

func (w StdoutWriter) Writer() io.Writer {
	return os.Stdout
}

func (w StdoutWriter) Close() error {
	return nil
}

type FileWriter struct {
	file os.File
}

func (w FileWriter) Writer() io.Writer {
	return &w.file
}

func (w FileWriter) Close() error {
	return w.file.Close()
}

// Initialize profiler from environment and set defaults

func init() {
	var usr1 []string
	var usr2 []string
	var out string

	usr1 = strings.Split(os.Getenv(`USR1_PROF`), ",")
	if len(usr1) == 0 {
		usr1 = append(usr1, "goroutine")
	}

	usr2 = strings.Split(os.Getenv(`USR2_PROF`), ",")
	if len(usr2) == 0 {
		usr2 = append(usr2, "heap")
	}

	out = os.Getenv(`SIG_PROF_OUT`)
	if out == "" {
		out = "file"
	}

	go profile(usr1, usr2, out)
}

// Wait for profile signal

func profile(usr1, usr2 []string, out string) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGUSR1, syscall.SIGUSR2)

	for {
		s := <-c
		switch s {
		case syscall.SIGUSR1:
			lookup(usr1, out)
		case syscall.SIGUSR2:
			lookup(usr2, out)
		}
	}
}

// Perform profile and write to `out` ProfileWriter

func lookup(profiles []string, out string) {
	for _, profileName := range profiles {
		w := writer(out)
		p := pprof.Lookup(profileName)
		if p != nil {
			p.WriteTo(w.Writer(), 1)
		}
		w.Close()
	}
}

// Build ProfileWriter to use for dumping the prof info

func writer(out string) ProfileWriter {
	switch out {
	case "file":
		file, err := os.Create(fmt.Sprintf("profile-%s.prof", time.Now()))
		if err != nil {
			return FileWriter{*file}
		} else {
			return StderrWriter{}
		}
	case "stdout":
		return StdoutWriter{}
	case "stderr":
		return StderrWriter{}
	default:
		return StderrWriter{}
	}
}
