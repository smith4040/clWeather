# clWeather — Aviation METAR/TAF CLI

[![CI](https://github.com/smith4040/clWeather/actions/workflows/ci.yml/badge.svg)](https://github.com/smith4040/clWeather/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/smith4040/clWeather)](https://goreportcard.com/report/github.com/smith4040/clWeather)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A fast CLI for aviation weather: METAR (observations) and TAF (forecasts) from [aviationweather.gov API](https://aviationweather.gov/data/api/).

## Features
- METAR: Current conditions, flight category (VFR/IFR/MVFR/LIFR), wind, vis, clouds
- TAF: 24-30hr forecasts with TEMPO changes
- Formats: Human-readable, raw text, JSON
- No API key; rate-limited to 100/min

## Installation

```bash
go install github.com/smith4040/clWeather/cmd/clweather@latest