package tool

import "time"

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func TimestampNow() time.Time {
	return time.Now().UTC()
}

//return the beginning and the end of t
func TimeDaySection(t time.Time) (time.Time, time.Time) {
	year, month, day := t.Date()
	start := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	end := time.Date(year, month, day, 23, 59, 59, 0, t.Location())
	return start, end
}

func TimestampNow2String() string {
	time := TimestampNow()
	newFormat := time.Format("2006-01-02 15:04:05")
	return newFormat
}

func TimeString2Time(timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", timeStr)
}
