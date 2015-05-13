FROM google/golang:1.3

# get dependency
RUN go get github.com/tam7t/sigprof

COPY . ${GOPATH}/src/github.com/tam7t/sigprof/example
WORKDIR ${GOPATH}/src/github.com/tam7t/sigprof/example

# run binary
RUN go run main.go

# Build & run the docker container
#   ~ $ docker build .
#
# Show running container
#   ~ $ docker ps
#   CONTAINER ID        IMAGE                                                                     COMMAND                CREATED             STATUS              PORTS                               NAMES
#   046bb6bd1e45        ae36a05487424e97637e18f1c29c3fb0370b0fd1a0b988616cfa91af10af280e:latest   "/bin/sh -c 'go run    11 seconds ago      Up 10 seconds                                           stoic_bartik
#
# Enter docker container shell
#   ~ $ docker exec -it 046bb6bd1e45 /bin/bash
#   root@ece227541990:/gopath/src/github.com/tam7t/sigprof/example#
#
# Find the pid of the process by iterating over /proc/<pid> (usually pid 1)
#   root@ece227541990:/gopath/src/github.com/tam7t/sigprof/example# echo "$(cat /proc/1/cmdline)"
#   /bin/sh-cgo run main.go
#
# Here pid 14 is the one we want since we used go run instead of directly running executable.
#   root@ece227541990:/gopath/src/github.com/tam7t/sigprof/example# echo "$(cat /proc/14/cmdline)"
#   /tmp/go-build163951639/command-line-arguments/_obj/exe/main
#
# Lets send the signal
#   root@ece227541990:/gopath/src/github.com/tam7t/sigprof/example# ls
#   Dockerfile  main  main.go
#   root@ece227541990:/gopath/src/github.com/tam7t/sigprof/example# kill -usr1 14
#   root@ece227541990:/gopath/src/github.com/tam7t/sigprof/example# ls
#   Dockerfile  main  main.go  profile-2015-05-12 14:43:08.72195103 +0000 UTC.prof
#
# And examine the file
#   root@ece227541990:/gopath/src/github.com/tam7t/sigprof/example# cat "profile-2015-05-12 14:43:08.72195103 +0000 UTC.prof"
#   goroutine profile: total 9
#   1 @ 0x45f26c 0x45f01b 0x45bfc3 0x42ae64 0x42abb1 0x413460
#   # 0x45f26c  runtime/pprof.writeRuntimeProfile+0xcc  /usr/local/go/src/pkg/runtime/pprof/pprof.go:540
#   # 0x45f01b  runtime/pprof.writeGoroutine+0x9b /usr/local/go/src/pkg/runtime/pprof/pprof.go:502
#   # 0x45bfc3  runtime/pprof.(*Profile).WriteTo+0xd3 /usr/local/go/src/pkg/runtime/pprof/pprof.go:229
#   # 0x42ae64  github.com/tam7t/sigprof.lookup+0x144 /gopath/src/github.com/tam7t/sigprof/sigprof.go:106
#   # 0x42abb1  github.com/tam7t/sigprof.profile+0x211  /gopath/src/github.com/tam7t/sigprof/sigprof.go:92
#
#   1 @ 0x4131c9 0x41cf59 0x400cb3 0x410c9a 0x413460
#   # 0x400cb3  main.main+0xb3    /gopath/src/github.com/tam7t/sigprof/example/main.go:32
#   # 0x410c9a  runtime.main+0x11a  /usr/local/go/src/pkg/runtime/proc.c:247
#
#   1 @ 0x4047f6 0x40c943 0x413460
#   # 0x4047f6  runtime.notetsleepg+0x46  /usr/local/go/src/pkg/runtime/lock_futex.c:198
#   # 0x40c943  runtime.MHeap_Scavenger+0xa3  /usr/local/go/src/pkg/runtime/mheap.c:532
#
#   1 @ 0x4131c9 0x41324b 0x4095ea 0x413460
#   # 0x4095ea  bgsweep+0x9a  /usr/local/go/src/pkg/runtime/mgc0.c:1993
#
#   1 @ 0x4131c9 0x41324b 0x40a90f 0x413460
#   # 0x40a90f  runfinq+0xcf  /usr/local/go/src/pkg/runtime/mgc0.c:2644
#
#   1 @ 0x4047f6 0x424e2a 0x45b3be 0x413460
#   # 0x45b3be  os/signal.loop+0x1e /usr/local/go/src/pkg/os/signal/signal_unix.go:21

#   1 @ 0x4131c9 0x41324b 0x41cc98 0x41cdc8 0x400d2d 0x413460
#   # 0x41cc98  chanrecv+0x4e8    /usr/local/go/src/pkg/runtime/chan.goc:268
#   # 0x41cdc8  runtime.chanrecv1+0x38  /usr/local/go/src/pkg/runtime/chan.goc:352
#   # 0x400d2d  main.func·001+0x6d  /gopath/src/github.com/tam7t/sigprof/example/main.go:17
#
#   1 @ 0x4131c9 0x41324b 0x42706a 0x426f61 0x400e0f 0x413460
#   # 0x426f61  time.Sleep+0x31   /usr/local/go/src/pkg/runtime/time.goc:39
#   # 0x400e0f  main.func·002+0x2f  /gopath/src/github.com/tam7t/sigprof/example/main.go:26
#
#   1 @ 0x4047f6 0x42730a 0x413460
#   # 0x42730a  timerproc+0xda  /usr/local/go/src/pkg/runtime/time.goc:260

