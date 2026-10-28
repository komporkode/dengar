//go:build windows

package main

import (
	"syscall"
	"time"
	"unsafe"
)

var (
	modUser32           = syscall.NewLazyDLL("user32.dll")
	procFindWindowW      = modUser32.NewProc("FindWindowW")
	procSendMessageW     = modUser32.NewProc("SendMessageW")
	procLoadImageW       = modUser32.NewProc("LoadImageW")
	procSetClassLongPtrW = modUser32.NewProc("SetClassLongPtrW")

	modKernel32          = syscall.NewLazyDLL("kernel32.dll")
	procGetModuleHandleW = modKernel32.NewProc("GetModuleHandleW")

	modShell32          = syscall.NewLazyDLL("shell32.dll")
	procSHChangeNotify  = modShell32.NewProc("SHChangeNotify")
)

const (
	wmSetIcon = 0x0080
	iconSmall = 0
	iconBig   = 1
	imageIcon = 1
	lrDefault = 0x0000

	shcneAssocChanged = 0x08000000
	shcnfIdList       = 0x0000
)

var (
	gclpHIcon   int32 = -14
	gclpHIconSm int32 = -34
)

// setupWindowsAppIcon guarantees that Windows Taskbar, Alt-Tab, and window title
// immediately receive the rich crimson gradient icon for both ICON_BIG and ICON_SMALL.
func setupWindowsAppIcon() {
	go func() {
		className, err := syscall.UTF16PtrFromString("wailsWindow")
		if err != nil {
			return
		}

		for i := 0; i < 30; i++ {
			time.Sleep(80 * time.Millisecond)

			hwnd, _, _ := procFindWindowW.Call(uintptr(unsafe.Pointer(className)), 0)
			if hwnd == 0 {
				continue
			}

			hInstance, _, _ := procGetModuleHandleW.Call(0)
			if hInstance == 0 {
				return
			}

			// Try Resource ID 3 (AppIconID used by Wails)
			hIconBig, _, _ := procLoadImageW.Call(
				hInstance,
				uintptr(3),
				imageIcon,
				32, 32,
				lrDefault,
			)

			hIconSmall, _, _ := procLoadImageW.Call(
				hInstance,
				uintptr(3),
				imageIcon,
				16, 16,
				lrDefault,
			)

			// Fallback to Resource ID 1 if ID 3 not found
			if hIconBig == 0 {
				hIconBig, _, _ = procLoadImageW.Call(
					hInstance,
					uintptr(1),
					imageIcon,
					32, 32,
					lrDefault,
				)
			}
			if hIconSmall == 0 {
				hIconSmall, _, _ = procLoadImageW.Call(
					hInstance,
					uintptr(1),
					imageIcon,
					16, 16,
					lrDefault,
				)
			}

			// Set ICON_BIG (crucial for Windows Taskbar button & Alt+Tab)
			if hIconBig != 0 {
				procSendMessageW.Call(hwnd, wmSetIcon, iconBig, hIconBig)
				procSetClassLongPtrW.Call(hwnd, uintptr(gclpHIcon), hIconBig)
			}

			// Set ICON_SMALL (for titlebar and small taskbar icons)
			if hIconSmall != 0 {
				procSendMessageW.Call(hwnd, wmSetIcon, iconSmall, hIconSmall)
				procSetClassLongPtrW.Call(hwnd, uintptr(gclpHIconSm), hIconSmall)
			}

			// Notify Windows Shell to refresh any cached taskbar icons
			procSHChangeNotify.Call(shcneAssocChanged, shcnfIdList, 0, 0)
			return
		}
	}()
}
