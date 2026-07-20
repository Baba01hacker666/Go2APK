package com.example.counter;

final class NativeBridge {
    private static final String LIBRARY_NAME = "go2apkapp";

    static {
        System.loadLibrary(LIBRARY_NAME);
    }

    private NativeBridge() {
    }

    static void sendEvent(String eventName) {
        sendEventToGo(eventName);
    }

    private static native void sendEventToGo(String eventName);
}
