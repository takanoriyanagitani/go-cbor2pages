// Code generated by "stringer -type=LogPageSize line.go"; DO NOT EDIT.

package log

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[LogPageXs-8]
	_ = x[LogPageSm-64]
	_ = x[LogPageMd-512]
	_ = x[LogPageLg-4096]
	_ = x[LogPageXl-32768]
}

const (
	_LogPageSize_name_0 = "LogPageXs"
	_LogPageSize_name_1 = "LogPageSm"
	_LogPageSize_name_2 = "LogPageMd"
	_LogPageSize_name_3 = "LogPageLg"
	_LogPageSize_name_4 = "LogPageXl"
)

func (i LogPageSize) String() string {
	switch {
	case i == 8:
		return _LogPageSize_name_0
	case i == 64:
		return _LogPageSize_name_1
	case i == 512:
		return _LogPageSize_name_2
	case i == 4096:
		return _LogPageSize_name_3
	case i == 32768:
		return _LogPageSize_name_4
	default:
		return "LogPageSize(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
