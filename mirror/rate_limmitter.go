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
