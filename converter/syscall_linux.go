package converter

import "syscall"

func getSysAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}
