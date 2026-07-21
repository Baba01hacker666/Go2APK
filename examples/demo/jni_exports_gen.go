package main

// #include <jni.h>
import "C"

import (
	"unsafe"
	"github.com/Baba01hacker666/Go2APK/android"
)

//export Java_com_example_go2apkapp_NativeBridge_sendEventToGo
func Java_com_example_go2apkapp_NativeBridge_sendEventToGo(env *C.JNIEnv, clazz C.jclass, eventName C.jstring) {
	android.Exported_sendEventToGo(unsafe.Pointer(env), unsafe.Pointer(clazz), unsafe.Pointer(eventName))
}

//export Java_com_example_go2apkapp_NativeBridge_onPermissionResult
func Java_com_example_go2apkapp_NativeBridge_onPermissionResult(env *C.JNIEnv, clazz C.jclass, permission C.jstring, granted C.jboolean) {
	android.Exported_onPermissionResult(unsafe.Pointer(env), unsafe.Pointer(clazz), unsafe.Pointer(permission), granted != 0)
}

//export Java_com_example_go2apkapp_NativeBridge_onVpnEstablished
func Java_com_example_go2apkapp_NativeBridge_onVpnEstablished(env *C.JNIEnv, clazz C.jclass, fd C.jint) {
	android.Exported_onVpnEstablished(unsafe.Pointer(env), unsafe.Pointer(clazz), int(fd))
}
