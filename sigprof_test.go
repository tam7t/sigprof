package sigprof

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"testing"
)

func setup() func() {
	stop()

	origNewSigChan := newSigChan
	origNewWriter := newWriter
	origNewProfiler := newProfiler

	origUsr1 := os.Getenv(`SIGPROF_USR1`)
	origUsr2 := os.Getenv(`SIGPROF_USR2`)
	origOut := os.Getenv(`SIGPROF_OUT`)

	return func() {
		newSigChan = origNewSigChan
		newWriter = origNewWriter
		newProfiler = origNewProfiler

		mustPutenv(`SIGPROF_USR1`, origUsr1)
		mustPutenv(`SIGPROF_USR2`, origUsr2)
		mustPutenv(`SIGPROF_OUT`, origOut)
	}
}

func mustPutenv(key, value string) {
	var err error
	if value == "" {
		err = os.Unsetenv(key)
	} else {
		err = os.Setenv(key, value)
	}
	if err != nil {
		panic(err)
	}
}

type bufferCloser struct {
	*bytes.Buffer
}

func (bufferCloser) Close() error { return nil }

type testProfiler struct{}

func (testProfiler) writeProfile(w io.Writer, profileName string) error {
	fmt.Fprintf(w, "test %s\n", profileName)
	return nil
}

func TestStubs(t *testing.T) {
	cleanup := setup()
	defer cleanup()

	// Send three signals
	newSigChan = func() <-chan (os.Signal) {
		c := make(chan os.Signal)
		go func() {
			c <- syscall.SIGUSR1
			c <- syscall.SIGUSR2
			c <- syscall.SIGHUP
			close(c)
		}()
		return c
	}

	outputs := map[string]*bytes.Buffer{}
	newWriter = func(profile string, out outputType) io.WriteCloser {
		if out != "orange" {
			t.Fatalf("unexpected output %q", out)
		}
		var buf bytes.Buffer
		outputs[profile] = &buf
		return bufferCloser{&buf}
	}

	newProfiler = func() profiler {
		return testProfiler{}
	}

	s := sigprof{
		usr1: []string{"foo", "bar"},
		usr2: []string{"baz", "quux"},
	}
	s.output = "orange"

	s.loop()

	if len(outputs) != 4 {
		t.Errorf("unexpected outputs len=%d", len(outputs))
	}

	for _, profile := range []string{"foo", "bar", "baz", "quux"} {
		buf, ok := outputs[profile]
		if !ok {
			t.Errorf("missing expected profile %q", profile)
		}
		if buf.String() != "test "+profile+"\n" {
			t.Errorf("unexpected profiler contents: %q", buf.String())
		}
	}
}

func TestPprof(t *testing.T) {
	cleanup := setup()
	defer cleanup()

	s := sigprof{
		usr1:   []string{"goroutine"},
		usr2:   []string{"heap"},
		output: "file",
	}

	outputs := []*bytes.Buffer{}
	newWriter = func(profile string, out outputType) io.WriteCloser {
		var buf bytes.Buffer
		outputs = append(outputs, &buf)
		return bufferCloser{&buf}
	}

	s.profileSignal(syscall.SIGUSR1)
	s.profileSignal(syscall.SIGUSR2)

	if len(outputs) != 2 {
		t.Errorf("unexpected number of profiles: %d", len(outputs))
	}

	var hasHeap, hasGoroutine bool
	for _, output := range outputs {
		if strings.Contains(output.String(), "goroutine profile") {
			hasGoroutine = true
		} else if strings.Contains(output.String(), "heap profile") {
			hasHeap = true
		}
	}
	if !hasGoroutine {
		t.Error("missing goroutine profile")
	}
	if !hasHeap {
		t.Error("missing heap profile")
	}
}
