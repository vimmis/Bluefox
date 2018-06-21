/*
	Weatherapp with basic structures (v1)
*/

package main

import "time"

//Aids in mapping to the weather data in mysql database and request queries

type WeatherDataAdd struct {
	Timestamp	string
	Tempt	string
}

type WeatherDataGet struct {
	Timestamp	time.Time
	Tempt	float64
}