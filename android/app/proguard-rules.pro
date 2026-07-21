# Keep JNI bridge classes so Go can call Java methods
-keep class *.NativeBridge { *; }
-keep class *.MainActivity { *; }
-keep class *.Go2APKBroadcastReceiver { *; }
