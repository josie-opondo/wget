package utils

import "io"

type ProgressRecoder struct {
	Reader           io.Reader
	Total            int64
	Progress         int64
	ProgressFunction func(int64, int64)
}

type WgetValues struct {
	BackgroudMode   bool   // Flag -B
	OutputFile      string // Flag -O
	OutPutDirectory string // Flag -P
	RateLimitValue  string //Flag --rate-limit
	Reject          bool
	Exclude         string // Flag exclude || -X
	ConvertLinks    bool   // Flag --convert-links
	Mirror          bool   //Flag --mirror
	Url             string // --- url given
}

func WgetInstance() *WgetValues {
	return &WgetValues{
		BackgroudMode:   false,
		OutputFile:      "",
		OutPutDirectory: "",
		RateLimitValue:  "",
		Reject:          false,
		Exclude:         "",
		ConvertLinks:    false,
		Mirror:          false,
	}
}
