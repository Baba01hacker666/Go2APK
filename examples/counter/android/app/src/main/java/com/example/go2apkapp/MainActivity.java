package com.example.go2apkapp;

import android.app.Activity;
import android.graphics.Color;
import android.os.Bundle;
import android.view.Gravity;
import android.view.View;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.TextView;
import android.widget.EditText;

public class MainActivity extends Activity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setTitle("Go2APK Calculator");
        setContentView(createAppView());
    }

    private View createAppView() {
        LinearLayout rootView = new LinearLayout(this);
        return rootView;
    }
}
