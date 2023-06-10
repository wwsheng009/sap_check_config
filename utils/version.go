package utils

import (
	"errors"
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

type VS_FIXEDFILEINFO struct {
	Signature        uint32
	StrucVersion     uint32
	FileVersionMS    uint32
	FileVersionLS    uint32
	ProductVersionMS uint32
	ProductVersionLS uint32
	FileFlagsMask    uint32
	FileFlags        uint32
	FileOS           uint32
	FileType         uint32
	FileSubtype      uint32
	FileDateMS       uint32
	FileDateLS       uint32
}

type WinVersion struct {
	Major uint32
	Minor uint32
	Patch uint32
	Build uint32
}

// FileVersion concatenates FileVersionMS and FileVersionLS to a uint64 value.
func (fi VS_FIXEDFILEINFO) FileVersion() uint64 {
	return uint64(fi.FileVersionMS)<<32 | uint64(fi.FileVersionLS)
}

// VerQueryValueRoot calls VerQueryValue
// (https://msdn.microsoft.com/en-us/library/windows/desktop/ms647464(v=vs.85).aspx)
// with `\` (root) to retieve the VS_FIXEDFILEINFO.
func VerQueryValueRoot(block []byte) (VS_FIXEDFILEINFO, error) {
	var offset uintptr
	var length uint32
	blockStart := unsafe.Pointer(&block[0])
	err := windows.VerQueryValue(blockStart, `\`, unsafe.Pointer(&offset), &length)
	if err != nil {
		return VS_FIXEDFILEINFO{}, nil
	}
	// ret, _, _ := verQueryValue.Call(
	// 	uintptr(blockStart),
	// 	uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(`\`))),
	// 	uintptr(unsafe.Pointer(&offset)),
	// 	uintptr(unsafe.Pointer(&length)),
	// )
	// if ret == 0 {
	// 	return VS_FIXEDFILEINFO{}, errors.New("VerQueryValueRoot: verQueryValue failed")
	// }
	start := int(offset) - int(uintptr(blockStart))
	end := start + int(length)
	if start < 0 || start >= len(block) || end < start || end > len(block) {
		return VS_FIXEDFILEINFO{}, errors.New("VerQueryValueRoot: find failed")
	}
	data := block[start:end]
	info := *((*VS_FIXEDFILEINFO)(unsafe.Pointer(&data[0])))
	return info, nil
}

// https://github.com/keybase/client/blob/master/go/install/winversion.go
func GetExeVersionInfo(exe string) (WinVersion, error) {
	var result WinVersion
	// get file version info size
	var size uint32
	var handler windows.Handle

	_, err := os.Stat(exe)
	if os.IsNotExist(err) {
		return result, err
	}

	// buffer1 := make([]byte, 100)
	// syscall.GetFileVersionInfoSize(syscall.StringToUTF16Ptr(filename), &size)
	size, err = windows.GetFileVersionInfoSize(exe, &handler)
	if err != nil {
		return result, err
	}
	info := make([]byte, size)

	err = windows.GetFileVersionInfo(exe, 0, size, unsafe.Pointer(&info[0]))
	if err != nil {
		return result, err
	}

	fixed, err := VerQueryValueRoot(info)
	if err != nil {
		return result, err
	}
	version := fixed.FileVersion()

	result.Major = uint32(version & 0xFFFF000000000000 >> 48)
	result.Minor = uint32(version & 0x0000FFFF00000000 >> 32)
	result.Patch = uint32(version & 0x00000000FFFF0000 >> 16)
	result.Build = uint32(version & 0x000000000000FFFF)

	return result, nil
}

func GetExeVersion(exe string) (string, error) {
	if exe == "" {
		return "", errors.New("exe文件安装路径不能为空")
	}
	ver, err := GetExeVersionInfo(exe)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d.%d.%d", int(ver.Major), int(ver.Minor), int(ver.Patch), int(ver.Build)), nil
}
