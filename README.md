# bad_weather ⛈️

[![Go](https://img.shields.io/badge/Go-1.27-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg?style=flat-square)](LICENSE)

**A weather CLI that tells you the forecast and then makes fun of you for asking.**

Under the hood it's a completely ordinary [wttr.in](https://wttr.in) client. On top of that it runs your local conditions through a deadpan commentary engine, because "12°C, light rain" doesn't really convey the emotional weight of standing outside in it.

---

### `~ > ./bad_weather berlin`

```
Fetching weather for berlin...

=== berlin ===

Chilly. Jacket weather if you're sane. Some breeze. Manageable. It's raining. Obviously.

— Reported condition: Light rain. Humidity: 76%. Whatever that means.
```

---

## What it does

1. Hits `wttr.in`'s JSON API for the city you give it (defaults to Berlin if you give it nothing)
2. Reads temperature, feels-like temperature, wind speed, precipitation, and humidity
3. Runs each of those through its own sarcasm tier — cold snaps, heatwaves, wind that will ruin your hair, rain that isn't quite committing, humidity that's decided to be personal about it
4. Prints the verdict, in color, to your terminal

No API key, no config file, no accounts. Just weather, and consequences.

## Requirements

- [Go 1.21+](https://go.dev/dl/)
- An internet connection (`wttr.in` does the actual meteorology; this project just has opinions about the output)

## Installation & running

```bash
git clone git@github.com:etokiyra/bad_weather.git
cd bad_weather
go build -o bad_weather
```

Then run it against any city:

```bash
./bad_weather Tokyo
./bad_weather "New York"
./bad_weather   # defaults to Berlin
```

Or skip the build step entirely:

```bash
go run . Warsaw
```

## Disabling color

Terminal escape codes are on by default. If you're piping output somewhere that doesn't appreciate ANSI, or you just have taste:

```bash
NO_COLOR=1 ./bad_weather Berlin
```

## How the sass works

Each metric — temperature, wind, precipitation, humidity, and the gap between actual and "feels like" temperature — has its own tiered set of one-liners, picked at random within the tier that matches current conditions. Extreme weather gets proportionally more extreme commentary. This is, as far as the author is aware, the correct way to build a weather app.

## License

GPL-3.0 — see [LICENSE](LICENSE). Do whatever you want with it, just keep it open.
