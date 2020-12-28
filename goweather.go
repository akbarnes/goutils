package main

import (
	"flag"
	"fmt"
	"strings"
	"os"
	"log"

	"github.com/alsm/forecastio"
	"github.com/BurntSushi/toml"
)

type weatherConfig struct {
	ApiKey string
	Lon float64
	Lat float64
}

var (
	key      string
	exclude  string
	excludes []string
	lat, lon float64
	extend   bool
	hourly	 bool
	daily    bool
	units    string
)

const cfgPath = "config.toml"

func init() {
	var cfg weatherConfig

	f, err := os.Open(cfgPath)

	if err != nil {
		log.Fatal(fmt.Sprintf("Could not open menuconfig file %s", cfgPath))
	}

	if _, err = toml.DecodeReader(f, &cfg); err != nil {
		log.Fatal(fmt.Sprintf("Could not decode TOML in config file %s", cfgPath))
	}

	f.Close()	
	flag.Float64Var(&lat, "lat", cfg.Lat, "Latitude for requested location")
	flag.Float64Var(&lon, "lon", cfg.Lon, "Longitude for requested location")
	flag.BoolVar(&hourly, "hourly", false, "Request hourly data")
	flag.BoolVar(&daily, "daily", false, "Request daily data")
	flag.BoolVar(&extend, "extend", false, "Request extended hourly data")
	flag.StringVar(&key, "apikey", cfg.ApiKey, "API key for forecast.io")
	flag.StringVar(&exclude, "exclude", "", "comma separated list of fields to exclude")
	flag.StringVar(&units, "units", "auto", "Units to return values in")
}

func main() {
	flag.Parse()
	excludes = strings.Split(exclude, ",")
	c := forecastio.NewConnection(key)
	c.SetUnits(units)

	f, err := c.Forecast(lat, lon, excludes, extend)
	if err != nil {
		panic(err)
	}
	f.ParseTimes()

	// fmt.Printf("API Calls made today: %d\n", c.APICalls())
	// fmt.Printf("Latitude: %.2f  Logitude: %.2f  Timezone: %s\n", f.Latitude, f.Longitude, f.Timezone)
	fmt.Printf("Report Time: %s  Summary: %s\n", f.Currently.Time.Format("02/Jan/2006 15:04"), f.Currently.Summary)
	fmt.Printf("Temperature: %.0f°  Pressure: %.0fmb  Humidity %.0f%%\n", f.Currently.Temperature, f.Currently.Pressure, f.Currently.Humidity*100)
	
	if hourly && len(f.Hourly.Data) > 0 {
		fmt.Printf("\nSummary for next hour: %s\n", f.Hourly.Summary)
		for _, h := range f.Hourly.Data {
			fmt.Printf("Time: %s  Temperature: %2.0f°  Pressure: %4.0fmb  - %s\n", h.Time.Format("02/Jan/2006 - 15:04"), h.Temperature, h.Pressure, h.Summary)
		}
	}
	
	if daily && len(f.Daily.Data) > 0 {
		fmt.Printf("\nSummary for next 7 days: %s\n", f.Daily.Summary)
		for _, d := range f.Daily.Data {
			fmt.Printf("Time: %s  Temperature (Min/Max): %2.0f/%2.0f°  Pressure: %4.0fmb  - %s\n", d.Time.Format("02/Jan/2006"), d.TemperatureMin, d.TemperatureMax, d.Pressure, d.Summary)
		}
	}
}
