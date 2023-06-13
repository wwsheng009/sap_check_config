package utils

import (
	"log"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const processEntrySize = 568

// const PROCESS_ALL_ACCESS = 0x1F0FFF

func CheckProcessIsRunning(process string) (bool, windows.HWND, error) {
	h, e := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if e != nil {
		return false, 0, e
	}
	defer windows.CloseHandle(h)
	p := windows.ProcessEntry32{Size: processEntrySize}

	e1 := windows.Process32First(h, &p)

	for e1 == nil {

		// if e != nil {
		// 	return false, nil
		// }
		s := windows.UTF16ToString(p.ExeFile[:])
		// hProcess, e := windows.OpenProcess(PROCESS_ALL_ACCESS, false, p.ProcessID)
		// if e == nil {
		// 	ret, e := windows.GetPriorityClass(hProcess)
		// 	if e == nil {
		// 		windows.CloseHandle(hProcess)
		// 	}
		// 	println("return classs", ret)

		// }
		//如果主窗口类存在，说明在运行
		if s == process {
			class := ""
			if process == "chrome.exe" || process == "msedge.exe" {
				class = "Chrome_WidgetWin_1"
			} else if process == "firefox.exe" {
				class = "MozillaWindowClass"
			}
			hwnd := getMainWindow(p.ProcessID, class)
			// println(hwnd)
			return true, hwnd, nil
		}
		e1 = windows.Process32Next(h, &p)
	}
	return false, 0, nil
}

// func getLowPartFromLPARAM(lParam uintptr) uint32 {
// 	// convert LPARAM to uintptr
// 	uip := uintptr(lParam)

// 	// extract low part using bitwise AND operation
// 	lowPart := uint32(uip & 0xffff)

//		return lowPart
//	}
//
//	func highPartFromLParam(lParam uintptr) uint32 {
//		return uint32((lParam >> 32) & 0xFFFFFFFF)
//	}
func getMainWindow(ProcessID uint32, windows_class string) windows.HWND {
	if windows_class == "" {
		return 0
	}
	var mainWindow windows.HWND
	windows.EnumWindows(windows.NewCallback(func(hwnd windows.HWND, lParam uintptr) uintptr {

		// pid := getLowPartFromLPARAM(lParam)
		var processId uint32
		windows.GetWindowThreadProcessId(hwnd, &processId)
		// fmt.Println(hwnd)
		className := make([]uint16, 256)
		windows.GetClassName(hwnd, &className[0], 256)
		class := windows.UTF16ToString(className[:])

		if processId == ProcessID {
			if class == windows_class { //"Chrome_WidgetWin_1" {
				mainWindow = hwnd
				return 0
			}
		}
		// Continue enumerating
		return 1
	}), unsafe.Pointer(&ProcessID))
	return mainWindow
}
func CheckProcessIsInTasklist(process_name string) (bool, error) {
	cmd := exec.Command("tasklist")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	processName := process_name
	if strings.Contains(string(output), processName) {
		return true, nil
	} else {
		return false, nil
	}
}
func KillProcess(process_name string) (bool, error) {
	log.Println("结束进程:", process_name)
	ok, err := CheckProcessIsInTasklist(process_name)
	if err != nil {
		return false, err
	}
	if ok {
		cmd := exec.Command("taskkill", "/F", "/IM", process_name)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: 0x08000000,
		}
		err := cmd.Run()
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}

	// fmt.Println("Process killed successfully.")
	return false, nil
}
