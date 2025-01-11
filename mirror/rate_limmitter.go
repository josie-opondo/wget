package mirror

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type RateLimitedReader struct {
	Reader     io.Reader
	RateLimit  int64
	Bucket     int64
	LastFilled time.Time
}

func parseRateLimit(rateLimit string) (int64, error) {
	if len(rateLimit) < 2 {
		return 0, fmt.Errorf("invalid rate limit")
	}

	multiplier := 1
	switch rateLimit[len(rateLimit)-1] {
	case 'k', 'K':
		multiplier = 1024
		rateLimit = rateLimit[:len(rateLimit)-1]
	case 'M':
		multiplier = 1024 * 1024
		rateLimit = rateLimit[:len(rateLimit)-1]
	}

	rate, err := strconv.Atoi(rateLimit)
	if err != nil {
		return 0, err
	}
	return int64(rate * multiplier), nil
}

func NewRateLimitedReader(reader io.Reader, limit string) *RateLimitedReader {
	// Convert limit to bytes per second (rateLimit)
	rateLimit, _ := parseRateLimit(limit)
	return &RateLimitedReader{Reader: reader, RateLimit: rateLimit, LastFilled: time.Now()}
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
