package utils

func (pr *ProgressRecoder) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	CheckError(err)
	pr.Progress += int64(n)
	if pr.ProgressFunction != nil {
		pr.ProgressFunction(pr.Progress, pr.Total, pr.startTime)
	}
	return
}
