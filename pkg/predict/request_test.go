package predict

import "testing"

func TestPost(t *testing.T) {
	type args struct {
		server  string
		host    string
		auth    string
		data    []byte
		timeout int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test-post-timeout",
			args: args{
				server:  "139.155.92.20",
				host:    "cnvdxnmodel-service.tbdsversion.com",
				auth:    "e70e0dc4c3f74c7380bc3dfba91bf1d3",
				data:    []byte(`{"instances": [{"x1":6.2, "x2":2.2, "x3":1.1, "x4":1.2}]}`),
				timeout: 50,
			},
			want:    `{"predictions": [{"probability(1)":"0.08642621987903544","probability(-1)":"0.9135737801209646","y":"-1"}]}`,
			wantErr: true,
		},
		{
			name: "test-post-success",
			args: args{
				server:  "139.155.92.20",
				host:    "cnvdxnmodel-service.tbdsversion.com",
				auth:    "e70e0dc4c3f74c7380bc3dfba91bf1d3",
				data:    []byte(`{"instances": [{"x1":6.2, "x2":2.2, "x3":1.1, "x4":1.2}]}`),
				timeout: 500,
			},
			want:    `{"predictions": [{"probability(1)":"0.08642621987903544","probability(-1)":"0.9135737801209646","y":"-1"}]}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Post(tt.args.server, tt.args.host, tt.args.auth, tt.args.timeout, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("Post() got = %v, want %v", got, tt.want)
			}
		})
	}
}
