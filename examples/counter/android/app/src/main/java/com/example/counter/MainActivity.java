package com.example.counter;

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
    private LinearLayout main_layout;
    private TextView title;
    private TextView counter_text;
    private Button increment_button;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setTitle("Counter");
        setContentView(createAppView());
    }

    private View createAppView() {
        LinearLayout view1 = new LinearLayout(this);
        view1.setOrientation(LinearLayout.VERTICAL);
        view1.setGravity(Gravity.CENTER_HORIZONTAL);
        view1.setPadding(24, 24, 24, 24);
        this.main_layout = view1;
        LinearLayout rootView = view1;
        TextView view2 = new TextView(this);
        view2.setText("Counter App");
        view2.setTextSize(24);
        this.title = view2;
        view1.addView(view2);
        TextView view3 = new TextView(this);
        view3.setText("0");
        view3.setTextSize(24);
        this.counter_text = view3;
        view1.addView(view3);
        Button view4 = new Button(this);
        view4.setText("Increment");
        view4.setOnClickListener(v -> {
            NativeBridge.sendEvent("increment_button_onclick");
        });
        this.increment_button = view4;
        view1.addView(view4);
        return rootView;
    }
}
