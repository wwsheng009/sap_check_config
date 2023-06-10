package utils

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

var (
	moduser32                    = syscall.NewLazyDLL("user32.dll")
	procEnumWindows              = moduser32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = moduser32.NewProc("GetWindowThreadProcessId")
)

func getWindowThreadProcessId(hwnd syscall.Handle) (uint32, uint32) {
	var processId uint32
	r0, _, _ := syscall.SyscallN(procGetWindowThreadProcessId.Addr(), 2, uintptr(hwnd), uintptr(unsafe.Pointer(&processId)), 0)
	threadId := uint32(r0)

	return threadId, processId
}

func enumWindowsCallback(hwnd syscall.Handle, lParam uintptr) uintptr {
	pid := uintptr(lParam)
	_, processId := getWindowThreadProcessId(hwnd)
	if processId == uint32(pid) {
		fmt.Println(hwnd)
		// Return 0 to continue enumeration
		return 0
	}
	// Return 1 to stop enumeration
	return 1
}

func GetMainWindow(processId int) (syscall.Handle, error) {
	result, _, err := procEnumWindows.Call(
		uintptr(syscall.NewCallback(enumWindowsCallback)),
		uintptr(processId),
	)
	if result == 0 {
		return 0, err
	}
	return syscall.Handle(result), nil
}

func GetAppPath(app string) string {
	if app == "" {
		return ""
	}
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\`+app, registry.QUERY_VALUE)
	loc := ""
	if err == nil {
		// Check whether the "InstallLocation" value exists in the registry key
		loc, _, _ = key.GetStringValue("")
	}
	defer key.Close()
	if _, err := os.Stat(loc); err == nil {
		return loc
	}
	return ""
}
