package main

/*
#include <jni.h>
#include <stdlib.h>

// Forward declaration of the native method we expose to Java.
// The package name is com.example.go2apkapp and class is NativeBridge.
JNIEXPORT void JNICALL Java_com_example_go2apkapp_NativeBridge_sendEventToGo(JNIEnv* env, jclass clazz, jstring eventName);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

//export Java_com_example_go2apkapp_NativeBridge_sendEventToGo
func Java_com_example_go2apkapp_NativeBridge_sendEventToGo(env *C.JNIEnv, clazz C.jclass, eventName C.jstring) {
	name := javaString(env, eventName)
	// In the future, this will route to the user's Go event handlers.
	fmt.Printf("Received event from Java: %s\n", name)
}

func javaString(env *C.JNIEnv, str C.jstring) string {
	if str == 0 {
		return ""
	}
	var isCopy C.jboolean
	chars := (*C.JNIEnv)(env).GetStringUTFChars(str, &isCopy)
	if chars == nil {
		return ""
	}
	defer (*C.JNIEnv)(env).ReleaseStringUTFChars(str, chars)
	return C.GoString(chars)
}

func main() {}
