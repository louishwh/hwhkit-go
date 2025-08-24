package utils

import (
	"fmt"
	"time"
)

// TimeUtils 时间工具集合
type TimeUtils struct{}

// NewTimeUtils 创建时间工具实例
func NewTimeUtils() *TimeUtils {
	return &TimeUtils{}
}

// 常用时间格式常量
const (
	DateTimeFormat     = "2006-01-02 15:04:05"
	DateFormat         = "2006-01-02"
	TimeFormat         = "15:04:05"
	ISO8601Format      = "2006-01-02T15:04:05Z07:00"
	RFC3339Format      = time.RFC3339
	CompactDateFormat  = "20060102"
	CompactTimeFormat  = "150405"
	ChineseDateFormat  = "2006年01月02日"
	ChineseTimeFormat  = "2006年01月02日 15时04分05秒"
)

// Now 获取当前时间
func (t *TimeUtils) Now() time.Time {
	return time.Now()
}

// NowUnix 获取当前Unix时间戳（秒）
func (t *TimeUtils) NowUnix() int64 {
	return time.Now().Unix()
}

// NowUnixMilli 获取当前Unix时间戳（毫秒）
func (t *TimeUtils) NowUnixMilli() int64 {
	return time.Now().UnixMilli()
}

// NowUnixMicro 获取当前Unix时间戳（微秒）
func (t *TimeUtils) NowUnixMicro() int64 {
	return time.Now().UnixMicro()
}

// NowUnixNano 获取当前Unix时间戳（纳秒）
func (t *TimeUtils) NowUnixNano() int64 {
	return time.Now().UnixNano()
}

// FormatNow 格式化当前时间
func (t *TimeUtils) FormatNow(layout string) string {
	return time.Now().Format(layout)
}

// FormatNowDateTime 格式化当前时间为日期时间字符串
func (t *TimeUtils) FormatNowDateTime() string {
	return time.Now().Format(DateTimeFormat)
}

// FormatNowDate 格式化当前时间为日期字符串
func (t *TimeUtils) FormatNowDate() string {
	return time.Now().Format(DateFormat)
}

// FormatNowTime 格式化当前时间为时间字符串
func (t *TimeUtils) FormatNowTime() string {
	return time.Now().Format(TimeFormat)
}

// Format 格式化时间
func (t *TimeUtils) Format(time time.Time, layout string) string {
	return time.Format(layout)
}

// Parse 解析时间字符串
func (t *TimeUtils) Parse(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

// ParseDateTime 解析日期时间字符串
func (t *TimeUtils) ParseDateTime(value string) (time.Time, error) {
	return time.Parse(DateTimeFormat, value)
}

// ParseDate 解析日期字符串
func (t *TimeUtils) ParseDate(value string) (time.Time, error) {
	return time.Parse(DateFormat, value)
}

// ParseTime 解析时间字符串
func (t *TimeUtils) ParseTime(value string) (time.Time, error) {
	return time.Parse(TimeFormat, value)
}

// ParseUnix 从Unix时间戳创建时间对象
func (t *TimeUtils) ParseUnix(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// ParseUnixMilli 从Unix时间戳（毫秒）创建时间对象
func (t *TimeUtils) ParseUnixMilli(timestamp int64) time.Time {
	return time.UnixMilli(timestamp)
}

// ParseUnixMicro 从Unix时间戳（微秒）创建时间对象
func (t *TimeUtils) ParseUnixMicro(timestamp int64) time.Time {
	return time.UnixMicro(timestamp)
}

// ParseUnixNano 从Unix时间戳（纳秒）创建时间对象
func (t *TimeUtils) ParseUnixNano(timestamp int64) time.Time {
	return time.Unix(0, timestamp)
}

// AddDays 添加天数
func (t *TimeUtils) AddDays(time time.Time, days int) time.Time {
	return time.AddDate(0, 0, days)
}

// AddHours 添加小时数
func (t *TimeUtils) AddHours(time time.Time, hours int) time.Time {
	return time.Add(time.Duration(hours) * time.Hour)
}

// AddMinutes 添加分钟数
func (t *TimeUtils) AddMinutes(time time.Time, minutes int) time.Time {
	return time.Add(time.Duration(minutes) * time.Minute)
}

// AddSeconds 添加秒数
func (t *TimeUtils) AddSeconds(time time.Time, seconds int) time.Time {
	return time.Add(time.Duration(seconds) * time.Second)
}

// SubtractDays 减去天数
func (t *TimeUtils) SubtractDays(time time.Time, days int) time.Time {
	return time.AddDate(0, 0, -days)
}

// SubtractHours 减去小时数
func (t *TimeUtils) SubtractHours(time time.Time, hours int) time.Time {
	return time.Add(time.Duration(-hours) * time.Hour)
}

// SubtractMinutes 减去分钟数
func (t *TimeUtils) SubtractMinutes(time time.Time, minutes int) time.Time {
	return time.Add(time.Duration(-minutes) * time.Minute)
}

// SubtractSeconds 减去秒数
func (t *TimeUtils) SubtractSeconds(time time.Time, seconds int) time.Time {
	return time.Add(time.Duration(-seconds) * time.Second)
}

// BeginOfDay 获取一天的开始时间
func (t *TimeUtils) BeginOfDay(time time.Time) time.Time {
	year, month, day := time.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Location())
}

// EndOfDay 获取一天的结束时间
func (t *TimeUtils) EndOfDay(time time.Time) time.Time {
	year, month, day := time.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, time.Location())
}

// BeginOfMonth 获取月份的开始时间
func (t *TimeUtils) BeginOfMonth(time time.Time) time.Time {
	year, month, _ := time.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, time.Location())
}

// EndOfMonth 获取月份的结束时间
func (t *TimeUtils) EndOfMonth(time time.Time) time.Time {
	year, month, _ := time.Date()
	// 下个月的第一天减去1纳秒
	nextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.Location())
	return nextMonth.Add(-time.Nanosecond)
}

// BeginOfYear 获取年份的开始时间
func (t *TimeUtils) BeginOfYear(time time.Time) time.Time {
	year, _, _ := time.Date()
	return time.Date(year, 1, 1, 0, 0, 0, 0, time.Location())
}

// EndOfYear 获取年份的结束时间
func (t *TimeUtils) EndOfYear(time time.Time) time.Time {
	year, _, _ := time.Date()
	return time.Date(year, 12, 31, 23, 59, 59, 999999999, time.Location())
}

// DiffDays 计算两个时间相差的天数
func (t *TimeUtils) DiffDays(from, to time.Time) int {
	return int(to.Sub(from).Hours() / 24)
}

// DiffHours 计算两个时间相差的小时数
func (t *TimeUtils) DiffHours(from, to time.Time) int {
	return int(to.Sub(from).Hours())
}

// DiffMinutes 计算两个时间相差的分钟数
func (t *TimeUtils) DiffMinutes(from, to time.Time) int {
	return int(to.Sub(from).Minutes())
}

// DiffSeconds 计算两个时间相差的秒数
func (t *TimeUtils) DiffSeconds(from, to time.Time) int {
	return int(to.Sub(from).Seconds())
}

// IsLeapYear 判断是否为闰年
func (t *TimeUtils) IsLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

// DaysInMonth 获取指定年月的天数
func (t *TimeUtils) DaysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// IsToday 判断是否为今天
func (t *TimeUtils) IsToday(time time.Time) bool {
	now := time.Now()
	return t.IsSameDay(time, now)
}

// IsYesterday 判断是否为昨天
func (t *TimeUtils) IsYesterday(time time.Time) bool {
	yesterday := time.Now().AddDate(0, 0, -1)
	return t.IsSameDay(time, yesterday)
}

// IsTomorrow 判断是否为明天
func (t *TimeUtils) IsTomorrow(time time.Time) bool {
	tomorrow := time.Now().AddDate(0, 0, 1)
	return t.IsSameDay(time, tomorrow)
}

// IsSameDay 判断两个时间是否为同一天
func (t *TimeUtils) IsSameDay(time1, time2 time.Time) bool {
	y1, m1, d1 := time1.Date()
	y2, m2, d2 := time2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// IsSameMonth 判断两个时间是否为同一月
func (t *TimeUtils) IsSameMonth(time1, time2 time.Time) bool {
	y1, m1, _ := time1.Date()
	y2, m2, _ := time2.Date()
	return y1 == y2 && m1 == m2
}

// IsSameYear 判断两个时间是否为同一年
func (t *TimeUtils) IsSameYear(time1, time2 time.Time) bool {
	return time1.Year() == time2.Year()
}

// IsBefore 判断时间1是否在时间2之前
func (t *TimeUtils) IsBefore(time1, time2 time.Time) bool {
	return time1.Before(time2)
}

// IsAfter 判断时间1是否在时间2之后
func (t *TimeUtils) IsAfter(time1, time2 time.Time) bool {
	return time1.After(time2)
}

// IsBetween 判断时间是否在指定范围内
func (t *TimeUtils) IsBetween(time, start, end time.Time) bool {
	return time.After(start) && time.Before(end)
}

// IsWorkday 判断是否为工作日（周一到周五）
func (t *TimeUtils) IsWorkday(time time.Time) bool {
	weekday := time.Weekday()
	return weekday >= time.Monday && weekday <= time.Friday
}

// IsWeekend 判断是否为周末
func (t *TimeUtils) IsWeekend(time time.Time) bool {
	weekday := time.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// Age 计算年龄
func (t *TimeUtils) Age(birthday time.Time) int {
	now := time.Now()
	age := now.Year() - birthday.Year()
	
	// 如果今年的生日还没到，年龄减1
	if now.Month() < birthday.Month() ||
		(now.Month() == birthday.Month() && now.Day() < birthday.Day()) {
		age--
	}
	
	return age
}

// TimeAgo 返回时间距离现在的描述（如：2小时前）
func (t *TimeUtils) TimeAgo(time time.Time) string {
	now := time.Now()
	diff := now.Sub(time)
	
	if diff < time.Minute {
		return "刚刚"
	}
	
	if diff < time.Hour {
		minutes := int(diff.Minutes())
		return fmt.Sprintf("%d分钟前", minutes)
	}
	
	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%d小时前", hours)
	}
	
	if diff < 30*24*time.Hour {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d天前", days)
	}
	
	if diff < 365*24*time.Hour {
		months := int(diff.Hours() / (24 * 30))
		return fmt.Sprintf("%d个月前", months)
	}
	
	years := int(diff.Hours() / (24 * 365))
	return fmt.Sprintf("%d年前", years)
}

// Sleep 休眠指定时间
func (t *TimeUtils) Sleep(duration time.Duration) {
	time.Sleep(duration)
}

// SleepSeconds 休眠指定秒数
func (t *TimeUtils) SleepSeconds(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

// SleepMilliseconds 休眠指定毫秒数
func (t *TimeUtils) SleepMilliseconds(milliseconds int) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}

// ToLocation 转换时区
func (t *TimeUtils) ToLocation(time time.Time, location *time.Location) time.Time {
	return time.In(location)
}

// ToUTC 转换为UTC时间
func (t *TimeUtils) ToUTC(time time.Time) time.Time {
	return time.UTC()
}

// ToLocal 转换为本地时间
func (t *TimeUtils) ToLocal(time time.Time) time.Time {
	return time.Local()
}

// GetLocation 获取时区
func (t *TimeUtils) GetLocation(name string) (*time.Location, error) {
	return time.LoadLocation(name)
}

// 全局时间工具实例
var Time = NewTimeUtils()