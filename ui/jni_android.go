//go:build android && cgo

package ui

/*
#include <jni.h>
#include <stdlib.h>

// Forward declaration
JNIEXPORT jstring JNICALL Java_com_example_go2apkapp_NativeBridge_sendEventToGo(JNIEnv* env, jclass clazz, jstring eventName);

static const char* GetString(JNIEnv* env, jstring str) {
    return (*env)->GetStringUTFChars(env, str, NULL);
}

static void ReleaseString(JNIEnv* env, jstring str, const char* chars) {
    (*env)->ReleaseStringUTFChars(env, str, chars);
}

static jstring CreateJavaString(JNIEnv* env, const char* str) {
    return (*env)->NewStringUTF(env, str);
}

static jclass NewGlobalRef(JNIEnv* env, jclass cls) {
    return (jclass)(*env)->NewGlobalRef(env, (jobject)cls);
}

static JavaVM* ExtractJavaVM(JNIEnv* env) {
    JavaVM* vm;
    if ((*env)->GetJavaVM(env, &vm) != JNI_OK) {
        return NULL;
    }
    return vm;
}

static JNIEnv* GetJNIEnv(JavaVM* vm) {
    JNIEnv* env;
    int status = (*vm)->GetEnv(vm, (void**)&env, JNI_VERSION_1_6);
    if (status == JNI_EDETACHED) {
        if ((*vm)->AttachCurrentThread(vm, &env, NULL) != JNI_OK) {
            return NULL;
        }
    }
    return env;
}

static void CallUpdateText(JNIEnv* env, jclass clazz, const char* id, const char* text) {
    jmethodID mid = (*env)->GetStaticMethodID(env, clazz, "updateText", "(Ljava/lang/String;Ljava/lang/String;)V");
    if (mid == NULL) return;
    jstring jId = (*env)->NewStringUTF(env, id);
    jstring jText = (*env)->NewStringUTF(env, text);
    (*env)->CallStaticVoidMethod(env, clazz, mid, jId, jText);
    (*env)->DeleteLocalRef(env, jId);
    (*env)->DeleteLocalRef(env, jText);
}
*/
import "C"
import "unsafe"

// globalVM stores the Java VM pointer for background thread attachment.
var globalVM *C.JavaVM
var globalBridgeClass C.jclass

// currentEnv holds the JNIEnv for the currently executing UI callback.
var currentEnv *C.JNIEnv

//export JNI_OnLoad
func JNI_OnLoad(vm *C.JavaVM, reserved unsafe.Pointer) C.jint {
	globalVM = vm
	return C.JNI_VERSION_1_6
}

//export Java_com_example_go2apkapp_NativeBridge_sendEventToGo
func Java_com_example_go2apkapp_NativeBridge_sendEventToGo(env *C.JNIEnv, clazz C.jclass, eventName C.jstring) C.jstring {
	currentEnv = env
	defer func() { currentEnv = nil }()

	if globalVM == nil {
		globalVM = C.ExtractJavaVM(env)
	}
	if globalBridgeClass == 0 {
		globalBridgeClass = C.NewGlobalRef(env, clazz)
	}
	name := javaString(env, eventName)
	handleEvent(name)

	msg := C.CString("Go callback finished: " + name)
	defer C.free(unsafe.Pointer(msg))
	return C.CreateJavaString(env, msg)
}

func javaString(env *C.JNIEnv, str C.jstring) string {
	if str == 0 {
		return ""
	}
	chars := C.GetString(env, str)
	if chars == nil {
		return ""
	}
	defer C.ReleaseString(env, str, chars)
	return C.GoString(chars)
}

func updateTextNative(id string, text string) {
	if globalBridgeClass == 0 {
		return
	}

	env := currentEnv
	if env == nil {
		if globalVM == nil {
			return
		}
		env = C.GetJNIEnv(globalVM)
	}
	if env == nil {
		return
	}

	cId := C.CString(id)
	defer C.free(unsafe.Pointer(cId))

	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	C.CallUpdateText(env, globalBridgeClass, cId, cText)
}
