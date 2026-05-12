# Cornealius-Eyeworth — Project Context

## Overview

A small Windows desktop application (WinForms) that reminds users to rest their eyes on a schedule. It uses Quartz for scheduling and Windows toast notifications for reminders. The app is written in C# and targets .NET 10 (net10.0-windows10.0.17763.0).

## Quick facts

- **Project name:** Cornealius-Eyeworth
- **Primary language:** C# (.NET)
- **Target framework:** net10.0-windows10.0.17763.0
- **UI framework:** Windows Forms
- **Entrypoint:** `Program.cs` → `AppHost.RunAsync()`
- **Primary config:** `config.json` (fields: `CronExpression`, `Language`)
- **Key packages:** Microsoft.Toolkit.Uwp.Notifications (7.1.3), Quartz (3.18.1)

## Layout & responsibilities

- **Application/** — `AppHost.cs`: composition root and app lifecycle.
- **Configuration/** — config model + `JsonConfigRepository` for reading/writing `config.json`.
- **Parsing/** — cron and minutes parsing (`CronExpressionParser`, `MinutesInputParser`).
- **Scheduling/** — Quartz integration and scheduling logic (`SchedulerService`, `EyeworthJob`, `NextTriggerProvider`).
- **Notifications/** — toast notifications (`NotificationService`).
- **UI/** — WinForms UI (`MainForm`, dialogs for schedule and language).
- **Localization/** — ResX resources and generated designer files.

## Important files (entry points)

- `Program.cs`
- `CornealiusEyeworth.csproj`
- `config.json`
- `Application/AppHost.cs`
- `Scheduling/SchedulerService.cs`
- `Notifications/NotificationService.cs`
- `Parsing/CronExpressionParser.cs`
- `UI/MainForm.cs`

## Build & run

- Build: `dotnet build -c Debug` or `dotnet build -c Release`.
- Run: execute the built binary from `bin/Debug/net10.0-windows10.0.17763.0/`.
- Note: `config.json` is copied to the output; edit and restart to change runtime behavior.

## Config

- `config.json` contains `CronExpression` and `Language`.
- Default CronExpression example: `0 20,40,55 * * * ?` (triggers at 20, 40 and 55 minutes each hour).

## Search tokens & quick tasks

- Search tokens: `AppHost`, `SchedulerService`, `NotificationService`, `CronExpressionParser`, `JsonConfigRepository`, `MainForm`, `config.json`.
- Common edits:
  - Change schedule: update `config.json` or `Scheduling/` code.
  - Change notifications: edit `Notifications/NotificationService.cs`.
  - Add language: add `.resx` in `Localization/` and update `LanguageDialog`.

## Known constraints & gaps

- Windows-only (WinForms + Windows toasts).
- No automated tests in the repo.
- global.json pins SDK `9.0.100` — may require matching SDK locally.

---
This file is a short, focused map of the repository to help future agents and contributors quickly find where to change scheduling, notifications, UI, and localization.
