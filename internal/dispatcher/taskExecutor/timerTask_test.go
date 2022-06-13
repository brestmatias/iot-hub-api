package taskExecutor

import (
	"testing"
	"time"
)

type inTimeTest struct {
	from, duration, now string
	expected            bool
}

var inTimeTests = []inTimeTest{
	{"21:00:00", "0h10m0s", "21:05:00", true},
	{"21:00:00", "0h1m1s", "21:01:01", true},
	{"21:00:00", "0h1m0s", "21:01:01", false},
	{"21:00:00", "0h1m0s", "20:00:00", false},

	{"12:00:00", "23h0m0s", "20:00:00", true},
	{"12:00:00", "23h0m0s", "01:00:00", true},
	{"12:00:00", "23h0m0s", "11:00:00", true},
	{"12:00:00", "23h0m0s", "11:00:01", false},

	{"01:00:00", "0h0m0s", "01:00:00", true},
	{"01:00:00", "0h1m0s", "01:01:00", true},
	{"01:00:00", "24h0m0s", "01:00:00", true},
	{"01:00:00", "23h0m0s", "11:00:01", true},
}

func TestIsInTimeSpan(t *testing.T) {
	for _, test := range inTimeTests {
		duration, _ := time.ParseDuration(test.duration)
		from, _ := time.Parse(HMSLayout, test.from)
		check, _ := time.Parse(HMSLayout, test.now)
		to, _ := time.Parse(HMSLayout, from.Add(duration).Format(HMSLayout))

		t.Logf("Running test isInTimeSpan for ➡️ From: %v Duration: %v Check: %v", test.from, test.duration, test.now)
		if result := isInTimeSpan(from, to, check); result != test.expected {
			t.Errorf("IsInTimeSpan function result %v not equal to expected %v \n\t From: %v \n\t Duration: %v \n\t Check: %v",
				result, test.expected, test.from, test.duration, test.now)
		}
	}
}
