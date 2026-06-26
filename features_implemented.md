# Features Implemented

This document consolidates all the features and improvements implemented across all development paths of Termodoro (zen-cli).

---

## 1. Unified Configuration & Presets System
* **Interactive TUI Dashboard**: Combined the separate, sequential start prompt screens (time input, session label, visual options) into a single, unified configuration dashboard.
* **Predefined Profiles**: Includes four start profiles that can be quickly toggled:
  * **Classic Pomodoro**: 25m Focus / 5m Break, "Melody" notification, "Tree" animation.
  * **Short Focus**: 15m Focus / 3m Break, "High Beep" notification, "Coffee" animation.
  * **Long Focus**: 50m Focus / 10m Break, "Double Beep" notification, "Flow" animation.
  * **Custom**: Dynamically active if you manually adjust any setting.
* **Controls**: 
  * <kbd>↑</kbd> / <kbd>↓</kbd> (or <kbd>j</kbd> / <kbd>k</kbd> / <kbd>Tab</kbd>) to navigate setting rows.
  * <kbd>←</kbd> / <kbd>→</kbd> (or <kbd>h</kbd> / <kbd>l</kbd>) to cycle selections and increment/decrement durations.
  * <kbd>Enter</kbd> to launch the session immediately.

---

## 2. Customizable Notification Sounds & Previews
* **Multiple Sound Profiles**: Incorporated six notification tones:
  1. **Default Alert** (Standard system ping)
  2. **High Beep** (High-pitched tone)
  3. **Low Beep** (Low-pitched tone)
  4. **Double Beep** (Two rapid notes)
  5. **Melody** (An ascending retro chime)
  6. **Silent** (Desktop banner only, no sound)
* **Real-time Previews**: Changing the sound setting in the configuration menu plays a brief preview of the audio so you can sample your alert choice.
* **Non-blocking Playback**: Sounds play in a background goroutine so that TUI rendering and timer precision are unaffected.

---

## 3. Focus & Break Session Cycle (State Machine)
* **Automated Pomodoro Loop**: The app automatically transitions between work and rest states:
  * **Focus Session Ends**: Triggers your custom sound, fires a desktop notification, and transitions to an ending screen prompting you to rest.
  * **Rest Transition**: Pressing <kbd>Space</kbd> starts the Break timer with the configured duration.
  * **Break Session Ends**: Triggers the sound, fires a notification, and prompts you to return to the focus phase.
* **Manual Rest Launch**: The rest timer waits for you to press Space, letting you step away from your keyboard without losing break time.
* **Menu Return**: Pressing <kbd>n</kbd> on the timer screen stops the active session and returns to the configuration dashboard.

---

## 4. Todoist-Style Todo List
* **Timer Overlay View**: Press <kbd>t</kbd> while the timer is running to open a keyboard-driven task list.
* **Controls**:
  * <kbd>a</kbd> — Add a new task (opens inline text input; <kbd>Enter</kbd> confirms, <kbd>Esc</kbd> cancels).
  * <kbd>Space</kbd> — Toggle a task done (strikethrough & faded).
  * <kbd>d</kbd> — Delete the selected task.
  * <kbd>↑</kbd> / <kbd>↓</kbd> (or <kbd>j</kbd> / <kbd>k</kbd>) — Navigate task items.
* **Priority Levels**: Quick keys to set urgency:
  * <kbd>1</kbd> — P1 (High priority, displayed in red)
  * <kbd>2</kbd> — P2 (Medium priority, displayed in orange)
  * <kbd>3</kbd> — P3 (Low priority, displayed in blue)
  * <kbd>0</kbd> — Clear priority
* **State Persistence**: Tasks are saved automatically to `~/.aimssh_todos.json` (or `~/.aimssh_todos_<username>.json` if running in multi-user/SSH mode).
* **Timer Integration**: The countdown timer continues running in the background. If the timer ends while you are managing tasks, you are automatically notified and returned to the timer display.

---

## 5. Per-Session Notes Taking
* **Access**: Press <kbd>o</kbd> from the Timer screen to open the scratchpad.
* **Persistent Scratchpad**: Powered by a multi-line `textarea` component (14 lines). Notes persist in memory during the active session, allowing you to cycle between the timer, todo list, and notes without losing your text.
* **Header**: Displays the session name as the header (`Notes — <session_name>:`).
* **Return**: Press <kbd>Ctrl+B</kbd> to return to the timer.

---

## 6. Personal Focus Stats Dashboard
* **Access**:
  * Press <kbd>Shift+T</kbd> from the Timer view.
  * Press <kbd>Ctrl+T</kbd> from the Config/Home screen.
* **Key Metrics**:
  * 🔥 **Streak**: Number of consecutive days (ending today or yesterday) with at least one completed focus session.
  * ⏱ **Total**: Cumulative logged focus time across all history (formatted as `Xh Ym`).
  * ✅ **Sessions**: Total completed focus runs.
  * 📅 **Today**: Total focus runs and minutes logged today.
* **7-Day Bar Chart**: A horizontal ascii bar chart displaying focus minutes for the trailing week. Today's bar is highlighted in a distinct custom color.
* **Local Persistence**: Logged focus sessions (Focus session name, duration, completion timestamp) are saved to `~/.aimssh/stats.json`.
* **State Safety**: Shows placeholder messages for empty stats and when run over SSH (where stats are disabled to prevent anonymous user cross-contamination).

---

## 7. Clean Codebase & Go Vet Compliance
* **Mutex Optimization**: Fixed a warning in `ascii_generator/tree.go` where `sync.Mutex` was copied by value during `initCanvas` and `GenerateTree` returns.
* **Pointers**: Refactored the methods to return pointers (`*Tree`), ensuring the repository fully passes `go vet ./...` with zero issues.
