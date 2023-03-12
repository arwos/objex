package providers

import "testing"

func TestUnit_ParseFTP(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name    string
		args    args
		host    string
		login   string
		passwd  string
		dir     string
		wantErr bool
	}{
		{
			name:    "Case0",
			args:    args{uri: "ftp://demo:p\\wd@ftp.server:21/folder"},
			host:    "ftp.server:21",
			login:   "demo",
			passwd:  "pwd",
			dir:     "/folder",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3, err := ParseFTP(tt.args.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.host {
				t.Errorf("ParseFTP() got = %v, want %v", got, tt.host)
			}
			if got1 != tt.login {
				t.Errorf("ParseFTP() got1 = %v, want %v", got1, tt.login)
			}
			if got2 != tt.passwd {
				t.Errorf("ParseFTP() got2 = %v, want %v", got2, tt.passwd)
			}
			if got3 != tt.dir {
				t.Errorf("ParseFTP() got3 = %v, want %v", got3, tt.dir)
			}
		})
	}
}
