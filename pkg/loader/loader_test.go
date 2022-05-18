/*
 * Copyright © 2022 photowey (photowey@gmail.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package loader

import (
	"testing"

	"github.com/spf13/viper"
)

type (
	config struct {
		Author `yml:"author" toml:"author" properties:"author"`
	}
	Author struct {
		Name  string `yml:"name" toml:"name" properties:"name"`
		Email string `yml:"email" toml:"email" properties:"email"`
	}
)

func TestBind_yml(t *testing.T) {
	type args struct {
		fileName   string
		fileType   string
		dst        any
		searchPath []string
	}

	var dst config

	tests := []struct {
		name      string
		args      args
		wantName  string
		wantEmail string
		wantErr   bool
	}{
		{
			name: "test load config file-yml",
			args: args{
				fileName:   "iocgo_yml_test",
				fileType:   "yml",
				dst:        &dst,
				searchPath: []string{"./testdata/"},
			},
			wantName:  "photowey$yml",
			wantEmail: "photowey.yml@gmail.com",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Bind(tt.args.fileName, tt.args.fileType, tt.args.dst, tt.args.searchPath...); (err != nil) != tt.wantErr {
				t.Errorf("Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if dst.Name != tt.wantName {
				t.Errorf("Bind() error got = %v, want = %v", dst.Name, tt.wantName)
			}
			if dst.Email != tt.wantEmail {
				t.Errorf("Bind() error got = %v, want = %v", dst.Email, tt.wantEmail)
			}

			t.Log("ok")
		})
	}
}

func TestBind_toml(t *testing.T) {
	type args struct {
		fileName   string
		fileType   string
		dst        any
		searchPath []string
	}

	var dst config

	tests := []struct {
		name      string
		args      args
		wantName  string
		wantEmail string
		wantErr   bool
	}{
		{
			name: "test load config file-toml",
			args: args{
				fileName:   "iocgo_toml_test",
				fileType:   "toml",
				dst:        &dst,
				searchPath: []string{"./testdata/"},
			},
			wantName:  "photowey.toml",
			wantEmail: "photowey.toml@gmail.com",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Bind(tt.args.fileName, tt.args.fileType, tt.args.dst, tt.args.searchPath...); (err != nil) != tt.wantErr {
				t.Errorf("Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if dst.Name != tt.wantName {
				t.Errorf("Bind() error got = %v, want = %v", dst.Name, tt.wantName)
			}
			if dst.Email != tt.wantEmail {
				t.Errorf("Bind() error got = %v, want = %v", dst.Email, tt.wantEmail)
			}

			t.Log("ok")
		})
	}
}

func TestBind_properties(t *testing.T) {
	type args struct {
		fileName   string
		fileType   string
		dst        any
		searchPath []string
	}

	var dst config

	tests := []struct {
		name      string
		args      args
		wantName  string
		wantEmail string
		wantErr   bool
	}{
		{
			name: "test load config file-properties",
			args: args{
				fileName:   "iocgo_properties_test",
				fileType:   "properties",
				dst:        &dst,
				searchPath: []string{"./testdata/"},
			},
			wantName:  "photowey#properties",
			wantEmail: "photowey.properties@gmail.com",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Bind(tt.args.fileName, tt.args.fileType, tt.args.dst, tt.args.searchPath...); (err != nil) != tt.wantErr {
				t.Errorf("Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if viper.GetString("author.name") != tt.wantName {
				t.Errorf("Bind() error got = %v, want = %v", dst.Name, tt.wantName)
			}
			if dst.Name != tt.wantName {
				t.Errorf("Bind() error got = %v, want = %v", dst.Name, tt.wantName)
			}
			if dst.Email != tt.wantEmail {
				t.Errorf("Bind() error got = %v, want = %v", dst.Email, tt.wantEmail)
			}

			t.Log("ok")
		})
	}
}

func TestLoad(t *testing.T) {
	type args struct {
		fileName   string
		fileType   string
		searchPath []string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test load config file-yml",
			args: args{
				fileName:   "iocgo_yml_test",
				fileType:   "yml",
				searchPath: []string{"./testdata/"},
			},
			want:    "photowey$yml",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Load(tt.args.fileName, tt.args.fileType, tt.args.searchPath...); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			if viper.GetString("author.name") != tt.want {
				t.Errorf("Bind() error got = %v, want = %v", viper.GetString("author.name"), tt.want)
			}
		})
	}
}
