package timefmt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRelativeTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		diff     time.Duration
		expected string
	}{
		{
			name:     "now",
			diff:     30 * time.Second,
			expected: "now",
		},
		{
			name:     "minutes",
			diff:     5 * time.Minute,
			expected: "5m",
		},
		{
			name:     "hours",
			diff:     2 * time.Hour,
			expected: "2h",
		},
		{
			name:     "days",
			diff:     3 * 24 * time.Hour,
			expected: "3d",
		},
		{
			name:     "months",
			diff:     40 * 24 * time.Hour,
			expected: "1mo",
		},
		{
			name:     "years",
			diff:     400 * 24 * time.Hour,
			expected: "1y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			past := now.Add(-tt.diff)
			got := RelativeTime(past)
			assert.Equal(t, tt.expected, got)
		})
	}
}
