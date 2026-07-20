package com.example.go2apkapp;

import android.app.Activity;
import android.content.Intent;
import android.net.Uri;
import android.Manifest;
import android.content.pm.PackageManager;

final class NativeBridge {
    private static final String LIBRARY_NAME = "go2apkapp";

    static {
        System.loadLibrary(LIBRARY_NAME);
    }

    private NativeBridge() {}

    // ── JNI entry points (called from Go) ────────────────────────────────────

    private static native void sendEventToGo(String eventName);
    private static native void onPermissionResult(String permission, boolean granted);

    // ── Java→Go dispatch ─────────────────────────────────────────────────────

    static void sendEvent(String eventName) {
        sendEventToGo(eventName);
    }

    // ── Activity & Context helpers ────────────────────────────────────────────

    private static Activity currentActivity;

    public static void setActivity(Activity activity) {
        currentActivity = activity;
    }

    // ── UI Updates ────────────────────────────────────────────────────────────

    /** Called from Go to update any widget text on the UI thread. */
    public static void updateText(String id, String text) {
        if (currentActivity instanceof MainActivity) {
            currentActivity.runOnUiThread(() ->
                ((MainActivity) currentActivity).updateWidgetText(id, text));
        }
    }

    /** Called from Go to read the current text of a widget. */
    public static String getText(String id) {
        if (currentActivity instanceof MainActivity) {
            return ((MainActivity) currentActivity).getWidgetText(id);
        }
        return "";
    }

    // ── Intents ───────────────────────────────────────────────────────────────

    /**
     * Called from Go to start an Android activity.
     * @param action  Intent action string (e.g. "android.intent.action.VIEW")
     * @param data    URI data string, or "" if none
     * @param pkg     explicit package name for the target app, or "" if implicit
     */
    public static void startActivity(String action, String data, String pkg) {
        if (currentActivity == null) return;
        Intent intent = new Intent(action);
        if (!data.isEmpty()) intent.setData(Uri.parse(data));
        if (!pkg.isEmpty())  intent.setPackage(pkg);
        currentActivity.startActivity(intent);
    }

    // ── Broadcasts ────────────────────────────────────────────────────────────

    /** Called from Go to send a local broadcast. */
    public static void sendBroadcast(String action) {
        if (currentActivity == null) return;
        Intent intent = new Intent(action);
        currentActivity.sendBroadcast(intent);
    }

    // ── Permissions ───────────────────────────────────────────────────────────

    private static final int PERMISSION_REQUEST_CODE = 1001;

    /**
     * Called from Go to request a dangerous runtime permission.
     * The result is delivered back via onPermissionResult on the Go side.
     */
    public static void requestPermission(String permission) {
        if (currentActivity == null) return;
        if (currentActivity.checkSelfPermission(permission) == PackageManager.PERMISSION_GRANTED) {
            onPermissionResult(permission, true);
            return;
        }
        currentActivity.requestPermissions(new String[]{permission}, PERMISSION_REQUEST_CODE);
    }

    /** Called by MainActivity.onRequestPermissionsResult to relay the result to Go. */
    static void deliverPermissionResult(String permission, boolean granted) {
        onPermissionResult(permission, granted);
    }
}
