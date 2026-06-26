<img src="./assets/logo.png" width="100" />

# Termodoro
A terminal Pomodoro timer for focused work sessions.

![Demo](./assets/demo-small-screen.png)

Introducing a fresh take on productivity — a minimalist terminal app for the tech-savvy professional. Keep your focus on the task without the clutter of traditional apps.

### Features
- ✅ Pomodoro timer with focus + break sessions
- ✅ 7 visual animations to choose from
- ✅ 6 sound options (including silent)
- ✅ Config presets (Classic, Short, Long, Custom)
- ✅ Desktop notifications
- ❌ Todo List
- ❌ Notes Taking
- ❌ Stats


---

## Visual Options

Visual effects make every session interesting. Choose one on the config screen.


<details open>
<summary>🌲 Tree</summary>
<br>

<img src="./assets/tree.gif" width="400" />

A random procedurally generated tree every time you start a session.
</details>

<details open>
<summary>🛶 Flow</summary>
<br>

<img src="./assets/flow.gif" width="400" />

A rower crossing the _"Time River"_.
</details>

<details open>
<summary>☕ Coffee</summary>
<br>

<img src="./assets/coffee.gif" width="400" />

A coffee mug that fills up as time passes.
</details>

<details open>
<summary>🔥 Campfire</summary>
<br>

Flickering ASCII flames with a glowing log. The fire cycles through yellow, orange, and red — updates every 250ms.
</details>

<details open>
<summary>🌧️ Rain</summary>
<br>

Raindrops falling down the screen in shades of blue. A calming, ambient effect — scrolls every 100ms.
</details>

<details open>
<summary>🌅 Sunrise</summary>
<br>

A sun that rises from the horizon as your session progresses (0% = dawn red, 100% = noon yellow). Tied to timer percentage.
</details>

<details open>
<summary>🕐 BigClock</summary>
<br>

Full-screen pixel-art digit display with a live progress bar and session label. Bypasses the normal UI box for a clean, immersive view.

```
 ██████  ██████   ██   ██████  ██████
 ██  ██  ██  ██  ████  ██      ██
 ██  ██  ██  ██   ██   ████    ████
 ██████  ██████   ██       ██      ██
                  ██   ██  ██  ██  ██
                  ██   ██████  ██████
```

Digits dim when the timer is paused.
</details>


---

## How to install Termodoro

__Linux/Mac:__

```sh
git clone https://github.com/hrushik98/termodoro
cd termodoro
go build -o binary/termodoro
```

__Run:__
```sh
./binary/termodoro
```

Or with make:

```sh
make run
```

---

## Controls

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate config options |
| `←` / `→` | Adjust selected value |
| `Enter` | Start session |
| `Space` | Pause / Resume / Next session |
| `r` | Restart timer |
| `n` | New session (back to config) |
| `q` | Quit |
| `Ctrl+B` | Back to config from timer |

<br/>

### Credits

This project is a fork of [aimssh](https://github.com/sairash/aimssh) by [Sairash](https://sairashgautam.com.np).
