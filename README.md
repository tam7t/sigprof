# sigprof

[![Join the chat at https://gitter.im/tam7t/sigprof](https://badges.gitter.im/tam7t/sigprof.svg)](https://gitter.im/tam7t/sigprof?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
Golang package for inspecting running processes. Similar to [net/http/pprof](https://golang.org/pkg/net/http/pprof/) but using `USR1` and `USR2` signals instead of HTTP server routes.

# Usage
Link the package:

```go
import _ "github.com/tam7t/sigprof"
```

Send the `USR1` or `USR2` signal to inspect the process.

```bash
kill -usr1 <golang process pid>
```

The default `USR1` profile is [goroutine](https://golang.org/pkg/runtime/pprof/#Profile). By default, `sigprof` will save results to timestamped files.

```bash
go tool pprof profile-<timestamp>.prof
```

# Configuration

`sigprof` loads its configuration from the following environment variables.

* `USR1_PROF` - Profile executed on the `USR1` signal. Default: `goroutine`
* `USR2_PROF` - Profile executed on the `USR2` signal. Default: `heap`
* `SIG_PROF_OUT` - Specify the output location, either `file`, `stderr`, or
  `stdout`. Default: `file`.
