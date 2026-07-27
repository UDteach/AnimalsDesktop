//go:build windows && !animalsdesktop_nonetwork

package main

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

const (
	inputMonitoringEnabled = true
	whKeyboardLL           = 13
	whMouseLL              = 14
)

var (
	procSetWindowsHookExW  = user32.NewProc("SetWindowsHookExW")
	procUnhookWindowsHook  = user32.NewProc("UnhookWindowsHookEx")
	procCallNextHookExProc = user32.NewProc("CallNextHookEx")
)

type mouseHookStruct struct {
	pt          win.POINT
	mouseData   uint32
	flags       uint32
	time        uint32
	dwExtraInfo uintptr
}

func (a *petApp) installInputMonitoring() {
	a.installKeyboardHook()
	a.installMouseHook()
}

func (a *petApp) cleanupInputMonitoring() {
	if a.keyHook != 0 {
		unhookWindowsHookEx(a.keyHook)
	}
	if a.mouseHook != 0 {
		unhookWindowsHookEx(a.mouseHook)
	}
}

func (a *petApp) installKeyboardHook() {
	cb := syscall.NewCallback(func(code int, wParam uintptr, lParam uintptr) uintptr {
		if code >= 0 && (wParam == win.WM_KEYDOWN || wParam == win.WM_SYSKEYDOWN) {
			a.postTypingFromHook()
		}
		return callNextHookEx(0, code, wParam, lParam)
	})
	a.keyHook = setWindowsHookEx(whKeyboardLL, cb, a.hinst, 0)
	a.keyHookFailed = a.keyHook == 0
}

func (a *petApp) installMouseHook() {
	cb := syscall.NewCallback(func(code int, wParam uintptr, lParam uintptr) uintptr {
		if code >= 0 && wParam == win.WM_LBUTTONDOWN {
			a.postMouseClickFromHook(lParam)
		}
		return callNextHookEx(0, code, wParam, lParam)
	})
	a.mouseHook = setWindowsHookEx(whMouseLL, cb, a.hinst, 0)
	a.mouseHookFailed = a.mouseHook == 0
}

func (a *petApp) postTypingFromHook() {
	defer recoverHookCallback()
	win.PostMessage(a.hwnd, wmTyping, 0, 0)
}

func (a *petApp) postMouseClickFromHook(lParam uintptr) {
	defer recoverHookCallback()
	pt := mouseHookPoint(lParam)
	win.PostMessage(a.hwnd, wmMouseClick, uintptr(uint32(pt.X)), uintptr(uint32(pt.Y)))
}

func recoverHookCallback() {
	_ = recover()
}

func mouseHookPoint(lParam uintptr) win.POINT {
	var hook mouseHookStruct
	if lParam != 0 {
		procRtlMoveMemory.Call(uintptr(unsafe.Pointer(&hook)), lParam, unsafe.Sizeof(hook))
	}
	return hook.pt
}

func setWindowsHookEx(idHook int, callback uintptr, module win.HINSTANCE, threadID uint32) uintptr {
	ret, _, _ := procSetWindowsHookExW.Call(uintptr(idHook), callback, uintptr(module), uintptr(threadID))
	return ret
}

func unhookWindowsHookEx(hook uintptr) bool {
	ret, _, _ := procUnhookWindowsHook.Call(hook)
	return ret != 0
}

func callNextHookEx(hook uintptr, code int, wParam uintptr, lParam uintptr) uintptr {
	ret, _, _ := procCallNextHookExProc.Call(hook, uintptr(code), wParam, lParam)
	return ret
}
