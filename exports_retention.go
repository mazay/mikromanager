package main

import (
	"time"

	"github.com/mazay/mikromanager/internal"
)

func exportInSlice(a *internal.Export, list []*internal.Export) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// getLatestExport returns latest export object in the list
func getLatestExport(slice []*internal.Export) *internal.Export {
	var latest *internal.Export
	for _, export := range slice {
		if latest == nil || export.LastModified.After(*latest.LastModified) {
			latest = export
		}
	}
	return latest
}

// exportsToKeep finds a list of export objects within the given time slots that we'd want to keep
// typically latest backup for a timeWindow after each item in the timeSlice
func exportsToKeep(exports []*internal.Export, timeSlice []time.Time, timeWindow time.Duration) []*internal.Export {
	var exportList []*internal.Export

	for _, t := range timeSlice {
		var tmpList []*internal.Export
		t2 := t.Add(timeWindow)
		for _, export := range exports {
			if export.LastModified.After(t) && export.LastModified.Before(t2) {
				tmpList = append(tmpList, export)
			}
		}
		earliestExport := getLatestExport(tmpList)
		if earliestExport != nil {
			exportList = append(exportList, earliestExport)
		}
	}
	return exportList
}

// rotateHourlyExports return a list of hourly exports that should be kept
func rotateHourlyExports(exports []*internal.Export, number int64) []*internal.Export {
	var exportsList []*internal.Export

	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	start := end.Add(-time.Hour * time.Duration(number))
	slice := timeSliceBy(start, end, time.Hour)

	timeWindow, err := time.ParseDuration("59m59s")
	if err != nil {
		logger.Fatal(err.Error())
		return exportsList
	}
	exportsList = exportsToKeep(exports, slice, timeWindow)

	return exportsList
}

// rotateDailyExports return a list of daily exports that should be kept
func rotateDailyExports(exports []*internal.Export, number int64) []*internal.Export {
	var exportsList []*internal.Export

	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	start := end.Add(-time.Hour * 24 * time.Duration(number))
	slice := timeSliceBy(start, end, time.Hour*24)

	timeWindow, err := time.ParseDuration("23h59m59s")
	if err != nil {
		logger.Fatal(err.Error())
		return exportsList
	}
	exportsList = exportsToKeep(exports, slice, timeWindow)

	return exportsList
}

// rotateWeeklyExports return a list of weekly exports that should be kept
func rotateWeeklyExports(exports []*internal.Export, number int64) []*internal.Export {
	var exportsList []*internal.Export

	now := time.Now()
	weekDayDiff := 7 - now.Weekday()
	end := time.Date(now.Year(), now.Month(), now.Day()+int(weekDayDiff), 0, 0, 0, 0, now.Location())
	start := end.Add(-time.Hour * 168 * 30)
	slice := timeSliceBy(start, end, time.Hour*168)

	// 6 days 23 hours 59 minutes 59 seconds
	timeWindow, err := time.ParseDuration("167h59m59s")
	if err != nil {
		logger.Fatal(err.Error())
		return exportsList
	}
	exportsList = exportsToKeep(exports, slice, timeWindow)

	return exportsList
}
