package utils

import (
	"fmt"
	"time"
)

func ParseDateTime(dt string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.000Z", dt)
}

func DateThaiFormatMonth(month time.Month) string {
	thaiMonths := []string{
		"มกราคม",
		"กุมภาพันธ์",
		"มีนาคม",
		"เมษายน",
		"พฤษภาคม",
		"มิถุนายน",
		"กรกฎาคม",
		"สิงหาคม",
		"กันยายน",
		"ตุลาคม",
		"พฤศจิกายน",
		"ธันวาคม",
	}
	return thaiMonths[month-1]
}

func DateThaiFormatShortMonth(month time.Month) string {
	thaiShortMonths := []string{
		"ม.ค.",  // January
		"ก.พ.",  // February
		"มี.ค.", // March
		"เม.ย.", // April
		"พ.ค.",  // May
		"มิ.ย.", // June
		"ก.ค.",  // July
		"ส.ค.",  // August
		"ก.ย.",  // September
		"ต.ค.",  // October
		"พ.ย.",  // November
		"ธ.ค.",  // December
	}
	return thaiShortMonths[month-1]
}

func FormatThaiDatetime(t time.Time) string {
	day := t.Day()
	month := DateThaiFormatShortMonth(t.Month())
	year := t.Year() + 543 // Buddhist Era
	hour, min := t.Hour(), t.Minute()
	return fmt.Sprintf("%02d %s %d, %02d:%02d", day, month, year, hour, min)
}
