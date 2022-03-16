package types

import (
	"testing"
)

func TestSdsAddress_Unmarshal(t *testing.T) {

	tests := []struct {
		name    string
		aa      string
		args    string
		wantErr bool
	}{
		{"test1", "stsds14c3em44vlh276cujnr2ez802uyjyeqrrsu9fuh", "stsds14c3em44vlh276cujnr2ez802uyjyeqrrsu9fuh", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := SdsAddressFromBech32(tt.aa)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			aa := &SdsAddress{}
			bz, err := addr.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			err = aa.UnmarshalJSON(bz)
			if !aa.Equals(addr) || (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		t.Run(tt.name, func(t *testing.T) {
			addr, err := SdsAddressFromBech32(tt.aa)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			aa := &SdsAddress{}
			bz, err := addr.MarshalYAML()
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			err = aa.UnmarshalYAML(bz.([]byte))
			if !aa.Equals(addr) || (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
