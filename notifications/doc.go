// Package notifications sends the eye-care and stand-up reminders to the
// user. notifier.go (non-Windows) wraps fyne.App.SendNotification;
// notifier_windows.go calls Shell_NotifyIcon directly, since
// SendNotification shells out to PowerShell there, which AV blocks on
// locked-down machines.
package notifications
