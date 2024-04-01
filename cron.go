package cronroutine

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	nanosecond = 0
	second     = 0
)

type Cron struct {
	Minute     []int
	Hour       []int
	DayOfMonth []int
	Month      []int
	DayOfWeek  []int
}

func ParseCron(cronConfig string) (*Cron, error) {
	fields := strings.Fields(cronConfig)
	if len(fields) != 5 {
		return nil, fmt.Errorf("given %d but need 5 fields for valid cron config: %s", len(fields), cronConfig)
	}

	return NewCron(fields[0], fields[1], fields[2], fields[3], fields[4])
}

func NewCron(minute, hour, dayOfMonth, month, dayOfWeek string) (*Cron, error) {
	c := &Cron{}
	m, err := parseField(minute, 0, 59)
	if err != nil {
		return nil, fmt.Errorf("failed to parse minute: %w", err)
	}
	c.Minute = m

	h, err := parseField(hour, 0, 23)
	if err != nil {
		return nil, fmt.Errorf("failed to parse hour: %w", err)
	}
	c.Hour = h

	dm, err := parseField(dayOfMonth, 1, 31)
	if err != nil {
		return nil, fmt.Errorf("failed to parse day of month: %w", err)
	}
	c.DayOfMonth = dm

	mo, err := parseField(month, 1, 12)
	if err != nil {
		return nil, fmt.Errorf("failed to parse month: %w", err)
	}
	c.Month = mo

	dw, err := parseField(dayOfWeek, 0, 6)
	if err != nil {
		return nil, fmt.Errorf("failed to parse day of week: %w", err)
	}
	c.DayOfWeek = dw

	// weird logic
	if dayOfMonth == "*" && dayOfWeek != "*" {
		c.DayOfMonth = nil
	} else if dayOfMonth != "*" && dayOfWeek == "*" {
		c.DayOfWeek = nil
	}

	return c, nil
}

func parseField(field string, minValue int, maxValue int) ([]int, error) {
	if field == "*" {
		return sliceWithStep(minValue, maxValue, 1), nil
	}

	values := strings.Split(field, ",")
	numbers := make([]int, 0, len(values))
	for _, value := range values {
		number, err := parseNumber(value, minValue, maxValue)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number: %w", err)
		}
		numbers = append(numbers, number...)
	}

	return sortUnique(numbers), nil
}

func parseNumber(value string, minValue int, maxValue int) ([]int, error) {
	if value == "*" {
		return sliceWithStep(minValue, maxValue, 1), nil
	}

	if strings.Contains(value, "/") {
		ret, err := parseSlash(value, minValue, maxValue)
		if err != nil {
			return nil, fmt.Errorf("failed to parse slash %q: %w", value, err)
		}
		return ret, nil
	}

	if strings.Contains(value, "-") {
		ret, err := parseRange(value, minValue, maxValue)
		if err != nil {
			return nil, fmt.Errorf("failed to parse range %q: %w", value, err)
		}
		return ret, nil
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("failed to parse int: %w", err)
	}

	if number < minValue || number > maxValue {
		return nil, fmt.Errorf("value %d is not in range %d-%d", number, minValue, maxValue)
	}

	return []int{number}, nil
}

func parseRange(value string, minValue int, maxValue int) ([]int, error) {
	parts := strings.Split(value, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected 2 parts but got %d when parsing %s", len(parts), value)
	}

	startNum, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse start int: %w", err)
	}

	endNum, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse end int: %w", err)
	}

	if startNum < minValue || startNum > maxValue {
		return nil, fmt.Errorf("start value %d is not in range %d-%d", startNum, minValue, maxValue)
	}

	if endNum < minValue || endNum > maxValue {
		return nil, fmt.Errorf("end value %d is not in range %d-%d", endNum, minValue, maxValue)
	}

	if startNum > endNum {
		return nil, fmt.Errorf("start value %d is greater than end value %d", startNum, endNum)
	}

	return sliceWithStep(startNum, endNum, 1), nil
}

func parseSlash(value string, minValue int, maxValue int) ([]int, error) {
	parts := strings.Split(value, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected 2 parts but got %d when parsing %s", len(parts), value)
	}

	step, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse step int: %w", err)
	}

	if step < 1 {
		return nil, fmt.Errorf("step value %d is less than 1", step)
	}

	if strings.Contains(parts[0], "-") {
		return parseSlashRange(parts[0], parts[1], minValue, maxValue)
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse number: %w", err)
	}

	if start < minValue || start > maxValue {
		return nil, fmt.Errorf("value %d is not in range %d-%d", start, minValue, maxValue)
	}

	return sliceWithStep(start, maxValue, step), nil
}

func parseSlashRange(value string, step string, minValue int, maxValue int) ([]int, error) {
	parts := strings.Split(value, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected 2 parts but got %d when parsing %s", len(parts), value)
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse start int: %w", err)
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse end int: %w", err)
	}

	if start < minValue || start > maxValue {
		return nil, fmt.Errorf("start value %d is not in range %d-%d", start, minValue, maxValue)
	}

	if end < minValue || end > maxValue {
		return nil, fmt.Errorf("end value %d is not in range %d-%d", end, minValue, maxValue)
	}

	if start > end {
		return nil, fmt.Errorf("start value %d is greater than end value %d", start, end)
	}

	stepNum, err := strconv.Atoi(step)
	if err != nil {
		return nil, fmt.Errorf("failed to parse step int: %w", err)
	}

	if stepNum < 1 {
		return nil, fmt.Errorf("step value %d is less than 1", stepNum)
	}

	return sliceWithStep(start, end, stepNum), nil
}

func (c *Cron) Next() time.Time {
	return c.next(time.Now().UTC())
}

func (c *Cron) next(t time.Time) time.Time {
	ct := newCTime(t)

	ct.minute = getFirstElementGreaterThan(c.Minute, ct.minute)

	if !ct.time().After(t) {
		ct.hour = getFirstElementGreaterThan(c.Hour, ct.hour)
	}

	if !ct.time().After(t) {
		ct.dayOfMonth = nextDay(c.DayOfWeek, c.DayOfMonth, ct.dayOfWeek, ct.dayOfMonth)
	}

	if !ct.time().After(t) {
		ct.month = getFirstElementGreaterThan(c.Month, ct.month)
	}

	if !ct.time().After(t) {
		ct.year++
	}

	return ct.time()
}

func (c *Cron) NextFor(t time.Duration) []time.Time {
	return c.nextFor(time.Now().UTC(), t)
}

func (c *Cron) nextFor(start time.Time, t time.Duration) []time.Time {
	times := []time.Time{}
	for next := c.next(start); next.Before(start.UTC().Add(t)); next = c.next(next) {
		times = append(times, next)
	}

	return times
}

func nextDay(cronDayOfWeek []int, cronDayOfMonth []int, currentDayOfWeek int, currentDayOfMonth int) int {
	// 0 1 * * *
	if cronDayOfMonth == nil && cronDayOfWeek == nil {
		return currentDayOfMonth + 1
	}

	o := []int{currentDayOfMonth}
	// 0 1 1 * *
	if cronDayOfMonth != nil {
		o = append(o, getFirstElementGreaterThan(cronDayOfMonth, currentDayOfMonth))
	}

	// 0 1 * * 1
	if cronDayOfWeek != nil {
		daysUntilJob := (getFirstElementGreaterThan(cronDayOfWeek, currentDayOfWeek) - currentDayOfWeek + 7) % 7
		o = append(o, (currentDayOfMonth+daysUntilJob)%31)
	}

	slices.Sort(o)

	return o[(slices.Index(o, currentDayOfMonth)+1)%len(o)]
}
