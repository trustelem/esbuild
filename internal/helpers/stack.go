package helpers

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"
)

var dumpMutex sync.Mutex
var StackStripPrefix string

func Stack(skip int) string {
	var pcs [200]uintptr
	// add 2 to skip count, to skip frames for 'runtime.Callers()' and 'Stack()'
	n := runtime.Callers(skip+2, pcs[:])

	frames := runtime.CallersFrames(pcs[:n])

	buf := bytes.NewBuffer(nil)
	buf.Grow(4096)

	f, more := frames.Next()
	for more {
		file := strings.TrimPrefix(f.File, StackStripPrefix)
		fmt.Fprintf(buf, "%s\n\t%s:%d\n", f.Function, file, f.Line)

		f, more = frames.Next()
	}

	return buf.String()
}

func DumpStackOpts(w io.Writer, skip int, msg string) {
	dumpMutex.Lock()
	defer dumpMutex.Unlock()

	if msg != "" {
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		w.Write([]byte(msg))
	}

	w.Write([]byte(Stack(skip + 1)))
}

func DumpStack(w io.Writer, msg string) {
	DumpStackOpts(w, 1, msg)
}
