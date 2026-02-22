package xdm1041

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMeasurePower(t *testing.T) {
	tests := []struct {
		name      string
		voltageDC float64
		currentDC float64
		wantPower float64
	}{
		{"Positive values", 5.0, 2.0, 10.0},
		{"Zero voltage", 0.0, 2.0, 0.0},
		{"Zero current", 5.0, 0.0, 0.0},
		{"Both zero", 0.0, 0.0, 0.0},
		{"Negative current", 5.0, -1.0, -5.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Measure{
				VoltageDC: tt.voltageDC,
				CurrentDC: tt.currentDC,
			}
			got := m.Power()
			assert.Equal(t, tt.wantPower, got)
		})
	}
}

func TestSCPICommands(t *testing.T) {
	tests := []struct {
		name     string
		got      []byte
		expected string
	}{
		{"IDENTITY", IDETITY, "*IDN?"},
		{"RATE query", RATE, "RATE?"},
		{"RANGE query", RANGE, "RANGE?"},
		{"MEAS query", MEASS, "MEAS?"},
		{"VOLT DC AUTO", VOLT_DC, "CONF:VOLT:DC AUTO"},
		{"VOLT DC 50mV", VOLT_DC_50e, "CONF:VOLT:DC 50E-3"},
		{"VOLT DC 500mV", VOLT_DC_500e, "CONF:VOLT:DC 500E-3"},
		{"VOLT DC 5V", VOLT_DC_5, "CONF:VOLT:DC 5"},
		{"VOLT DC 50V", VOLT_DC_50, "CONF:VOLT:DC 50"},
		{"VOLT DC 500V", VOLT_DC_500, "CONF:VOLT:DC 500"},
		{"VOLT DC 1000V", VOLT_DC_1000, "CONF:VOLT:DC 1000"},
		{"VOLT AC AUTO", VOLT_AC, "CONF:VOLT:AC AUTO"},
		{"CURR DC 500uA", CURR_DC_500u, "CONF:CURR:DC 500E-6"},
		{"CURR DC 5mA", CURR_DC_5m, "CONF:CURR:DC 5E-3"},
		{"CURR DC 5A", CURR_DC_5, "CONF:CURR:DC 5"},
		{"CURR DC 10A", CURR_DC_10, "CONF:CURR:DC 10"},
		{"CURR AC AUTO", CURR_AC, "CONF:CURR:AC AUTO"},
		{"RES AUTO", RES_AUTO, "CONF:RES AUTO"},
		{"RES 500Ω", RES_500, "CONF:RES 500"},
		{"RES 5kΩ", RES_5k, "CONF:RES 5E3"},
		{"RES 50kΩ", RES_50k, "CONF:RES 50E3"},
		{"CAP AUTO", CAP_AUTO, "CONF:CAP AUTO"},
		{"CAP 50nF", CAP_50nF, "CONF:CAP 50E-9"},
		{"CAP 500nF", CAP_500nF, "CONF:CAP 500E-9"},
		{"CAP 5μF", CAP_5uF, "CONF:CAP 5E-6"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.got))
		})
	}
}

func TestTimeouts(t *testing.T) {
	assert.Equal(t, 350*time.Millisecond, COMMUNICATION_TIMEOUT)
	assert.Equal(t, time.Second, SWITH_RANGE_TIMEOUT)
}
