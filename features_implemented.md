# Features Implemented

Here is a summary of the new features and improvements implemented in Termodoro.

## 1. Professional Repository Credits
* **Fork Acknowledgment**: Updated the footer in [README.md](file:///Users/hrushikreddy/Desktop/termodoro/README.md) to explicitly state that the project is a fork of [aimssh](https://github.com/sairash/aimssh) by [Sairash](https://sairashgautam.com.np).
* **Removed Clutter**: Deleted the informal "Made with ❤️" branding in favor of a clean, structured **Credits** section.

## 2. Interactive Configuration & Presets System
* **Unified Screen**: Combined the separate, sequential inputs (time, visual options) into a single, cohesive **Configuration & Presets** dashboard on startup.
* **Predefined Presets**: Added four start profiles:
  * **Classic Pomodoro**: 25-minute focus session, 5-minute break, "Melody" sound, "Tree" animation.
  * **Short Focus**: 15-minute focus session, 3-minute break, "High Beep" sound, "Coffee" animation.
  * **Long Focus**: 50-minute focus session, 10-minute break, "Double Beep" sound, "Flow" animation.
  * **Custom**: Automatically activated if any parameter is adjusted individually.
* **Controls**:
  * <kbd>↑</kbd> / <kbd>↓</kbd> (or <kbd>j</kbd> / <kbd>k</kbd> / <kbd>Tab</kbd>) to navigate fields.
  * <kbd>←</kbd> / <kbd>→</kbd> (or <kbd>h</kbd> / <kbd>l</kbd>) to cycle settings and increment/decrement durations.
  * <kbd>Enter</kbd> to launch the session.

## 3. Customizable Notification Sounds
* **Sound Profiles**: Integrated six different sound options:
  1. **Default Alert** (Standard system beep)
  2. **High Beep** (High-pitched tone)
  3. **Low Beep** (Low-pitched tone)
  4. **Double Beep** (Two rapid notes)
  5. **Melody** (An ascending retro chime)
  6. **Silent** (Desktop banner only, no sound)
* **Live Previews**: Changing the sound setting in the configuration menu plays a brief preview of the audio so you know what it sounds like.
* **Asynchronous Execution**: Audio runs on a separate goroutine so it doesn't freeze the TUI display.

## 4. Focus & Break Session transitions
* **Pomodoro Loop**: Rather than terminating when the timer runs out, Termodoro now loops between focus and break:
  * **Focus Session Ends**: Plays the chosen sound, fires a desktop notification, and prompts you to start the break.
  * **Break Transition**: Pressing <kbd>Space</kbd> starts the Break timer with the configured duration.
  * **Break Session Ends**: Plays the sound, triggers a notification, and prompts you to return to the focus session.
* **Manual Control**: Allows you to step away from your computer between sessions without the break timer starting automatically.
* **Menu Return**: Pressing <kbd>n</kbd> at any point returns you to the configuration screen.

## 5. Code Health & Bugfixes
* **Mutex Copy Fix**: Fixed a compiler warning in `ascii_generator/tree.go` where `sync.Mutex` was copied by value. Updated `initCanvas` and `GenerateTree` to return pointers (`*Tree`), rendering the repository 100% compliant with `go vet`.

## 6. Rename: AimSSH → Termodoro
* Renamed project from **AimSSH** → **Termodoro** throughout (module name, binary name, UI text, README)
* Removed the SSH server component entirely (`app/server.go`, `app/config.go` deleted)
* Removed SSH-related Go dependencies (`charmbracelet/ssh`, `charmbracelet/wish`, `charmbracelet/log`) from `go.mod`
* Removed the `ssh` subcommand from the CLI
* Updated all `aimssh/…` import paths to `termodoro/…`

## 7. New Visual Animations

### 🔥 Campfire (`ascii_generator/campfire.go`)
* 4-frame flickering ASCII fire — flames cycle through yellow tip → orange mid → red base
* Glowing log rendered below in brown / ember-orange
* Updates every **250 ms** via a background goroutine

### 🌧️ Rain (`ascii_generator/rain.go`)
* Raindrops (`|`, `/`, `'`, `.`) scroll downward across a 40×15 grid
* Three shades of blue alternate row-by-row for a depth effect
* Updates every **100 ms** via a background goroutine

### 🌅 Sunrise (`ascii_generator/sunrise.go`)
* Sun rises from the horizon (0% elapsed) to the top of the sky (100% elapsed)
* Color shifts red → orange → yellow as the session progresses
* Horizon line (`~~~`) and ground line (`===`) anchor the scene
* Purely percentage-driven, no goroutine needed

### 🕐 BigClock (`app/bigclock_view.go`)
* Full-screen 5×7 **pixel-art digit display** — each pixel rendered as `██` (2 terminal chars)
* MM:SS readout is 58 chars wide; centers automatically to any terminal width
* Digits **dim** (`#3d3878`) when paused, **brighten** (`#7C6AE8`) when running
* Animated **progress bar** below: filled `█` in purple + empty `░` in dark + `%` label
* Session label shows "focus session" / "break session" / "paused" / end-of-session messages
* Bypasses the normal boxed AppStyle and renders full-screen via `lipgloss.Place`
