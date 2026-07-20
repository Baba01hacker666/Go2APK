//go:build android && cgo

package android

/*
#include <jni.h>
#include <stdlib.h>

// Forward declarations
JNIEXPORT void JNICALL Java_com_example_go2apkapp_NativeBridge_sendEventToGo(JNIEnv* env, jclass clazz, jstring eventName);
JNIEXPORT void JNICALL Java_com_example_go2apkapp_NativeBridge_onPermissionResult(JNIEnv* env, jclass clazz, jstring permission, jboolean granted);

static const char* GetString(JNIEnv* env, jstring str) {
    return (*env)->GetStringUTFChars(env, str, NULL);
}

static void ReleaseString(JNIEnv* env, jstring str, const char* chars) {
    (*env)->ReleaseStringUTFChars(env, str, chars);
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
    jstring jId   = (*env)->NewStringUTF(env, id);
    jstring jText = (*env)->NewStringUTF(env, text);
    (*env)->CallStaticVoidMethod(env, clazz, mid, jId, jText);
    (*env)->DeleteLocalRef(env, jId);
    (*env)->DeleteLocalRef(env, jText);
}

static jstring CallGetText(JNIEnv* env, jclass clazz, const char* id) {
    jmethodID mid = (*env)->GetStaticMethodID(env, clazz, "getText", "(Ljava/lang/String;)Ljava/lang/String;");
    if (mid == NULL) return NULL;
    jstring jId = (*env)->NewStringUTF(env, id);
    jstring result = (jstring)(*env)->CallStaticObjectMethod(env, clazz, mid, jId);
    (*env)->DeleteLocalRef(env, jId);
    return result;
}

static void CallStartActivity(JNIEnv* env, jclass clazz, const char* action, const char* data, const char* pkg) {
    jmethodID mid = (*env)->GetStaticMethodID(env, clazz, "startActivity",
        "(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V");
    if (mid == NULL) return;
    jstring jAction = (*env)->NewStringUTF(env, action);
    jstring jData   = (*env)->NewStringUTF(env, data);
    jstring jPkg    = (*env)->NewStringUTF(env, pkg);
    (*env)->CallStaticVoidMethod(env, clazz, mid, jAction, jData, jPkg);
    (*env)->DeleteLocalRef(env, jAction);
    (*env)->DeleteLocalRef(env, jData);
    (*env)->DeleteLocalRef(env, jPkg);
}

static void CallSendBroadcast(JNIEnv* env, jclass clazz, const char* action) {
    jmethodID mid = (*env)->GetStaticMethodID(env, clazz, "sendBroadcast", "(Ljava/lang/String;)V");
    if (mid == NULL) return;
    jstring jAction = (*env)->NewStringUTF(env, action);
    (*env)->CallStaticVoidMethod(env, clazz, mid, jAction);
    (*env)->DeleteLocalRef(env, jAction);
}

static void CallRequestPermission(JNIEnv* env, jclass clazz, const char* permission) {
    jmethodID mid = (*env)->GetStaticMethodID(env, clazz, "requestPermission", "(Ljava/lang/String;)V");
    if (mid == NULL) return;
    jstring jPerm = (*env)->NewStringUTF(env, permission);
    (*env)->CallStaticVoidMethod(env, clazz, mid, jPerm);
    (*env)->DeleteLocalRef(env, jPerm);
}

static void CallAnimate(JNIEnv* env, jclass clazz, const char* id, const char* property, float to, int durationMs) {
    jmethodID mid = (*env)->GetStaticMethodID(env, clazz, "animate", "(Ljava/lang/String;Ljava/lang/String;FI)V");
    if (mid == NULL) return;
    jstring jId = (*env)->NewStringUTF(env, id);
    jstring jProp = (*env)->NewStringUTF(env, property);
    (*env)->CallStaticVoidMethod(env, clazz, mid, jId, jProp, (jfloat)to, (jint)durationMs);
    (*env)->DeleteLocalRef(env, jId);
    (*env)->DeleteLocalRef(env, jProp);
}

static void CallStartVpn(JNIEnv* env, jclass clazz, const char* configJson) {
    jmethodID mid = (*env)->GetStaticMethodID(env, clazz, "startVpn", "(Ljava/lang/String;)V");
    if (mid == NULL) return;
    jstring jConfig = (*env)->NewStringUTF(env, configJson);
    (*env)->CallStaticVoidMethod(env, clazz, mid, jConfig);
    (*env)->DeleteLocalRef(env, jConfig);
}
*/
import "C"
import "unsafe"

// globalVM stores the Java VM pointer for background thread attachment.
var globalVM *C.JavaVM
var globalBridgeClass C.jclass

// currentEnv holds the JNIEnv for the currently executing UI callback.
var currentEnv *C.JNIEnv

// permissionCallbacks holds pending permission request callbacks keyed by permission name.
var permissionCallbacks = map[string]func(bool){}

//export JNI_OnLoad
func JNI_OnLoad(vm *C.JavaVM, reserved unsafe.Pointer) C.jint {
	globalVM = vm
	return C.JNI_VERSION_1_6
}

//export Java_com_example_go2apkapp_NativeBridge_sendEventToGo
func Java_com_example_go2apkapp_NativeBridge_sendEventToGo(env *C.JNIEnv, clazz C.jclass, eventName C.jstring) {
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
}

//export Java_com_example_go2apkapp_NativeBridge_onPermissionResult
func Java_com_example_go2apkapp_NativeBridge_onPermissionResult(env *C.JNIEnv, clazz C.jclass, permission C.jstring, granted C.jboolean) {
	perm := javaString(env, permission)
	if cb, ok := permissionCallbacks[perm]; ok {
		cb(granted != 0)
		delete(permissionCallbacks, perm)
	}
}

//export Java_com_example_go2apkapp_NativeBridge_onVpnEstablished
func Java_com_example_go2apkapp_NativeBridge_onVpnEstablished(env *C.JNIEnv, clazz C.jclass, fd C.jint) {
	if VpnCallback != nil {
		VpnCallback(int(fd))
	}
}

// VpnCallback is called when the VPN interface is established.
var VpnCallback func(fd int)

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

func getEnv() *C.JNIEnv {
	if currentEnv != nil {
		return currentEnv
	}
	if globalVM == nil {
		return nil
	}
	return C.GetJNIEnv(globalVM)
}

func updateTextNative(id string, text string) {
	if globalBridgeClass == 0 {
		return
	}
	env := getEnv()
	if env == nil {
		return
	}
	cId := C.CString(id)
	defer C.free(unsafe.Pointer(cId))
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.CallUpdateText(env, globalBridgeClass, cId, cText)
}

func getTextNative(id string) string {
	if globalBridgeClass == 0 {
		return ""
	}
	env := getEnv()
	if env == nil {
		return ""
	}
	cId := C.CString(id)
	defer C.free(unsafe.Pointer(cId))
	jstr := C.CallGetText(env, globalBridgeClass, cId)
	if jstr == 0 {
		return ""
	}
	return javaString(env, jstr)
}

func startActivityNative(intent Intent) {
	if globalBridgeClass == 0 {
		return
	}
	env := getEnv()
	if env == nil {
		return
	}
	cAction := C.CString(intent.Action)
	defer C.free(unsafe.Pointer(cAction))
	cData := C.CString(intent.Data)
	defer C.free(unsafe.Pointer(cData))
	cPkg := C.CString(intent.Package)
	defer C.free(unsafe.Pointer(cPkg))
	C.CallStartActivity(env, globalBridgeClass, cAction, cData, cPkg)
}

func startActivityForResultNative(intent Intent, onResult func(resultCode int, data map[string]string)) {
	// For now, just launch without result capture.
	startActivityNative(intent)
}

func sendBroadcastNative(action string) {
	if globalBridgeClass == 0 {
		return
	}
	env := getEnv()
	if env == nil {
		return
	}
	cAction := C.CString(action)
	defer C.free(unsafe.Pointer(cAction))
	C.CallSendBroadcast(env, globalBridgeClass, cAction)
}

func sendBroadcastWithExtrasNative(action string, extras map[string]string) {
	// Delegate to basic sendBroadcast for now.
	sendBroadcastNative(action)
}

func requestPermissionNative(permission string, onResult func(granted bool)) {
	if globalBridgeClass == 0 {
		if onResult != nil {
			onResult(false)
		}
		return
	}
	env := getEnv()
	if env == nil {
		if onResult != nil {
			onResult(false)
		}
		return
	}
	if onResult != nil {
		permissionCallbacks[permission] = onResult
	}
	cPerm := C.CString(permission)
	defer C.free(unsafe.Pointer(cPerm))
	C.CallRequestPermission(env, globalBridgeClass, cPerm)
}

func animateNative(id string, property string, to float32, durationMs int) {
	if globalBridgeClass == 0 {
		return
	}
	env := getEnv()
	if env == nil {
		return
	}
	cId := C.CString(id)
	defer C.free(unsafe.Pointer(cId))
	cProp := C.CString(property)
	defer C.free(unsafe.Pointer(cProp))
	C.CallAnimate(env, globalBridgeClass, cId, cProp, C.float(to), C.int(durationMs))
}

// StartVPN requests Android VPN permissions. If granted, it starts the VPN
// service with the given config. The established TUN file descriptor (fd)
// is sent back to the onEstablished callback.
func StartVPN(config VpnConfig, onEstablished func(fd int)) {
	VpnCallback = onEstablished

	env := getEnv()
	if env == nil {
		return
	}

	cConfig := C.CString(config.toJSON())
	defer C.free(unsafe.Pointer(cConfig))

	C.CallStartVpn(env, globalBridgeClass, cConfig)
}
