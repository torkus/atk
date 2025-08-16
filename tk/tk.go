// Copyright 2018 visualfc. All rights reserved.

package tk

import (
	"fmt"
	"runtime"

	"github.com/visualfc/atk/tk/interp"
)

var (
	tkHasInit            bool
	tkWindowInitAutoHide bool
	mainInterp           *interp.Interp
	rootWindow           *Window
	fnDebugHandle        func(string) = func(script string) {}
	fnErrorHandle        func(error)  = func(err error) {}
)

func Init() error {
	return InitEx(true, "", "")
}

func InitEx(tk_window_init_hide bool, tcl_library string, tk_library string) (err error) {
	mainInterp, err = interp.NewInterp()

	mainInterp.FnDebugHandle = fnDebugHandle
	mainInterp.FnErrorHandle = fnErrorHandle

	if err != nil {
		return err
	}

	err = mainInterp.InitTcl(tcl_library)
	if err != nil {
		return err
	}

	err = mainInterp.InitTk(tk_library)
	if err != nil {
		return err
	}

	tkWindowInitAutoHide = tk_window_init_hide
	//hide console for macOS bundle
	mainInterp.Eval("if {[info commands console] == \"console\"} {console hide}")

	for _, fn := range init_func_list {
		fn()
	}
	rootWindow = &Window{}
	rootWindow.Attach(".")
	if tkWindowInitAutoHide {
		rootWindow.Hide()
	}
	//hide wish menu on macos
	rootWindow.SetMenu(NewMenu(nil))
	tkHasInit = true
	return nil
}

var (
	init_func_list []func()
)

func registerInit(fn func()) {
	init_func_list = append(init_func_list, fn)
}

func SetDebugHandle(fn func(string)) {
	fnDebugHandle = fn
}

func SetErrorHandle(fn func(error)) {
	fnErrorHandle = fn
}

func MainInterp() *interp.Interp {
	return mainInterp
}

func TclVersion() (ver string) {
	return mainInterp.TclVersion()
}

func TkVersion() (ver string) {
	return mainInterp.TkVersion()
}

func TclLibary() (path string) {
	path, _ = evalAsString("set tcl_library")
	return
}

func TkLibrary() (path string) {
	path, _ = evalAsString("set tk_library")
	return
}

func AutoPath() (path string) {
	path, _ = evalAsString("puts $auto_path")
	return
}

func SetAutoPath(path string) string {
	r, _ := evalAsString(fmt.Sprintf("set auto_path [linsert $auto_path 0 %s]", path))
	return r
}

func init() {
	runtime.LockOSThread()
}

func MainLoop(fn func()) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if !tkHasInit {
		err := Init()
		if err != nil {
			return err
		}
	}
	interp.MainLoop(fn)
	return nil
}

func Async(fn func()) {
	interp.Async(fn)
}

func Update() {
	eval("update")
}

func Quit() {
	Async(func() {
		DestroyWidget(rootWindow)
	})
}

func eval(script string) error {
	return mainInterp.Eval(script)
}

func evalAsString(script string) (string, error) {
	return mainInterp.EvalAsString(script)
}

func evalAsInt(script string) (int, error) {
	return mainInterp.EvalAsInt(script)
}

func evalAsUint(script string) (uint, error) {
	return mainInterp.EvalAsUint(script)
}

func evalAsFloat64(script string) (float64, error) {
	return mainInterp.EvalAsFloat64(script)
}

func evalAsBool(script string) (bool, error) {
	return mainInterp.EvalAsBool(script)
}

func evalAsStringList(script string) ([]string, error) {
	return mainInterp.EvalAsStringList(script)
}

func evalAsIntList(script string) ([]int, error) {
	return mainInterp.EvalAsIntList(script)
}

func setObjText(obj string, text string) {
	mainInterp.SetStringVar(obj, text, false)
}

func setObjTextList(obj string, list []string) {
	mainInterp.SetStringList(obj, list, false)
}
