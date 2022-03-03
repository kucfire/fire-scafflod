package utils

import "time"

var cst *time.Location

// CSTLayout China Standard Time Layout
const CSTLayout = "2006-01-02 15:04:05"

func init() {
	var err error
	if cst, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		panic(err)
	}

	// 默认设置为中国时区
	time.Local = cst
}

func CSTLayoutString() string {
	now := time.Now()

	return now.In(cst).Format(CSTLayout)
}
