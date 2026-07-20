package com.example.go2apkapp;

import android.app.Activity;

final class NativeBridge {
    private static final String LIBRARY_NAME = "go2apkapp";

    static {
        System.loadLibrary(LIBRARY_NAME);
    }

    private NativeBridge() {
    }

    static void sendEvent(String eventName) {
        updateText("display", "Event: " + eventName);
        String res = sendEventToGo(eventName);
        if (res != null && !res.isEmpty()) {
            updateText("display", res);
        }
    }

    private static native String sendEventToGo(String eventName);

    public static void setActivity(Activity activity) {
        currentActivity = activity;
    }

    private static Activity currentActivity;

    public static void updateText(String id, String text) {
        if (currentActivity instanceof MainActivity) {
            currentActivity.runOnUiThread(() -> {
                ((MainActivity) currentActivity).updateWidgetText(id, text);
            });
        }
    }
}
