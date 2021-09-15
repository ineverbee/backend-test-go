package currency_test

import (
	"testing"
	"time"

	"github.com/ineverbee/backend-test-go/internal/app/currency"
)

func TestCurrency_changeCurrency(t *testing.T) {
	testData := []struct {
		s        string
		n        int
		expected float64
	}{
		{"USD", 0, 0},
		{"EUR", 0, 0},
		{"USD", 10, 0.1366},
		{"EUR", 630, 7.2663},
		{"USD", 365901088, 4997200.2266},
		{"EUR", 209452600, 2415796.9287},
	}
	for _, tc := range testData {
		r, err := currency.ChangeCurrency(tc.s, tc.n)
		if r != tc.expected {
			t.Fatalf("encode expected '%v', but got '%v' with '%v'", tc.expected, r, err)
		}
		time.Sleep(1 * time.Second)
	}
}
