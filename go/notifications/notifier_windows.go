//go:build windows

package notifications

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"fyne.io/fyne/v2"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
)

func pickRandom(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[rand.Intn(len(items))]
}

const notificationTemplate = `$title = "%s"
$content = "%s"
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null
$template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastText02)
$toastXml = [xml] $template.GetXml()
$toastXml.GetElementsByTagName("text")[0].AppendChild($toastXml.CreateTextNode($title)) > $null
$toastXml.GetElementsByTagName("text")[1].AppendChild($toastXml.CreateTextNode($content)) > $null
$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($toastXml.OuterXml)
$toast = [Windows.UI.Notifications.ToastNotification]::new($xml)
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("%s").Show($toast);`

func escapePS(s string) string {
	s = strings.ReplaceAll(s, "`", "``")
	return strings.ReplaceAll(s, "\"", "`\"")
}

func sendToast(a fyne.App, title, content string) {
	appID := a.UniqueID()
	if appID == "" || strings.HasPrefix(appID, "missing-id") {
		appID = a.Metadata().Name
	}
	script := fmt.Sprintf(notificationTemplate, escapePS(title), escapePS(content), escapePS(appID))

	tmp, err := os.CreateTemp("", "fyne-notify-*.ps1")
	if err != nil {
		fyne.LogError("notification: cannot create temp script", err)
		return
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err := tmp.WriteString(script); err != nil {
		tmp.Close()
		fyne.LogError("notification: cannot write script", err)
		return
	}
	tmp.Close()

	launch := "(Get-Content -Encoding UTF8 -Path " + tmpPath + " -Raw) | Invoke-Expression"
	cmd := exec.Command("PowerShell", "-ExecutionPolicy", "Bypass", launch)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		fyne.LogError("notification: toast script failed", err)
	}
}

// SendStartup delivers the "on duty" notification shown when the app starts.
func SendStartup(a fyne.App) {
	go sendToast(a, i18n.Active.NotificationOnDuty, pickRandom(i18n.Active.NotificationQuips))
}

// SendMinimizedToTray delivers a "still running" hint when the window is hidden.
func SendMinimizedToTray(a fyne.App) {
	go sendToast(a, i18n.Active.AppName, i18n.Active.NotificationMinimizedToTray)
}

// SendReminder delivers an eye-rest reminder notification.
func SendReminder(a fyne.App) {
	go sendToast(a, pickRandom(i18n.Active.NotificationReminders), pickRandom(i18n.Active.NotificationQuips))
}
