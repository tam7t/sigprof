// Package sigprof provides signal-triggered profiling.
package sigprof

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"sync"
	"syscall"
	"time"
)

func init() {
	s := newSigprof()
	go s.loop()
}

var (
	stopMu   sync.Mutex
	stopChan = make(chan struct{})
)

type stderrWriter struct{}

// Write implements io.Writer.
func (w stderrWriter) Write(p []byte) (int, error) {
	return os.Stderr.Write(p)
}

// Close implements io.Closer.
func (w stderrWriter) Close() error {
	return nil
}

type stdoutWriter struct{}

// Write implements io.Writer.
func (w stdoutWriter) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

// Close implements io.Closer.
func (w stdoutWriter) Close() error {
	return nil
}

type outputType string

const (
	stdoutOutput = outputType("stdout")
	stderrOutput = outputType("stderr")
	fileOutput   = outputType("file")
)

type sigprof struct {
	usr1, usr2 []string
	output     outputType
}

func newSigprof() sigprof {
	s := sigprof{}

	usr1EnvStr := os.Getenv(`SIGPROF_USR1`)
	if usr1EnvStr == "" {
		usr1EnvStr = "goroutine"
	}
	s.usr1 = strings.Split(usr1EnvStr, ",")

	usr2EnvStr := os.Getenv(`SIGPROF_USR2`)
	if usr2EnvStr == "" {
		usr2EnvStr = "heap"
	}
	s.usr2 = strings.Split(usr2EnvStr, ",")

	output := os.Getenv(`SIGPROF_OUT`)
	if output == "" {
		output = "file"
	}
	s.output = outputType(output)

	return s
}

func stop() {
	stopMu.Lock()
	if stopChan != nil {
		close(stopChan)
		stopChan = nil
	}
	stopMu.Unlock()
}

// loop handles signals and writes profiles.
func (s *sigprof) loop() {
	c := newSigChan()
	for {
		select {
		case sig, ok := <-c:
			if !ok {
				return
			}
			s.profileSignal(sig)
		case _, ok := <-stopChan:
			if !ok {
				return
			}
		}
	}
}

var newSigChan = func() <-chan (os.Signal) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGUSR1, syscall.SIGUSR2)
	return c
}

// profileSignal writes the profiles for the given signal.
func (s *sigprof) profileSignal(sig os.Signal) {
	var profiles []string
	switch sig {
	case syscall.SIGUSR1:
		profiles = s.usr1
	case syscall.SIGUSR2:
		profiles = s.usr2
	default:
		return
	}

	for _, profile := range profiles {
		w := s.writer(profile)
		s.profile(profile, w)
	}
}

// writer returns an io.WriteCloser to where the profile should be written.
func (s *sigprof) writer(profile string) io.WriteCloser {
	return newWriter(profile, s.output)
}

var newWriter = func(profile string, output outputType) io.WriteCloser {
	switch output {
	case "file":
		file, err := os.Create(fmt.Sprintf("%s-%s.prof", profile, time.Now()))
		if err != nil {
			log.Println("failed to create file for %s profile: %v", profile, err)
			return stderrWriter{}
		} else {
			return file
		}
	case "stdout":
		return stdoutWriter{}
	case "stderr":
		return stderrWriter{}
	default:
		return stderrWriter{}
	}
}

type profiler interface {
	writeProfile(w io.Writer, profileName string) error
}

type pprofiler struct{}

func (pprofiler) writeProfile(w io.Writer, profileName string) error {
	p := pprof.Lookup(profileName)
	if p == nil {
		return fmt.Errorf("failed to lookup profile %q", profileName)
	}
	return p.WriteTo(w, 1)
}

var newProfiler = func() profiler {
	return pprofiler{}
}

func (s *sigprof) profile(profileName string, w io.WriteCloser) {
	defer w.Close()
	p := newProfiler()
	err := p.writeProfile(w, profileName)
	if err != nil {
		log.Printf("failed to write %s profile: %v", profileName, err)
	}
}
