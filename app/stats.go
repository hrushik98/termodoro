package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// SessionRecord is a single completed focus session.
type SessionRecord struct {
	Name        string    `json:"name"`
	Minutes     int       `json:"minutes"`
	CompletedAt time.Time `json:"completed_at"`
}

// Stats holds the history of completed sessions.
type Stats struct {
	Records []SessionRecord `json:"records"`
}

// DayStat aggregates focus minutes for a single calendar day.
type DayStat struct {
	Label   string
	Minutes int
	Today   bool
}

// statsFilePath returns the on-disk location of the stats file (~/.aimssh/stats.json).
func statsFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".aimssh", "stats.json"), nil
}

// LoadStats reads stats from disk. A missing or unreadable file yields empty stats.
func LoadStats() Stats {
	path, err := statsFilePath()
	if err != nil {
		return Stats{}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Stats{}
	}

	var s Stats
	if err := json.Unmarshal(data, &s); err != nil {
		return Stats{}
	}
	return s
}

// Save persists stats to disk, creating the parent directory if needed.
func (s *Stats) Save() error {
	path, err := statsFilePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

// RecordSession appends a completed session and persists it.
func RecordSession(name string, minutes int) {
	s := LoadStats()
	s.Records = append(s.Records, SessionRecord{
		Name:        name,
		Minutes:     minutes,
		CompletedAt: time.Now(),
	})
	_ = s.Save()
}

// TotalSessions returns the number of completed sessions.
func (s Stats) TotalSessions() int {
	return len(s.Records)
}

// TotalMinutes returns the cumulative focus minutes across all sessions.
func (s Stats) TotalMinutes() int {
	total := 0
	for _, r := range s.Records {
		total += r.Minutes
	}
	return total
}

// dayKey normalizes a time to a YYYY-MM-DD key in local time.
func dayKey(t time.Time) string {
	return t.Local().Format("2006-01-02")
}

// minutesByDay buckets total focus minutes per calendar day.
func (s Stats) minutesByDay() map[string]int {
	byDay := map[string]int{}
	for _, r := range s.Records {
		byDay[dayKey(r.CompletedAt)] += r.Minutes
	}
	return byDay
}

// TodayMinutes returns the focus minutes logged today.
func (s Stats) TodayMinutes() int {
	return s.minutesByDay()[dayKey(time.Now())]
}

// TodaySessions returns the number of sessions completed today.
func (s Stats) TodaySessions() int {
	today := dayKey(time.Now())
	count := 0
	for _, r := range s.Records {
		if dayKey(r.CompletedAt) == today {
			count++
		}
	}
	return count
}

// Streak returns the number of consecutive days (ending today or yesterday)
// that have at least one completed session.
func (s Stats) Streak() int {
	byDay := s.minutesByDay()
	if len(byDay) == 0 {
		return 0
	}

	// Anchor at today; if nothing logged today, the streak can still be
	// "alive" from yesterday.
	cursor := time.Now()
	if _, ok := byDay[dayKey(cursor)]; !ok {
		cursor = cursor.AddDate(0, 0, -1)
	}

	streak := 0
	for {
		if _, ok := byDay[dayKey(cursor)]; !ok {
			break
		}
		streak++
		cursor = cursor.AddDate(0, 0, -1)
	}
	return streak
}

// Last7Days returns focus minutes per day for the trailing week, oldest first.
func (s Stats) Last7Days() []DayStat {
	byDay := s.minutesByDay()
	today := dayKey(time.Now())

	days := make([]DayStat, 0, 7)
	for i := 6; i >= 0; i-- {
		d := time.Now().AddDate(0, 0, -i)
		key := dayKey(d)
		days = append(days, DayStat{
			Label:   d.Format("Mon"),
			Minutes: byDay[key],
			Today:   key == today,
		})
	}
	return days
}
