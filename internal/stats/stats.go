package stats

import "time"

// Result holds statistics collected during a slice operation.
type Result struct {
	TotalLines    int
	MatchedLines  int
	SkippedLines  int
	FilteredLines int
	EarliestTime  *time.Time
	LatestTime    *time.Time
	Duration      time.Duration
}

// Collector accumulates stats during processing.
type Collector struct {
	start         time.Time
	totalLines    int
	matchedLines  int
	skippedLines  int
	filteredLines int
	earliestTime  *time.Time
	latestTime    *time.Time
}

// NewCollector creates a new Collector and records the start time.
func NewCollector() *Collector {
	return &Collector{start: time.Now()}
}

// RecordTotal increments the total line count.
func (c *Collector) RecordTotal() {
	c.totalLines++
}

// RecordMatched increments matched lines and updates time bounds.
func (c *Collector) RecordMatched(t *time.Time) {
	c.matchedLines++
	if t != nil {
		if c.earliestTime == nil || t.Before(*c.earliestTime) {
			copy := *t
			c.earliestTime = &copy
		}
		if c.latestTime == nil || t.After(*c.latestTime) {
			copy := *t
			c.latestTime = &copy
		}
	}
}

// RecordSkipped increments the skipped line count.
func (c *Collector) RecordSkipped() {
	c.skippedLines++
}

// RecordFiltered increments the filtered-out line count.
func (c *Collector) RecordFiltered() {
	c.filteredLines++
}

// Finalise stops the timer and returns the accumulated Result.
func (c *Collector) Finalise() Result {
	return Result{
		TotalLines:    c.totalLines,
		MatchedLines:  c.matchedLines,
		SkippedLines:  c.skippedLines,
		FilteredLines: c.filteredLines,
		EarliestTime:  c.earliestTime,
		LatestTime:    c.latestTime,
		Duration:      time.Since(c.start),
	}
}
