package utils

import (
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

const thisFile string = "logging.go"

func LogStart(parent *log.Logger, method string) (log.Logger, func(time.Time)) {
	logger := log.With(*parent, "method", method)
	level.Info(logger).Log("msg", "Starting method")

	return logger, func(start time.Time) {
		level.Info(logger).Log("msg", "Method completed.", "took", time.Since(start))
	}
}

func Caller(depth int) log.Valuer {
	return func() interface{} {
		mydepth := depth
		for {
			_, file, line, _ := runtime.Caller(mydepth)
			idx := strings.LastIndexByte(file, '/')
			// using idx+1 below handles both of following cases:
			// idx == -1 because no "/" was found, or
			// idx >= 0 and we want to start at the character after the found "/".
			if file[idx+1:] != thisFile {
				return file[idx+1:] + ":" + strconv.Itoa(line)
			}

			// If caller is this file, try the caller of the caller
			mydepth = mydepth + 1
		}
	}
}
