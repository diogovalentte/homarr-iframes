package alarms

import "time"

type Alarm struct {
	// Time: time related to the alarm
	Time time.Time
	// Summary: like "Low disk space on /dev/sda1"
	Summary string
	// URL: an URL to be used in the link of the alarm
	URL string
	// Status: like "CLEAR", "WARNING", "ERROR", or "CRITICAL"
	// Prefer to uppercase the status
	Status string
	// Value: a value related to the alarm, like "12GB free"
	Value string
	// Property: custom property to be used in the alarm card
	Property string
	// Source: like "Netdata", "Radarr", etc.
	Source string
	// BackgroundImgURL: URL to an image to be used as background of the alarm card
	BackgroundImgURL string
	// BackgroundColor: Color to be used as background of the alarm card if no BackgroundImgURL is set
	BackgroundColor string
	// BackgroundImgSize: Size of the background image in %, like 80 or 102.5
	BackgroundImgSize float32
}
