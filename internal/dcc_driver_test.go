package conductor

import (
	"reflect"
	"testing"
)

func TestParsePacket(t *testing.T) {
	tests := []struct {
		input   string
		expect  interface{}
		wantErr bool
	}{
		{
			input:   "<",
			wantErr: true,
		},
		{
			input:   ">",
			wantErr: true,
		},
		{
			input:   "<>",
			wantErr: true,
		},
		{
			input: "<iDCC-EX V-5.0.0 / MEGA / STANDARD_MOTOR_SHIELD G-3bddf4d>",
			expect: NewEvent("status", StatusEvent{
				StatusRaw: "V-5.0.0 / MEGA / STANDARD_MOTOR_SHIELD G-3bddf4d",
			}),
			wantErr: false,
		},
		{
			input:   "<p0>",
			expect:  NewEvent("power", PowerEvent{On: false}),
			wantErr: false,
		},
		{
			input:   "<p1>",
			expect:  NewEvent("power", PowerEvent{On: true}),
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			actual, err := parsePacket(test.input)

			if (err != nil) != test.wantErr {
				t.Errorf("parsePacket err = %v, wantErr = %v", err, test.wantErr)
			}

			if actual != nil {
				if !reflect.DeepEqual(actual, test.expect) {
					t.Errorf("parsePacket actual = %v, expect = %v", actual, test.expect)
				}
			}
		})
	}
}
