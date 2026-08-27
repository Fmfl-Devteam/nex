//go:build linux

package process

import (
	"os"
	"strconv"
	"syscall"
)

func fileUID(info os.FileInfo) string {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return ""
	}
	return strconv.FormatUint(uint64(stat.Uid), 10)
}
