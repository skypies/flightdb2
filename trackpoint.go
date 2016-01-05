package flightdb2

import (
	"fmt"
	"time"

	"github.com/skypies/adsb"
	"github.com/skypies/geo"
)

// Trackpoint is a data point that locates an aircraft in space and time, etc
type Trackpoint struct {
	DataSource   string    // What kind of trackpoint is this; flightaware radar, local ADSB, etc
	ReceiverName string    // For local ADSB
	TimestampUTC time.Time // Always in UTC, to make life SIMPLE

	geo.Latlong            // Embedded type, so we can call all the geo stuff directly on trackpoints

	Altitude     float64   // In feet
	GroundSpeed  float64   // In knots
	Heading      float64   // [0.0, 360.0) degrees. Direction plane is pointing in. Mag or real north?
	VerticalRate float64   // In feet per minute (multiples of 64)
	Squawk       string    // Generally, a string of four digits.
}

func (tp Trackpoint)String() string {
	return fmt.Sprintf("[%s] %s %.0fft, %.0fkts", tp.TimestampUTC, tp.Latlong,
		tp.Altitude, tp.GroundSpeed)
}

func (tp Trackpoint)ToJSString() string {
	return fmt.Sprintf("source:%q, receiver:%q, pos:{lat:%.6f,lng:%.6f}, "+
		"alt:%.0f, speed:%.0f, track:%.0f, vert:%.0f, t:\"%s\"",
		tp.DataSource, tp.ReceiverName, tp.Lat, tp.Long,
		tp.Altitude, tp.GroundSpeed, tp.Heading, tp.VerticalRate, tp.TimestampUTC)
}

func (tp Trackpoint)LongSource() string {
	switch tp.DataSource {
	case "":      return "(none specified)"
	case "FA:TZ": return "FlightAware, Radar (TZ)"
	case "FA:TA": return "FlightAware, ADS-B Mode-ES (TA)"
	case "ADSB":  return "Private receiver, ADS-B Mode-ES ("+tp.ReceiverName+")"
	}
	return tp.DataSource
}

func TrackpointFromADSB(m *adsb.CompositeMsg) Trackpoint {
	return Trackpoint{
		DataSource: "ADSB",
		ReceiverName: m.ReceiverName,
		TimestampUTC: m.GeneratedTimestampUTC,
		Latlong: m.Position,
		Altitude: float64(m.Altitude),
		GroundSpeed: float64(m.GroundSpeed),
		Heading: float64(m.Track),
		VerticalRate: float64(m.VerticalRate),
		Squawk: m.Squawk,
	}
}
