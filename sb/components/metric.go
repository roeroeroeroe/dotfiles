package components

import (
	"fmt"

	"roe/sb/util"
)

type UsageMetric uint8

const (
	MetricUsed UsageMetric = iota
	MetricUsedPerc
	// "used/total"
	MetricUsedOutOf
	MetricFree
)

func (m UsageMetric) Validate(name string) {
	switch m {
	case MetricUsed, MetricUsedPerc, MetricUsedOutOf, MetricFree:
	default:
		panic(name + ": unknown metric")
	}
}

func (m UsageMetric) Format(used, total uint64) string {
	switch m {
	case MetricUsed:
		return util.HumanBytes(used)
	case MetricUsedPerc:
		return fmt.Sprintf("%.1f%%", float64(used)/float64(total)*100.0)
	case MetricUsedOutOf:
		return fmt.Sprintf("%s/%s", util.HumanBytes(used), util.HumanBytes(total))
	case MetricFree:
		return util.HumanBytes(total - used)
	}
	panic("unreachable")
}
