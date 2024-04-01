package cronroutine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testCaseForNext struct {
	name       string
	Minute     []int
	Hour       []int
	DayOfMonth []int
	Month      []int
	DayOfWeek  []int

	currentTime time.Time
	expected    time.Time
}

var (
	nextAfterNoWeekdayCases = []testCaseForNext{
		{
			name:        "every minute",
			Minute:      allMinutes,
			Hour:        allHours,
			DayOfMonth:  allDaysInMonth,
			Month:       allMonths,
			DayOfWeek:   allDaysOfWeek,
			currentTime: time.Date(2021, time.January, 3, 14, 35, second, nanosecond, time.UTC),
			expected:    time.Date(2021, time.January, 3, 14, 36, second, nanosecond, time.UTC),
		},
		{
			name:        "every hour at 35 minutes",
			Minute:      []int{35},
			Hour:        allHours,
			DayOfMonth:  allDaysInMonth,
			Month:       allMonths,
			DayOfWeek:   allDaysOfWeek,
			currentTime: time.Date(2021, time.January, 3, 14, 15, second, nanosecond, time.UTC),
			expected:    time.Date(2021, time.January, 3, 14, 35, second, nanosecond, time.UTC),
		},
		{
			name:        "every hour at 35 minutes but already past the hour",
			Minute:      []int{35},
			Hour:        allHours,
			DayOfMonth:  allDaysInMonth,
			Month:       allMonths,
			DayOfWeek:   allDaysOfWeek,
			currentTime: time.Date(2021, time.January, 3, 14, 36, second, nanosecond, time.UTC),
			expected:    time.Date(2021, time.January, 3, 15, 35, second, nanosecond, time.UTC),
		},
		{
			name:        "At 03:00 on the first of every month (0 3 1 * *)",
			Minute:      []int{0},
			Hour:        []int{3},
			DayOfMonth:  []int{1},
			Month:       allMonths,
			DayOfWeek:   nil,
			currentTime: time.Date(2021, time.January, 3, 14, 36, second, nanosecond, time.UTC),
			expected:    time.Date(2021, time.February, 1, 3, 0, second, nanosecond, time.UTC),
		},
		{
			name:        "March 3rd every year at 3:32",
			Minute:      []int{32},
			Hour:        []int{3},
			DayOfMonth:  []int{3},
			Month:       []int{3},
			DayOfWeek:   []int{6},
			currentTime: time.Date(2021, time.January, 3, 14, 36, second, nanosecond, time.UTC),
			expected:    time.Date(2021, time.March, 3, 3, 32, second, nanosecond, time.UTC),
		},
		{
			name:        "March 3rd every year at 3:32, but check next year",
			Minute:      []int{32},
			Hour:        []int{3},
			DayOfMonth:  []int{3},
			Month:       []int{3},
			DayOfWeek:   nil,
			currentTime: time.Date(2021, time.March, 4, 14, 36, second, nanosecond, time.UTC),
			expected:    time.Date(2022, time.March, 3, 3, 32, second, nanosecond, time.UTC),
		},
		{
			name:        "every five hours at the top of the hour (0 0/5 * * *)",
			Minute:      []int{0},
			Hour:        []int{0, 5, 10, 15, 20},
			DayOfMonth:  allDaysInMonth,
			Month:       allMonths,
			DayOfWeek:   allDaysOfWeek,
			currentTime: time.Date(2021, time.March, 4, 14, 36, second, nanosecond, time.UTC),
			expected:    time.Date(2021, time.March, 4, 15, 00, second, nanosecond, time.UTC),
		},
		{
			name:        "every Sunday at 2:30",
			Minute:      []int{30},
			Hour:        []int{2},
			DayOfMonth:  nil,
			Month:       allMonths,
			DayOfWeek:   []int{0},
			currentTime: time.Date(2024, time.March, 27, 4, 35, second, nanosecond, time.UTC),
			expected:    time.Date(2024, time.March, 31, 2, 30, second, nanosecond, time.UTC),
		},
		{
			name:        "next month, monday morning (0 9-17/2 * * 1-5)",
			Minute:      []int{0},
			Hour:        []int{9, 11, 13, 15, 17},
			DayOfMonth:  nil,
			Month:       allMonths,
			DayOfWeek:   []int{1, 2, 3, 4, 5},
			currentTime: time.Date(2024, time.March, 29, 17, 36, second, nanosecond, time.UTC),
			expected:    time.Date(2024, time.April, 1, 9, 0, second, nanosecond, time.UTC),
		},
		{
			name:        "same day (0 9-17/2 * * 1-5)",
			Minute:      []int{0},
			Hour:        []int{9, 11, 13, 15, 17},
			DayOfMonth:  nil,
			Month:       allMonths,
			DayOfWeek:   []int{1, 2, 3, 4, 5},
			currentTime: time.Date(2024, time.March, 29, 10, 36, second, nanosecond, time.UTC),
			expected:    time.Date(2024, time.March, 29, 11, 0, second, nanosecond, time.UTC),
		},
	}

	allMinutes = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
		31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59}
	allDaysOfWeek  = []int{0, 1, 2, 3, 4, 5, 6}
	allDaysInMonth = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	allMonths      = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	allHours       = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
)

func TestCron_nextAfterWeekday(t *testing.T) {
	for _, tt := range nextAfterNoWeekdayCases {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cron{
				Minute:     tt.Minute,
				Hour:       tt.Hour,
				DayOfMonth: tt.DayOfMonth,
				Month:      tt.Month,
				DayOfWeek:  tt.DayOfWeek,
			}
			actual := c.next(tt.currentTime)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCron_NextFor(t *testing.T) {
	tests := []struct {
		name        string
		cronConfig  string
		start       time.Time
		forDuration time.Duration
		expected    []time.Time
	}{
		{
			name:        "every minute",
			cronConfig:  "* * * * *",
			start:       time.Date(2021, time.January, 3, 6, 35, second, nanosecond, time.UTC),
			forDuration: 5 * time.Minute,
			expected: []time.Time{
				time.Date(2021, time.January, 3, 6, 36, second, nanosecond, time.UTC),
				time.Date(2021, time.January, 3, 6, 37, second, nanosecond, time.UTC),
				time.Date(2021, time.January, 3, 6, 38, second, nanosecond, time.UTC),
				time.Date(2021, time.January, 3, 6, 39, second, nanosecond, time.UTC),
			},
		},
		{
			name:        "every 2 hours at the 30 on weekdays during working hours (30 9-17/2 * * 1-5)",
			cronConfig:  "30 9-17/2 * * 1-5",
			start:       time.Date(2024, time.March, 29, 10, 36, second, nanosecond, time.UTC),
			forDuration: 24 * 3 * time.Hour,
			expected: []time.Time{
				time.Date(2024, time.March, 29, 11, 30, second, nanosecond, time.UTC),
				time.Date(2024, time.March, 29, 13, 30, second, nanosecond, time.UTC),
				time.Date(2024, time.March, 29, 15, 30, second, nanosecond, time.UTC),
				time.Date(2024, time.March, 29, 17, 30, second, nanosecond, time.UTC),
				time.Date(2024, time.April, 1, 9, 30, second, nanosecond, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cron, err := ParseCron(tt.cronConfig)
			assert.NoError(t, err)
			actual := cron.nextFor(tt.start, tt.forDuration)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParseCron(t *testing.T) {
	tests := []struct {
		name       string
		cronConfig string
		expected   *Cron
		errMsg     string
	}{
		{
			name:       "every day",
			cronConfig: "0,30 0 * * *",
			expected: &Cron{
				Minute:     []int{0, 30},
				Hour:       []int{0},
				DayOfMonth: allDaysInMonth,
				Month:      allMonths,
				DayOfWeek:  allDaysOfWeek,
			},
		},
		{
			name:       "every month on the 2nd at midnight",
			cronConfig: "0 0 2 * *",
			expected: &Cron{
				Minute:     []int{0},
				Hour:       []int{0},
				DayOfMonth: []int{2},
				Month:      allMonths,
				DayOfWeek:  nil,
			},
		},
		{
			name:       "every month on the 2nd at midnight and Sunday",
			cronConfig: "0 0 2 * 0",
			expected: &Cron{
				Minute:     []int{0},
				Hour:       []int{0},
				DayOfMonth: []int{2},
				Month:      allMonths,
				DayOfWeek:  []int{0},
			},
		},
		{
			name:       "every sunday at 3:32",
			cronConfig: "32 3 * * 0",
			expected: &Cron{
				Minute:     []int{32},
				Hour:       []int{3},
				DayOfMonth: nil,
				Month:      allMonths,
				DayOfWeek:  []int{0},
			},
		},
		{
			name:       "every weekday at 3:32",
			cronConfig: "32 3 * * 1-5",
			expected: &Cron{
				Minute:     []int{32},
				Hour:       []int{3},
				DayOfMonth: nil,
				Month:      allMonths,
				DayOfWeek:  []int{1, 2, 3, 4, 5},
			},
		},
		{
			name:       "every fifteen minutes on weekdays during working hours",
			cronConfig: "0/15 9-17 * * 1-5",
			expected: &Cron{
				Minute:     []int{0, 15, 30, 45},
				Hour:       []int{9, 10, 11, 12, 13, 14, 15, 16, 17},
				DayOfMonth: nil,
				Month:      allMonths,
				DayOfWeek:  []int{1, 2, 3, 4, 5},
			},
		},
		{
			name:       "too few fields",
			cronConfig: "0 0 2 *",
			expected:   nil,
			errMsg:     "given 4 but need 5 fields for valid cron config: 0 0 2 *",
		},
		{
			name:       "too many fields",
			cronConfig: "0 0 2 1 1 0",
			expected:   nil,
			errMsg:     "given 6 but need 5 fields for valid cron config: 0 0 2 1 1 0",
		},
		{
			name:       "too many fields",
			cronConfig: "0 0 2 hello what",
			expected:   nil,
			errMsg:     "failed to parse month: failed to parse number: failed to parse int: strconv.Atoi: parsing \"hello\": invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ParseCron(tt.cronConfig)
			if tt.errMsg != "" {
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParse_parseNumber(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		minValue int
		maxValue int
		expected []int
		errMsg   string
	}{
		// Happy Path
		{
			name:     "parse *",
			value:    "*",
			minValue: 0,
			maxValue: 59,
			expected: allMinutes,
		},
		{
			name:     "parse minute",
			value:    "7",
			minValue: 0,
			maxValue: 59,
			expected: []int{7},
		},
		{
			name:     "parse hour",
			value:    "7",
			minValue: 0,
			maxValue: 23,
			expected: []int{7},
		},
		{
			name:     "parse dayOfMonth",
			value:    "7",
			minValue: 1,
			maxValue: 31,
			expected: []int{7},
		},
		{
			name:     "parse month",
			value:    "7",
			minValue: 1,
			maxValue: 12,
			expected: []int{7},
		},
		{
			name:     "parse dayOfWeek",
			value:    "4",
			minValue: 0,
			maxValue: 6,
			expected: []int{4},
		},
		// Error Path
		{
			name:     "parse minute",
			value:    "60",
			minValue: 0,
			maxValue: 59,
			expected: nil,
			errMsg:   "value 60 is not in range 0-59",
		},
		{
			name:     "missing start value",
			value:    "-3",
			minValue: 0,
			maxValue: 23,
			expected: nil,
			errMsg:   "failed to parse range \"-3\": failed to parse start int: strconv.Atoi: parsing \"\": invalid syntax",
		},
		{
			name:     "parse dayOfMonth",
			value:    "0",
			minValue: 1,
			maxValue: 31,
			expected: nil,
			errMsg:   "value 0 is not in range 1-31",
		},
		{
			name:     "parse month",
			value:    "0",
			minValue: 1,
			maxValue: 12,
			expected: nil,
			errMsg:   "value 0 is not in range 1-12",
		},
		{
			name:     "parse dayOfWeek",
			value:    "7",
			minValue: 0,
			maxValue: 6,
			expected: nil,
			errMsg:   "value 7 is not in range 0-6",
		},
		{
			name:     "parse weird string",
			value:    "@forty",
			minValue: 0,
			maxValue: 6,
			expected: nil,
			errMsg:   "failed to parse int: strconv.Atoi: parsing \"@forty\": invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := parseNumber(tt.value, tt.minValue, tt.maxValue)
			if tt.errMsg != "" {
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParse_parseRange(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		minValue int
		maxValue int
		expected []int
		errMsg   string
	}{
		{
			name:     "parse 2-5",
			value:    "2-5",
			minValue: 0,
			maxValue: 6,
			expected: []int{2, 3, 4, 5},
		},
		{
			name:     "parse 5-2",
			value:    "5-2",
			minValue: 0,
			maxValue: 6,
			expected: nil,
			errMsg:   "start value 5 is greater than end value 2",
		},
		{
			name:     "parse 0-0",
			value:    "0-0",
			minValue: 0,
			maxValue: 15,
			expected: []int{0},
		},
		{
			name:     "parse 5-5",
			value:    "5-5",
			minValue: 0,
			maxValue: 10,
			expected: []int{5},
		},
		{
			name:     "out of range parse 5-5",
			value:    "5-5",
			minValue: 0,
			maxValue: 4,
			expected: nil,
			errMsg:   "start value 5 is not in range 0-4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := parseRange(tt.value, tt.minValue, tt.maxValue)
			if tt.errMsg != "" {
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParse_parseSlash(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		minValue int
		maxValue int
		expected []int
		errMsg   string
	}{
		{
			name:     "parse /",
			value:    "0/15",
			minValue: 0,
			maxValue: 59,
			expected: []int{0, 15, 30, 45},
		},
		{
			name:     "parse 2/6",
			value:    "2/6",
			minValue: 0,
			maxValue: 45,
			expected: []int{2, 8, 14, 20, 26, 32, 38, 44},
		},
		{
			name:     "parse 1-6/6",
			value:    "1-2/6",
			minValue: 0,
			maxValue: 59,
			expected: []int{1},
		},
		{
			name:     "parse 1-7/6",
			value:    "1-10/6",
			minValue: 0,
			maxValue: 59,
			expected: []int{1, 7},
		},
		{
			name:     "parse */6",
			value:    "*/6",
			minValue: 0,
			maxValue: 45,
			expected: nil,
			errMsg:   "failed to parse number: strconv.Atoi: parsing \"*\": invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := parseSlash(tt.value, tt.minValue, tt.maxValue)
			if tt.errMsg != "" {
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParse_parseSlashRange(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		minValue int
		maxValue int
		expected []int
		errMsg   string
	}{
		{
			name:     "parse /",
			value:    "0-12/2",
			minValue: 0,
			maxValue: 59,
			expected: []int{0, 2, 4, 6, 8, 10, 12},
		},
		{
			name:     "parse 2/6",
			value:    "0-30/6",
			minValue: 0,
			maxValue: 45,
			expected: []int{0, 6, 12, 18, 24, 30},
		},
		{
			name:     "parse 0-45/6",
			value:    "1-45/6",
			minValue: 0,
			maxValue: 59,
			expected: []int{1, 7, 13, 19, 25, 31, 37, 43},
		},
		{
			name:     "parse */6",
			value:    "*/6",
			minValue: 0,
			maxValue: 45,
			expected: nil,
			errMsg:   "failed to parse number: strconv.Atoi: parsing \"*\": invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := parseSlash(tt.value, tt.minValue, tt.maxValue)
			if tt.errMsg != "" {
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_nextDay(t *testing.T) {
	type args struct {
		cronDayOfWeek     []int
		cronDayOfMonth    []int
		currentDayOfWeek  int
		currentDayOfMonth int
	}
	tests := []struct {
		name     string
		args     args
		expected int
	}{
		{
			name: "next day is the same day",
			args: args{
				cronDayOfWeek:     []int{0, 1, 2, 3, 4, 5, 6},
				cronDayOfMonth:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				currentDayOfWeek:  0,
				currentDayOfMonth: 1,
			},
			expected: 2,
		},
		{
			name: "weekdays",
			args: args{
				cronDayOfWeek:     []int{1, 2, 3, 4, 5},
				cronDayOfMonth:    nil,
				currentDayOfWeek:  0,
				currentDayOfMonth: 6,
			},
			expected: 7,
		},
		{
			name: "everyday",
			args: args{
				cronDayOfWeek:     allDaysOfWeek,
				cronDayOfMonth:    nil,
				currentDayOfWeek:  0,
				currentDayOfMonth: 6,
			},
			expected: 7,
		},
		{
			name: "sundays and the first of the month",
			args: args{
				cronDayOfWeek:     []int{int(time.Sunday)},
				cronDayOfMonth:    []int{1},
				currentDayOfWeek:  int(time.Friday),
				currentDayOfMonth: 6,
			},
			expected: 8,
		},
		{
			name: "first of the month",
			args: args{
				cronDayOfWeek:     nil,
				cronDayOfMonth:    []int{1},
				currentDayOfWeek:  int(time.Friday),
				currentDayOfMonth: 6,
			},
			expected: 1,
		},
		{
			name: "3rd of the month and Saturdays",
			args: args{
				cronDayOfWeek:     []int{int(time.Saturday)},
				cronDayOfMonth:    []int{3},
				currentDayOfWeek:  int(time.Sunday),
				currentDayOfMonth: 30,
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := nextDay(tt.args.cronDayOfWeek, tt.args.cronDayOfMonth, tt.args.currentDayOfWeek, tt.args.currentDayOfMonth)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
