package utils

import (
	"io"
	"time"
)

type ProgressRecoder struct {
	Reader           io.Reader
	Total            int64
	Progress         int64
	startTime        time.Time
	ProgressFunction func(int64, int64, time.Time)
}

type WgetValues struct {
	BackgroudMode   bool   // Flag -B
	OutputFile      string // Flag -O
	OutPutDirectory string // Flag -P
	RateLimitValue  int    //Flag --rate-limit
	Reject          bool
	Exclude         string   // Flag exclude || -X
	ConvertLinks    bool     // Flag --convert-links
	MirrorMode      bool     //Flag --mirror
	MirrorStarted   bool	 // Flag --mirror start time
	Url             string   // --- url given
	RejectSuffixes  []string // Flag rejects suffixes
}

// Rate limiter Struct
// RateLimitedReader limits the read speed to a specified rate in bytes per second
type RateLimitedReader struct {
	Reader io.Reader
	Rate   int64 // bytes per second
	Ticker *time.Ticker
}

// Read implements io.Reader.
func (r *RateLimitedReader) Read(p []byte) (n int, err error) {
	toRead := len(p)
	if int64(toRead) > r.Rate {
		toRead = int(r.Rate)
	}
	n, err = r.Reader.Read(p[:toRead])
	if n > 0 {
		<-r.Ticker.C // Wait for the next tick
	}
	return n, err
}

func WgetInstance() *WgetValues {
	return &WgetValues{
		BackgroudMode:   false,
		OutputFile:      "",
		OutPutDirectory: ".",
		RateLimitValue:  0,
		Reject:          false,
		Exclude:         "",
		ConvertLinks:    false,
		MirrorMode:      false,
		MirrorStarted:   false,
	}
}
