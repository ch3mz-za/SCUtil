package logmon

import (
	"strings"
	"time"
)

type EventType string

const (
	ActorDeath         EventType = "Actor Death"
	VehicleDestruction EventType = "Vehicle Destruction"
)

// Domain parser
type LogItem struct {
	Attacker string
	Driver   string
	Location string
	Time     *time.Time
	Type     EventType
	Vehicle  string
	Victim   string
	Weapon   string
}

type IndicesAD struct{ Time, Victim, Attacker, Weapon int }
type IndicesVD struct{ Time, Vehicle, Location, Driver, Attacker, Weapon int }

func DefaultAD() IndicesAD { return IndicesAD{0, 5, 12, 15} }
func DefaultVD() IndicesVD { return IndicesVD{0, 6, 10, 27, 38, 41} }

func GameParser(ad IndicesAD, vd IndicesVD) Parser {
	return func(line string) (*LogItem, bool) {
		if line == "" {
			return nil, false
		}

		var fields []string
		get := func(i int) string {
			if i < 0 || i >= len(fields) {
				return ""
			}
			return fields[i]
		}

		// Actor Death
		if strings.Contains(line, string(ActorDeath)) {
			fields = strings.Fields(line)
			return &LogItem{
				Time:     roundTimeToSeconds(trimAngleBrackets(get(ad.Time))),
				Attacker: normalize(get(ad.Attacker)),
				Victim:   normalize(get(ad.Victim)),
				Weapon:   normalize(get(ad.Weapon)),
				Type:     ActorDeath,
			}, true
		}
		// Vehicle Death
		if strings.Contains(line, string(VehicleDestruction)) {
			fields = strings.Fields(line)
			return &LogItem{
				Time:     roundTimeToSeconds(trimAngleBrackets(get(ad.Time))),
				Vehicle:  normalize(get(vd.Vehicle)),
				Attacker: normalize(get(vd.Attacker)),
				Location: trimQuotes(get(vd.Location)),
				Driver:   get(vd.Driver),
				Weapon:   trimQuotes(get(vd.Weapon)),
				Type:     VehicleDestruction,
			}, true
		}
		return nil, false
	}
}
