package konfig

import (
	"flag"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type config struct {
	unexported         string
	SkipFlag           string          `flag:"-"`
	SkipFlagEnv        string          `flag:"-" env:"-"`
	SkipFlagEnvFile    string          `flag:"-" env:"-" fileenv:"-"`
	FieldString        string          // `flag:"field.string" env:"FIELD_STRING" fileenv:"FIELD_STRING_FILE"`
	FieldBool          bool            // `flag:"field.bool" env:"FIELD_BOOL" fileenv:"FIELD_BOOL_FILE"`
	FieldFloat32       float32         // `flag:"field.float32" env:"FIELD_FLOAT32" fileenv:"FIELD_FLOAT32_FILE"`
	FieldFloat64       float64         // `flag:"field.float64" env:"FIELD_FLOAT64" fileenv:"FIELD_FLOAT64_FILE"`
	FieldInt           int             // `flag:"field.int" env:"FIELD_INT" fileenv:"FIELD_INT_FILE"`
	FieldInt8          int8            // `flag:"field.int8" env:"FIELD_INT8" fileenv:"FIELD_INT8_FILE"`
	FieldInt16         int16           // `flag:"field.int16" env:"FIELD_INT16" fileenv:"FIELD_INT16_FILE"`
	FieldInt32         int32           // `flag:"field.int32" env:"FIELD_INT32" fileenv:"FIELD_INT32_FILE"`
	FieldInt64         int64           // `flag:"field.int64" env:"FIELD_INT64" fileenv:"FIELD_INT64_FILE"`
	FieldUint          uint            // `flag:"field.uint" env:"FIELD_UINT" fileenv:"FIELD_UINT_FILE"`
	FieldUint8         uint8           // `flag:"field.uint8" env:"FIELD_UINT8" fileenv:"FIELD_UINT8_FILE"`
	FieldUint16        uint16          // `flag:"field.uint16" env:"FIELD_UINT16" fileenv:"FIELD_UINT16_FILE"`
	FieldUint32        uint32          // `flag:"field.uint32" env:"FIELD_UINT32" fileenv:"FIELD_UINT32_FILE"`
	FieldUint64        uint64          // `flag:"field.uint64" env:"FIELD_UINT64" fileenv:"FIELD_UINT64_FILE"`
	FieldDuration      time.Duration   // `flag:"field.duration" env:"FIELD_DURATION" fileenv:"FIELD_DURATION_FILE"`
	FieldURL           url.URL         // `flag:"field.url" env:"FIELD_URL" fileenv:"FIELD_URL_FILE"`
	FieldStringArray   []string        // `flag:"field.string.array" env:"FIELD_STRING_ARRAY" fileenv:"FIELD_STRING_ARRAY_FILE" sep:","`
	FieldBoolArray     []bool          // `flag:"field.bool.array" env:"FIELD_BOOL_ARRAY" fileenv:"FIELD_BOOL_ARRAY_FILE" sep:","`
	FieldFloat32Array  []float32       // `flag:"field.float32.array" env:"FIELD_FLOAT32_ARRAY" fileenv:"FIELD_FLOAT32_ARRAY_FILE" sep:","`
	FieldFloat64Array  []float64       // `flag:"field.float64.array" env:"FIELD_FLOAT64_ARRAY" fileenv:"FIELD_FLOAT64_ARRAY_FILE" sep:","`
	FieldIntArray      []int           // `flag:"field.int.array" env:"FIELD_INT_ARRAY" fileenv:"FIELD_INT_ARRAY_FILE" sep:","`
	FieldInt8Array     []int8          // `flag:"field.int8.array" env:"FIELD_INT8_ARRAY" fileenv:"FIELD_INT8_ARRAY_FILE" sep:","`
	FieldInt16Array    []int16         // `flag:"field.int16.array" env:"FIELD_INT16_ARRAY" fileenv:"FIELD_INT16_ARRAY_FILE" sep:","`
	FieldInt32Array    []int32         // `flag:"field.int32.array" env:"FIELD_INT32_ARRAY" fileenv:"FIELD_INT32_ARRAY_FILE" sep:","`
	FieldInt64Array    []int64         // `flag:"field.int64.array" env:"FIELD_INT64_ARRAY" fileenv:"FIELD_INT64_ARRAY_FILE" sep:","`
	FieldUintArray     []uint          // `flag:"field.uint.array" env:"FIELD_UINT_ARRAY" fileenv:"FIELD_UINT_ARRAY_FILE" sep:","`
	FieldUint8Array    []uint8         // `flag:"field.uint8.array" env:"FIELD_UINT8_ARRAY" fileenv:"FIELD_UINT8_ARRAY_FILE" sep:","`
	FieldUint16Array   []uint16        // `flag:"field.uint16.array" env:"FIELD_UINT16_ARRAY" fileenv:"FIELD_UINT16_ARRAY_FILE" sep:","`
	FieldUint32Array   []uint32        // `flag:"field.uint32.array" env:"FIELD_UINT32_ARRAY" fileenv:"FIELD_UINT32_ARRAY_FILE" sep:","`
	FieldUint64Array   []uint64        // `flag:"field.uint64.array" env:"FIELD_UINT64_ARRAY" fileenv:"FIELD_UINT64_ARRAY_FILE" sep:","`
	FieldDurationArray []time.Duration // `flag:"field.duration.array" env:"FIELD_DURATION_ARRAY" fileenv:"FIELD_DURATION_ARRAY_FILE" sep:","`
	FieldURLArray      []url.URL       // `flag:"field.url.array" env:"FIELD_URL_ARRAY" fileenv:"FIELD_URL_ARRAY_FILE" sep:","`
}

func TestDefaultOptions(t *testing.T) {
	o := defaultOptions()

	assert.False(t, o.debug)
	assert.False(t, o.telepresence)
}

func TestDebug(t *testing.T) {
	tests := []struct {
		options         *options
		expectedOptions *options
	}{
		{
			&options{},
			&options{
				debug: true,
			},
		},
	}

	for _, tc := range tests {
		opt := Debug()
		opt.apply(tc.options)

		assert.Equal(t, tc.expectedOptions, tc.options)
	}
}

func TestTelepresence(t *testing.T) {
	tests := []struct {
		options         *options
		expectedOptions *options
	}{
		{
			&options{},
			&options{
				telepresence: true,
			},
		},
	}

	for _, tc := range tests {
		opt := Telepresence()
		opt.apply(tc.options)

		assert.Equal(t, tc.expectedOptions, tc.options)
	}
}

func TestOptionsString(t *testing.T) {
	tests := []struct {
		name           string
		o              options
		expectedString string
	}{
		{
			"NoOption",
			options{},
			"",
		},
		{
			"WithDebug",
			options{
				debug: true,
			},
			"Debug",
		},
		{
			"WithTelepresence",
			options{
				telepresence: true,
			},
			"Telepresence",
		},
		{
			"WithAll",
			options{
				debug:        true,
				telepresence: true,
			},
			"Debug + Telepresence",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			str := tc.o.String()

			assert.Equal(t, tc.expectedString, str)
		})
	}
}

func TestPrint(t *testing.T) {
	tests := []struct {
		name string
		o    options
		msg  string
		args []interface{}
	}{
		{
			"WithoutDebug",
			options{},
			"debugging ...",
			nil,
		},
		{
			"WithDebug",
			options{
				debug: true,
			},
			"debugging ...",
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.o.print(tc.msg, tc.args...)
		})
	}
}

func TestGetFieldValue(t *testing.T) {
	type env struct {
		varName string
		value   string
	}

	type file struct {
		varName string
		value   string
	}

	tests := []struct {
		name                   string
		args                   []string
		envConfig              env
		fileConfig             file
		field, flag, env, file string
		o                      options
		expectedValue          string
		expectedFromFile       bool
	}{
		{
			"SkipFlag",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "-", "LOG_LEVEL", "LOG_LEVEL_FILE",
			options{},
			"info",
			false,
		},
		{
			"SkipFlagAndEnv",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "-", "-", "LOG_LEVEL_FILE",
			options{},
			"error",
			true,
		},
		{
			"SkipFlagAndEnvAndFile",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "-", "-", "-",
			options{},
			"",
			false,
		},
		{
			"FromFlag",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			options{},
			"debug",
			false,
		},
		{
			"FromFlag",
			[]string{"/path/to/executable", "--log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			options{},
			"debug",
			false,
		},
		{
			"FromFlag",
			[]string{"/path/to/executable", "-log.level", "debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			options{},
			"debug",
			false,
		},
		{
			"FromFlag",
			[]string{"/path/to/executable", "--log.level", "debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			options{},
			"debug",
			false,
		},
		{
			"FromEnvironmentVariable",
			[]string{"/path/to/executable"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			options{},
			"info",
			false,
		},
		{
			"FromFiles",
			[]string{"/path/to/executable"},
			env{"LOG_LEVEL", ""},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			options{},
			"error",
			true,
		},
		{
			"FromFilesWithTelepresenceOption",
			[]string{"/path/to/executable"},
			env{"LOG_LEVEL", ""},
			file{"LOG_LEVEL_FILE", "info"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			options{telepresence: true},
			"info",
			true,
		},
	}

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set value using a flag
			os.Args = tc.args

			// Set value in an environment variable
			err := os.Setenv(tc.envConfig.varName, tc.envConfig.value)
			assert.NoError(t, err)
			defer os.Unsetenv(tc.envConfig.varName)

			// Testing Telepresence option
			if tc.o.telepresence {
				err := os.Setenv(telepresenceEnvVar, "/")
				assert.NoError(t, err)
				defer os.Unsetenv(telepresenceEnvVar)
			}

			// Write value in a temporary config file

			tmpfile, err := ioutil.TempFile("", "gotest_")
			assert.NoError(t, err)
			defer os.Remove(tmpfile.Name())

			_, err = tmpfile.WriteString(tc.fileConfig.value)
			assert.NoError(t, err)

			err = tmpfile.Close()
			assert.NoError(t, err)

			err = os.Setenv(tc.fileConfig.varName, tmpfile.Name())
			assert.NoError(t, err)
			defer os.Unsetenv(tc.fileConfig.varName)

			// Verify

			value, fromFile := tc.o.getFieldValue(tc.field, tc.flag, tc.env, tc.file)
			assert.Equal(t, tc.expectedValue, value)
			assert.Equal(t, tc.expectedFromFile, fromFile)
		})
	}
}

func TestSetString(t *testing.T) {
	tests := []struct {
		name       string
		o          options
		field      string
		fieldName  string
		fieldValue string
		expected   string
	}{
		{
			name:       "NewValue",
			o:          options{},
			field:      "",
			fieldName:  "Field",
			fieldValue: "milad",
			expected:   "milad",
		},
		{
			name:       "NoNewValue",
			o:          options{},
			field:      "milad",
			fieldName:  "Field",
			fieldValue: "milad",
			expected:   "milad",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setString(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetBool(t *testing.T) {
	tests := []struct {
		name       string
		o          options
		field      bool
		fieldName  string
		fieldValue string
		expected   bool
	}{
		{
			name:       "NewValue",
			o:          options{},
			field:      false,
			fieldName:  "Field",
			fieldValue: "true",
			expected:   true,
		},
		{
			name:       "NoNewValue",
			o:          options{},
			field:      true,
			fieldName:  "Field",
			fieldValue: "true",
			expected:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setBool(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetFloat32(t *testing.T) {
	tests := []struct {
		name       string
		o          options
		field      float32
		fieldName  string
		fieldValue string
		expected   float32
	}{
		{
			name:       "NewValue",
			o:          options{},
			field:      0,
			fieldName:  "Field",
			fieldValue: "3.1415",
			expected:   3.1415,
		},
		{
			name:       "NoNewValue",
			o:          options{},
			field:      3.1415,
			fieldName:  "Field",
			fieldValue: "3.1415",
			expected:   3.1415,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setFloat(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetFloat64(t *testing.T) {
	tests := []struct {
		name       string
		o          options
		field      float64
		fieldName  string
		fieldValue string
		expected   float64
	}{
		{
			name:       "NewValue",
			o:          options{},
			field:      0,
			fieldName:  "Field",
			fieldValue: "3.14159265359",
			expected:   3.14159265359,
		},
		{
			name:       "NoNewValue",
			o:          options{},
			field:      3.14159265359,
			fieldName:  "Field",
			fieldValue: "3.14159265359",
			expected:   3.14159265359,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setFloat(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetInt(t *testing.T) {
	tests := []struct {
		name       string
		o          options
		field      int
		fieldName  string
		fieldValue string
		expected   int
	}{
		{
			name:       "NewValue",
			o:          options{},
			field:      0,
			fieldName:  "Field",
			fieldValue: "-27",
			expected:   -27,
		},
		{
			name:       "NoNewValue",
			o:          options{},
			field:      -27,
			fieldName:  "Field",
			fieldValue: "-27",
			expected:   -27,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setInt(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetInt8(t *testing.T) {
	tests := []struct {
		name       string
		o          options
		field      int8
		fieldName  string
		fieldValue string
		expected   int8
	}{
		{
			name:       "NewValue",
			o:          options{},
			field:      0,
			fieldName:  "Field",
			fieldValue: "-128",
			expected:   -128,
		},
		{
			name:       "NoNewValue",
			o:          options{},
			field:      -128,
			fieldName:  "Field",
			fieldValue: "-128",
			expected:   -128,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setInt(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetInt16(t *testing.T) {
	tests := []struct {
		name       string
		o          options
		field      int16
		fieldName  string
		fieldValue string
		expected   int16
	}{
		{
			name:       "NewValue",
			o:          options{},
			field:      0,
			fieldName:  "Field",
			fieldValue: "-32768",
			expected:   -32768,
		},
		{
			name:       "NoNewValue",
			o:          options{},
			field:      -32768,
			fieldName:  "Field",
			fieldValue: "-32768",
			expected:   -32768,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setInt(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetInt32(t *testing.T) {
	tests := []struct {
		name       string
		o          options
		field      int32
		fieldName  string
		fieldValue string
		expected   int32
	}{
		{
			name:       "NewValue",
			o:          options{},
			field:      0,
			fieldName:  "Field",
			fieldValue: "-2147483648",
			expected:   -2147483648,
		},
		{
			name:       "NoNewValue",
			o:          options{},
			field:      -2147483648,
			fieldName:  "Field",
			fieldValue: "-2147483648",
			expected:   -2147483648,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setInt(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetInt64(t *testing.T) {
	tests := []struct {
		name       string
		o          options
		field      int64
		fieldName  string
		fieldValue string
		expected   int64
	}{
		{
			name:       "NewValue",
			o:          options{},
			field:      0,
			fieldName:  "Field",
			fieldValue: "-9223372036854775808",
			expected:   -9223372036854775808,
		},
		{
			name:       "NoNewValue",
			o:          options{},
			field:      -9223372036854775808,
			fieldName:  "Field",
			fieldValue: "-9223372036854775808",
			expected:   -9223372036854775808,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setInt(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetDuration(t *testing.T) {
	tests := []struct {
		name       string
		o          options
		field      time.Duration
		fieldName  string
		fieldValue string
		expected   time.Duration
	}{
		{
			name:       "NewValue",
			o:          options{},
			field:      0,
			fieldName:  "Field",
			fieldValue: "1h0m0s",
			expected:   time.Hour,
		},
		{
			name:       "NoNewValue",
			o:          options{},
			field:      time.Hour,
			fieldName:  "Field",
			fieldValue: "1h0m0s",
			expected:   time.Hour,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setInt(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetUint(t *testing.T) {
	tests := []struct {
		name       string
		o          options
		field      uint
		fieldName  string
		fieldValue string
		expected   uint
	}{
		{
			name:       "NewValue",
			o:          options{},
			field:      0,
			fieldName:  "Field",
			fieldValue: "27",
			expected:   27,
		},
		{
			name:       "NoNewValue",
			o:          options{},
			field:      27,
			fieldName:  "Field",
			fieldValue: "27",
			expected:   27,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setUint(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetUint8(t *testing.T) {
	tests := []struct {
		name       string
		o          options
		field      uint8
		fieldName  string
		fieldValue string
		expected   uint8
	}{
		{
			name:       "NewValue",
			o:          options{},
			field:      0,
			fieldName:  "Field",
			fieldValue: "255",
			expected:   255,
		},
		{
			name:       "NoNewValue",
			o:          options{},
			field:      255,
			fieldName:  "Field",
			fieldValue: "255",
			expected:   255,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setUint(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetUint16(t *testing.T) {
	tests := []struct {
		name       string
		o          options
		field      uint16
		fieldName  string
		fieldValue string
		expected   uint16
	}{
		{
			name:       "NewValue",
			o:          options{},
			field:      0,
			fieldName:  "Field",
			fieldValue: "65535",
			expected:   65535,
		},
		{
			name:       "NoNewValue",
			o:          options{},
			field:      65535,
			fieldName:  "Field",
			fieldValue: "65535",
			expected:   65535,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setUint(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetUint32(t *testing.T) {
	tests := []struct {
		name       string
		o          options
		field      uint32
		fieldName  string
		fieldValue string
		expected   uint32
	}{
		{
			name:       "NewValue",
			o:          options{},
			field:      0,
			fieldName:  "Field",
			fieldValue: "4294967295",
			expected:   4294967295,
		},
		{
			name:       "NoNewValue",
			o:          options{},
			field:      4294967295,
			fieldName:  "Field",
			fieldValue: "4294967295",
			expected:   4294967295,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setUint(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetUint64(t *testing.T) {
	tests := []struct {
		name       string
		o          options
		field      uint64
		fieldName  string
		fieldValue string
		expected   uint64
	}{
		{
			name:       "NewValue",
			o:          options{},
			field:      0,
			fieldName:  "Field",
			fieldValue: "18446744073709551615",
			expected:   18446744073709551615,
		},
		{
			name:       "NoNewValue",
			o:          options{},
			field:      18446744073709551615,
			fieldName:  "Field",
			fieldValue: "18446744073709551615",
			expected:   18446744073709551615,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setUint(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetURL(t *testing.T) {
	u, _ := url.Parse("example.com")

	tests := []struct {
		name       string
		o          options
		field      url.URL
		fieldName  string
		fieldValue string
		expected   url.URL
	}{
		{
			name:       "NewValue",
			o:          options{},
			field:      url.URL{},
			fieldName:  "Field",
			fieldValue: "example.com",
			expected:   *u,
		},
		{
			name:       "NoNewValue",
			o:          options{},
			field:      *u,
			fieldName:  "Field",
			fieldValue: "example.com",
			expected:   *u,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setStruct(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetStringSlice(t *testing.T) {
	tests := []struct {
		name        string
		o           options
		field       []string
		fieldName   string
		fieldValues []string
		expected    []string
	}{
		{
			name:        "NewValue",
			o:           options{},
			field:       []string{},
			fieldName:   "Field",
			fieldValues: []string{"milad", "mona"},
			expected:    []string{"milad", "mona"},
		},
		{
			name:        "NoNewValue",
			o:           options{},
			field:       []string{"milad", "mona"},
			fieldName:   "Field",
			fieldValues: []string{"milad", "mona"},
			expected:    []string{"milad", "mona"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setStringSlice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetBoolSlice(t *testing.T) {
	tests := []struct {
		name        string
		o           options
		field       []bool
		fieldName   string
		fieldValues []string
		expected    []bool
	}{
		{
			name:        "NewValue",
			o:           options{},
			field:       []bool{},
			fieldName:   "Field",
			fieldValues: []string{"false", "true"},
			expected:    []bool{false, true},
		},
		{
			name:        "NoNewValue",
			o:           options{},
			field:       []bool{false, true},
			fieldName:   "Field",
			fieldValues: []string{"false", "true"},
			expected:    []bool{false, true},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setBoolSlice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetFloat32Slice(t *testing.T) {
	tests := []struct {
		name        string
		o           options
		field       []float32
		fieldName   string
		fieldValues []string
		expected    []float32
	}{
		{
			name:        "NewValue",
			o:           options{},
			field:       []float32{},
			fieldName:   "Field",
			fieldValues: []string{"3.1415", "2.7182"},
			expected:    []float32{3.1415, 2.7182},
		},
		{
			name:        "NoNewValue",
			o:           options{},
			field:       []float32{3.1415, 2.7182},
			fieldName:   "Field",
			fieldValues: []string{"3.1415", "2.7182"},
			expected:    []float32{3.1415, 2.7182},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setFloat32Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetFloat64Slice(t *testing.T) {
	tests := []struct {
		name        string
		o           options
		field       []float64
		fieldName   string
		fieldValues []string
		expected    []float64
	}{
		{
			name:        "NewValue",
			o:           options{},
			field:       []float64{},
			fieldName:   "Field",
			fieldValues: []string{"3.14159265", "2.71828182"},
			expected:    []float64{3.14159265, 2.71828182},
		},
		{
			name:        "NoNewValue",
			o:           options{},
			field:       []float64{3.14159265, 2.71828182},
			fieldName:   "Field",
			fieldValues: []string{"3.14159265", "2.71828182"},
			expected:    []float64{3.14159265, 2.71828182},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setFloat64Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetIntSlice(t *testing.T) {
	tests := []struct {
		name        string
		o           options
		field       []int
		fieldName   string
		fieldValues []string
		expected    []int
	}{
		{
			name:        "NewValue",
			o:           options{},
			field:       []int{},
			fieldName:   "Field",
			fieldValues: []string{"27", "69"},
			expected:    []int{27, 69},
		},
		{
			name:        "NoNewValue",
			o:           options{},
			field:       []int{27, 69},
			fieldName:   "Field",
			fieldValues: []string{"27", "69"},
			expected:    []int{27, 69},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setIntSlice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetInt8Slice(t *testing.T) {
	tests := []struct {
		name        string
		o           options
		field       []int8
		fieldName   string
		fieldValues []string
		expected    []int8
	}{
		{
			name:        "NewValue",
			o:           options{},
			field:       []int8{},
			fieldName:   "Field",
			fieldValues: []string{"-128", "127"},
			expected:    []int8{-128, 127},
		},
		{
			name:        "NoNewValue",
			o:           options{},
			field:       []int8{-128, 127},
			fieldName:   "Field",
			fieldValues: []string{"-128", "127"},
			expected:    []int8{-128, 127},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setInt8Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetInt16Slice(t *testing.T) {
	tests := []struct {
		name        string
		o           options
		field       []int16
		fieldName   string
		fieldValues []string
		expected    []int16
	}{
		{
			name:        "NewValue",
			o:           options{},
			field:       []int16{},
			fieldName:   "Field",
			fieldValues: []string{"-32768", "32767"},
			expected:    []int16{-32768, 32767},
		},
		{
			name:        "NoNewValue",
			o:           options{},
			field:       []int16{-32768, 32767},
			fieldName:   "Field",
			fieldValues: []string{"-32768", "32767"},
			expected:    []int16{-32768, 32767},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setInt16Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetInt32Slice(t *testing.T) {
	tests := []struct {
		name        string
		o           options
		field       []int32
		fieldName   string
		fieldValues []string
		expected    []int32
	}{
		{
			name:        "NewValue",
			o:           options{},
			field:       []int32{},
			fieldName:   "Field",
			fieldValues: []string{"-2147483648", "2147483647"},
			expected:    []int32{-2147483648, 2147483647},
		},
		{
			name:        "NoNewValue",
			o:           options{},
			field:       []int32{-2147483648, 2147483647},
			fieldName:   "Field",
			fieldValues: []string{"-2147483648", "2147483647"},
			expected:    []int32{-2147483648, 2147483647},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setInt32Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetInt64Slice(t *testing.T) {
	tests := []struct {
		name        string
		o           options
		field       []int64
		fieldName   string
		fieldValues []string
		expected    []int64
	}{
		{
			name:        "NewValue",
			o:           options{},
			field:       []int64{},
			fieldName:   "Field",
			fieldValues: []string{"-9223372036854775808", "9223372036854775807"},
			expected:    []int64{-9223372036854775808, 9223372036854775807},
		},
		{
			name:        "NoNewValue",
			o:           options{},
			field:       []int64{-9223372036854775808, 9223372036854775807},
			fieldName:   "Field",
			fieldValues: []string{"-9223372036854775808", "9223372036854775807"},
			expected:    []int64{-9223372036854775808, 9223372036854775807},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setInt64Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetDurationSlice(t *testing.T) {
	tests := []struct {
		name        string
		o           options
		field       []time.Duration
		fieldName   string
		fieldValues []string
		expected    []time.Duration
	}{
		{
			name:        "NewValue",
			o:           options{},
			field:       []time.Duration{},
			fieldName:   "Field",
			fieldValues: []string{"1h0m0s", "1m0s"},
			expected:    []time.Duration{time.Hour, time.Minute},
		},
		{
			name:        "NoNewValue",
			o:           options{},
			field:       []time.Duration{time.Hour, time.Minute},
			fieldName:   "Field",
			fieldValues: []string{"1h0m0s", "1m0s"},
			expected:    []time.Duration{time.Hour, time.Minute},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setInt64Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetUintSlice(t *testing.T) {
	tests := []struct {
		name        string
		o           options
		field       []uint
		fieldName   string
		fieldValues []string
		expected    []uint
	}{
		{
			name:        "NewValue",
			o:           options{},
			field:       []uint{},
			fieldName:   "Field",
			fieldValues: []string{"27", "69"},
			expected:    []uint{27, 69},
		},
		{
			name:        "NoNewValue",
			o:           options{},
			field:       []uint{27, 69},
			fieldName:   "Field",
			fieldValues: []string{"27", "69"},
			expected:    []uint{27, 69},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setUintSlice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetUint8Slice(t *testing.T) {
	tests := []struct {
		name        string
		o           options
		field       []uint8
		fieldName   string
		fieldValues []string
		expected    []uint8
	}{
		{
			name:        "NewValue",
			o:           options{},
			field:       []uint8{},
			fieldName:   "Field",
			fieldValues: []string{"128", "255"},
			expected:    []uint8{128, 255},
		},
		{
			name:        "NoNewValue",
			o:           options{},
			field:       []uint8{128, 255},
			fieldName:   "Field",
			fieldValues: []string{"128", "255"},
			expected:    []uint8{128, 255},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setUint8Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetUint16Slice(t *testing.T) {
	tests := []struct {
		name        string
		o           options
		field       []uint16
		fieldName   string
		fieldValues []string
		expected    []uint16
	}{
		{
			name:        "NewValue",
			o:           options{},
			field:       []uint16{},
			fieldName:   "Field",
			fieldValues: []string{"32768", "65535"},
			expected:    []uint16{32768, 65535},
		},
		{
			name:        "NoNewValue",
			o:           options{},
			field:       []uint16{32768, 65535},
			fieldName:   "Field",
			fieldValues: []string{"32768", "65535"},
			expected:    []uint16{32768, 65535},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setUint16Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetUint32Slice(t *testing.T) {
	tests := []struct {
		name        string
		o           options
		field       []uint32
		fieldName   string
		fieldValues []string
		expected    []uint32
	}{
		{
			name:        "NewValue",
			o:           options{},
			field:       []uint32{},
			fieldName:   "Field",
			fieldValues: []string{"2147483648", "4294967295"},
			expected:    []uint32{2147483648, 4294967295},
		},
		{
			name:        "NoNewValue",
			o:           options{},
			field:       []uint32{2147483648, 4294967295},
			fieldName:   "Field",
			fieldValues: []string{"2147483648", "4294967295"},
			expected:    []uint32{2147483648, 4294967295},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setUint32Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetUint64Slice(t *testing.T) {
	tests := []struct {
		name        string
		o           options
		field       []uint64
		fieldName   string
		fieldValues []string
		expected    []uint64
	}{
		{
			name:        "NewValue",
			o:           options{},
			field:       []uint64{},
			fieldName:   "Field",
			fieldValues: []string{"9223372036854775808", "18446744073709551615"},
			expected:    []uint64{9223372036854775808, 18446744073709551615},
		},
		{
			name:        "NoNewValue",
			o:           options{},
			field:       []uint64{9223372036854775808, 18446744073709551615},
			fieldName:   "Field",
			fieldValues: []string{"9223372036854775808", "18446744073709551615"},
			expected:    []uint64{9223372036854775808, 18446744073709551615},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setUint64Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestSetURLSlice(t *testing.T) {
	u1, _ := url.Parse("localhost")
	u2, _ := url.Parse("example.com")

	tests := []struct {
		name        string
		o           options
		field       []url.URL
		fieldName   string
		fieldValues []string
		expected    []url.URL
	}{
		{
			name:        "NewValue",
			o:           options{},
			field:       []url.URL{},
			fieldName:   "Field",
			fieldValues: []string{"localhost", "example.com"},
			expected:    []url.URL{*u1, *u2},
		},
		{
			name:        "NoNewValue",
			o:           options{},
			field:       []url.URL{*u1, *u2},
			fieldName:   "Field",
			fieldValues: []string{"localhost", "example.com"},
			expected:    []url.URL{*u1, *u2},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			tc.o.setURLSlice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expected, tc.field)
		})
	}
}

func TestPickError(t *testing.T) {
	tests := []struct {
		name          string
		config        interface{}
		opts          []Option
		expectedError string
	}{
		{
			"NonPointer",
			config{},
			nil,
			"a non-pointer type is passed",
		},
		{
			"NonStruct",
			new(string),
			nil,
			"a non-struct type is passed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Pick(tc.config, tc.opts...)
			assert.Equal(t, tc.expectedError, err.Error())
		})
	}
}

func TestPick(t *testing.T) {
	type env struct {
		varName string
		value   string
	}

	type file struct {
		varName string
		value   string
	}

	d90m := 90 * time.Minute
	d120m := 120 * time.Minute
	service1URL, _ := url.Parse("service-1:8080")
	service2URL, _ := url.Parse("service-2:8080")

	tests := []struct {
		name           string
		args           []string
		envs           []env
		files          []file
		config         config
		opts           []Option
		expectedConfig config
	}{
		{
			"Empty",
			[]string{"path/to/binary"},
			[]env{},
			[]file{},
			config{},
			nil,
			config{},
		},
		{
			"AllFromDefaults",
			[]string{"path/to/binary"},
			[]env{},
			[]file{},
			config{
				unexported:         "internal",
				SkipFlag:           "default",
				SkipFlagEnv:        "default",
				SkipFlagEnvFile:    "default",
				FieldString:        "default",
				FieldBool:          false,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			nil,
			config{
				unexported:         "internal",
				SkipFlag:           "default",
				SkipFlagEnv:        "default",
				SkipFlagEnvFile:    "default",
				FieldString:        "default",
				FieldBool:          false,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"AllFromFlags",
			[]string{
				"path/to/binary",
				"-field.string", "content",
				"-field.bool",
				"-field.float32", "3.1415",
				"-field.float64", "3.14159265359",
				"-field.int", "-2147483648",
				"-field.int8", "-128",
				"-field.int16", "-32768",
				"-field.int32", "-2147483648",
				"-field.int64", "-9223372036854775808",
				"-field.uint", "4294967295",
				"-field.uint8", "255",
				"-field.uint16", "65535",
				"-field.uint32", "4294967295",
				"-field.uint64", "18446744073709551615",
				"-field.duration", "90m",
				"-field.url", "service-1:8080",
				"-field.string.array", "milad,mona",
				"-field.bool.array", "false,true",
				"-field.float32.array", "3.1415,2.7182",
				"-field.float64.array", "3.14159265359,2.71828182845",
				"-field.int.array", "-2147483648,2147483647",
				"-field.int8.array", "-128,127",
				"-field.int16.array", "-32768,32767",
				"-field.int32.array", "-2147483648,2147483647",
				"-field.int64.array", "-9223372036854775808,9223372036854775807",
				"-field.uint.array", "0,4294967295",
				"-field.uint8.array", "0,255",
				"-field.uint16.array", "0,65535",
				"-field.uint32.array", "0,4294967295",
				"-field.uint64.array", "0,18446744073709551615",
				"-field.duration.array", "90m,120m",
				"-field.url.array", "service-1:8080,service-2:8080",
			},
			[]env{},
			[]file{},
			config{},
			nil,
			config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"AllFromFlags",
			[]string{
				"path/to/binary",
				"--field.string", "content",
				"--field.bool",
				"--field.float32", "3.1415",
				"--field.float64", "3.14159265359",
				"--field.int", "-2147483648",
				"--field.int8", "-128",
				"--field.int16", "-32768",
				"--field.int32", "-2147483648",
				"--field.int64", "-9223372036854775808",
				"--field.uint", "4294967295",
				"--field.uint8", "255",
				"--field.uint16", "65535",
				"--field.uint32", "4294967295",
				"--field.uint64", "18446744073709551615",
				"--field.duration", "90m",
				"--field.url", "service-1:8080",
				"--field.string.array", "milad,mona",
				"--field.bool.array", "false,true",
				"--field.float32.array", "3.1415,2.7182",
				"--field.float64.array", "3.14159265359,2.71828182845",
				"--field.int.array", "-2147483648,2147483647",
				"--field.int8.array", "-128,127",
				"--field.int16.array", "-32768,32767",
				"--field.int32.array", "-2147483648,2147483647",
				"--field.int64.array", "-9223372036854775808,9223372036854775807",
				"--field.uint.array", "0,4294967295",
				"--field.uint8.array", "0,255",
				"--field.uint16.array", "0,65535",
				"--field.uint32.array", "0,4294967295",
				"--field.uint64.array", "0,18446744073709551615",
				"--field.duration.array", "90m,120m",
				"--field.url.array", "service-1:8080,service-2:8080",
			},
			[]env{},
			[]file{},
			config{},
			nil,
			config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"AllFromFlags",
			[]string{
				"path/to/binary",
				"-field.string=content",
				"-field.bool",
				"-field.float32=3.1415",
				"-field.float64=3.14159265359",
				"-field.int=-2147483648",
				"-field.int8=-128",
				"-field.int16=-32768",
				"-field.int32=-2147483648",
				"-field.int64=-9223372036854775808",
				"-field.uint=4294967295",
				"-field.uint8=255",
				"-field.uint16=65535",
				"-field.uint32=4294967295",
				"-field.uint64=18446744073709551615",
				"-field.duration=90m",
				"-field.url=service-1:8080",
				"-field.string.array=milad,mona",
				"-field.bool.array=false,true",
				"-field.float32.array=3.1415,2.7182",
				"-field.float64.array=3.14159265359,2.71828182845",
				"-field.int.array=-2147483648,2147483647",
				"-field.int8.array=-128,127",
				"-field.int16.array=-32768,32767",
				"-field.int32.array=-2147483648,2147483647",
				"-field.int64.array=-9223372036854775808,9223372036854775807",
				"-field.uint.array=0,4294967295",
				"-field.uint8.array=0,255",
				"-field.uint16.array=0,65535",
				"-field.uint32.array=0,4294967295",
				"-field.uint64.array=0,18446744073709551615",
				"-field.duration.array=90m,120m",
				"-field.url.array=service-1:8080,service-2:8080",
			},
			[]env{},
			[]file{},
			config{},
			nil,
			config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"AllFromFlags",
			[]string{
				"path/to/binary",
				"--field.string=content",
				"--field.bool",
				"--field.float32=3.1415",
				"--field.float64=3.14159265359",
				"--field.int=-2147483648",
				"--field.int8=-128",
				"--field.int16=-32768",
				"--field.int32=-2147483648",
				"--field.int64=-9223372036854775808",
				"--field.uint=4294967295",
				"--field.uint8=255",
				"--field.uint16=65535",
				"--field.uint32=4294967295",
				"--field.uint64=18446744073709551615",
				"--field.duration=90m",
				"--field.url=service-1:8080",
				"--field.string.array=milad,mona",
				"--field.bool.array=false,true",
				"--field.float32.array=3.1415,2.7182",
				"--field.float64.array=3.14159265359,2.71828182845",
				"--field.int.array=-2147483648,2147483647",
				"--field.int8.array=-128,127",
				"--field.int16.array=-32768,32767",
				"--field.int32.array=-2147483648,2147483647",
				"--field.int64.array=-9223372036854775808,9223372036854775807",
				"--field.uint.array=0,4294967295",
				"--field.uint8.array=0,255",
				"--field.uint16.array=0,65535",
				"--field.uint32.array=0,4294967295",
				"--field.uint64.array=0,18446744073709551615",
				"--field.duration.array=90m,120m",
				"--field.url.array=service-1:8080,service-2:8080",
			},
			[]env{},
			[]file{},
			config{},
			nil,
			config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"AllFromEnvironmentVariables",
			[]string{"path/to/binary"},
			[]env{
				{"SKIP_FLAG", "from_env"},
				{"SKIP_FLAG_ENV", "from_env"},
				{"SKIP_FLAG_ENV_FILE", "from_env"},
				{"FIELD_STRING", "content"},
				{"FIELD_BOOL", "true"},
				{"FIELD_FLOAT32", "3.1415"},
				{"FIELD_FLOAT64", "3.14159265359"},
				{"FIELD_INT", "-2147483648"},
				{"FIELD_INT8", "-128"},
				{"FIELD_INT16", "-32768"},
				{"FIELD_INT32", "-2147483648"},
				{"FIELD_INT64", "-9223372036854775808"},
				{"FIELD_UINT", "4294967295"},
				{"FIELD_UINT8", "255"},
				{"FIELD_UINT16", "65535"},
				{"FIELD_UINT32", "4294967295"},
				{"FIELD_UINT64", "18446744073709551615"},
				{"FIELD_DURATION", "90m"},
				{"FIELD_URL", "service-1:8080"},
				{"FIELD_STRING_ARRAY", "milad,mona"},
				{"FIELD_BOOL_ARRAY", "false,true"},
				{"FIELD_FLOAT32_ARRAY", "3.1415,2.7182"},
				{"FIELD_FLOAT64_ARRAY", "3.14159265359,2.71828182845"},
				{"FIELD_INT_ARRAY", "-2147483648,2147483647"},
				{"FIELD_INT8_ARRAY", "-128,127"},
				{"FIELD_INT16_ARRAY", "-32768,32767"},
				{"FIELD_INT32_ARRAY", "-2147483648,2147483647"},
				{"FIELD_INT64_ARRAY", "-9223372036854775808,9223372036854775807"},
				{"FIELD_UINT_ARRAY", "0,4294967295"},
				{"FIELD_UINT8_ARRAY", "0,255"},
				{"FIELD_UINT16_ARRAY", "0,65535"},
				{"FIELD_UINT32_ARRAY", "0,4294967295"},
				{"FIELD_UINT64_ARRAY", "0,18446744073709551615"},
				{"FIELD_DURATION_ARRAY", "90m,120m"},
				{"FIELD_URL_ARRAY", "service-1:8080,service-2:8080"},
			},
			[]file{},
			config{},
			nil,
			config{
				unexported:         "",
				SkipFlag:           "from_env",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"AllFromFromFiles",
			[]string{"path/to/binary"},
			[]env{},
			[]file{
				{"SKIP_FLAG_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE_FILE", "from_file"},
				{"FIELD_STRING_FILE", "content"},
				{"FIELD_BOOL_FILE", "true"},
				{"FIELD_FLOAT32_FILE", "3.1415"},
				{"FIELD_FLOAT64_FILE", "3.14159265359"},
				{"FIELD_INT_FILE", "-2147483648"},
				{"FIELD_INT8_FILE", "-128"},
				{"FIELD_INT16_FILE", "-32768"},
				{"FIELD_INT32_FILE", "-2147483648"},
				{"FIELD_INT64_FILE", "-9223372036854775808"},
				{"FIELD_UINT_FILE", "4294967295"},
				{"FIELD_UINT8_FILE", "255"},
				{"FIELD_UINT16_FILE", "65535"},
				{"FIELD_UINT32_FILE", "4294967295"},
				{"FIELD_UINT64_FILE", "18446744073709551615"},
				{"FIELD_DURATION_FILE", "90m"},
				{"FIELD_URL_FILE", "service-1:8080"},
				{"FIELD_STRING_ARRAY_FILE", "milad,mona"},
				{"FIELD_BOOL_ARRAY_FILE", "false,true"},
				{"FIELD_FLOAT32_ARRAY_FILE", "3.1415,2.7182"},
				{"FIELD_FLOAT64_ARRAY_FILE", "3.14159265359,2.71828182845"},
				{"FIELD_INT_ARRAY_FILE", "-2147483648,2147483647"},
				{"FIELD_INT8_ARRAY_FILE", "-128,127"},
				{"FIELD_INT16_ARRAY_FILE", "-32768,32767"},
				{"FIELD_INT32_ARRAY_FILE", "-2147483648,2147483647"},
				{"FIELD_INT64_ARRAY_FILE", "-9223372036854775808,9223372036854775807"},
				{"FIELD_UINT_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT8_ARRAY_FILE", "0,255"},
				{"FIELD_UINT16_ARRAY_FILE", "0,65535"},
				{"FIELD_UINT32_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT64_ARRAY_FILE", "0,18446744073709551615"},
				{"FIELD_DURATION_ARRAY_FILE", "90m,120m"},
				{"FIELD_URL_ARRAY_FILE", "service-1:8080,service-2:8080"},
			},
			config{},
			nil,
			config{
				unexported:         "",
				SkipFlag:           "from_file",
				SkipFlagEnv:        "from_file",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"FromMixedSources",
			[]string{
				"path/to/binary",
				"-field.float32=3.1415",
				"--field.float64=3.14159265359",
				"-field.float32.array", "3.1415,2.7182",
				"--field.float64.array", "3.14159265359,2.71828182845",
			},
			[]env{
				{"SKIP_FLAG", "from_env"},
				{"SKIP_FLAG_ENV", "from_env"},
				{"SKIP_FLAG_ENV_FILE", "from_env"},
				{"FIELD_INT", "-2147483648"},
				{"FIELD_INT8", "-128"},
				{"FIELD_INT16", "-32768"},
				{"FIELD_INT32", "-2147483648"},
				{"FIELD_INT64", "-9223372036854775808"},
				{"FIELD_INT_ARRAY", "-2147483648,2147483647"},
				{"FIELD_INT8_ARRAY", "-128,127"},
				{"FIELD_INT16_ARRAY", "-32768,32767"},
				{"FIELD_INT32_ARRAY", "-2147483648,2147483647"},
				{"FIELD_INT64_ARRAY", "-9223372036854775808,9223372036854775807"},
			},
			[]file{
				{"SKIP_FLAG_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE_FILE", "from_file"},
				{"FIELD_UINT_FILE", "4294967295"},
				{"FIELD_UINT8_FILE", "255"},
				{"FIELD_UINT16_FILE", "65535"},
				{"FIELD_UINT32_FILE", "4294967295"},
				{"FIELD_UINT64_FILE", "18446744073709551615"},
				{"FIELD_UINT_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT8_ARRAY_FILE", "0,255"},
				{"FIELD_UINT16_ARRAY_FILE", "0,65535"},
				{"FIELD_UINT32_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT64_ARRAY_FILE", "0,18446744073709551615"},
			},
			config{
				FieldString:        "default",
				FieldStringArray:   []string{"milad", "mona"},
				FieldBool:          true,
				FieldBoolArray:     []bool{false, true},
				FieldDuration:      d90m,
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURL:           *service1URL,
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			nil,
			config{
				unexported:         "",
				SkipFlag:           "from_env",
				SkipFlagEnv:        "from_file",
				SkipFlagEnvFile:    "",
				FieldString:        "default",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"WithTelepresenceOption",
			[]string{"path/to/binary"},
			[]env{},
			[]file{
				{"SKIP_FLAG_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE_FILE", "from_file"},
				{"FIELD_STRING_FILE", "content"},
				{"FIELD_BOOL_FILE", "true"},
				{"FIELD_FLOAT32_FILE", "3.1415"},
				{"FIELD_FLOAT64_FILE", "3.14159265359"},
				{"FIELD_INT_FILE", "-2147483648"},
				{"FIELD_INT8_FILE", "-128"},
				{"FIELD_INT16_FILE", "-32768"},
				{"FIELD_INT32_FILE", "-2147483648"},
				{"FIELD_INT64_FILE", "-9223372036854775808"},
				{"FIELD_UINT_FILE", "4294967295"},
				{"FIELD_UINT8_FILE", "255"},
				{"FIELD_UINT16_FILE", "65535"},
				{"FIELD_UINT32_FILE", "4294967295"},
				{"FIELD_UINT64_FILE", "18446744073709551615"},
				{"FIELD_DURATION_FILE", "90m"},
				{"FIELD_URL_FILE", "service-1:8080"},
				{"FIELD_STRING_ARRAY_FILE", "milad,mona"},
				{"FIELD_BOOL_ARRAY_FILE", "false,true"},
				{"FIELD_FLOAT32_ARRAY_FILE", "3.1415,2.7182"},
				{"FIELD_FLOAT64_ARRAY_FILE", "3.14159265359,2.71828182845"},
				{"FIELD_INT_ARRAY_FILE", "-2147483648,2147483647"},
				{"FIELD_INT8_ARRAY_FILE", "-128,127"},
				{"FIELD_INT16_ARRAY_FILE", "-32768,32767"},
				{"FIELD_INT32_ARRAY_FILE", "-2147483648,2147483647"},
				{"FIELD_INT64_ARRAY_FILE", "-9223372036854775808,9223372036854775807"},
				{"FIELD_UINT_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT8_ARRAY_FILE", "0,255"},
				{"FIELD_UINT16_ARRAY_FILE", "0,65535"},
				{"FIELD_UINT32_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT64_ARRAY_FILE", "0,18446744073709551615"},
				{"FIELD_DURATION_ARRAY_FILE", "90m,120m"},
				{"FIELD_URL_ARRAY_FILE", "service-1:8080,service-2:8080"},
			},
			config{},
			[]Option{
				Telepresence(),
			},
			config{
				unexported:         "",
				SkipFlag:           "from_file",
				SkipFlagEnv:        "from_file",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
	}

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			os.Args = tc.args

			// Set environment variables
			for _, e := range tc.envs {
				err := os.Setenv(e.varName, e.value)
				assert.NoError(t, err)
				defer os.Unsetenv(e.varName)
			}

			o := defaultOptions()
			for _, opt := range tc.opts {
				opt.apply(o)
			}

			// Testing Telepresence option
			if o.telepresence {
				err := os.Setenv(telepresenceEnvVar, "/")
				assert.NoError(t, err)
				defer os.Unsetenv(telepresenceEnvVar)
			}

			// Write configuration files
			for _, f := range tc.files {
				tmpfile, err := ioutil.TempFile("", "gotest_")
				assert.NoError(t, err)
				defer os.Remove(tmpfile.Name())

				_, err = tmpfile.WriteString(f.value)
				assert.NoError(t, err)

				err = tmpfile.Close()
				assert.NoError(t, err)

				err = os.Setenv(f.varName, tmpfile.Name())
				assert.NoError(t, err)
				defer os.Unsetenv(f.varName)
			}

			err := Pick(&tc.config, tc.opts...)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedConfig, tc.config)
		})
	}

	// flag.Parse() can be called only once
	flag.Parse()
}
