package com.example.go2apkapp;

final class NativeCalculator {
    private static final String LIBRARY_NAME = "go2apkcalc";

    static {
        System.loadLibrary(LIBRARY_NAME);
    }

    private NativeCalculator() {
    }

    static String calculate(String left, String operator, String right) {
        return calculateWithGo(left, operator, right);
    }

    private static native String calculateWithGo(String left, String operator, String right);
}
