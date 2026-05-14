//go:build !windows

package passwd

import "syscall"
import "unsafe"

type termios struct {
	Iflag, Oflag, Cflag, Lflag uint32
	Line                       uint8
	Cc                         [19]uint8
	Ispeed, Ospeed             uint32
}

func Read() ([]byte, error) {
	const ECHO, ICANON, TCGETS, TCSETS = 0x8, 0x2, uintptr(0x5401), uintptr(0x5402)
	var old termios
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), TCGETS, uintptr(unsafe.Pointer(&old)))
	n := old
	n.Lflag &^= ECHO | ICANON
	n.Cc[6] = 1
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), TCSETS, uintptr(unsafe.Pointer(&n)))
	defer syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), TCSETS, uintptr(unsafe.Pointer(&old)))
	var result []byte
	buf := make([]byte, 1)
	for {
		cnt, err := syscall.Read(syscall.Stdin, buf)
		if err != nil || cnt == 0 || buf[0] == '\n' || buf[0] == '\r' {
			break
		}
		if buf[0] == 127 || buf[0] == 8 {
			if len(result) > 0 {
				result = result[:len(result)-1]
			}
			continue
		}
		result = append(result, buf[0])
	}
	return result, nil
}
