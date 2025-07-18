package ipresolver

import (
	"testing"
)

func Test_v6Ipshudi_Resolve(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name    string
		r       *v6Ipshudi
		args    args
		want    string
		wantErr bool
	}{
		{name: "1", r: &v6Ipshudi{}, args: args{ip: "2001:0db8:85a3:0000:0000:8a2e:0370:7334"}, want: "保留地址", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &v6Ipshudi{}
			got, err := r.Resolve(tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("v6Ipshudi.Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("v6Ipshudi.Resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_v6Ip77_Resolve(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name    string
		r       *v6Ip77
		args    args
		want    string
		wantErr bool
	}{
		{name: "1", r: &v6Ip77{}, args: args{ip: "2001:0db8:85a3:0000:0000:8a2e:0370:7334"}, want: "保留", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &v6Ip77{}
			got, err := r.Resolve(tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("v6Ip77.Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("v6Ip77.Resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}
