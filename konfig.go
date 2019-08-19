// Package konfig is a minimal and unopinionated library for reading configuration values in Go applications
// based on The 12-Factor App (https://12factor.net/config).
package konfig

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	flagTag    = "flag"
	envTag     = "env"
	fileEnvTag = "fileenv"
	sepTag     = "sep"
	skipValue  = "-"

	telepresenceEnvVar = "TELEPRESENCE_ROOT"

	separatorLog = "----------------------------------------------------------------------------------------------------"
)

// Option configures how configuration values are read
type Option interface {
	apply(*options)
}

// option implements Option interface
type option struct {
	function func(*options)
}

func (o *option) apply(s *options) {
	o.function(s)
}

func newOption(function func(*options)) *option {
	return &option{
		function: function,
	}
}

// Debug is the option for enabling logs for debugging purposes.
// You should not use this option in Production.
func Debug() Option {
	return newOption(func(o *options) {
		o.debug = true
	})
}

// Telepresence is the option for reading files when running in a Telepresence shell.
// If the TELEPRESENCE_ROOT environment variable exist, files will be read from mounted volume.
// See https://telepresence.io/howto/volumes.html for details.
func Telepresence() Option {
	return newOption(func(o *options) {
		o.telepresence = true
	})
}

// options contains all the options
type options struct {
	debug        bool
	telepresence bool
}

// defaultOptions creates options with default values
func defaultOptions() *options {
	return &options{
		debug:        false,
		telepresence: false,
	}
}

// String is used for printing debugging information.
// The output should fit in one line.
func (o options) String() string {
	opts := []string{}

	if o.debug {
		opts = append(opts, "Debug")
	}

	if o.telepresence {
		opts = append(opts, "Telepresence")
	}

	return strings.Join(opts, " + ")
}

func (o *options) print(msg string, args ...interface{}) {
	if o.debug {
		log.Printf(msg+"\n", args...)
	}
}

/*
 * getFieldValue reads and returns the string value for a field from either
 *   - command-line flags,
 *   - environment variables,
 *   - or configuration files
 * If the value is read from a configuration file, the second return will be true.
 */
func (o *options) getFieldValue(field, flag, env, fileenv string) (string, bool) {
	var value string
	var fromFile bool

	// First, try reading from flag
	if value == "" && flag != skipValue {
		value = getFlagValue(flag)
		o.print("[%s] value read from flag %s: %s", field, flag, value)
	}

	// Second, try reading from environment variable
	if value == "" && env != skipValue {
		value = os.Getenv(env)
		o.print("[%s] value read from environment variable %s: %s", field, env, value)
	}

	// Third, try reading from file
	if value == "" && fileenv != skipValue {
		// Read file environment variable
		val := os.Getenv(fileenv)
		o.print("[%s] value read from file environment variable %s: %s", field, fileenv, val)

		if val != "" {
			root := "/"

			// Check for Telepresence
			// See https://telepresence.io/howto/volumes.html for details
			if o.telepresence {
				if tr := os.Getenv(telepresenceEnvVar); tr != "" {
					root = tr
					o.print("[%s] telepresence root path: %s", field, tr)
				}
			}

			// Read config file
			file := filepath.Join(root, val)
			content, err := ioutil.ReadFile(file)
			if err == nil {
				value = string(content)
				fromFile = true
				o.print("[%s] value read from file %s: %s", field, file, value)
			}
		}
	}

	if value == "" {
		o.print("[%s] falling back to default value", field)
	}

	return value, fromFile
}

func (o *options) setString(v reflect.Value, name, val string) {
	if v.String() != val {
		o.print("[%s] setting string value: %s", name, val)
		v.SetString(val)
	}
}

func (o *options) setBool(v reflect.Value, name, val string) {
	if b, err := strconv.ParseBool(val); err == nil {
		if v.Bool() != b {
			o.print("[%s] setting boolean value: %t", name, b)
			v.SetBool(b)
		}
	}
}

func (o *options) setFloat(v reflect.Value, name, val string) {
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		if v.Float() != f {
			o.print("[%s] setting float value: %f", name, f)
			v.SetFloat(f)
		}
	}
}

func (o *options) setInt(v reflect.Value, name, val string) {
	if t := v.Type(); t.PkgPath() == "time" && t.Name() == "Duration" {
		// time.Duration
		if d, err := time.ParseDuration(val); err == nil {
			if v.Interface() != d {
				o.print("[%s] setting duration value: %s", name, d)
				v.Set(reflect.ValueOf(d))
			}
		}
	} else if i, err := strconv.ParseInt(val, 10, 64); err == nil {
		if v.Int() != i {
			o.print("[%s] setting integer value: %d", name, i)
			v.SetInt(i)
		}
	}
}

func (o *options) setUint(v reflect.Value, name, val string) {
	if u, err := strconv.ParseUint(val, 10, 64); err == nil {
		if v.Uint() != u {
			o.print("[%s] setting unsigned integer value: %d", name, u)
			v.SetUint(u)
		}
	}
}

func (o *options) setStruct(v reflect.Value, name, val string) {
	if t := v.Type(); t.PkgPath() == "net/url" && t.Name() == "URL" {
		// url.URL
		if u, err := url.Parse(val); err == nil {
			// u is a pointer
			if !reflect.DeepEqual(v.Interface(), *u) {
				o.print("[%s] setting url value: %s", name, val)
				v.Set(reflect.ValueOf(u).Elem())
			}
		}
	}
}

func (o *options) setStringSlice(v reflect.Value, name string, vals []string) {
	if !reflect.DeepEqual(v.Interface(), vals) {
		o.print("[%s] setting string slice: %v", name, vals)
		v.Set(reflect.ValueOf(vals))
	}
}

func (o *options) setBoolSlice(v reflect.Value, name string, vals []string) {
	bools := []bool{}
	for _, val := range vals {
		if b, err := strconv.ParseBool(val); err == nil {
			bools = append(bools, b)
		}
	}

	if !reflect.DeepEqual(v.Interface(), bools) {
		o.print("[%s] setting boolean slice: %v", name, bools)
		v.Set(reflect.ValueOf(bools))
	}
}

func (o *options) setFloat32Slice(v reflect.Value, name string, vals []string) {
	floats := []float32{}
	for _, val := range vals {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			floats = append(floats, float32(f))
		}
	}

	if !reflect.DeepEqual(v.Interface(), floats) {
		o.print("[%s] setting float32 slice: %v", name, floats)
		v.Set(reflect.ValueOf(floats))
	}
}

func (o *options) setFloat64Slice(v reflect.Value, name string, vals []string) {
	floats := []float64{}
	for _, val := range vals {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			floats = append(floats, f)
		}
	}

	if !reflect.DeepEqual(v.Interface(), floats) {
		o.print("[%s] setting float64 slice: %v", name, floats)
		v.Set(reflect.ValueOf(floats))
	}
}

func (o *options) setIntSlice(v reflect.Value, name string, vals []string) {
	ints := []int{}
	for _, val := range vals {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			ints = append(ints, int(i))
		}
	}

	if !reflect.DeepEqual(v.Interface(), ints) {
		o.print("[%s] setting int slice: %v", name, ints)
		v.Set(reflect.ValueOf(ints))
	}
}

func (o *options) setInt8Slice(v reflect.Value, name string, vals []string) {
	ints := []int8{}
	for _, val := range vals {
		if i, err := strconv.ParseInt(val, 10, 8); err == nil {
			ints = append(ints, int8(i))
		}
	}

	if !reflect.DeepEqual(v.Interface(), ints) {
		o.print("[%s] setting int8 slice: %v", name, ints)
		v.Set(reflect.ValueOf(ints))
	}
}

func (o *options) setInt16Slice(v reflect.Value, name string, vals []string) {
	ints := []int16{}
	for _, val := range vals {
		if i, err := strconv.ParseInt(val, 10, 16); err == nil {
			ints = append(ints, int16(i))
		}
	}

	if !reflect.DeepEqual(v.Interface(), ints) {
		o.print("[%s] setting int16 slice: %v", name, ints)
		v.Set(reflect.ValueOf(ints))
	}
}

func (o *options) setInt32Slice(v reflect.Value, name string, vals []string) {
	ints := []int32{}
	for _, val := range vals {
		if i, err := strconv.ParseInt(val, 10, 32); err == nil {
			ints = append(ints, int32(i))
		}
	}

	if !reflect.DeepEqual(v.Interface(), ints) {
		o.print("[%s] setting int32 slice: %v", name, ints)
		v.Set(reflect.ValueOf(ints))
	}
}

func (o *options) setInt64Slice(v reflect.Value, name string, vals []string) {
	if t := reflect.TypeOf(v.Interface()).Elem(); t.PkgPath() == "time" && t.Name() == "Duration" {
		durations := []time.Duration{}
		for _, val := range vals {
			if d, err := time.ParseDuration(val); err == nil {
				durations = append(durations, d)
			}
		}

		// []time.Duration
		if !reflect.DeepEqual(v.Interface(), durations) {
			o.print("[%s] setting duration slice: %v", name, durations)
			v.Set(reflect.ValueOf(durations))
		}
	} else {
		ints := []int64{}
		for _, val := range vals {
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				ints = append(ints, i)
			}
		}

		if !reflect.DeepEqual(v.Interface(), ints) {
			o.print("[%s] setting int64 slice: %v", name, ints)
			v.Set(reflect.ValueOf(ints))
		}
	}
}

func (o *options) setUintSlice(v reflect.Value, name string, vals []string) {
	uints := []uint{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 64); err == nil {
			uints = append(uints, uint(u))
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		o.print("[%s] setting uint slice: %v", name, uints)
		v.Set(reflect.ValueOf(uints))
	}
}

func (o *options) setUint8Slice(v reflect.Value, name string, vals []string) {
	uints := []uint8{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 8); err == nil {
			uints = append(uints, uint8(u))
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		o.print("[%s] setting uint8 slice: %v", name, uints)
		v.Set(reflect.ValueOf(uints))
	}
}

func (o *options) setUint16Slice(v reflect.Value, name string, vals []string) {
	uints := []uint16{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 16); err == nil {
			uints = append(uints, uint16(u))
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		o.print("[%s] setting uint16 slice: %v", name, uints)
		v.Set(reflect.ValueOf(uints))
	}
}

func (o *options) setUint32Slice(v reflect.Value, name string, vals []string) {
	uints := []uint32{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 32); err == nil {
			uints = append(uints, uint32(u))
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		o.print("[%s] setting uint32 slice: %v", name, uints)
		v.Set(reflect.ValueOf(uints))
	}
}

func (o *options) setUint64Slice(v reflect.Value, name string, vals []string) {
	uints := []uint64{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 64); err == nil {
			uints = append(uints, u)
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		o.print("[%s] setting uint64 slice: %v", name, uints)
		v.Set(reflect.ValueOf(uints))
	}
}

func (o *options) setURLSlice(v reflect.Value, name string, vals []string) {
	t := reflect.TypeOf(v.Interface()).Elem()

	if t.PkgPath() == "net/url" && t.Name() == "URL" {
		urls := []url.URL{}
		for _, val := range vals {
			if u, err := url.Parse(val); err == nil {
				urls = append(urls, *u)
			}
		}

		// []url.URL
		if !reflect.DeepEqual(v.Interface(), urls) {
			o.print("[%s] setting url slice: %v", name, urls)
			v.Set(reflect.ValueOf(urls))
		}
	}
}

func (o *options) read(config interface{}) error {
	v := reflect.ValueOf(config) // reflect.Value --> v.Type(), v.Kind(), v.NumField()
	t := reflect.TypeOf(config)  // reflect.Type --> t.Name(), t.Kind(), t.NumField()

	// A pointer to a struct should be passed
	if t.Kind() != reflect.Ptr {
		o.print("a non-pointer type is passed")
		return errors.New("a non-pointer type is passed")
	}

	// Navigate to the pointer value
	v = v.Elem()
	t = t.Elem()

	if t.Kind() != reflect.Struct {
		o.print("a non-struct type is passed")
		return errors.New("a non-struct type is passed")
	}

	// Iterate over struct fields
	for i := 0; i < v.NumField(); i++ {
		vField := v.Field(i) // reflect.Value --> vField.Kind(), vField.Type().Name(), vField.Type().Kind(), vField.Interface()
		tField := t.Field(i) // reflect.StructField --> tField.Name, tField.Type.Name(), tField.Type.Kind(), tField.Tag.Get(tag)

		// Skip unexported fields
		if !vField.CanSet() {
			continue
		}

		o.print(separatorLog)

		name := tField.Name

		// `flag:"..."`
		flagName := tField.Tag.Get(flagTag)
		if flagName == "" {
			flagName = getFlagName(name)
		}

		o.print("[%s] expecting flag name: %s", name, flagName)

		// `env:"..."`
		envName := tField.Tag.Get(envTag)
		if envName == "" {
			envName = getEnvVarName(name)
		}
		o.print("[%s] expecting environment variable name: %s", name, envName)

		// `fileenv:"..."`
		fileEnvName := tField.Tag.Get(fileEnvTag)
		if fileEnvName == "" {
			fileEnvName = getFileEnvVarName(name)
		}

		o.print("[%s] expecting file environment variable name: %s", name, fileEnvName)

		// `sep:"..."`
		sep := tField.Tag.Get(sepTag)
		if sep == "" {
			sep = ","
		}

		o.print("[%s] expecting list separator: %s", name, sep)

		// Define a flag for the field, so flag.Parse() can be called
		defaultValue := fmt.Sprintf("%v", vField.Interface())
		defineFlag(flagName, defaultValue, envName, fileEnvName)

		// Try reading the configuration value for current field
		// If no value, skip this field
		str, _ := o.getFieldValue(name, flagName, envName, fileEnvName)
		if str == "" {
			continue
		}

		switch vField.Kind() {
		case reflect.String:
			o.setString(vField, name, str)
		case reflect.Bool:
			o.setBool(vField, name, str)
		case reflect.Float32, reflect.Float64:
			o.setFloat(vField, name, str)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			o.setInt(vField, name, str)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			o.setUint(vField, name, str)
		case reflect.Struct:
			o.setStruct(vField, name, str)

		case reflect.Slice:
			tSlice := reflect.TypeOf(vField.Interface()).Elem()
			strs := strings.Split(str, sep)

			switch tSlice.Kind() {
			case reflect.String:
				o.setStringSlice(vField, name, strs)
			case reflect.Bool:
				o.setBoolSlice(vField, name, strs)
			case reflect.Float32:
				o.setFloat32Slice(vField, name, strs)
			case reflect.Float64:
				o.setFloat64Slice(vField, name, strs)
			case reflect.Int:
				o.setIntSlice(vField, name, strs)
			case reflect.Int8:
				o.setInt8Slice(vField, name, strs)
			case reflect.Int16:
				o.setInt16Slice(vField, name, strs)
			case reflect.Int32:
				o.setInt32Slice(vField, name, strs)
			case reflect.Int64:
				o.setInt64Slice(vField, name, strs)
			case reflect.Uint:
				o.setUintSlice(vField, name, strs)
			case reflect.Uint8:
				o.setUint8Slice(vField, name, strs)
			case reflect.Uint16:
				o.setUint16Slice(vField, name, strs)
			case reflect.Uint32:
				o.setUint32Slice(vField, name, strs)
			case reflect.Uint64:
				o.setUint64Slice(vField, name, strs)
			case reflect.Struct:
				o.setURLSlice(vField, name, strs)
			}
		}
	}

	o.print(separatorLog)

	return nil
}

// Pick reads values for exported fields of a struct from either command-line flags, environment variables, or configuration files.
// You can also specify default values.
func Pick(config interface{}, opts ...Option) error {
	// Create settings
	o := defaultOptions()
	for _, opt := range opts {
		opt.apply(o)
	}

	o.print("pick options: %s", o)

	err := o.read(config)
	if err != nil {
		return err
	}

	return nil
}
