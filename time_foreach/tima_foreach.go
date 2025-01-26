package main

import "time"

func main() {
	forWeek("2023-01-23 00:00:00", "2023-05-01 00:00:00")
}

func forDay(start, end string) {
	startTime, _ := time.Parse("2006-01-02 15:04:05", start)
	endTime, _ := time.Parse("2006-01-02 15:04:05", end)
	startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.Local)
	endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
	for timeRotate := startTime.AddDate(0, 0, 1); timeRotate.Before(endTime) || timeRotate.Equal(endTime); timeRotate = timeRotate.AddDate(0, 0, 1) {
		println(timeRotate.AddDate(0, 0, -1).Format("2006-01-02"))
	}
}

func forWeek(start, end string) {
	startTime, _ := time.Parse("2006-01-02 15:04:05", start)
	endTime, _ := time.Parse("2006-01-02 15:04:05", end)
	startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.Local)
	endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, time.Local)
	tmp := startTime
	for {
		tmp = tmp.AddDate(0, 0, 7)
		if endTime.Before(tmp) {
			endTime = tmp
			break
		}
	}
	for timeRotate := startTime.AddDate(0, 0, 7); timeRotate.Before(endTime) || timeRotate.Equal(endTime); timeRotate = timeRotate.AddDate(0, 0, 7) {
		println(timeRotate.AddDate(0, 0, -7).Format("2006-01-02") + "~" + timeRotate.AddDate(0, 0, -1).Format("2006-01-02"))
	}
}

func forMonth(start, end string) {
	startTime, _ := time.Parse("2006-01-02 15:04:05", start)
	endTime, _ := time.Parse("2006-01-02 15:04:05", end)
	startTime = time.Date(startTime.Year(), startTime.Month(), 1, 0, 0, 0, 0, time.Local)
	endTime = time.Date(endTime.Year(), endTime.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0)
	for timeRotate := startTime.AddDate(0, 1, 0); timeRotate.Before(endTime) || timeRotate.Equal(endTime); timeRotate = timeRotate.AddDate(0, 1, 0) {
		println(timeRotate.AddDate(0, -1, 0).Format("2006-01-02") + "~" + timeRotate.AddDate(0, 0, -1).Format("2006-01-02"))
	}
}
