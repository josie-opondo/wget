package mirror

import (
	"io"
	"time"
)

type RateLimitedReader struct {
	Reader     io.Reader
	RateLimit  int64
	Bucket     int64
	LastFilled time.Time
}

func (r *RateLimitedReader) Read(p []byte) (n int, err error) {
	if r.Bucket <= 0 {
		time.Sleep(time.Second)
		r.Bucket = r.RateLimit
		r.LastFilled = time.Now()
	}

	toRead := int64(len(p))
	if toRead > r.Bucket {
		toRead = r.Bucket
	}

	n, err = r.Reader.Read(p[:toRead])
	r.Bucket -= int64(n)

	return n, err
}

// ProgressRecoder tracks the download progress
type ProgressRecoder struct {
	Reader           io.Reader
	Total            int64
	Progress         int64
	StartTime        time.Time
	ProgressFunction func(downloaded, total int64, start time.Time)
}

// Read updates the progress bar
func (p *ProgressRecoder) Read(b []byte) (int, error) {
	n, err := p.Reader.Read(b)
	if n > 0 {
		p.Progress += int64(n)
		p.ProgressFunction(p.Progress, p.Total, p.StartTime)
	}
	return n, err
}
