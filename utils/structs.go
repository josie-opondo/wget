package utils

import "io"

type ProgressRecoder struct {
	Reader           io.Reader
	Total            int64
	Progress         int64
	ProgressFunction func(int64, int64)
}
