package com.example.go2apkapp;

import android.app.Activity;
import android.os.Bundle;
import android.widget.TextView;

public class MainActivity extends Activity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        TextView view = new TextView(this);
        view.setText("Hello from Go2APK");
        view.setTextSize(24);
        view.setPadding(32, 32, 32, 32);
        setContentView(view);
    }
}
