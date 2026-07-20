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
        sendEventToGo(eventName);
    }

    private static native void sendEventToGo(String eventName);

    public static void setActivity(Activity activity) {
        currentActivity = activity;
    }

    private static Activity currentActivity;

    public static void updateText(String id, String text) {
        if (currentActivity == null) return;
        currentActivity.runOnUiThread(() -> {
            try {
                java.lang.reflect.Field field = currentActivity.getClass().getDeclaredField(id);
                field.setAccessible(true);
                Object view = field.get(currentActivity);
                if (view instanceof android.widget.TextView) {
                    ((android.widget.TextView) view).setText(text);
                }
            } catch (Exception e) {
                e.printStackTrace();
            }
        });
    }
}
