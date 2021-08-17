package protocol

import "testing"

func TestGetProtocol(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name    string
		args    args
		want    Type
		wantErr bool
	}{
		{name: "vmess test", args: args{"vmess://d6Npz6dxdHxmKY@127.0.0.1:51507#this+is+comment"}, want: Vmess},
		{name: "shadowsocks test", args: args{"ss://d6Npz6dxdHxmKY@127.0.0.1:51507#this+is+comment"}, want: Shadowsocks},
		{name: "error test", args: args{"trojan://d6Npz6dxdHxmKY@127.0.0.1:51507#this+is+comment"}, want: None, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetProtocol(tt.args.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProtocol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetProtocol() got = %v, want %v", got, tt.want)
			}
		})
	}
}
