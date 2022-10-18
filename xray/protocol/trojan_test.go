package protocol

import (
	"reflect"
	"testing"
)

func TestParseTrojanUri(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name    string
		args    args
		want    *TrojanConfig
		wantErr bool
	}{
		{
			name: "Trojan",
			args: args{"trojan://a6a647d2-1234-4c19-a343-beeec21a66ac@127.0.0.1:51507"},
			want: &TrojanConfig{
				Password: "a6a647d2-1234-4c19-a343-beeec21a66ac",
				Host:     "127.0.0.1",
				Port:     51507,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTrojanUri(tt.args.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTrojanUri() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTrojanUri() got = %v, want %v", got, tt.want)
			}
		})
	}
}
