//go:build windows

package passwd

import (
	"syscall"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	getConsoleMode     = kernel32.NewProc("GetConsoleMode")
	setConsoleMode     = kernel32.NewProc("SetConsoleMode")
	readConsoleInputW  = kernel32.NewProc("ReadConsoleInputW")
)

const (
	enableEchoInput = 0x0004
	enableLineInput = 0x0002
)

type keyEventRecord struct {
	KeyDown         int32
	RepeatCount     uint16
	VirtualKeyCode  uint16
	VirtualScanCode uint16
	UnicodeChar     uint16
	ControlKeyState uint32
}

type inputRecord struct {
	EventType uint16
	_         [2]byte
	Event     [16]byte
}

func Read() ([]byte, error) {
	h := syscall.Handle(uintptr(syscall.Stdin))
	var old uint32
	getConsoleMode.Call(uintptr(h), uintptr(unsafe.Pointer(&old)))
	setConsoleMode.Call(uintptr(h), uintptr(old&^(enableEchoInput|enableLineInput)))
	defer setConsoleMode.Call(uintptr(h), uintptr(old))

	var result []byte
	var rec inputRecord
	var n uint32
	for {
		readConsoleInputW.Call(uintptr(h), uintptr(unsafe.Pointer(&rec)), 1, uintptr(unsafe.Pointer(&n)))
		if rec.EventType != 1 {
			continue
		}
		key := (*keyEventRecord)(unsafe.Pointer(&rec.Event))
		if key.KeyDown == 0 {
			continue
		}
		ch := key.UnicodeChar
		if ch == '\r' || ch == '\n' {
			break
		}
		if ch == 8 || ch == 127 {
			if len(result) > 0 {
				result = result[:len(result)-1]
			}
			continue
		}
		if ch == 0 {
			continue
		}
		if ch < 0x80 {
			result = append(result, byte(ch))
		} else if ch < 0x800 {
			result = append(result, byte(0xC0|(ch>>6)), byte(0x80|(ch&0x3F)))
		} else {
			result = append(result, byte(0xE0|(ch>>12)), byte(0x80|((ch>>6)&0x3F)), byte(0x80|(ch&0x3F)))
		}
	}
	return result, nil
}
