package com.example.go2apkapp;

import android.content.Intent;
import android.app.Activity;
import android.graphics.Color;
import android.os.Bundle;
import android.view.Gravity;
import android.view.View;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.TextView;
import android.widget.EditText;
import android.widget.ImageView;
import android.widget.VideoView;
import android.media.MediaPlayer;
import android.net.Uri;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import java.net.URL;
import java.util.concurrent.Executors;
import android.animation.ObjectAnimator;

public class MainActivity extends Activity implements Go2ApkActivity {
    private LinearLayout main_layout;
    private TextView welcome_text;
    private Button btn_calc;
    private Button btn_http;
    private TextView http_result;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setTitle("Go2APK Calculator");
        setContentView(createAppView());
        NativeBridge.setActivity(this);
    }

    private View createAppView() {
        LinearLayout view1 = new LinearLayout(this);
        view1.setOrientation(LinearLayout.VERTICAL);
        view1.setGravity(Gravity.CENTER);
        LinearLayout.LayoutParams lp_view1 = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.MATCH_PARENT, 0.000000f);
        view1.setLayoutParams(lp_view1);
        view1.setPadding(24, 24, 24, 24);
        view1.setBackgroundColor(Color.parseColor("#1E1E2E"));
        this.main_layout = view1;
        TextView view2 = new TextView(this);
        view2.setText(android.text.Html.fromHtml("Welcome to Go2APK Multi-Page Demo", android.text.Html.FROM_HTML_MODE_COMPACT));
        view2.setGravity(Gravity.END);
        LinearLayout.LayoutParams lp_view2 = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT, 0.000000f);
        lp_view2.setMargins(16, 16, 16, 16);
        view2.setLayoutParams(lp_view2);
        view2.setTextSize(24);
        view2.setTextColor(Color.parseColor("#CBA6F7"));
        this.welcome_text = view2;
        view1.addView(view2);
        Button view3 = new Button(this);
        view3.setText(android.text.Html.fromHtml("Open Calculator", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view3 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view3.setMargins(16, 16, 16, 16);
        view3.setLayoutParams(lp_view3);
        view3.setBackgroundColor(Color.parseColor("#89B4FA"));
        view3.setTextSize(20);
        view3.setTextColor(Color.parseColor("#11111B"));
        this.btn_calc = view3;
        view1.addView(view3);
        Button view4 = new Button(this);
        view4.setText(android.text.Html.fromHtml("Test HTTP GET", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view4 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view4.setMargins(16, 16, 16, 16);
        view4.setLayoutParams(lp_view4);
        view4.setBackgroundColor(Color.parseColor("#A6E3A1"));
        view4.setTextSize(20);
        view4.setTextColor(Color.parseColor("#11111B"));
        view4.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_http_onclick");
        });
        this.btn_http = view4;
        view1.addView(view4);
        TextView view5 = new TextView(this);
        view5.setText(android.text.Html.fromHtml("Result will appear here", android.text.Html.FROM_HTML_MODE_COMPACT));
        view5.setGravity(Gravity.END);
        LinearLayout.LayoutParams lp_view5 = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT, 0.000000f);
        lp_view5.setMargins(16, 16, 16, 16);
        view5.setLayoutParams(lp_view5);
        view5.setTextSize(14);
        view5.setTextColor(Color.parseColor("#CDD6F4"));
        this.http_result = view5;
        view1.addView(view5);
        return view1;
    }

    public void updateWidgetText(String id, String text) {
        if (id.equals("welcome_text")) { if (this.welcome_text != null) this.welcome_text.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_calc")) { if (this.btn_calc != null) this.btn_calc.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_http")) { if (this.btn_http != null) this.btn_http.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("http_result")) { if (this.http_result != null) this.http_result.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
    }

    public void animateWidget(String id, String property, float to, int durationMs) {
        if (id.equals("welcome_text")) { if (this.welcome_text != null) ObjectAnimator.ofFloat(this.welcome_text, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_calc")) { if (this.btn_calc != null) ObjectAnimator.ofFloat(this.btn_calc, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_http")) { if (this.btn_http != null) ObjectAnimator.ofFloat(this.btn_http, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("http_result")) { if (this.http_result != null) ObjectAnimator.ofFloat(this.http_result, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("main_layout")) { if (this.main_layout != null) ObjectAnimator.ofFloat(this.main_layout, property, to).setDuration(durationMs).start(); return; }
    }

    public String getWidgetText(String id) {
        if (id.equals("welcome_text") && this.welcome_text != null) return this.welcome_text.getText().toString();
        if (id.equals("http_result") && this.http_result != null) return this.http_result.getText().toString();
        return "";
    }

    private void loadImage(ImageView view, String src) {
        if (src.startsWith("http")) {
            Executors.newSingleThreadExecutor().execute(() -> {
                try {
                    Bitmap bmp = BitmapFactory.decodeStream(new URL(src).openConnection().getInputStream());
                    runOnUiThread(() -> view.setImageBitmap(bmp));
                } catch (Exception e) { e.printStackTrace(); }
            });
        }
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);
        NativeBridge.handleActivityResult(requestCode, resultCode, data);
    }

    @Override
    public void onRequestPermissionsResult(int requestCode, String[] permissions, int[] grantResults) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults);
        for (int i = 0; i < permissions.length; i++) {
            boolean granted = (grantResults[i] == android.content.pm.PackageManager.PERMISSION_GRANTED);
            NativeBridge.deliverPermissionResult(permissions[i], granted);
        }
    }
}
