package time

import (
	"time"
)

type TimeClass struct {
	//radius float64 // Private field
}

// ============================================================================
// 1. CÁC HÀM ĐỊNH DẠNG THỜI GIAN (FORMATTING)
// ============================================================================

// FormatDate định dạng time.Time thành "yyyy-mm-dd"
func (tc *TimeClass) FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDateTime định dạng time.Time thành "yyyy-mm-dd H:i:s" (24h)
func (tc *TimeClass) FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// ============================================================================
// 2. CÁC HÀM LẤY THỜI GIAN TƯƠNG LAI (FUTURE TIME)
// ============================================================================

// AddHours lấy thời gian sau n giờ
func (tc *TimeClass) AddHours(baseTime time.Time, n int) time.Time {
	return baseTime.Add(time.Duration(n) * time.Hour)
}

// AddDays lấy thời gian sau n ngày
func (tc *TimeClass) AddDays(baseTime time.Time, n int) time.Time {
	return baseTime.AddDate(0, 0, n)
}

// AddMonths lấy thời gian sau n tháng
func (tc *TimeClass) AddMonths(baseTime time.Time, n int) time.Time {
	return baseTime.AddDate(0, n, 0)
}

// AddYears lấy thời gian sau n năm
func (tc *TimeClass) AddYears(baseTime time.Time, n int) time.Time {
	return baseTime.AddDate(n, 0, 0)
}
