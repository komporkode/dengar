//go:build windows

package main

import (
	"syscall"
	"time"
	"unsafe"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	modUser32                   = syscall.NewLazyDLL("user32.dll")
	procFindWindowW             = modUser32.NewProc("FindWindowW")
	procSendMessageW            = modUser32.NewProc("SendMessageW")
	procLoadImageW              = modUser32.NewProc("LoadImageW")
	procSetClassLongPtrW        = modUser32.NewProc("SetClassLongPtrW")
	procSetWindowLongPtrW       = modUser32.NewProc("SetWindowLongPtrW")
	procCallWindowProcW         = modUser32.NewProc("CallWindowProcW")
	procCreatePopupMenu         = modUser32.NewProc("CreatePopupMenu")
	procAppendMenuW             = modUser32.NewProc("AppendMenuW")
	procTrackPopupMenu          = modUser32.NewProc("TrackPopupMenu")
	procDestroyMenu             = modUser32.NewProc("DestroyMenu")
	procGetCursorPos            = modUser32.NewProc("GetCursorPos")
	procSetForegroundWindow     = modUser32.NewProc("SetForegroundWindow")
	procShowWindow              = modUser32.NewProc("ShowWindow")
	procPostMessageW            = modUser32.NewProc("PostMessageW")
	procCreateIconIndirect      = modUser32.NewProc("CreateIconIndirect")
	procDestroyIcon             = modUser32.NewProc("DestroyIcon")
	procRegisterWindowMessageW = modUser32.NewProc("RegisterWindowMessageW")

	modGdi32         = syscall.NewLazyDLL("gdi32.dll")
	procCreateBitmap = modGdi32.NewProc("CreateBitmap")
	procDeleteObject = modGdi32.NewProc("DeleteObject")

	modKernel32          = syscall.NewLazyDLL("kernel32.dll")
	procGetModuleHandleW = modKernel32.NewProc("GetModuleHandleW")

	modShell32           = syscall.NewLazyDLL("shell32.dll")
	procSHChangeNotify   = modShell32.NewProc("SHChangeNotify")
	procShellNotifyIconW = modShell32.NewProc("Shell_NotifyIconW")

	modOle32             = syscall.NewLazyDLL("ole32.dll")
	procCoInitialize     = modOle32.NewProc("CoInitialize")
	procCoCreateInstance = modOle32.NewProc("CoCreateInstance")
)

const (
	wmSetIcon = 0x0080
	iconSmall = 0
	iconBig   = 1
	imageIcon = 1
	lrDefault = 0x0000

	shcneAssocChanged = 0x08000000
	shcnfIdList       = 0x0000

	nimAdd    = 0x00000000
	nimDelete = 0x00000002

	nifMessage = 0x00000001
	nifIcon    = 0x00000002
	nifTip     = 0x00000004
	wmApp      = 0x8000

	swRestore      uintptr = 9
	mfString       uintptr = 0x00000000
	mfSeparator    uintptr = 0x00000800
	tpmReturnCmd   uintptr = 0x0100
	tpmRightButton uintptr = 0x0002

	wmLButtonUp     = 0x0202
	wmLButtonDblClk = 0x0203
	wmRButtonUp     = 0x0205
	wmContextMenu   = 0x007B

	wmCommand   = 0x0111
	thbnClicked = 0x1800

	thbIcon    = 0x0002
	thbTooltip = 0x0004
	thbFlags   = 0x0008

	thbfEnabled = 0x0000

	thumbBtnPrev      = 3001
	thumbBtnPlayPause = 3002
	thumbBtnNext      = 3003
)

var (
	gclpHIcon   int32 = -14
	gclpHIconSm int32 = -34
	gwlpWndProc int32 = -4
)

type point struct {
	X int32
	Y int32
}

type notifyIconDataW struct {
	cbSize           uint32
	hWnd             uintptr
	uID              uint32
	uFlags           uint32
	uCallbackMessage uint32
	hIcon            uintptr
	szTip            [128]uint16
}

type iconInfo struct {
	fIcon    int32
	xHotspot uint32
	yHotspot uint32
	hbmMask  uintptr
	hbmColor uintptr
}

type guid struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

var (
	clsidTaskbarList = guid{
		Data1: 0x56FDF344,
		Data2: 0xFD6D,
		Data3: 0x11D0,
		Data4: [8]byte{0x95, 0x8A, 0x00, 0x60, 0x97, 0xC9, 0xA0, 0x90},
	}
	iidITaskbarList3 = guid{
		Data1: 0xEA1AFB91,
		Data2: 0x9E28,
		Data3: 0x4B86,
		Data4: [8]byte{0x90, 0xE9, 0x9E, 0x9F, 0x8A, 0x5E, 0xEF, 0xAF},
	}
)

type taskbarList3 struct {
	lpVtbl *taskbarList3Vtbl
}

type taskbarList3Vtbl struct {
	QueryInterface        uintptr
	AddRef                uintptr
	Release               uintptr
	HrInit                uintptr
	AddTab                uintptr
	DeleteTab             uintptr
	ActivateTab           uintptr
	SetActiveAlt          uintptr
	MarkFullscreenWindow  uintptr
	SetProgressValue      uintptr
	SetProgressState      uintptr
	RegisterTab           uintptr
	UnregisterTab         uintptr
	SetTabOrder           uintptr
	SetTabActive          uintptr
	ThumbBarAddButtons    uintptr
	ThumbBarUpdateButtons uintptr
	ThumbBarSetImageList  uintptr
	SetOverlayIcon        uintptr
	SetThumbnailTooltip   uintptr
	SetThumbnailClip      uintptr
}

type thumbButton struct {
	dwMask  uint32
	iId     uint32
	iBitmap uint32
	_pad0   uint32
	hIcon   uintptr
	szTip   [260]uint16
	dwFlags uint32
	_pad1   uint32
}

var (
	globalHwnd          uintptr
	origWndProc         uintptr
	globalTaskbarList3  *taskbarList3
	taskbarButtonsAdded bool
	wmTaskbarBtnCreated uintptr
	wmInitThumbbar      uintptr = wmApp + 102

	hIconPrev  uintptr
	hIconPlay  uintptr
	hIconPause uintptr
	hIconNext  uintptr
)

// createMonoIcon generates a crisp 16x16 monochrome icon with transparent background and white glyph
func createMonoIcon(xorMask [16]uint16) uintptr {
	var bits [32]uint16
	for i := 0; i < 16; i++ {
		bits[i] = ^xorMask[i]  // AND mask: 0 where glyph is drawn, 1 where transparent
		bits[16+i] = xorMask[i] // XOR mask: 1 where glyph is drawn (white), 0 where transparent
	}

	hbm, _, _ := procCreateBitmap.Call(16, 32, 1, 1, uintptr(unsafe.Pointer(&bits[0])))
	if hbm == 0 {
		return 0
	}
	defer procDeleteObject.Call(hbm)

	var ii iconInfo
	ii.fIcon = 1
	ii.hbmMask = hbm
	ii.hbmColor = 0

	hIcon, _, _ := procCreateIconIndirect.Call(uintptr(unsafe.Pointer(&ii)))
	return hIcon
}

func initThumbnailIcons() {
	if hIconPlay != 0 {
		return
	}

	// 16x16 glyphs
	prevMask := [16]uint16{
		0x0000, 0x0000, 0x0000,
		0x3010, 0x3030, 0x3070, 0x30F0, 0x31F0,
		0x31F0, 0x30F0, 0x3070, 0x3030, 0x3010,
		0x0000, 0x0000, 0x0000,
	}
	playMask := [16]uint16{
		0x0000, 0x0000, 0x0000,
		0x0800, 0x0C00, 0x0E00, 0x0F00, 0x0F80,
		0x0F80, 0x0F00, 0x0E00, 0x0C00, 0x0800,
		0x0000, 0x0000, 0x0000,
	}
	pauseMask := [16]uint16{
		0x0000, 0x0000, 0x0000,
		0x0E70, 0x0E70, 0x0E70, 0x0E70, 0x0E70,
		0x0E70, 0x0E70, 0x0E70, 0x0E70, 0x0E70,
		0x0000, 0x0000, 0x0000,
	}
	nextMask := [16]uint16{
		0x0000, 0x0000, 0x0000,
		0x1018, 0x1818, 0x1C18, 0x1E18, 0x1F18,
		0x1F18, 0x1E18, 0x1C18, 0x1818, 0x1018,
		0x0000, 0x0000, 0x0000,
	}

	hIconPrev = createMonoIcon(prevMask)
	hIconPlay = createMonoIcon(playMask)
	hIconPause = createMonoIcon(pauseMask)
	hIconNext = createMonoIcon(nextMask)
}

func initTaskbarThumbbar(hwnd uintptr) {
	if taskbarButtonsAdded || hwnd == 0 {
		return
	}

	initThumbnailIcons()

	if globalTaskbarList3 == nil {
		procCoInitialize.Call(0)
		var p *taskbarList3
		hr, _, _ := procCoCreateInstance.Call(
			uintptr(unsafe.Pointer(&clsidTaskbarList)),
			0,
			1, // CLSCTX_INPROC_SERVER
			uintptr(unsafe.Pointer(&iidITaskbarList3)),
			uintptr(unsafe.Pointer(&p)),
		)
		if hr != 0 || p == nil {
			return
		}
		globalTaskbarList3 = p
		syscall.SyscallN(globalTaskbarList3.lpVtbl.HrInit, uintptr(unsafe.Pointer(globalTaskbarList3)))
		syscall.SyscallN(globalTaskbarList3.lpVtbl.AddTab, uintptr(unsafe.Pointer(globalTaskbarList3)), hwnd)
	}

	var btns [3]thumbButton

	// Button 0: Prev
	btns[0].dwMask = thbIcon | thbTooltip | thbFlags
	btns[0].iId = thumbBtnPrev
	btns[0].hIcon = hIconPrev
	btns[0].dwFlags = thbfEnabled
	tipPrev, _ := syscall.UTF16FromString("Lagu Sebelumnya")
	copy(btns[0].szTip[:], tipPrev)

	// Button 1: Play/Pause
	btns[1].dwMask = thbIcon | thbTooltip | thbFlags
	btns[1].iId = thumbBtnPlayPause
	btns[1].hIcon = hIconPlay
	btns[1].dwFlags = thbfEnabled
	tipPlay, _ := syscall.UTF16FromString("Putar")
	copy(btns[1].szTip[:], tipPlay)

	// Button 2: Next
	btns[2].dwMask = thbIcon | thbTooltip | thbFlags
	btns[2].iId = thumbBtnNext
	btns[2].hIcon = hIconNext
	btns[2].dwFlags = thbfEnabled
	tipNext, _ := syscall.UTF16FromString("Lagu Selanjutnya")
	copy(btns[2].szTip[:], tipNext)

	hr, _, _ := syscall.SyscallN(
		globalTaskbarList3.lpVtbl.ThumbBarAddButtons,
		uintptr(unsafe.Pointer(globalTaskbarList3)),
		hwnd,
		3,
		uintptr(unsafe.Pointer(&btns[0])),
	)
	if hr == 0 {
		taskbarButtonsAdded = true
	}
}

func updateTaskbarThumbButtons(isPlaying bool) {
	if globalTaskbarList3 == nil || !taskbarButtonsAdded || globalHwnd == 0 {
		return
	}

	initThumbnailIcons()

	var btn thumbButton
	btn.dwMask = thbIcon | thbTooltip | thbFlags
	btn.iId = thumbBtnPlayPause
	btn.dwFlags = thbfEnabled

	if isPlaying {
		btn.hIcon = hIconPause
		tip, _ := syscall.UTF16FromString("Jeda")
		copy(btn.szTip[:], tip)
	} else {
		btn.hIcon = hIconPlay
		tip, _ := syscall.UTF16FromString("Putar")
		copy(btn.szTip[:], tip)
	}

	syscall.SyscallN(
		globalTaskbarList3.lpVtbl.ThumbBarUpdateButtons,
		uintptr(unsafe.Pointer(globalTaskbarList3)),
		globalHwnd,
		1,
		uintptr(unsafe.Pointer(&btn)),
	)
}

// trayWndProc intercepts Windows messages directed to the app window,
// specifically catching Tray Icon callbacks (WM_APP+101) and Taskbar Thumbnail Toolbar clicks (WM_COMMAND).
func trayWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	// Taskbar thumbnail toolbar button click
	if msg == wmCommand {
		wmId := wParam & 0xFFFF
		wmNotifyCode := (wParam >> 16) & 0xFFFF
		if wmNotifyCode == thbnClicked || wmNotifyCode == 0 {
			switch wmId {
			case thumbBtnPrev:
				if globalApp != nil && globalApp.ctx != nil {
					wailsRuntime.EventsEmit(globalApp.ctx, "tray:prev")
				}
				return 0
			case thumbBtnPlayPause:
				if globalApp != nil && globalApp.ctx != nil {
					wailsRuntime.EventsEmit(globalApp.ctx, "tray:play_pause")
				}
				return 0
			case thumbBtnNext:
				if globalApp != nil && globalApp.ctx != nil {
					wailsRuntime.EventsEmit(globalApp.ctx, "tray:next")
				}
				return 0
			}
		}
	}

	// Dynamic taskbar created notification from Windows Explorer
	if wmTaskbarBtnCreated != 0 && msg == wmTaskbarBtnCreated {
		initTaskbarThumbbar(hwnd)
		return 0
	}

	if msg == wmInitThumbbar {
		initTaskbarThumbbar(hwnd)
		return 0
	}

	if msg == wmApp+101 {
		switch lParam {
		case wmLButtonUp, wmLButtonDblClk:
			restoreAndShowApp(hwnd)
			return 0
		case wmRButtonUp, wmContextMenu:
			showTrayContextMenu(hwnd)
			return 0
		}
	}
	r, _, _ := procCallWindowProcW.Call(origWndProc, hwnd, msg, wParam, lParam)
	return r
}

func restoreAndShowApp(hwnd uintptr) {
	procShowWindow.Call(hwnd, swRestore)
	procSetForegroundWindow.Call(hwnd)
	if globalApp != nil {
		globalApp.ShowAppWindow()
	}
}

func showTrayContextMenu(hwnd uintptr) {
	hMenu, _, _ := procCreatePopupMenu.Call()
	if hMenu == 0 {
		return
	}
	defer procDestroyMenu.Call(hMenu)

	titleText, _ := syscall.UTF16PtrFromString("Dengerin - Buka Aplikasi")
	playText, _ := syscall.UTF16PtrFromString("Putar / Jeda")
	nextText, _ := syscall.UTF16PtrFromString("Lagu Selanjutnya")
	prevText, _ := syscall.UTF16PtrFromString("Lagu Sebelumnya")
	exitText, _ := syscall.UTF16PtrFromString("Keluar dari Dengerin")

	procAppendMenuW.Call(hMenu, mfString, 2001, uintptr(unsafe.Pointer(titleText)))
	procAppendMenuW.Call(hMenu, mfSeparator, 0, 0)
	procAppendMenuW.Call(hMenu, mfString, 2002, uintptr(unsafe.Pointer(playText)))
	procAppendMenuW.Call(hMenu, mfString, 2003, uintptr(unsafe.Pointer(nextText)))
	procAppendMenuW.Call(hMenu, mfString, 2004, uintptr(unsafe.Pointer(prevText)))
	procAppendMenuW.Call(hMenu, mfSeparator, 0, 0)
	procAppendMenuW.Call(hMenu, mfString, 2005, uintptr(unsafe.Pointer(exitText)))

	var pt point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))

	// Bring window to foreground before TrackPopupMenu so clicking away dismisses it cleanly
	procSetForegroundWindow.Call(hwnd)

	cmd, _, _ := procTrackPopupMenu.Call(
		hMenu,
		tpmReturnCmd|tpmRightButton,
		uintptr(pt.X),
		uintptr(pt.Y),
		0,
		hwnd,
		0,
	)

	// Post message to clear menu focus
	procPostMessageW.Call(hwnd, 0, 0, 0)

	switch cmd {
	case 2001:
		restoreAndShowApp(hwnd)
	case 2002:
		if globalApp != nil && globalApp.ctx != nil {
			wailsRuntime.EventsEmit(globalApp.ctx, "tray:play_pause")
		}
	case 2003:
		if globalApp != nil && globalApp.ctx != nil {
			wailsRuntime.EventsEmit(globalApp.ctx, "tray:next")
		}
	case 2004:
		if globalApp != nil && globalApp.ctx != nil {
			wailsRuntime.EventsEmit(globalApp.ctx, "tray:prev")
		}
	case 2005:
		if globalApp != nil && globalApp.ctx != nil {
			wailsRuntime.Quit(globalApp.ctx)
		}
	}
}

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

			// Save global handle and initialize tray + taskbar thumbnail toolbar
			globalHwnd = hwnd
			if hIconSmall != 0 {
				setupWindowsTray(hwnd, hIconSmall)
			}

			// Notify Windows Shell to refresh any cached taskbar icons
			procSHChangeNotify.Call(shcneAssocChanged, shcnfIdList, 0, 0)
			return
		}
	}()
}

func setupWindowsTray(hwnd uintptr, hIcon uintptr) {
	// Subclass window procedure once to intercept tray and taskbar messages
	if origWndProc == 0 {
		cb := syscall.NewCallback(trayWndProc)
		old, _, _ := procSetWindowLongPtrW.Call(hwnd, uintptr(gwlpWndProc), cb)
		origWndProc = old
	}

	var nid notifyIconDataW
	nid.cbSize = uint32(unsafe.Sizeof(nid))
	nid.hWnd = hwnd
	nid.uID = 1001
	nid.uFlags = nifMessage | nifIcon | nifTip
	nid.uCallbackMessage = wmApp + 101
	nid.hIcon = hIcon

	tip, _ := syscall.UTF16FromString("Dengerin Desktop")
	copy(nid.szTip[:], tip)

	procShellNotifyIconW.Call(uintptr(nimAdd), uintptr(unsafe.Pointer(&nid)))

	// Register TaskbarButtonCreated window message
	if wmTaskbarBtnCreated == 0 {
		taskbarMsg, _ := syscall.UTF16PtrFromString("TaskbarButtonCreated")
		msgId, _, _ := procRegisterWindowMessageW.Call(uintptr(unsafe.Pointer(taskbarMsg)))
		wmTaskbarBtnCreated = msgId
	}

	// Initialize Taskbar Thumbnail Toolbar on the window's UI thread
	procSendMessageW.Call(hwnd, wmInitThumbbar, 0, 0)
}

func removeWindowsTray() {
	if globalHwnd != 0 {
		var nid notifyIconDataW
		nid.cbSize = uint32(unsafe.Sizeof(nid))
		nid.hWnd = globalHwnd
		nid.uID = 1001
		procShellNotifyIconW.Call(uintptr(nimDelete), uintptr(unsafe.Pointer(&nid)))
		if origWndProc != 0 {
			procSetWindowLongPtrW.Call(globalHwnd, uintptr(gwlpWndProc), origWndProc)
			origWndProc = 0
		}
		if globalTaskbarList3 != nil {
			syscall.SyscallN(globalTaskbarList3.lpVtbl.Release, uintptr(unsafe.Pointer(globalTaskbarList3)))
			globalTaskbarList3 = nil
		}
		if hIconPlay != 0 {
			procDestroyIcon.Call(hIconPlay)
			procDestroyIcon.Call(hIconPause)
			procDestroyIcon.Call(hIconPrev)
			procDestroyIcon.Call(hIconNext)
			hIconPlay = 0
			hIconPause = 0
			hIconPrev = 0
			hIconNext = 0
		}
	}
}
