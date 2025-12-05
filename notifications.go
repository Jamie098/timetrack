package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

func sendNotification(title, message string) {
	if runtime.GOOS == "windows" {
		// PowerShell toast notification
		script := fmt.Sprintf(`
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null
$template = @"
<toast>
    <visual>
        <binding template="ToastText02">
            <text id="1">%s</text>
            <text id="2">%s</text>
        </binding>
    </visual>
</toast>
"@
$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($template)
$toast = [Windows.UI.Notifications.ToastNotification]::new($xml)
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("TimeTrack").Show($toast)
`, title, message)
		cmd := exec.Command("powershell", "-Command", script)
		cmd.Run()
	} else if runtime.GOOS == "darwin" {
		// macOS notification
		script := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
		exec.Command("osascript", "-e", script).Run()
	} else {
		// Linux notification
		exec.Command("notify-send", title, message).Run()
	}
}
