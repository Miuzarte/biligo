package biligo

import (
	"fmt"
	"strconv"
	"time"
)

type number interface {
	integer | ~float32 | ~float64
}

type integer interface {
	~int | ~int32 | ~int64 |
		~uint | ~uint32 | ~uint64
}

func itoa[T number](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

func formatNumber[T number](num T) string {
	const (
		numA = 1
		numB = numA * 10 // 十
		numC = numB * 10 // 百
		numD = numC * 10 // 千
		numE = numD * 10 // 万
		numF = numE * 10 // 十万
		numG = numF * 10 // 百万
		numH = numG * 10 // 千万
		numI = numH * 10 // 亿
	)

	switch {
	case num > numI: // 1.xx亿
		hundredMillion := int64(num) / numI
		tenMillion := (int64(num) - hundredMillion*numI) / numH
		million := (int64(num) - hundredMillion*numI - tenMillion*numH) / numG
		if million == 0 && tenMillion == 0 {
			return fmt.Sprintf("%d亿", hundredMillion)
		} else if tenMillion == 0 {
			return fmt.Sprintf("%d.%d亿", hundredMillion, million)
		} else {
			return fmt.Sprintf("%d.%d%d亿", hundredMillion, tenMillion, million)
		}

	case num > numH: // 1.x千万
		tenMillion := int64(num) / numH
		million := (int64(num) - tenMillion*numH) / numG
		if million == 0 {
			return fmt.Sprintf("%d千万", tenMillion)
		} else {
			return fmt.Sprintf("%d.%d千万", tenMillion, million)
		}

	case num > numE: // 1.x万
		tenTh := int64(num) / numE
		thousand := (int64(num) - tenTh*numE) / numD
		if thousand == 0 {
			return fmt.Sprintf("%d万", tenTh)
		} else {
			return fmt.Sprintf("%d.%d万", tenTh, thousand)
		}

	case num > numD: // 1kx
		thousand := int64(num) / numD
		hundred := (int64(num) - thousand*numD) / numC
		if hundred == 0 {
			return fmt.Sprintf("%dk", thousand)
		} else {
			return fmt.Sprintf("%dk%d", thousand, hundred)
		}

	default:
		return itoa(int64(num))
	}
}

func formatTimeInterval(start, end time.Time) (string, string) {
	if start.Format("2006") == end.Format("2006") { // 结束日期同年 不显示年份
		if end.Format("01") == time.Now().Format("01") { // 结束日期同月 不显示月份
			return start.Format(TIME_LAYOUT_M24), end.Format(TIME_LAYOUT_S24)
		}
		return start.Format(TIME_LAYOUT_M24), end.Format(TIME_LAYOUT_M24)
	}
	return start.Format(TIME_LAYOUT_L24), end.Format(TIME_LAYOUT_L24)
}

// formatDuration 格式化秒级时间戳至 时分秒 x:x:x
func formatDuration[T integer](duration T) (format string) {
	h := (duration / (60 * 60)) % 24
	m := (duration / 60) % 60
	s := duration % 60
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

// formatInterval 格式化秒级时间戳至 x天x小时x分钟x秒
func formatInterval[T integer](interval T) (format string) {
	if interval == 0 {
		return "0秒"
	}
	days := interval / (24 * 60 * 60)
	hours := (interval / (60 * 60)) % 24
	minutes := (interval / 60) % 60
	seconds := interval % 60
	switch {
	case days > 0:
		format += itoa(days) + "天"
		fallthrough
	case hours > 0:
		format += itoa(hours) + "小时"
		fallthrough
	case minutes > 0:
		format += itoa(minutes) + "分"
		fallthrough
	default:
		format += itoa(seconds) + "秒"
	}
	return format
}

// formatTime 根据当前时间格式化时间戳
// TODO: 实现"昨天"、"前天"、"明天"、"后天"...
func formatTime[T integer](timestamp T) string {
	tn := time.Now()
	t := time.Unix(int64(timestamp), 0)
	sameYear := tn.Year() == t.Year()
	sameMonth := tn.Month() == t.Month()
	sameDay := tn.Day() == t.Day()
	var f string
	if sameYear && sameMonth && sameDay {
		f = "15:04"
		// } else if sameYear && sameMonth {
		// 	f = "02 15:04"
	} else if sameYear {
		f = "01/02 15:04"
	} else {
		f = "2006/01/02 15:04"
	}
	return t.Format(f)
}

func formatPercent[T number](part, total T) string {
	return fmt.Sprintf("%.0f%%", float64(part)/float64(total)*100.0)
}

func liveStatusText[T integer](status T) string {
	switch status {
	case 0:
		return "未开播"
	case 1:
		return "直播中"
	case 2:
		return "轮播中"
	default:
		return "未知状态"
	}
}
