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
import android.widget.ScrollView;
import android.widget.ProgressBar;
import android.widget.Switch;
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
        view3.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_calc_onclick");
        });
        this.btn_calc = view3;
        view1.addView(view3);
        return view1;
    }

    public void updateWidgetText(String id, String text) {
        if (id.equals("welcome_text")) { if (this.welcome_text != null) this.welcome_text.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_calc")) { if (this.btn_calc != null) this.btn_calc.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
    }

    public void animateWidget(String id, String property, float to, int durationMs) {
        if (id.equals("welcome_text")) { if (this.welcome_text != null) ObjectAnimator.ofFloat(this.welcome_text, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_calc")) { if (this.btn_calc != null) ObjectAnimator.ofFloat(this.btn_calc, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("main_layout")) { if (this.main_layout != null) ObjectAnimator.ofFloat(this.main_layout, property, to).setDuration(durationMs).start(); return; }
    }

    public String getWidgetText(String id) {
        if (id.equals("welcome_text") && this.welcome_text != null) return this.welcome_text.getText().toString();
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
