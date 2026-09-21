package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type WttrResponse struct {
	CurrentCondition []struct {
		TempC       string `json:"temp_C"`
		FeelsLikeC  string `json:"FeelsLikeC"`
		WeatherDesc []struct {
			Value string `json:"value"`
		} `json:"weatherDesc"`
		WindspeedKmph string `json:"windspeedKmph"`
		PrecipMM      string `json:"precipMM"`
		Humidity      string `json:"humidity"`
	} `json:"current_condition"`
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

// ANSI colors. Disabled automatically if NO_COLOR is set — kept simple
// rather than pulling in a terminal detection library for a joke program.
const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
)

var colorsEnabled = os.Getenv("NO_COLOR") == ""

func colorize(code, s string) string {
	if !colorsEnabled {
		return s
	}
	return code + s + colorReset
}

func fetchWeather(city string) (*WttrResponse, error) {
	endpoint := fmt.Sprintf("https://wttr.in/%s?format=j1", url.PathEscape(city))

	resp, err := httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wttr.in returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var w WttrResponse
	if err := json.Unmarshal(body, &w); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &w, nil
}

func pick(options []string) string {
	return options[rand.Intn(len(options))]
}

// tempSass covers the actual temperature.
func tempSass(tempC float64) (string, string) {
	switch {
	case tempC < -10:
		return colorBlue, pick([]string{
			"It's cold enough to question your life choices.",
			"Negative double digits. Bold of you to go outside.",
			"Your breath is now a weapon.",
			"This is soup weather, and by soup I mean your blood if you stay out long enough.",
		})
	case tempC < 0:
		return colorBlue, pick([]string{
			"Below freezing. Congrats on the visible breath.",
			"Cold enough to matter, not cold enough to be interesting.",
			"Ice is forming somewhere and it's judging you.",
		})
	case tempC < 10:
		return colorCyan, pick([]string{
			"Chilly. Jacket weather if you're sane.",
			"Cool enough that you'll regret that t-shirt.",
			"Bring a hoodie. Trust me.",
			"Sweater season, no negotiation.",
		})
	case tempC < 20:
		return colorGray, pick([]string{
			"Mild. Suspiciously pleasant.",
			"Perfectly average. Like a Tuesday.",
			"This is fine. Everything is fine.",
			"Weather so unremarkable it doesn't even deserve a complaint.",
		})
	case tempC < 30:
		return colorYellow, pick([]string{
			"Warm. Shorts weather if you're brave.",
			"Nice out. Don't waste it indoors.",
			"Vitamin D opportunity. Take it.",
			"Prime patio weather. Act accordingly.",
		})
	case tempC < 35:
		return colorRed, pick([]string{
			"Hot. Stay hydrated or perish.",
			"It's melting season. Good luck.",
			"The sun has chosen violence today.",
		})
	default:
		return colorRed, pick([]string{
			"This isn't weather, it's a warning label.",
			"The asphalt is basically lava now.",
			"Everything outside is a hazard. Stay in. Drink water. Pray.",
		})
	}
}

// feelsLikeSass compares actual vs. feels-like temperature — the gap is
// usually where the real complaint lives.
func feelsLikeSass(tempC, feelsC float64) string {
	diff := feelsC - tempC
	switch {
	case diff <= -5:
		return fmt.Sprintf("Feels like %.0f°C — the wind is personally offended by you.", feelsC)
	case diff >= 5:
		return fmt.Sprintf("Feels like %.0f°C — humidity is cosplaying as a sauna.", feelsC)
	default:
		return "" // not different enough to bother mentioning
	}
}

func windSass(windF float64) string {
	switch {
	case windF < 5:
		return "Barely any wind. Your hair is safe."
	case windF < 20:
		return "Some breeze. Manageable."
	case windF < 40:
		return "Windy. Your hair is already a lost cause."
	case windF < 60:
		return "Hold onto something. Or someone."
	default:
		return "That's not wind, that's an eviction notice for anything not bolted down."
	}
}

func precipSass(precipF float64) string {
	if precipF <= 0 {
		return pick([]string{
			"No rain. Yet.",
			"Dry for now. Enjoy the illusion of safety.",
			"Sky is holding it in. For now.",
		})
	}
	switch {
	case precipF < 2.5:
		return pick([]string{
			"Light rain. Barely counts, but bring a hood.",
			"A light drizzle, just enough to be annoying.",
		})
	case precipF < 10:
		return pick([]string{
			"It's raining. Obviously.",
			"Precipitation detected. Bring an umbrella or don't.",
			"Wet stuff falling from the sky. Shocking.",
		})
	default:
		return pick([]string{
			"This is less 'rain' and more 'the sky giving up.'",
			"Ark-building weather. Start recruiting animals.",
		})
	}
}

func humiditySass(humidityF float64) string {
	switch {
	case humidityF >= 80:
		return "Humidity's high enough to make the air feel personally hostile."
	case humidityF <= 25:
		return "Bone-dry air. Static shocks incoming."
	default:
		return ""
	}
}

func sass(tempC, feelsC float64, desc, wind, precip, humidity string) string {
	var b strings.Builder

	tColor, tLine := tempSass(tempC)
	b.WriteString(colorize(tColor+colorBold, tLine))
	b.WriteString(" ")

	if fl := feelsLikeSass(tempC, feelsC); fl != "" {
		b.WriteString(colorize(colorGray, fl))
		b.WriteString(" ")
	}

	windF, _ := strconv.ParseFloat(wind, 64)
	b.WriteString(windSass(windF))
	b.WriteString(" ")

	precipF, _ := strconv.ParseFloat(precip, 64)
	b.WriteString(precipSass(precipF))

	humidityF, _ := strconv.ParseFloat(humidity, 64)
	if hLine := humiditySass(humidityF); hLine != "" {
		b.WriteString(" ")
		b.WriteString(hLine)
	}

	b.WriteString(fmt.Sprintf("\n\n%s Reported condition: %s. Humidity: %s%%. Whatever that means.",
		colorize(colorGray, "—"), desc, humidity))

	return b.String()
}

func main() {
	city := "Berlin"
	if len(os.Args) > 1 {
		city = strings.Join(os.Args[1:], " ")
	}

	fmt.Printf("Fetching weather for %s...\n\n", city)

	w, err := fetchWeather(city)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not fetch weather: %v\n", err)
		os.Exit(1)
	}

	if len(w.CurrentCondition) == 0 {
		fmt.Fprintln(os.Stderr, "No weather data returned. Weird city?")
		os.Exit(1)
	}

	c := w.CurrentCondition[0]

	tempC, err := strconv.ParseFloat(c.TempC, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse temperature: %v\n", err)
		os.Exit(1)
	}
	feelsC, err := strconv.ParseFloat(c.FeelsLikeC, 64)
	if err != nil {
		feelsC = tempC // degrade gracefully rather than crash over a cosmetic field
	}

	desc := "unknown"
	if len(c.WeatherDesc) > 0 {
		desc = c.WeatherDesc[0].Value
	}

	fmt.Printf("%s\n\n", colorize(colorBold, fmt.Sprintf("=== %s ===", city)))
	fmt.Printf("%s\n", sass(tempC, feelsC, desc, c.WindspeedKmph, c.PrecipMM, c.Humidity))
}
