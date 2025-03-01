package ui

import (
	"testing"
	"time"
)

func TestAllEqual(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name     string
		input    []versionInfo
		expected bool
	}{
		// success
		{
			name: "all versions equal",
			input: []versionInfo{
				{
					version:      "v0.1.0",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
				{
					version:      "v0.1.0",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
				{
					version:      "v0.1.0",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
			},
			expected: true,
		},
		{
			name: "different versions",
			input: []versionInfo{
				{
					version:      "v0.3.0",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
				{
					version:      "v0.2.0",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
				{
					version:      "v0.1.0",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
			},
			expected: false,
		},
		{
			name: "first version empty",
			input: []versionInfo{
				{
					version:      "",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
				{
					version:      "v0.1.0",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
				{
					version:      "v0.1.0",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
			},
			expected: true,
		},
		{
			name: "second version empty",
			input: []versionInfo{
				{
					version:      "v0.1.0",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
				{
					version:      "",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
				{
					version:      "v0.1.0",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
			},
			expected: true,
		},
		{
			name: "two empty",
			input: []versionInfo{
				{
					version:      "v0.1.0",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
				{
					version:      "",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
				{
					version:      "",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
			},
			expected: true,
		},
		{
			name: "all empty",
			input: []versionInfo{
				{
					version:      "",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
				{
					version:      "",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
				{
					version:      "",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
			},
			expected: false,
		},
		{
			name: "one error",
			input: []versionInfo{
				{
					version:      "",
					errMsg:       "some error",
					registeredAt: &now,
					notFound:     false,
				},
				{
					version:      "v0.1.0",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
				{
					version:      "v0.1.0",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
			},
			expected: false,
		},
		{
			name: "one not found",
			input: []versionInfo{
				{
					version:      "v0.1.0",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
				{
					version:      "",
					errMsg:       "",
					registeredAt: &now,
					notFound:     true,
				},
				{
					version:      "v0.1.0",
					errMsg:       "",
					registeredAt: &now,
					notFound:     false,
				},
			},
			expected: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := allEqual(tt.input)

			if got != tt.expected {
				if got != tt.expected {
					t.Errorf("got: %v, expected: %v", got, tt.expected)
				}
			}
		})
	}
}
