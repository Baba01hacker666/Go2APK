package android

import (
	"fmt"
	"strings"
)

// RenderJNIExports generates a CGo file with the correct JNI exported functions for the user's package.
func RenderJNIExports(pkgName string, androidImport string) string {
	if androidImport == "" {
		androidImport = "github.com/Baba01hacker666/Go2APK/android"
	}
	
	jniPrefix := "Java_" + strings.ReplaceAll(pkgName, ".", "_") + "_NativeBridge_"
	
	return fmt.Sprintf(`//go:build android && cgo
package main

// #include <jni.h>
import "C"

import (
	"unsafe"
	"%s"
)

//export %ssendEventToGo
func %ssendEventToGo(env *C.JNIEnv, clazz C.jclass, eventName C.jstring) {
	android.Exported_sendEventToGo(unsafe.Pointer(env), unsafe.Pointer(clazz), unsafe.Pointer(eventName))
}

//export %sonPermissionResult
func %sonPermissionResult(env *C.JNIEnv, clazz C.jclass, permission C.jstring, granted C.jboolean) {
	android.Exported_onPermissionResult(unsafe.Pointer(env), unsafe.Pointer(clazz), unsafe.Pointer(permission), granted != 0)
}

//export %sonVpnEstablished
func %sonVpnEstablished(env *C.JNIEnv, clazz C.jclass, fd C.jint) {
	android.Exported_onVpnEstablished(unsafe.Pointer(env), unsafe.Pointer(clazz), int(fd))
}
`, androidImport, jniPrefix, jniPrefix, jniPrefix, jniPrefix, jniPrefix, jniPrefix)
}
