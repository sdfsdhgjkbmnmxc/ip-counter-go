package counter

import "testing"

func TestParseIPv4(t *testing.T) {
	tests := []struct {
		input   string
		want    uint32
		wantErr bool
	}{
		{"0.0.0.0", 0, false},
		{"255.255.255.255", 0xFFFFFFFF, false},
		{"192.168.1.1", 0xC0A80101, false},
		{"10.0.0.1", 0x0A000001, false},
		{"127.0.0.1", 0x7F000001, false},
		{"8.8.8.8", 0x08080808, false},
		{"1.2.3.4", 0x01020304, false},

		{"256.0.0.0", 0, true},
		{"1.256.1.1", 0, true},
		{"1.1.1.256", 0, true},
		{"1.1.1", 0, true},
		{"1.1.1.1.1", 0, true},
		{"abc.def.ghi.jkl", 0, true},
		{"", 0, true},
		{"1.2.3.4.5", 0, true},
		{"192.168.1", 0, true},
		{"192.168..1", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseIPv4(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIPv4(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseIPv4(%q) = 0x%X, want 0x%X", tt.input, got, tt.want)
			}
		})
	}
}
