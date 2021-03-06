package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadConfig(t *testing.T) {
	type args struct {
		path *string
	}

	var (
		translationOK         = "../tests/data/config/translation-ok.yml"
		translationBadCode    = "../tests/data/config/translation-bad-code.yml"
		translationBadDefault = "../tests/data/config/translation-bad-default.yml"
		translationNoDefault  = "../tests/data/config/translation-no-default.yml"
		translationNoEnabled  = "../tests/data/config/translation-no-enabled.yml"
	)

	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "Enabled Languages Present",
			args: args{
				path: &translationOK,
			},
			want: Config{
				TranslationEnabled: true,
				DefaultLanguage:    "en",
			},
			wantErr: false,
		},
		{
			name: "Default Language Missing",
			args: args{
				path: &translationNoDefault,
			},
			wantErr: true,
			errMsg:  "Translation enabled but no default language specified",
		},
		{
			name: "Enabled Languages Missing",
			args: args{
				path: &translationNoEnabled,
			},
			wantErr: true,
			errMsg:  "Translation enabled but no languages enabled",
		},
		{
			name: "Default Language Not Defined",
			args: args{
				path: &translationBadDefault,
			},
			wantErr: true,
			errMsg:  "Default language 'de' not found",
		},
		{
			name: "Language Code Missing",
			args: args{
				path: &translationBadCode,
			},
			wantErr: true,
			errMsg:  "Language code 'de' not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			c, err := loadConfig(tt.args.path)

			if tt.wantErr {
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.Equal(t, tt.want.TranslationEnabled, c.TranslationEnabled)
				assert.Equal(t, tt.want.DefaultLanguage, c.DefaultLanguage)
			}

		})
	}
}
