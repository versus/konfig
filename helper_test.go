package konfig

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlagValue(t *testing.T) {
	tests := []struct {
		fv               *flagValue
		expectedString   string
		setString        string
		expectedSetError error
	}{
		{
			fv:               &flagValue{},
			expectedString:   "",
			setString:        "anything",
			expectedSetError: nil,
		},
	}

	for _, tc := range tests {
		str := tc.fv.String()
		assert.Equal(t, tc.expectedString, str)

		err := tc.fv.Set(tc.setString)
		assert.Equal(t, tc.expectedSetError, err)
	}
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		fieldName      string
		expectedTokens []string
	}{
		{"c", []string{"c"}},
		{"C", []string{"C"}},
		{"camel", []string{"camel"}},
		{"Camel", []string{"Camel"}},
		{"camelCase", []string{"camel", "Case"}},
		{"CamelCase", []string{"Camel", "Case"}},
		{"OneTwoThree", []string{"One", "Two", "Three"}},
		{"DatabaseURL", []string{"Database", "URL"}},
		{"DBEndpoints", []string{"DB", "Endpoints"}},
	}

	for _, tc := range tests {
		tokens := tokenize(tc.fieldName)
		assert.Equal(t, tc.expectedTokens, tokens)
	}
}

func TestGetFlagName(t *testing.T) {
	tests := []struct {
		fieldName        string
		expectedFlagName string
	}{
		{"c", "c"},
		{"C", "c"},
		{"camel", "camel"},
		{"Camel", "camel"},
		{"camelCase", "camel.case"},
		{"CamelCase", "camel.case"},
		{"OneTwoThree", "one.two.three"},
		{"DatabaseURL", "database.url"},
		{"DBEndpoints", "db.endpoints"},
	}

	for _, tc := range tests {
		flagName := getFlagName(tc.fieldName)
		assert.Equal(t, tc.expectedFlagName, flagName)
	}
}

func TestGetEnvVarName(t *testing.T) {
	tests := []struct {
		fieldName          string
		expectedEnvVarName string
	}{
		{"c", "C"},
		{"C", "C"},
		{"camel", "CAMEL"},
		{"Camel", "CAMEL"},
		{"camelCase", "CAMEL_CASE"},
		{"CamelCase", "CAMEL_CASE"},
		{"OneTwoThree", "ONE_TWO_THREE"},
		{"DatabaseURL", "DATABASE_URL"},
		{"DBEndpoints", "DB_ENDPOINTS"},
	}

	for _, tc := range tests {
		envVarName := getEnvVarName(tc.fieldName)
		assert.Equal(t, tc.expectedEnvVarName, envVarName)
	}
}

func TestGetFileVarName(t *testing.T) {
	tests := []struct {
		fieldName              string
		expectedFileEnvVarName string
	}{
		{"c", "C_FILE"},
		{"C", "C_FILE"},
		{"camel", "CAMEL_FILE"},
		{"Camel", "CAMEL_FILE"},
		{"camelCase", "CAMEL_CASE_FILE"},
		{"CamelCase", "CAMEL_CASE_FILE"},
		{"OneTwoThree", "ONE_TWO_THREE_FILE"},
		{"DatabaseURL", "DATABASE_URL_FILE"},
		{"DBEndpoints", "DB_ENDPOINTS_FILE"},
	}

	for _, tc := range tests {
		fileEnvVarName := getFileEnvVarName(tc.fieldName)
		assert.Equal(t, tc.expectedFileEnvVarName, fileEnvVarName)
	}
}

func TestDefineFlag(t *testing.T) {
	tests := []struct {
		name             string
		flagName         string
		defaultValue     string
		envName          string
		fileName         string
		expectedFlagName string
	}{
		{"SkipFlag", "-", "default", "SKIP_FLAG", "SKIP_FLAG_FILE", ""},
		{"ExampleFlag", "example.flag", "default", "EXAMPLE_FLAG", "EXAMPLE_FLAG_FILE", "example.flag"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defineFlag(tc.flagName, tc.defaultValue, tc.envName, tc.fileName)

			if tc.expectedFlagName != "" {
				fl := flag.Lookup(tc.expectedFlagName)
				assert.NotEmpty(t, fl)
			}
		})
	}
}

func TestGetFlagValue(t *testing.T) {
	tests := []struct {
		args              []string
		flagName          string
		expectedFlagValue string
	}{
		{[]string{"exe", "-enabled"}, "enabled", "true"},
		{[]string{"exe", "--enabled"}, "enabled", "true"},
		{[]string{"exe", "-enabled=false"}, "enabled", "false"},
		{[]string{"exe", "--enabled=false"}, "enabled", "false"},
		{[]string{"exe", "-enabled", "false"}, "enabled", "false"},
		{[]string{"exe", "--enabled", "false"}, "enabled", "false"},

		{[]string{"exe", "-port=-10"}, "port", "-10"},
		{[]string{"exe", "--port=-10"}, "port", "-10"},
		{[]string{"exe", "-port", "-10"}, "port", "-10"},
		{[]string{"exe", "--port", "-10"}, "port", "-10"},

		{[]string{"exe", "-text=content"}, "text", "content"},
		{[]string{"exe", "--text=content"}, "text", "content"},
		{[]string{"exe", "-text", "content"}, "text", "content"},
		{[]string{"exe", "--text", "content"}, "text", "content"},

		{[]string{"exe", "-enabled", "-text", "content"}, "enabled", "true"},
		{[]string{"exe", "--enabled", "--text", "content"}, "enabled", "true"},

		{[]string{"exec", "-service.name=go-service"}, "service.name", "go-service"},
		{[]string{"exec", "--service.name=go-service"}, "service.name", "go-service"},
		{[]string{"exec", "-service.name", "go-service"}, "service.name", "go-service"},
		{[]string{"exec", "--service.name", "go-service"}, "service.name", "go-service"},
	}

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	for _, tc := range tests {
		os.Args = tc.args
		flagValue := getFlagValue(tc.flagName)

		assert.Equal(t, tc.expectedFlagValue, flagValue)
	}
}
