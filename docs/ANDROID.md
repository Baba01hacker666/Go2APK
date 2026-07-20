# Android Platform APIs in Go2APK

Go2APK provides access to core Android platform features (Intents, Broadcasts, Permissions) directly from Go.

These APIs are available in the `github.com/Baba01hacker666/Go2APK/android` package.

## 1. Permissions

### Declaring Manifest Permissions
If your app needs a system permission, declare it at the package level or in `main()`. Go2APK will automatically inject it into the generated `AndroidManifest.xml`.

```go
android.Permission("android.permission.CAMERA")
android.Permission("android.permission.INTERNET")
```

### Requesting Runtime Permissions
For dangerous permissions (like Camera or Location), you must request them at runtime:

```go
android.RequestPermission("android.permission.CAMERA", func(granted bool) {
	if granted {
		// Start camera
	} else {
		// Show error message
	}
})
```

## 2. Intents and Activities

### Launching External Apps
You can construct and fire Android Intents to open other apps, URLs, or system screens:

```go
// Open a website in the default browser
android.StartActivity(android.Intent{
	Action: "android.intent.action.VIEW",
	Data:   "https://github.com",
})

// Open a specific app package
android.StartActivity(android.Intent{
	Package: "com.example.otherapp",
})
```

## 3. Broadcast Receivers

### Listening to System Broadcasts
You can register a broadcast receiver to listen for system-wide events (like battery low, screen unlock, etc). Go2APK will automatically generate the Java `BroadcastReceiver` and wire it up in the Manifest.

```go
android.BroadcastReceiver("android.intent.action.BATTERY_LOW", "battery_low", func() {
	println("Battery is low!")
})
```

### Sending Broadcasts
You can also broadcast your own events to other apps or system components:

```go
// Simple broadcast
android.SendBroadcast("com.my.app.CUSTOM_EVENT")

// Broadcast with string extras
android.SendBroadcastWithExtras("com.my.app.DATA_EVENT", map[string]string{
	"user_id": "12345",
	"status":  "active",
})
```
