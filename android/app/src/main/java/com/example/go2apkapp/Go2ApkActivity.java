package com.example.go2apkapp;

public interface Go2ApkActivity {
    void updateWidgetText(String id, String text);
    void animateWidget(String id, String property, float to, int durationMs);
    String getWidgetText(String id);
}
