package protocol

import (
	"reflect"
	"testing"
)

func TestParseShadowsocksUri(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name    string
		args    args
		want    *ShadowsocksConfig
		wantErr bool
	}{
		{
			name: "shadowsocks without comment",
			args: args{"ss://YWVzLTI1Ni1nY206dGVzdHBhc3N3b3Jk@127.0.0.1:51507"},
			want: &ShadowsocksConfig{
				Method:   "aes-256-gcm",
				Password: "testpassword",
				Hostname: "127.0.0.1",
				Port:     51507,
			},
			wantErr: false,
		}, {
			name: "shadowsocks with comment",
			args: args{"ss://YWVzLTI1Ni1nY206dGVzdHBhc3N3b3Jk@127.0.0.1:51507#test%20comment"},
			want: &ShadowsocksConfig{
				Method:   "aes-256-gcm",
				Password: "testpassword",
				Hostname: "127.0.0.1",
				Port:     51507,
				Comment:  "test comment",
			},
			wantErr: false,
		}, {
			name:    "shadowsocks incorrect",
			args:    args{"ss://illegalbase64@127.0.0.1:51507"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseShadowsocksUri(tt.args.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseShadowsocksUri() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseShadowsocksUri() got = %v, want %v", got, tt.want)
			}
		})
	}
}
