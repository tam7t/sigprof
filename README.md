# sigprof
Golang user signal based package for collecting pprof information. Similar to
[net/http/pprof](https://golang.org/pkg/net/http/pprof/) but with the USR1 and
USR2 signal instead of an HTTP server.

# Usage
Link the package into your program:

```
import _ "github.com/tam7t/sigprof"
```

When you want to inspect a running process send the `USR1` or `USR2` signal.

```
kill -usr1 <golang process pid>
```

By default the profile information for USR1 is [goroutine](https://golang.org/pkg/runtime/pprof/#Profile)
and the information will be dumped to a file `profile-<current time>.prof`

```
go tool pprof <prof file>
```

# Configuration

Configuration is done through environment variables.

* `USR1_PROF` - Profile to lookup on the `USR1` signal. Default: `goroutine`
* `USR2_PROF` - Profile to lookup on the `USR2` signal. Default: `heap`
* `SIG_PROF_OUT` - Specify the output location, either `file`, `stderr`, or
  `stdout`. Default: `file`.
