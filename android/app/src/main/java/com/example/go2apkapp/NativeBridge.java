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
    static native void onVpnEstablished(int fd);

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
        if (currentActivity instanceof Go2ApkActivity) {
            currentActivity.runOnUiThread(() ->
                ((Go2ApkActivity) currentActivity).updateWidgetText(id, text));
        }
    }

    /** Called from Go to animate a widget on the UI thread. */
    public static void animate(String id, String property, float to, int durationMs) {
        if (currentActivity instanceof Go2ApkActivity) {
            currentActivity.runOnUiThread(() ->
                ((Go2ApkActivity) currentActivity).animateWidget(id, property, to, durationMs));
        }
    }

    /** Called from Go to read the current text of a widget. */
    public static String getText(String id) {
        if (currentActivity instanceof Go2ApkActivity) {
            return ((Go2ApkActivity) currentActivity).getWidgetText(id);
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
    public static void startActivity(String action, String data, String pkg, String extrasJson) {
        if (currentActivity == null) return;
        Intent intent = new Intent(action);
        if (!data.isEmpty()) intent.setData(Uri.parse(data));
        if (!pkg.isEmpty())  intent.setPackage(pkg);
        if (extrasJson != null && !extrasJson.isEmpty()) {
            try {
                org.json.JSONObject obj = new org.json.JSONObject(extrasJson);
                java.util.Iterator<String> keys = obj.keys();
                while (keys.hasNext()) {
                    String key = keys.next();
                    intent.putExtra(key, obj.getString(key));
                }
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
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

    // ── VPN ───────────────────────────────────────────────────────────────────

    private static String pendingVpnConfig = null;
    private static final int VPN_REQUEST_CODE = 1002;

    public static void startVpn(String configJson) {
        if (currentActivity == null) return;
        Intent intent = android.net.VpnService.prepare(currentActivity);
        if (intent != null) {
            pendingVpnConfig = configJson;
            currentActivity.startActivityForResult(intent, VPN_REQUEST_CODE);
        } else {
            // Already authorized
            launchVpnService(configJson);
        }
    }

    static void handleActivityResult(int requestCode, int resultCode, Intent data) {
        if (requestCode == VPN_REQUEST_CODE && resultCode == Activity.RESULT_OK) {
            launchVpnService(pendingVpnConfig);
            pendingVpnConfig = null;
        }
    }

    public static void navigate(String target) {
        if (currentActivity == null) return;
        try {
            Class<?> targetClass = Class.forName(currentActivity.getPackageName() + "." + target);
            Intent intent = new Intent(currentActivity, targetClass);
            currentActivity.startActivity(intent);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private static void launchVpnService(String configJson) {
        if (currentActivity == null) return;
        try {
            Class<?> vpnClass = Class.forName(currentActivity.getPackageName() + ".Go2ApkVpnService");
            Intent intent = new Intent(currentActivity, vpnClass);
            intent.putExtra("config", configJson);
            currentActivity.startService(intent);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
