package utils

import "io"

type ProgressRecoder struct {
	Reader           io.Reader
	Total            int64
	Progress         int64
	ProgressFunction func(int64, int64)
}

func (pr *ProgressRecoder) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	CheckError(err)
	pr.Progress += int64(n)
	if pr.ProgressFunction != nil {
		pr.ProgressFunction(pr.Progress, pr.Total)
	}
	return
}