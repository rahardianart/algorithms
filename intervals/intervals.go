// Package intervals holds the runnable templates from intervals-study-plan.md.
//
// Intervals are closed [start, end] unless a function says otherwise.
package intervals

import "sort"

func mergeIntervals(intervals [][]int) [][]int {
	sort.Slice(intervals, func(i, j int) bool { return intervals[i][0] < intervals[j][0] })

	var merged [][]int
	for _, iv := range intervals {
		if len(merged) > 0 && iv[0] <= merged[len(merged)-1][1] {
			last := merged[len(merged)-1]
			last[1] = max(last[1], iv[1])
		} else {
			merged = append(merged, []int{iv[0], iv[1]})
		}
	}
	return merged
}

// intervalIntersection intersects two lists that are each sorted and internally disjoint.
func intervalIntersection(a, b [][]int) [][]int {
	var result [][]int
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		lo := max(a[i][0], b[j][0])
		hi := min(a[i][1], b[j][1])
		if lo <= hi {
			result = append(result, []int{lo, hi})
		}
		if a[i][1] < b[j][1] {
			i++
		} else {
			j++
		}
	}
	return result
}

// maxNonOverlapping returns how many intervals can be kept with no two overlapping.
// Intervals that only touch at an endpoint ([1,2] and [2,3]) do not overlap here.
func maxNonOverlapping(intervals [][]int) int {
	sort.Slice(intervals, func(i, j int) bool { return intervals[i][1] < intervals[j][1] })

	kept := 0
	lastEnd := 0
	for i, iv := range intervals {
		if i == 0 || iv[0] >= lastEnd {
			kept++
			lastEnd = iv[1]
		}
	}
	return kept
}

// minIntervalsToCover returns the fewest intervals whose union covers [start, end], or -1.
func minIntervalsToCover(intervals [][]int, start, end int) int {
	sort.Slice(intervals, func(i, j int) bool { return intervals[i][0] < intervals[j][0] })

	count, covered, i := 0, start, 0
	for covered < end {
		farthest := covered
		for i < len(intervals) && intervals[i][0] <= covered {
			farthest = max(farthest, intervals[i][1])
			i++
		}
		if farthest == covered {
			return -1
		}
		count++
		covered = farthest
	}
	return count
}

// maxConcurrent returns the largest number of half-open [start, end) intervals active at once.
func maxConcurrent(intervals [][]int) int {
	type event struct{ time, delta int }
	events := make([]event, 0, 2*len(intervals))
	for _, iv := range intervals {
		events = append(events, event{iv[0], +1}, event{iv[1], -1})
	}
	sort.Slice(events, func(i, j int) bool {
		if events[i].time != events[j].time {
			return events[i].time < events[j].time
		}
		return events[i].delta < events[j].delta // ends before starts at the same time
	})

	active, best := 0, 0
	for _, e := range events {
		active += e.delta
		best = max(best, active)
	}
	return best
}
