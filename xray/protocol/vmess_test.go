package protocol

import (
	"reflect"
	"testing"
)

func TestParseVmessUri(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name    string
		args    args
		want    *VmessConfig
		wantErr bool
	}{
		{
			name: "vmess correct",
			args: args{"vmess://eyJ2IjoiMiIsInBzIjoiIOWkh+azqOaIluWIq+WQjSAgIiwiYWRkIjoiMTExLjExMS4xMTEuMTExIiwicG9ydCI6IjMyMDAwIiwiaWQiOiIxMzg2Zjg1ZS02NTdiLTRkNmUtOWQ1Ni03OGJhZGI3NWUxZmQiLCJhaWQiOiIxMDAiLCJzY3kiOiJ6ZXJvIiwibmV0IjoidGNwIiwidHlwZSI6Im5vbmUiLCJob3N0Ijoid3d3LmJiYi5jb20iLCJwYXRoIjoiLyIsInRscyI6InRscyIsInNuaSI6Ind3dy5jY2MuY29tIn0="},
			want: &VmessConfig{
				V:    "2",
				Ps:   " 备注或别名  ",
				Add:  "111.111.111.111",
				Port: "32000",
				Id:   "1386f85e-657b-4d6e-9d56-78badb75e1fd",
				Aid:  "100",
				Scy:  "zero",
				Net:  "tcp",
				Type: "none",
				Host: "www.bbb.com",
				Path: "/",
				Tls:  "tls",
				Sni:  "www.ccc.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVmessUri(tt.args.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVmessUri() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseVmessUri() got = %v, want %v", got, tt.want)
			}
		})
	}
}
