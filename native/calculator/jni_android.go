//go:build android && cgo

package main

/*
#include <stdlib.h>
#include <jni.h>

static const char* go2apk_get_string_utf_chars(JNIEnv* env, jstring value) {
    return (*env)->GetStringUTFChars(env, value, 0);
}

static void go2apk_release_string_utf_chars(JNIEnv* env, jstring value, const char* chars) {
    (*env)->ReleaseStringUTFChars(env, value, chars);
}

static jstring go2apk_new_string_utf(JNIEnv* env, const char* value) {
    return (*env)->NewStringUTF(env, value);
}
*/
import "C"
import "unsafe"

//export Java_com_example_go2apkapp_NativeCalculator_calculateWithGo
func Java_com_example_go2apkapp_NativeCalculator_calculateWithGo(env *C.JNIEnv, _ C.jclass, jLeft C.jstring, jOperator C.jstring, jRight C.jstring) C.jstring {
	left := javaString(env, jLeft)
	operator := javaString(env, jOperator)
	right := javaString(env, jRight)

	result, err := Calculate(left, operator, right)
	if err != nil {
		result = err.Error()
	}
	return newJavaString(env, result)
}

func javaString(env *C.JNIEnv, value C.jstring) string {
	chars := C.go2apk_get_string_utf_chars(env, value)
	defer C.go2apk_release_string_utf_chars(env, value, chars)
	return C.GoString(chars)
}

func newJavaString(env *C.JNIEnv, value string) C.jstring {
	chars := C.CString(value)
	defer C.free(unsafe.Pointer(chars))
	return C.go2apk_new_string_utf(env, chars)
}
