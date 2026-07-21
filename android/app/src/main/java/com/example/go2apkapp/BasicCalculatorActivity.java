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

public class BasicCalculatorActivity extends Activity implements Go2ApkActivity {
    private LinearLayout calc_layout;
    private TextView display;
    private Button btn_adv;
    private Button btn_clear;
    private Button btn_lp;
    private Button btn_rp;
    private Button btn_div;
    private Button btn_7;
    private Button btn_8;
    private Button btn_9;
    private Button btn_mul;
    private Button btn_4;
    private Button btn_5;
    private Button btn_6;
    private Button btn_sub;
    private Button btn_1;
    private Button btn_2;
    private Button btn_3;
    private Button btn_add;
    private Button btn_0;
    private Button btn_dot;
    private Button btn_del;
    private Button btn_eq;

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
        this.calc_layout = view1;
        TextView view2 = new TextView(this);
        view2.setText(android.text.Html.fromHtml("0", android.text.Html.FROM_HTML_MODE_COMPACT));
        view2.setGravity(Gravity.END);
        LinearLayout.LayoutParams lp_view2 = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT, 0.000000f);
        lp_view2.setMargins(16, 16, 16, 16);
        view2.setLayoutParams(lp_view2);
        view2.setPadding(32, 32, 32, 32);
        view2.setBackgroundColor(Color.parseColor("#181825"));
        view2.setTextSize(48);
        view2.setTextColor(Color.parseColor("#CBA6F7"));
        this.display = view2;
        view1.addView(view2);
        Button view3 = new Button(this);
        view3.setText(android.text.Html.fromHtml("Advanced Mode", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view3 = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view3.setMargins(8, 8, 8, 8);
        view3.setLayoutParams(lp_view3);
        view3.setBackgroundColor(Color.parseColor("#F5C2E7"));
        view3.setTextSize(20);
        view3.setTextColor(Color.parseColor("#11111B"));
        view3.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_adv_onclick");
        });
        this.btn_adv = view3;
        view1.addView(view3);
        LinearLayout view4 = new LinearLayout(this);
        view4.setOrientation(LinearLayout.HORIZONTAL);
        view4.setGravity(Gravity.CENTER);
        LinearLayout.LayoutParams lp_view4 = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view4.setMargins(4, 4, 4, 4);
        view4.setLayoutParams(lp_view4);
        view1.addView(view4);
        Button view5 = new Button(this);
        view5.setText(android.text.Html.fromHtml("C", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view5 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view5.setMargins(8, 8, 8, 8);
        view5.setLayoutParams(lp_view5);
        view5.setBackgroundColor(Color.parseColor("#F38BA8"));
        view5.setTextSize(28);
        view5.setTextColor(Color.parseColor("#11111B"));
        view5.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_clear_onclick");
        });
        this.btn_clear = view5;
        view4.addView(view5);
        Button view6 = new Button(this);
        view6.setText(android.text.Html.fromHtml("(", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view6 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view6.setMargins(8, 8, 8, 8);
        view6.setLayoutParams(lp_view6);
        view6.setBackgroundColor(Color.parseColor("#89B4FA"));
        view6.setTextSize(28);
        view6.setTextColor(Color.parseColor("#11111B"));
        view6.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_lp_onclick");
        });
        this.btn_lp = view6;
        view4.addView(view6);
        Button view7 = new Button(this);
        view7.setText(android.text.Html.fromHtml(")", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view7 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view7.setMargins(8, 8, 8, 8);
        view7.setLayoutParams(lp_view7);
        view7.setBackgroundColor(Color.parseColor("#89B4FA"));
        view7.setTextSize(28);
        view7.setTextColor(Color.parseColor("#11111B"));
        view7.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_rp_onclick");
        });
        this.btn_rp = view7;
        view4.addView(view7);
        Button view8 = new Button(this);
        view8.setText(android.text.Html.fromHtml("/", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view8 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view8.setMargins(8, 8, 8, 8);
        view8.setLayoutParams(lp_view8);
        view8.setBackgroundColor(Color.parseColor("#F9E2AF"));
        view8.setTextSize(28);
        view8.setTextColor(Color.parseColor("#11111B"));
        view8.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_div_onclick");
        });
        this.btn_div = view8;
        view4.addView(view8);
        LinearLayout view9 = new LinearLayout(this);
        view9.setOrientation(LinearLayout.HORIZONTAL);
        view9.setGravity(Gravity.CENTER);
        LinearLayout.LayoutParams lp_view9 = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view9.setMargins(4, 4, 4, 4);
        view9.setLayoutParams(lp_view9);
        view1.addView(view9);
        Button view10 = new Button(this);
        view10.setText(android.text.Html.fromHtml("7", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view10 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view10.setMargins(8, 8, 8, 8);
        view10.setLayoutParams(lp_view10);
        view10.setBackgroundColor(Color.parseColor("#313244"));
        view10.setTextSize(28);
        view10.setTextColor(Color.parseColor("#CDD6F4"));
        view10.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_7_onclick");
        });
        this.btn_7 = view10;
        view9.addView(view10);
        Button view11 = new Button(this);
        view11.setText(android.text.Html.fromHtml("8", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view11 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view11.setMargins(8, 8, 8, 8);
        view11.setLayoutParams(lp_view11);
        view11.setBackgroundColor(Color.parseColor("#313244"));
        view11.setTextSize(28);
        view11.setTextColor(Color.parseColor("#CDD6F4"));
        view11.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_8_onclick");
        });
        this.btn_8 = view11;
        view9.addView(view11);
        Button view12 = new Button(this);
        view12.setText(android.text.Html.fromHtml("9", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view12 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view12.setMargins(8, 8, 8, 8);
        view12.setLayoutParams(lp_view12);
        view12.setBackgroundColor(Color.parseColor("#313244"));
        view12.setTextSize(28);
        view12.setTextColor(Color.parseColor("#CDD6F4"));
        view12.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_9_onclick");
        });
        this.btn_9 = view12;
        view9.addView(view12);
        Button view13 = new Button(this);
        view13.setText(android.text.Html.fromHtml("*", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view13 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view13.setMargins(8, 8, 8, 8);
        view13.setLayoutParams(lp_view13);
        view13.setBackgroundColor(Color.parseColor("#F9E2AF"));
        view13.setTextSize(28);
        view13.setTextColor(Color.parseColor("#11111B"));
        view13.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_mul_onclick");
        });
        this.btn_mul = view13;
        view9.addView(view13);
        LinearLayout view14 = new LinearLayout(this);
        view14.setOrientation(LinearLayout.HORIZONTAL);
        view14.setGravity(Gravity.CENTER);
        LinearLayout.LayoutParams lp_view14 = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view14.setMargins(4, 4, 4, 4);
        view14.setLayoutParams(lp_view14);
        view1.addView(view14);
        Button view15 = new Button(this);
        view15.setText(android.text.Html.fromHtml("4", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view15 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view15.setMargins(8, 8, 8, 8);
        view15.setLayoutParams(lp_view15);
        view15.setBackgroundColor(Color.parseColor("#313244"));
        view15.setTextSize(28);
        view15.setTextColor(Color.parseColor("#CDD6F4"));
        view15.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_4_onclick");
        });
        this.btn_4 = view15;
        view14.addView(view15);
        Button view16 = new Button(this);
        view16.setText(android.text.Html.fromHtml("5", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view16 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view16.setMargins(8, 8, 8, 8);
        view16.setLayoutParams(lp_view16);
        view16.setBackgroundColor(Color.parseColor("#313244"));
        view16.setTextSize(28);
        view16.setTextColor(Color.parseColor("#CDD6F4"));
        view16.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_5_onclick");
        });
        this.btn_5 = view16;
        view14.addView(view16);
        Button view17 = new Button(this);
        view17.setText(android.text.Html.fromHtml("6", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view17 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view17.setMargins(8, 8, 8, 8);
        view17.setLayoutParams(lp_view17);
        view17.setBackgroundColor(Color.parseColor("#313244"));
        view17.setTextSize(28);
        view17.setTextColor(Color.parseColor("#CDD6F4"));
        view17.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_6_onclick");
        });
        this.btn_6 = view17;
        view14.addView(view17);
        Button view18 = new Button(this);
        view18.setText(android.text.Html.fromHtml("-", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view18 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view18.setMargins(8, 8, 8, 8);
        view18.setLayoutParams(lp_view18);
        view18.setBackgroundColor(Color.parseColor("#F9E2AF"));
        view18.setTextSize(28);
        view18.setTextColor(Color.parseColor("#11111B"));
        view18.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_sub_onclick");
        });
        this.btn_sub = view18;
        view14.addView(view18);
        LinearLayout view19 = new LinearLayout(this);
        view19.setOrientation(LinearLayout.HORIZONTAL);
        view19.setGravity(Gravity.CENTER);
        LinearLayout.LayoutParams lp_view19 = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view19.setMargins(4, 4, 4, 4);
        view19.setLayoutParams(lp_view19);
        view1.addView(view19);
        Button view20 = new Button(this);
        view20.setText(android.text.Html.fromHtml("1", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view20 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view20.setMargins(8, 8, 8, 8);
        view20.setLayoutParams(lp_view20);
        view20.setBackgroundColor(Color.parseColor("#313244"));
        view20.setTextSize(28);
        view20.setTextColor(Color.parseColor("#CDD6F4"));
        view20.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_1_onclick");
        });
        this.btn_1 = view20;
        view19.addView(view20);
        Button view21 = new Button(this);
        view21.setText(android.text.Html.fromHtml("2", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view21 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view21.setMargins(8, 8, 8, 8);
        view21.setLayoutParams(lp_view21);
        view21.setBackgroundColor(Color.parseColor("#313244"));
        view21.setTextSize(28);
        view21.setTextColor(Color.parseColor("#CDD6F4"));
        view21.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_2_onclick");
        });
        this.btn_2 = view21;
        view19.addView(view21);
        Button view22 = new Button(this);
        view22.setText(android.text.Html.fromHtml("3", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view22 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view22.setMargins(8, 8, 8, 8);
        view22.setLayoutParams(lp_view22);
        view22.setBackgroundColor(Color.parseColor("#313244"));
        view22.setTextSize(28);
        view22.setTextColor(Color.parseColor("#CDD6F4"));
        view22.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_3_onclick");
        });
        this.btn_3 = view22;
        view19.addView(view22);
        Button view23 = new Button(this);
        view23.setText(android.text.Html.fromHtml("+", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view23 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view23.setMargins(8, 8, 8, 8);
        view23.setLayoutParams(lp_view23);
        view23.setBackgroundColor(Color.parseColor("#F9E2AF"));
        view23.setTextSize(28);
        view23.setTextColor(Color.parseColor("#11111B"));
        view23.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_add_onclick");
        });
        this.btn_add = view23;
        view19.addView(view23);
        LinearLayout view24 = new LinearLayout(this);
        view24.setOrientation(LinearLayout.HORIZONTAL);
        view24.setGravity(Gravity.CENTER);
        LinearLayout.LayoutParams lp_view24 = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view24.setMargins(4, 4, 4, 4);
        view24.setLayoutParams(lp_view24);
        view1.addView(view24);
        Button view25 = new Button(this);
        view25.setText(android.text.Html.fromHtml("0", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view25 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view25.setMargins(8, 8, 8, 8);
        view25.setLayoutParams(lp_view25);
        view25.setBackgroundColor(Color.parseColor("#313244"));
        view25.setTextSize(28);
        view25.setTextColor(Color.parseColor("#CDD6F4"));
        view25.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_0_onclick");
        });
        this.btn_0 = view25;
        view24.addView(view25);
        Button view26 = new Button(this);
        view26.setText(android.text.Html.fromHtml(".", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view26 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view26.setMargins(8, 8, 8, 8);
        view26.setLayoutParams(lp_view26);
        view26.setBackgroundColor(Color.parseColor("#313244"));
        view26.setTextSize(28);
        view26.setTextColor(Color.parseColor("#CDD6F4"));
        view26.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_dot_onclick");
        });
        this.btn_dot = view26;
        view24.addView(view26);
        Button view27 = new Button(this);
        view27.setText(android.text.Html.fromHtml("DEL", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view27 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view27.setMargins(8, 8, 8, 8);
        view27.setLayoutParams(lp_view27);
        view27.setBackgroundColor(Color.parseColor("#F38BA8"));
        view27.setTextSize(28);
        view27.setTextColor(Color.parseColor("#11111B"));
        view27.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_del_onclick");
        });
        this.btn_del = view27;
        view24.addView(view27);
        Button view28 = new Button(this);
        view28.setText(android.text.Html.fromHtml("=", android.text.Html.FROM_HTML_MODE_COMPACT));
        LinearLayout.LayoutParams lp_view28 = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.000000f);
        lp_view28.setMargins(8, 8, 8, 8);
        view28.setLayoutParams(lp_view28);
        view28.setBackgroundColor(Color.parseColor("#A6E3A1"));
        view28.setTextSize(28);
        view28.setTextColor(Color.parseColor("#11111B"));
        view28.setOnClickListener(v -> {
            NativeBridge.sendEvent("btn_eq_onclick");
        });
        this.btn_eq = view28;
        view24.addView(view28);
        return view1;
    }

    public void updateWidgetText(String id, String text) {
        if (id.equals("display")) { if (this.display != null) this.display.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_adv")) { if (this.btn_adv != null) this.btn_adv.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_clear")) { if (this.btn_clear != null) this.btn_clear.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_lp")) { if (this.btn_lp != null) this.btn_lp.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_rp")) { if (this.btn_rp != null) this.btn_rp.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_div")) { if (this.btn_div != null) this.btn_div.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_7")) { if (this.btn_7 != null) this.btn_7.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_8")) { if (this.btn_8 != null) this.btn_8.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_9")) { if (this.btn_9 != null) this.btn_9.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_mul")) { if (this.btn_mul != null) this.btn_mul.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_4")) { if (this.btn_4 != null) this.btn_4.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_5")) { if (this.btn_5 != null) this.btn_5.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_6")) { if (this.btn_6 != null) this.btn_6.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_sub")) { if (this.btn_sub != null) this.btn_sub.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_1")) { if (this.btn_1 != null) this.btn_1.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_2")) { if (this.btn_2 != null) this.btn_2.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_3")) { if (this.btn_3 != null) this.btn_3.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_add")) { if (this.btn_add != null) this.btn_add.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_0")) { if (this.btn_0 != null) this.btn_0.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_dot")) { if (this.btn_dot != null) this.btn_dot.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_del")) { if (this.btn_del != null) this.btn_del.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
        if (id.equals("btn_eq")) { if (this.btn_eq != null) this.btn_eq.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }
    }

    public void animateWidget(String id, String property, float to, int durationMs) {
        if (id.equals("display")) { if (this.display != null) ObjectAnimator.ofFloat(this.display, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_adv")) { if (this.btn_adv != null) ObjectAnimator.ofFloat(this.btn_adv, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_clear")) { if (this.btn_clear != null) ObjectAnimator.ofFloat(this.btn_clear, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_lp")) { if (this.btn_lp != null) ObjectAnimator.ofFloat(this.btn_lp, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_rp")) { if (this.btn_rp != null) ObjectAnimator.ofFloat(this.btn_rp, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_div")) { if (this.btn_div != null) ObjectAnimator.ofFloat(this.btn_div, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_7")) { if (this.btn_7 != null) ObjectAnimator.ofFloat(this.btn_7, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_8")) { if (this.btn_8 != null) ObjectAnimator.ofFloat(this.btn_8, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_9")) { if (this.btn_9 != null) ObjectAnimator.ofFloat(this.btn_9, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_mul")) { if (this.btn_mul != null) ObjectAnimator.ofFloat(this.btn_mul, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_4")) { if (this.btn_4 != null) ObjectAnimator.ofFloat(this.btn_4, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_5")) { if (this.btn_5 != null) ObjectAnimator.ofFloat(this.btn_5, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_6")) { if (this.btn_6 != null) ObjectAnimator.ofFloat(this.btn_6, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_sub")) { if (this.btn_sub != null) ObjectAnimator.ofFloat(this.btn_sub, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_1")) { if (this.btn_1 != null) ObjectAnimator.ofFloat(this.btn_1, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_2")) { if (this.btn_2 != null) ObjectAnimator.ofFloat(this.btn_2, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_3")) { if (this.btn_3 != null) ObjectAnimator.ofFloat(this.btn_3, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_add")) { if (this.btn_add != null) ObjectAnimator.ofFloat(this.btn_add, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_0")) { if (this.btn_0 != null) ObjectAnimator.ofFloat(this.btn_0, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_dot")) { if (this.btn_dot != null) ObjectAnimator.ofFloat(this.btn_dot, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_del")) { if (this.btn_del != null) ObjectAnimator.ofFloat(this.btn_del, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("btn_eq")) { if (this.btn_eq != null) ObjectAnimator.ofFloat(this.btn_eq, property, to).setDuration(durationMs).start(); return; }
        if (id.equals("calc_layout")) { if (this.calc_layout != null) ObjectAnimator.ofFloat(this.calc_layout, property, to).setDuration(durationMs).start(); return; }
    }

    public String getWidgetText(String id) {
        if (id.equals("display") && this.display != null) return this.display.getText().toString();
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
