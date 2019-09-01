package konfig

import (
	"errors"
	"flag"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type config struct {
	sync.Mutex
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

func TestDebug(t *testing.T) {
	tests := []struct {
		c         *controller
		verbosity uint
		expected  *controller
	}{
		{
			&controller{},
			2,
			&controller{
				debug: 2,
			},
		},
	}

	for _, tc := range tests {
		opt := Debug(tc.verbosity)
		opt(tc.c)

		assert.Equal(t, tc.expected, tc.c)
	}
}

func TestTelepresence(t *testing.T) {
	tests := []struct {
		c        *controller
		expected *controller
	}{
		{
			&controller{},
			&controller{
				telepresence: true,
			},
		},
	}

	for _, tc := range tests {
		opt := Telepresence()
		opt(tc.c)

		assert.Equal(t, tc.expected, tc.c)
	}
}

func TestWatchInterval(t *testing.T) {
	tests := []struct {
		c        *controller
		d        time.Duration
		expected *controller
	}{
		{
			&controller{},
			10 * time.Second,
			&controller{
				watchInterval: 10 * time.Second,
			},
		},
	}

	for _, tc := range tests {
		opt := WatchInterval(tc.d)
		opt(tc.c)

		assert.Equal(t, tc.expected, tc.c)
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		expectedString string
	}{
		{
			"NoOption",
			&controller{},
			"",
		},
		{
			"WithDebug",
			&controller{
				debug: 2,
			},
			"Debug<2>",
		},
		{
			"WithTelepresence",
			&controller{
				telepresence: true,
			},
			"Telepresence",
		},
		{
			"WithWatchInterval",
			&controller{
				watchInterval: 5 * time.Second,
			},
			"Watch<5s>",
		},
		{
			"WithSubscribers",
			&controller{
				subscribers: []chan Update{
					make(chan Update),
					make(chan Update),
				},
			},
			"Subscribers<2>",
		},
		{
			"WithAll",
			&controller{
				debug:         2,
				telepresence:  true,
				watchInterval: 5 * time.Second,
				subscribers: []chan Update{
					make(chan Update),
					make(chan Update),
				},
			},
			"Debug<2> + Telepresence + Watch<5s> + Subscribers<2>",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			str := tc.c.String()

			assert.Equal(t, tc.expectedString, str)
		})
	}
}

func TestLog(t *testing.T) {
	tests := []struct {
		name string
		c    *controller
		v    uint
		msg  string
		args []interface{}
	}{
		{
			"WithoutDebug",
			&controller{},
			1,
			"testing ...",
			nil,
		},
		{
			"WithDebug",
			&controller{
				debug: 2,
			},
			2,
			"testing ...",
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.c.log(tc.v, tc.msg, tc.args...)
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
		name                                      string
		args                                      []string
		envConfig                                 env
		fileConfig                                file
		fieldName, flagName, envName, fileEnvName string
		c                                         *controller
		expectedValue                             string
		expectedFromFile                          bool
	}{
		{
			"SkipFlag",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "-", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{},
			"info",
			false,
		},
		{
			"SkipFlagAndEnv",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "-", "-", "LOG_LEVEL_FILE",
			&controller{},
			"error",
			true,
		},
		{
			"SkipFlagAndEnvAndFile",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "-", "-", "-",
			&controller{},
			"",
			false,
		},
		{
			"FromFlag",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{},
			"debug",
			false,
		},
		{
			"FromFlag",
			[]string{"/path/to/executable", "--log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{},
			"debug",
			false,
		},
		{
			"FromFlag",
			[]string{"/path/to/executable", "-log.level", "debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{},
			"debug",
			false,
		},
		{
			"FromFlag",
			[]string{"/path/to/executable", "--log.level", "debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{},
			"debug",
			false,
		},
		{
			"FromEnvironmentVariable",
			[]string{"/path/to/executable"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{},
			"info",
			false,
		},
		{
			"FromFiles",
			[]string{"/path/to/executable"},
			env{"LOG_LEVEL", ""},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{},
			"error",
			true,
		},
		{
			"FromFilesWithTelepresenceOption",
			[]string{"/path/to/executable"},
			env{"LOG_LEVEL", ""},
			file{"LOG_LEVEL_FILE", "info"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{telepresence: true},
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
			if tc.c.telepresence {
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

			value, fromFile := tc.c.getFieldValue(tc.fieldName, tc.flagName, tc.envName, tc.fileEnvName)
			assert.Equal(t, tc.expectedValue, value)
			assert.Equal(t, tc.expectedFromFile, fromFile)
		})
	}
}

func TestNotifySubscribers(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		fieldName      string
		fieldValue     interface{}
		expectedUpdate Update
	}{
		{
			"Nil",
			&controller{},
			"FieldBool", true,
			Update{},
		},
		{
			"NoChannel",
			&controller{
				subscribers: []chan Update{},
			},
			"FieldString", "value",
			Update{},
		},
		{
			"WithBlockingChannels",
			&controller{
				subscribers: []chan Update{
					make(chan Update),
					make(chan Update),
				},
			},
			"FieldInt", 27,
			Update{"FieldInt", 27},
		},
		{
			"WithBufferedChannels",
			&controller{
				subscribers: []chan Update{
					make(chan Update, 1),
					make(chan Update, 1),
				},
			},
			"FieldFloat", 3.1415,
			Update{"FieldFloat", 3.1415},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.c.notifySubscribers(tc.fieldName, tc.fieldValue)

			if tc.expectedUpdate != (Update{}) {
				for _, ch := range tc.c.subscribers {
					update := <-ch
					assert.Equal(t, tc.expectedUpdate, update)
				}
			}
		})
	}
}

func TestSetString(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          string
		fieldName      string
		fieldValue     string
		expectedValue  string
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          "",
			fieldName:      "Field",
			fieldValue:     "milad",
			expectedValue:  "milad",
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          "milad",
			fieldName:      "Field",
			fieldValue:     "milad",
			expectedValue:  "milad",
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setString(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetBool(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          bool
		fieldName      string
		fieldValue     string
		expectedValue  bool
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          false,
			fieldName:      "Field",
			fieldValue:     "true",
			expectedValue:  true,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          true,
			fieldName:      "Field",
			fieldValue:     "true",
			expectedValue:  true,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setBool(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetFloat32(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          float32
		fieldName      string
		fieldValue     string
		expectedValue  float32
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "3.1415",
			expectedValue:  3.1415,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          3.1415,
			fieldName:      "Field",
			fieldValue:     "3.1415",
			expectedValue:  3.1415,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setFloat32(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetFloat64(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          float64
		fieldName      string
		fieldValue     string
		expectedValue  float64
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "3.14159265359",
			expectedValue:  3.14159265359,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          3.14159265359,
			fieldName:      "Field",
			fieldValue:     "3.14159265359",
			expectedValue:  3.14159265359,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setFloat64(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetInt(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          int
		fieldName      string
		fieldValue     string
		expectedValue  int
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "-27",
			expectedValue:  -27,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          -27,
			fieldName:      "Field",
			fieldValue:     "-27",
			expectedValue:  -27,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetInt8(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          int8
		fieldName      string
		fieldValue     string
		expectedValue  int8
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "-128",
			expectedValue:  -128,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          -128,
			fieldName:      "Field",
			fieldValue:     "-128",
			expectedValue:  -128,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetInt16(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          int16
		fieldName      string
		fieldValue     string
		expectedValue  int16
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "-32768",
			expectedValue:  -32768,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          -32768,
			fieldName:      "Field",
			fieldValue:     "-32768",
			expectedValue:  -32768,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetInt32(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          int32
		fieldName      string
		fieldValue     string
		expectedValue  int32
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "-2147483648",
			expectedValue:  -2147483648,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          -2147483648,
			fieldName:      "Field",
			fieldValue:     "-2147483648",
			expectedValue:  -2147483648,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetInt64(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          int64
		fieldName      string
		fieldValue     string
		expectedValue  int64
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "-9223372036854775808",
			expectedValue:  -9223372036854775808,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          -9223372036854775808,
			fieldName:      "Field",
			fieldValue:     "-9223372036854775808",
			expectedValue:  -9223372036854775808,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetDuration(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          time.Duration
		fieldName      string
		fieldValue     string
		expectedValue  time.Duration
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "1h0m0s",
			expectedValue:  time.Hour,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          time.Hour,
			fieldName:      "Field",
			fieldValue:     "1h0m0s",
			expectedValue:  time.Hour,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUint(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          uint
		fieldName      string
		fieldValue     string
		expectedValue  uint
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "27",
			expectedValue:  27,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          27,
			fieldName:      "Field",
			fieldValue:     "27",
			expectedValue:  27,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUint(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUint8(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          uint8
		fieldName      string
		fieldValue     string
		expectedValue  uint8
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "255",
			expectedValue:  255,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          255,
			fieldName:      "Field",
			fieldValue:     "255",
			expectedValue:  255,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUint(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUint16(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          uint16
		fieldName      string
		fieldValue     string
		expectedValue  uint16
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "65535",
			expectedValue:  65535,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          65535,
			fieldName:      "Field",
			fieldValue:     "65535",
			expectedValue:  65535,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUint(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUint32(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          uint32
		fieldName      string
		fieldValue     string
		expectedValue  uint32
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "4294967295",
			expectedValue:  4294967295,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          4294967295,
			fieldName:      "Field",
			fieldValue:     "4294967295",
			expectedValue:  4294967295,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUint(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUint64(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          uint64
		fieldName      string
		fieldValue     string
		expectedValue  uint64
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "18446744073709551615",
			expectedValue:  18446744073709551615,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          18446744073709551615,
			fieldName:      "Field",
			fieldValue:     "18446744073709551615",
			expectedValue:  18446744073709551615,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUint(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetURL(t *testing.T) {
	u, _ := url.Parse("example.com")

	tests := []struct {
		name           string
		c              *controller
		field          url.URL
		fieldName      string
		fieldValue     string
		expectedValue  url.URL
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          url.URL{},
			fieldName:      "Field",
			fieldValue:     "example.com",
			expectedValue:  *u,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          *u,
			fieldName:      "Field",
			fieldValue:     "example.com",
			expectedValue:  *u,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setStruct(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetStringSlice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []string
		fieldName      string
		fieldValues    []string
		expectedValues []string
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []string{},
			fieldName:      "Field",
			fieldValues:    []string{"milad", "mona"},
			expectedValues: []string{"milad", "mona"},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []string{"milad", "mona"},
			fieldName:      "Field",
			fieldValues:    []string{"milad", "mona"},
			expectedValues: []string{"milad", "mona"},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setStringSlice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetBoolSlice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []bool
		fieldName      string
		fieldValues    []string
		expectedValues []bool
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []bool{},
			fieldName:      "Field",
			fieldValues:    []string{"false", "true"},
			expectedValues: []bool{false, true},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []bool{false, true},
			fieldName:      "Field",
			fieldValues:    []string{"false", "true"},
			expectedValues: []bool{false, true},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setBoolSlice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetFloat32Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []float32
		fieldName      string
		fieldValues    []string
		expectedValues []float32
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []float32{},
			fieldName:      "Field",
			fieldValues:    []string{"3.1415", "2.7182"},
			expectedValues: []float32{3.1415, 2.7182},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []float32{3.1415, 2.7182},
			fieldName:      "Field",
			fieldValues:    []string{"3.1415", "2.7182"},
			expectedValues: []float32{3.1415, 2.7182},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setFloat32Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetFloat64Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []float64
		fieldName      string
		fieldValues    []string
		expectedValues []float64
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []float64{},
			fieldName:      "Field",
			fieldValues:    []string{"3.14159265", "2.71828182"},
			expectedValues: []float64{3.14159265, 2.71828182},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []float64{3.14159265, 2.71828182},
			fieldName:      "Field",
			fieldValues:    []string{"3.14159265", "2.71828182"},
			expectedValues: []float64{3.14159265, 2.71828182},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setFloat64Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetIntSlice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []int
		fieldName      string
		fieldValues    []string
		expectedValues []int
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []int{},
			fieldName:      "Field",
			fieldValues:    []string{"27", "69"},
			expectedValues: []int{27, 69},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []int{27, 69},
			fieldName:      "Field",
			fieldValues:    []string{"27", "69"},
			expectedValues: []int{27, 69},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setIntSlice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetInt8Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []int8
		fieldName      string
		fieldValues    []string
		expectedValues []int8
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []int8{},
			fieldName:      "Field",
			fieldValues:    []string{"-128", "127"},
			expectedValues: []int8{-128, 127},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []int8{-128, 127},
			fieldName:      "Field",
			fieldValues:    []string{"-128", "127"},
			expectedValues: []int8{-128, 127},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt8Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetInt16Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []int16
		fieldName      string
		fieldValues    []string
		expectedValues []int16
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []int16{},
			fieldName:      "Field",
			fieldValues:    []string{"-32768", "32767"},
			expectedValues: []int16{-32768, 32767},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []int16{-32768, 32767},
			fieldName:      "Field",
			fieldValues:    []string{"-32768", "32767"},
			expectedValues: []int16{-32768, 32767},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt16Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetInt32Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []int32
		fieldName      string
		fieldValues    []string
		expectedValues []int32
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []int32{},
			fieldName:      "Field",
			fieldValues:    []string{"-2147483648", "2147483647"},
			expectedValues: []int32{-2147483648, 2147483647},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []int32{-2147483648, 2147483647},
			fieldName:      "Field",
			fieldValues:    []string{"-2147483648", "2147483647"},
			expectedValues: []int32{-2147483648, 2147483647},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt32Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetInt64Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []int64
		fieldName      string
		fieldValues    []string
		expectedValues []int64
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []int64{},
			fieldName:      "Field",
			fieldValues:    []string{"-9223372036854775808", "9223372036854775807"},
			expectedValues: []int64{-9223372036854775808, 9223372036854775807},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []int64{-9223372036854775808, 9223372036854775807},
			fieldName:      "Field",
			fieldValues:    []string{"-9223372036854775808", "9223372036854775807"},
			expectedValues: []int64{-9223372036854775808, 9223372036854775807},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt64Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetDurationSlice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []time.Duration
		fieldName      string
		fieldValues    []string
		expectedValues []time.Duration
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []time.Duration{},
			fieldName:      "Field",
			fieldValues:    []string{"1h0m0s", "1m0s"},
			expectedValues: []time.Duration{time.Hour, time.Minute},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []time.Duration{time.Hour, time.Minute},
			fieldName:      "Field",
			fieldValues:    []string{"1h0m0s", "1m0s"},
			expectedValues: []time.Duration{time.Hour, time.Minute},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt64Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUintSlice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []uint
		fieldName      string
		fieldValues    []string
		expectedValues []uint
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []uint{},
			fieldName:      "Field",
			fieldValues:    []string{"27", "69"},
			expectedValues: []uint{27, 69},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []uint{27, 69},
			fieldName:      "Field",
			fieldValues:    []string{"27", "69"},
			expectedValues: []uint{27, 69},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUintSlice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUint8Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []uint8
		fieldName      string
		fieldValues    []string
		expectedValues []uint8
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []uint8{},
			fieldName:      "Field",
			fieldValues:    []string{"128", "255"},
			expectedValues: []uint8{128, 255},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []uint8{128, 255},
			fieldName:      "Field",
			fieldValues:    []string{"128", "255"},
			expectedValues: []uint8{128, 255},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUint8Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUint16Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []uint16
		fieldName      string
		fieldValues    []string
		expectedValues []uint16
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []uint16{},
			fieldName:      "Field",
			fieldValues:    []string{"32768", "65535"},
			expectedValues: []uint16{32768, 65535},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []uint16{32768, 65535},
			fieldName:      "Field",
			fieldValues:    []string{"32768", "65535"},
			expectedValues: []uint16{32768, 65535},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUint16Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUint32Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []uint32
		fieldName      string
		fieldValues    []string
		expectedValues []uint32
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []uint32{},
			fieldName:      "Field",
			fieldValues:    []string{"2147483648", "4294967295"},
			expectedValues: []uint32{2147483648, 4294967295},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []uint32{2147483648, 4294967295},
			fieldName:      "Field",
			fieldValues:    []string{"2147483648", "4294967295"},
			expectedValues: []uint32{2147483648, 4294967295},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUint32Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUint64Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []uint64
		fieldName      string
		fieldValues    []string
		expectedValues []uint64
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []uint64{},
			fieldName:      "Field",
			fieldValues:    []string{"9223372036854775808", "18446744073709551615"},
			expectedValues: []uint64{9223372036854775808, 18446744073709551615},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []uint64{9223372036854775808, 18446744073709551615},
			fieldName:      "Field",
			fieldValues:    []string{"9223372036854775808", "18446744073709551615"},
			expectedValues: []uint64{9223372036854775808, 18446744073709551615},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUint64Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetURLSlice(t *testing.T) {
	u1, _ := url.Parse("localhost")
	u2, _ := url.Parse("example.com")

	tests := []struct {
		name           string
		c              *controller
		field          []url.URL
		fieldName      string
		fieldValues    []string
		expectedValues []url.URL
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []url.URL{},
			fieldName:      "Field",
			fieldValues:    []string{"localhost", "example.com"},
			expectedValues: []url.URL{*u1, *u2},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []url.URL{*u1, *u2},
			fieldName:      "Field",
			fieldValues:    []string{"localhost", "example.com"},
			expectedValues: []url.URL{*u1, *u2},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setURLSlice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestIterateOnFields(t *testing.T) {
	tests := []struct {
		name                 string
		c                    *controller
		config               interface{}
		expectedValues       []reflect.Value
		expectedFieldNames   []string
		expectedFlagNames    []string
		expectedEnvNames     []string
		expectedFileEnvNames []string
		expectedListSeps     []string
		expectedError        error
	}{
		{
			name:           "OK",
			c:              &controller{},
			config:         &config{},
			expectedValues: []reflect.Value{},
			expectedFieldNames: []string{
				"SkipFlag", "SkipFlagEnv", "SkipFlagEnvFile",
				"FieldString",
				"FieldBool",
				"FieldFloat32", "FieldFloat64",
				"FieldInt", "FieldInt8", "FieldInt16", "FieldInt32", "FieldInt64",
				"FieldUint", "FieldUint8", "FieldUint16", "FieldUint32", "FieldUint64",
				"FieldDuration", "FieldURL",
				"FieldStringArray",
				"FieldBoolArray",
				"FieldFloat32Array", "FieldFloat64Array",
				"FieldIntArray", "FieldInt8Array", "FieldInt16Array", "FieldInt32Array", "FieldInt64Array",
				"FieldUintArray", "FieldUint8Array", "FieldUint16Array", "FieldUint32Array", "FieldUint64Array",
				"FieldDurationArray", "FieldURLArray",
			},
			expectedFlagNames: []string{
				"-", "-", "-",
				"field.string",
				"field.bool",
				"field.float32", "field.float64",
				"field.int", "field.int8", "field.int16", "field.int32", "field.int64",
				"field.uint", "field.uint8", "field.uint16", "field.uint32", "field.uint64",
				"field.duration", "field.url",
				"field.string.array",
				"field.bool.array",
				"field.float32.array", "field.float64.array",
				"field.int.array", "field.int8.array", "field.int16.array", "field.int32.array", "field.int64.array",
				"field.uint.array", "field.uint8.array", "field.uint16.array", "field.uint32.array", "field.uint64.array",
				"field.duration.array", "field.url.array",
			},
			expectedEnvNames: []string{
				"SKIP_FLAG", "-", "-",
				"FIELD_STRING",
				"FIELD_BOOL",
				"FIELD_FLOAT32", "FIELD_FLOAT64",
				"FIELD_INT", "FIELD_INT8", "FIELD_INT16", "FIELD_INT32", "FIELD_INT64",
				"FIELD_UINT", "FIELD_UINT8", "FIELD_UINT16", "FIELD_UINT32", "FIELD_UINT64",
				"FIELD_DURATION", "FIELD_URL",
				"FIELD_STRING_ARRAY",
				"FIELD_BOOL_ARRAY",
				"FIELD_FLOAT32_ARRAY", "FIELD_FLOAT64_ARRAY",
				"FIELD_INT_ARRAY", "FIELD_INT8_ARRAY", "FIELD_INT16_ARRAY", "FIELD_INT32_ARRAY", "FIELD_INT64_ARRAY",
				"FIELD_UINT_ARRAY", "FIELD_UINT8_ARRAY", "FIELD_UINT16_ARRAY", "FIELD_UINT32_ARRAY", "FIELD_UINT64_ARRAY",
				"FIELD_DURATION_ARRAY", "FIELD_URL_ARRAY",
			},
			expectedFileEnvNames: []string{
				"SKIP_FLAG_FILE", "SKIP_FLAG_ENV_FILE", "-",
				"FIELD_STRING_FILE",
				"FIELD_BOOL_FILE",
				"FIELD_FLOAT32_FILE", "FIELD_FLOAT64_FILE",
				"FIELD_INT_FILE", "FIELD_INT8_FILE", "FIELD_INT16_FILE", "FIELD_INT32_FILE", "FIELD_INT64_FILE",
				"FIELD_UINT_FILE", "FIELD_UINT8_FILE", "FIELD_UINT16_FILE", "FIELD_UINT32_FILE", "FIELD_UINT64_FILE",
				"FIELD_DURATION_FILE", "FIELD_URL_FILE",
				"FIELD_STRING_ARRAY_FILE",
				"FIELD_BOOL_ARRAY_FILE",
				"FIELD_FLOAT32_ARRAY_FILE", "FIELD_FLOAT64_ARRAY_FILE",
				"FIELD_INT_ARRAY_FILE", "FIELD_INT8_ARRAY_FILE", "FIELD_INT16_ARRAY_FILE", "FIELD_INT32_ARRAY_FILE", "FIELD_INT64_ARRAY_FILE",
				"FIELD_UINT_ARRAY_FILE", "FIELD_UINT8_ARRAY_FILE", "FIELD_UINT16_ARRAY_FILE", "FIELD_UINT32_ARRAY_FILE", "FIELD_UINT64_ARRAY_FILE",
				"FIELD_DURATION_ARRAY_FILE", "FIELD_URL_ARRAY_FILE",
			},
			expectedListSeps: []string{
				",", ",", ",",
				",",
				",",
				",", ",",
				",", ",", ",", ",", ",",
				",", ",", ",", ",", ",",
				",", ",",
				",",
				",",
				",", ",",
				",", ",", ",", ",", ",",
				",", ",", ",", ",", ",",
				",", ",",
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// values := []reflect.Value{}
			fieldNames := []string{}
			flagNames := []string{}
			envNames := []string{}
			fileEnvNames := []string{}
			listSeps := []string{}

			vStruct, err := validateStruct(tc.config)
			assert.NoError(t, err)

			tc.c.iterateOnFields(vStruct, func(v reflect.Value, fieldName, flagName, envName, fileEnvName, listSep string) {
				// values = append(values, v)
				fieldNames = append(fieldNames, fieldName)
				flagNames = append(flagNames, flagName)
				envNames = append(envNames, envName)
				fileEnvNames = append(fileEnvNames, fileEnvName)
				listSeps = append(listSeps, listSep)
			})

			// assert.Equal(t, tc.expectedValues, values)
			assert.Equal(t, tc.expectedFieldNames, fieldNames)
			assert.Equal(t, tc.expectedFlagNames, flagNames)
			assert.Equal(t, tc.expectedEnvNames, envNames)
			assert.Equal(t, tc.expectedFileEnvNames, fileEnvNames)
			assert.Equal(t, tc.expectedListSeps, listSeps)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestRegisterFlags(t *testing.T) {
	tests := []struct {
		name          string
		c             *controller
		config        interface{}
		expectedError error
		expectedFlags []string
	}{
		{
			name:          "OK",
			c:             &controller{},
			config:        &config{},
			expectedError: nil,
			expectedFlags: []string{
				"field.string",
				"field.bool",
				"field.float32", "field.float64",
				"field.int", "field.int8", "field.int16", "field.int32", "field.int64",
				"field.uint", "field.uint8", "field.uint16", "field.uint32", "field.uint64",
				"field.duration", "field.url",
				"field.string.array",
				"field.bool.array",
				"field.float32.array", "field.float64.array",
				"field.int.array", "field.int8.array", "field.int16.array", "field.int32.array", "field.int64.array",
				"field.uint.array", "field.uint8.array", "field.uint16.array", "field.uint32.array", "field.uint64.array",
				"field.duration.array", "field.url.array",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vStruct, err := validateStruct(tc.config)
			assert.NoError(t, err)

			tc.c.registerFlags(vStruct)

			for _, expectedFlag := range tc.expectedFlags {
				f := flag.Lookup(expectedFlag)
				assert.NotEmpty(t, f)
			}
		})
	}
}

func TestReadConfig(t *testing.T) {
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
		c              *controller
		config         interface{}
		watchMode      bool
		expectedConfig interface{}
	}{
		{
			"Empty",
			[]string{"path/to/binary"},
			[]env{},
			[]file{},
			&controller{},
			&config{},
			false,
			&config{},
		},
		{
			"AllFromDefaults",
			[]string{"path/to/binary"},
			[]env{},
			[]file{},
			&controller{},
			&config{
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
			false,
			&config{
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
			&controller{},
			&config{},
			false,
			&config{
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
			&controller{},
			&config{},
			false,
			&config{
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
			&controller{},
			&config{},
			false,
			&config{
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
			&controller{},
			&config{},
			false,
			&config{
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
			&controller{},
			&config{},
			false,
			&config{
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
			&controller{},
			&config{},
			false,
			&config{
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
			&controller{},
			&config{
				FieldString:        "default",
				FieldStringArray:   []string{"milad", "mona"},
				FieldBool:          true,
				FieldBoolArray:     []bool{false, true},
				FieldDuration:      d90m,
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURL:           *service1URL,
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			false,
			&config{
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
			&controller{
				telepresence: true,
			},
			&config{},
			false,
			&config{
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
			"WatchMode",
			[]string{
				"path/to/binary",
				"-field.string", "content",
				"-field.bool",
				"-field.string.array", "milad,mona",
				"-field.bool.array", "false,true",
			},
			[]env{
				{"FIELD_DURATION", "90m"},
				{"FIELD_URL", "service-1:8080"},
				{"FIELD_DURATION_ARRAY", "90m,120m"},
				{"FIELD_URL_ARRAY", "service-1:8080,service-2:8080"},
			},
			[]file{
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
			},
			&controller{},
			&config{},
			true,
			&config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "",
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
				FieldDuration:      0,
				FieldURL:           url.URL{},
				FieldStringArray:   nil,
				FieldBoolArray:     nil,
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
				FieldDurationArray: nil,
				FieldURLArray:      nil,
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

			// Testing Telepresence option
			if tc.c.telepresence {
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

			vStruct, err := validateStruct(tc.config)
			assert.NoError(t, err)

			tc.c.readConfig(vStruct, tc.watchMode)
			assert.Equal(t, tc.expectedConfig, tc.config)
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
		config         interface{}
		opts           []Option
		expectedError  error
		expectedConfig *config
	}{
		{
			"NonStruct",
			[]string{"path/to/binary"},
			[]env{},
			[]file{},
			new(string),
			nil,
			errors.New("a non-struct type is passed"),
			&config{},
		},
		{
			"NonPointer",
			[]string{"path/to/binary"},
			[]env{},
			[]file{},
			config{},
			nil,
			errors.New("a non-pointer type is passed"),
			&config{},
		},
		{
			"Empty",
			[]string{"path/to/binary"},
			[]env{},
			[]file{},
			&config{},
			nil,
			nil,
			&config{},
		},
		{
			"AllFromDefaults",
			[]string{"path/to/binary"},
			[]env{},
			[]file{},
			&config{
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
			nil,
			&config{
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
			&config{},
			nil,
			nil,
			&config{
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
			&config{},
			nil,
			nil,
			&config{
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
			&config{},
			nil,
			nil,
			&config{
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
			&config{},
			nil,
			nil,
			&config{
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
			&config{},
			nil,
			nil,
			&config{
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
			&config{},
			nil,
			nil,
			&config{
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
			&config{
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
			nil,
			&config{
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
			&config{},
			[]Option{
				Telepresence(),
			},
			nil,
			&config{
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
			c := &controller{}
			for _, opt := range tc.opts {
				opt(c)
			}

			// Set arguments for flags
			os.Args = tc.args

			// Set environment variables
			for _, e := range tc.envs {
				err := os.Setenv(e.varName, e.value)
				assert.NoError(t, err)
				defer os.Unsetenv(e.varName)
			}

			// Testing Telepresence option
			if c.telepresence {
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

			err := Pick(tc.config, tc.opts...)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfig, tc.config)
			}
		})
	}

	// flag.Parse() can be called only once
	flag.Parse()
}
