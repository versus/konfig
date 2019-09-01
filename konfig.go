// Package konfig is a minimal and unopinionated library for reading configuration values in Go applications
// based on The 12-Factor App (https://12factor.net/config).
package konfig

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	flagTag    = "flag"
	envTag     = "env"
	fileEnvTag = "fileenv"
	sepTag     = "sep"
	skip       = "-"

	defaultInterval    = 10 * time.Second
	telepresenceEnvVar = "TELEPRESENCE_ROOT"

	line = "----------------------------------------------------------------------------------------------------"
)

// Update represents a configuration field that received a new value
type Update struct {
	Name  string
	Value interface{}
}

// controller controls how configuration values are read
type controller struct {
	debug         uint
	telepresence  bool
	watchInterval time.Duration
	subscribers   []chan Update
}

// Option sets optional parameters for controller
type Option func(*controller)

// Debug is the option for enabling logs for debugging purposes.
// verbosity is the verbosity level of logs.
// You should not use this option in production.
func Debug(verbosity uint) Option {
	return func(c *controller) {
		c.debug = verbosity
	}
}

// Telepresence is the option for reading files when running in a Telepresence shell.
// If the TELEPRESENCE_ROOT environment variable exist, files will be read from mounted volume.
// See https://telepresence.io/howto/volumes.html for details.
func Telepresence() Option {
	return func(c *controller) {
		c.telepresence = true
	}
}

// WatchInterval is the option for overriding the default interval (10s) for watching.
func WatchInterval(d time.Duration) Option {
	return func(c *controller) {
		c.watchInterval = d
	}
}

// String is used for printing debugging information.
// The output should fit in one line.
func (c *controller) String() string {
	strs := []string{}

	if c.debug > 0 {
		strs = append(strs, fmt.Sprintf("Debug<%d>", c.debug))
	}

	if c.telepresence {
		strs = append(strs, "Telepresence")
	}

	if c.watchInterval > 0 {
		strs = append(strs, fmt.Sprintf("Watch<%s>", c.watchInterval))
	}

	if len(c.subscribers) > 0 {
		strs = append(strs, fmt.Sprintf("Subscribers<%d>", len(c.subscribers)))
	}

	return strings.Join(strs, " + ")
}

func (c *controller) log(v uint, msg string, args ...interface{}) {
	if c.debug >= v {
		log.Printf(msg+"\n", args...)
	}
}

/*
 * getFieldValue reads and returns the string value for a field from either
 *   - command-line flags,
 *   - environment variables,
 *   - or configuration files
 * If the value is read from a file, the second returned value will be true.
 */
func (c *controller) getFieldValue(fieldName, flagName, envName, fileEnvName string) (string, bool) {
	var value string
	var fromFile bool

	// First, try reading from flag
	if value == "" && flagName != skip {
		value = getFlagValue(flagName)
		c.log(3, "[%s] value read from flag %s: %s", fieldName, flagName, value)
	}

	// Second, try reading from environment variable
	if value == "" && envName != skip {
		value = os.Getenv(envName)
		c.log(3, "[%s] value read from environment variable %s: %s", fieldName, envName, value)
	}

	// Third, try reading from file
	if value == "" && fileEnvName != skip {
		// Read file environment variable
		val := os.Getenv(fileEnvName)
		c.log(3, "[%s] value read from file environment variable %s: %s", fieldName, fileEnvName, val)

		if val != "" {
			root := "/"

			// Check for Telepresence
			// See https://telepresence.io/howto/volumes.html for details
			if c.telepresence {
				if tr := os.Getenv(telepresenceEnvVar); tr != "" {
					root = tr
					c.log(3, "[%s] telepresence root path: %s", fieldName, tr)
				}
			}

			// Read config file
			file := filepath.Join(root, val)
			content, err := ioutil.ReadFile(file)
			if err == nil {
				value = string(content)
				fromFile = true
				c.log(3, "[%s] value read from file %s: %s", fieldName, file, value)
			}
		}
	}

	return value, fromFile
}

func (c *controller) notifySubscribers(name string, value interface{}) {
	if len(c.subscribers) == 0 {
		return
	}

	c.log(1, "[%s] notifying %d subscribers ...", name, len(c.subscribers))

	update := Update{
		Name:  name,
		Value: value,
	}

	for i, sub := range c.subscribers {
		go func(id int, ch chan Update) {
			c.log(1, "[%s] notifying subscriber %d ...", name, id)
			ch <- update
			c.log(1, "[%s] subscriber %d notified", name, id)
		}(i, sub)
	}
}

func (c *controller) setString(v reflect.Value, name, val string) bool {
	if v.String() != val {
		c.log(2, "[%s] setting string value: %s", name, val)
		v.SetString(val)
		c.notifySubscribers(name, val)
		return true
	}

	return false
}

func (c *controller) setBool(v reflect.Value, name, val string) bool {
	if b, err := strconv.ParseBool(val); err == nil {
		if v.Bool() != b {
			c.log(2, "[%s] setting boolean value: %t", name, b)
			v.SetBool(b)
			c.notifySubscribers(name, b)
			return true
		}
	}

	return false
}

func (c *controller) setFloat32(v reflect.Value, name, val string) bool {
	if f, err := strconv.ParseFloat(val, 32); err == nil {
		if v.Float() != f {
			c.log(2, "[%s] setting float value: %f", name, f)
			v.SetFloat(f)
			c.notifySubscribers(name, float32(f))
			return true
		}
	}

	return false
}

func (c *controller) setFloat64(v reflect.Value, name, val string) bool {
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		if v.Float() != f {
			c.log(2, "[%s] setting float value: %f", name, f)
			v.SetFloat(f)
			c.notifySubscribers(name, f)
			return true
		}
	}

	return false
}

func (c *controller) setInt(v reflect.Value, name, val string) bool {
	// int size and range are platform-dependent
	if i, err := strconv.ParseInt(val, 10, 64); err == nil {
		if v.Int() != i {
			c.log(2, "[%s] setting integer value: %d", name, i)
			v.SetInt(i)
			c.notifySubscribers(name, int(i))
			return true
		}
	}

	return false
}

func (c *controller) setInt8(v reflect.Value, name, val string) bool {
	if i, err := strconv.ParseInt(val, 10, 8); err == nil {
		if v.Int() != i {
			c.log(2, "[%s] setting integer value: %d", name, i)
			v.SetInt(i)
			c.notifySubscribers(name, int8(i))
			return true
		}
	}

	return false
}

func (c *controller) setInt16(v reflect.Value, name, val string) bool {
	if i, err := strconv.ParseInt(val, 10, 16); err == nil {
		if v.Int() != i {
			c.log(2, "[%s] setting integer value: %d", name, i)
			v.SetInt(i)
			c.notifySubscribers(name, int16(i))
			return true
		}
	}

	return false
}

func (c *controller) setInt32(v reflect.Value, name, val string) bool {
	if i, err := strconv.ParseInt(val, 10, 32); err == nil {
		if v.Int() != i {
			c.log(2, "[%s] setting integer value: %d", name, i)
			v.SetInt(i)
			c.notifySubscribers(name, int32(i))
			return true
		}
	}

	return false
}

func (c *controller) setInt64(v reflect.Value, name, val string) bool {
	if t := v.Type(); t.PkgPath() == "time" && t.Name() == "Duration" {
		// time.Duration
		if d, err := time.ParseDuration(val); err == nil {
			if v.Interface() != d {
				c.log(2, "[%s] setting duration value: %s", name, d)
				v.Set(reflect.ValueOf(d))
				c.notifySubscribers(name, d)
				return true
			}
		}
	} else if i, err := strconv.ParseInt(val, 10, 64); err == nil {
		if v.Int() != i {
			c.log(2, "[%s] setting integer value: %d", name, i)
			v.SetInt(i)
			c.notifySubscribers(name, i)
			return true
		}
	}

	return false
}

func (c *controller) setUint(v reflect.Value, name, val string) bool {
	// uint size and range are platform-dependent
	if u, err := strconv.ParseUint(val, 10, 64); err == nil {
		if v.Uint() != u {
			c.log(2, "[%s] setting unsigned integer value: %d", name, u)
			v.SetUint(u)
			c.notifySubscribers(name, uint(u))
			return true
		}
	}

	return false
}

func (c *controller) setUint8(v reflect.Value, name, val string) bool {
	if u, err := strconv.ParseUint(val, 10, 8); err == nil {
		if v.Uint() != u {
			c.log(2, "[%s] setting unsigned integer value: %d", name, u)
			v.SetUint(u)
			c.notifySubscribers(name, uint8(u))
			return true
		}
	}

	return false
}

func (c *controller) setUint16(v reflect.Value, name, val string) bool {
	if u, err := strconv.ParseUint(val, 10, 16); err == nil {
		if v.Uint() != u {
			c.log(2, "[%s] setting unsigned integer value: %d", name, u)
			v.SetUint(u)
			c.notifySubscribers(name, uint16(u))
			return true
		}
	}

	return false
}

func (c *controller) setUint32(v reflect.Value, name, val string) bool {
	if u, err := strconv.ParseUint(val, 10, 32); err == nil {
		if v.Uint() != u {
			c.log(2, "[%s] setting unsigned integer value: %d", name, u)
			v.SetUint(u)
			c.notifySubscribers(name, uint32(u))
			return true
		}
	}

	return false
}

func (c *controller) setUint64(v reflect.Value, name, val string) bool {
	if u, err := strconv.ParseUint(val, 10, 64); err == nil {
		if v.Uint() != u {
			c.log(2, "[%s] setting unsigned integer value: %d", name, u)
			v.SetUint(u)
			c.notifySubscribers(name, u)
			return true
		}
	}

	return false
}

func (c *controller) setStruct(v reflect.Value, name, val string) bool {
	if t := v.Type(); t.PkgPath() == "net/url" && t.Name() == "URL" {
		// url.URL
		if u, err := url.Parse(val); err == nil {
			// u is a pointer
			if !reflect.DeepEqual(v.Interface(), *u) {
				c.log(2, "[%s] setting url value: %s", name, val)
				v.Set(reflect.ValueOf(u).Elem())
				c.notifySubscribers(name, *u)
				return true
			}
		}
	}

	return false
}

func (c *controller) setStringSlice(v reflect.Value, name string, vals []string) bool {
	if !reflect.DeepEqual(v.Interface(), vals) {
		c.log(2, "[%s] setting string slice: %v", name, vals)
		v.Set(reflect.ValueOf(vals))
		c.notifySubscribers(name, vals)
		return true
	}

	return false
}

func (c *controller) setBoolSlice(v reflect.Value, name string, vals []string) bool {
	bools := []bool{}
	for _, val := range vals {
		if b, err := strconv.ParseBool(val); err == nil {
			bools = append(bools, b)
		}
	}

	if !reflect.DeepEqual(v.Interface(), bools) {
		c.log(2, "[%s] setting boolean slice: %v", name, bools)
		v.Set(reflect.ValueOf(bools))
		c.notifySubscribers(name, bools)
		return true
	}

	return false
}

func (c *controller) setFloat32Slice(v reflect.Value, name string, vals []string) bool {
	floats := []float32{}
	for _, val := range vals {
		if f, err := strconv.ParseFloat(val, 32); err == nil {
			floats = append(floats, float32(f))
		}
	}

	if !reflect.DeepEqual(v.Interface(), floats) {
		c.log(2, "[%s] setting float32 slice: %v", name, floats)
		v.Set(reflect.ValueOf(floats))
		c.notifySubscribers(name, floats)
		return true
	}

	return false
}

func (c *controller) setFloat64Slice(v reflect.Value, name string, vals []string) bool {
	floats := []float64{}
	for _, val := range vals {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			floats = append(floats, f)
		}
	}

	if !reflect.DeepEqual(v.Interface(), floats) {
		c.log(2, "[%s] setting float64 slice: %v", name, floats)
		v.Set(reflect.ValueOf(floats))
		c.notifySubscribers(name, floats)
		return true
	}

	return false
}

func (c *controller) setIntSlice(v reflect.Value, name string, vals []string) bool {
	// int size and range are platform-dependent
	ints := []int{}
	for _, val := range vals {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			ints = append(ints, int(i))
		}
	}

	if !reflect.DeepEqual(v.Interface(), ints) {
		c.log(2, "[%s] setting int slice: %v", name, ints)
		v.Set(reflect.ValueOf(ints))
		c.notifySubscribers(name, ints)
		return true
	}

	return false
}

func (c *controller) setInt8Slice(v reflect.Value, name string, vals []string) bool {
	ints := []int8{}
	for _, val := range vals {
		if i, err := strconv.ParseInt(val, 10, 8); err == nil {
			ints = append(ints, int8(i))
		}
	}

	if !reflect.DeepEqual(v.Interface(), ints) {
		c.log(2, "[%s] setting int8 slice: %v", name, ints)
		v.Set(reflect.ValueOf(ints))
		c.notifySubscribers(name, ints)
		return true
	}

	return false
}

func (c *controller) setInt16Slice(v reflect.Value, name string, vals []string) bool {
	ints := []int16{}
	for _, val := range vals {
		if i, err := strconv.ParseInt(val, 10, 16); err == nil {
			ints = append(ints, int16(i))
		}
	}

	if !reflect.DeepEqual(v.Interface(), ints) {
		c.log(2, "[%s] setting int16 slice: %v", name, ints)
		v.Set(reflect.ValueOf(ints))
		c.notifySubscribers(name, ints)
		return true
	}

	return false
}

func (c *controller) setInt32Slice(v reflect.Value, name string, vals []string) bool {
	ints := []int32{}
	for _, val := range vals {
		if i, err := strconv.ParseInt(val, 10, 32); err == nil {
			ints = append(ints, int32(i))
		}
	}

	if !reflect.DeepEqual(v.Interface(), ints) {
		c.log(2, "[%s] setting int32 slice: %v", name, ints)
		v.Set(reflect.ValueOf(ints))
		c.notifySubscribers(name, ints)
		return true
	}

	return false
}

func (c *controller) setInt64Slice(v reflect.Value, name string, vals []string) bool {
	if t := reflect.TypeOf(v.Interface()).Elem(); t.PkgPath() == "time" && t.Name() == "Duration" {
		durations := []time.Duration{}
		for _, val := range vals {
			if d, err := time.ParseDuration(val); err == nil {
				durations = append(durations, d)
			}
		}

		// []time.Duration
		if !reflect.DeepEqual(v.Interface(), durations) {
			c.log(2, "[%s] setting duration slice: %v", name, durations)
			v.Set(reflect.ValueOf(durations))
			c.notifySubscribers(name, durations)
			return true
		}
	} else {
		ints := []int64{}
		for _, val := range vals {
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				ints = append(ints, i)
			}
		}

		if !reflect.DeepEqual(v.Interface(), ints) {
			c.log(2, "[%s] setting int64 slice: %v", name, ints)
			v.Set(reflect.ValueOf(ints))
			c.notifySubscribers(name, ints)
			return true
		}
	}

	return false
}

func (c *controller) setUintSlice(v reflect.Value, name string, vals []string) bool {
	// uint size and range are platform-dependent
	uints := []uint{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 64); err == nil {
			uints = append(uints, uint(u))
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		c.log(2, "[%s] setting uint slice: %v", name, uints)
		v.Set(reflect.ValueOf(uints))
		c.notifySubscribers(name, uints)
		return true
	}

	return false
}

func (c *controller) setUint8Slice(v reflect.Value, name string, vals []string) bool {
	uints := []uint8{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 8); err == nil {
			uints = append(uints, uint8(u))
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		c.log(2, "[%s] setting uint8 slice: %v", name, uints)
		v.Set(reflect.ValueOf(uints))
		c.notifySubscribers(name, uints)
		return true
	}

	return false
}

func (c *controller) setUint16Slice(v reflect.Value, name string, vals []string) bool {
	uints := []uint16{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 16); err == nil {
			uints = append(uints, uint16(u))
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		c.log(2, "[%s] setting uint16 slice: %v", name, uints)
		v.Set(reflect.ValueOf(uints))
		c.notifySubscribers(name, uints)
		return true
	}

	return false
}

func (c *controller) setUint32Slice(v reflect.Value, name string, vals []string) bool {
	uints := []uint32{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 32); err == nil {
			uints = append(uints, uint32(u))
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		c.log(2, "[%s] setting uint32 slice: %v", name, uints)
		v.Set(reflect.ValueOf(uints))
		c.notifySubscribers(name, uints)
		return true
	}

	return false
}

func (c *controller) setUint64Slice(v reflect.Value, name string, vals []string) bool {
	uints := []uint64{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 64); err == nil {
			uints = append(uints, u)
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		c.log(2, "[%s] setting uint64 slice: %v", name, uints)
		v.Set(reflect.ValueOf(uints))
		c.notifySubscribers(name, uints)
		return true
	}

	return false
}

func (c *controller) setURLSlice(v reflect.Value, name string, vals []string) bool {
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
			c.log(2, "[%s] setting url slice: %v", name, urls)
			v.Set(reflect.ValueOf(urls))
			c.notifySubscribers(name, urls)
			return true
		}
	}

	return false
}

func (c *controller) iterateOnFields(vStruct reflect.Value, handle func(v reflect.Value, fieldName, flagName, envName, fileEnvName, listSep string)) {
	// Iterate over struct fields
	for i := 0; i < vStruct.NumField(); i++ {
		v := vStruct.Field(i)        // reflect.Value --> vField.Kind(), vField.Type().Name(), vField.Type().Kind(), vField.Interface()
		f := vStruct.Type().Field(i) // reflect.StructField --> tField.Name, tField.Type.Name(), tField.Type.Kind(), tField.Tag.Get(tag)

		// Skip unexported and unsupported fields
		if !v.CanSet() || !isTypeSupported(v.Type()) {
			continue
		}

		// `flag:"..."`
		flagName := f.Tag.Get(flagTag)
		if flagName == "" {
			flagName = getFlagName(f.Name)
		}

		// `env:"..."`
		envName := f.Tag.Get(envTag)
		if envName == "" {
			envName = getEnvVarName(f.Name)
		}

		// `fileenv:"..."`
		fileEnvName := f.Tag.Get(fileEnvTag)
		if fileEnvName == "" {
			fileEnvName = getFileEnvVarName(f.Name)
		}

		// `sep:"..."`
		listSep := f.Tag.Get(sepTag)
		if listSep == "" {
			listSep = ","
		}

		handle(v, f.Name, flagName, envName, fileEnvName, listSep)
	}
}

func (c *controller) registerFlags(vStruct reflect.Value) {
	c.log(3, "Registering configuration flags ...")
	c.log(3, line)
	defer c.log(3, line)

	c.iterateOnFields(vStruct, func(v reflect.Value, fieldName, flagName, envName, fileEnvName, listSep string) {
		if flagName == skip {
			return
		}

		defaultValue := fmt.Sprintf("%v", v.Interface())
		usage := fmt.Sprintf(
			"%s:\t\t\t\t%s\n%s:\t\t\t%s\n%s:\t%s",
			"default value", defaultValue,
			"environment variable", envName,
			"environment variable for file path", fileEnvName,
		)

		// Define a flag for the field, so flag.Parse() can be called
		if flag.Lookup(flagName) == nil {
			switch v.Kind() {
			case reflect.Bool:
				flag.Bool(flagName, v.Bool(), usage)
			default:
				flag.Var(&flagValue{}, flagName, usage)
			}
		}

		c.log(3, "[%s] flag registered: %s", fieldName, flagName)
	})
}

func (c *controller) readConfig(vStruct reflect.Value, watchMode bool) {
	if watchMode {
		c.log(2, "watching for new configurations ...")
	} else {
		c.log(2, "Reading configuration values ...")
	}

	c.log(2, line)

	c.iterateOnFields(vStruct, func(v reflect.Value, fieldName, flagName, envName, fileEnvName, listSep string) {
		c.log(3, "[%s] expecting flag name: %s", fieldName, flagName)
		c.log(3, "[%s] expecting environment variable name: %s", fieldName, envName)
		c.log(3, "[%s] expecting file environment variable name: %s", fieldName, fileEnvName)
		c.log(3, "[%s] expecting list separator: %s", fieldName, listSep)
		defer c.log(2, line)

		// Try reading the configuration value for current field
		val, fromFile := c.getFieldValue(fieldName, flagName, envName, fileEnvName)

		// If no value, skip this field
		if val == "" {
			c.log(2, "[%s] falling back to default value: %v", fieldName, v.Interface())
			return
		}

		// Only those configuration values read from files can recieve new values
		// In watch mode, if the value for a field is not read from a file, skip the field
		if watchMode && !fromFile {
			return
		}

		switch v.Kind() {
		case reflect.String:
			c.setString(v, fieldName, val)
		case reflect.Bool:
			c.setBool(v, fieldName, val)
		case reflect.Float32:
			c.setFloat32(v, fieldName, val)
		case reflect.Float64:
			c.setFloat64(v, fieldName, val)
		case reflect.Int:
			c.setInt(v, fieldName, val)
		case reflect.Int8:
			c.setInt8(v, fieldName, val)
		case reflect.Int16:
			c.setInt16(v, fieldName, val)
		case reflect.Int32:
			c.setInt32(v, fieldName, val)
		case reflect.Int64:
			c.setInt64(v, fieldName, val)
		case reflect.Uint:
			c.setUint(v, fieldName, val)
		case reflect.Uint8:
			c.setUint8(v, fieldName, val)
		case reflect.Uint16:
			c.setUint16(v, fieldName, val)
		case reflect.Uint32:
			c.setUint32(v, fieldName, val)
		case reflect.Uint64:
			c.setUint64(v, fieldName, val)
		case reflect.Struct:
			c.setStruct(v, fieldName, val)

		case reflect.Slice:
			tSlice := reflect.TypeOf(v.Interface()).Elem()
			vals := strings.Split(val, listSep)

			switch tSlice.Kind() {
			case reflect.String:
				c.setStringSlice(v, fieldName, vals)
			case reflect.Bool:
				c.setBoolSlice(v, fieldName, vals)
			case reflect.Float32:
				c.setFloat32Slice(v, fieldName, vals)
			case reflect.Float64:
				c.setFloat64Slice(v, fieldName, vals)
			case reflect.Int:
				c.setIntSlice(v, fieldName, vals)
			case reflect.Int8:
				c.setInt8Slice(v, fieldName, vals)
			case reflect.Int16:
				c.setInt16Slice(v, fieldName, vals)
			case reflect.Int32:
				c.setInt32Slice(v, fieldName, vals)
			case reflect.Int64:
				c.setInt64Slice(v, fieldName, vals)
			case reflect.Uint:
				c.setUintSlice(v, fieldName, vals)
			case reflect.Uint8:
				c.setUint8Slice(v, fieldName, vals)
			case reflect.Uint16:
				c.setUint16Slice(v, fieldName, vals)
			case reflect.Uint32:
				c.setUint32Slice(v, fieldName, vals)
			case reflect.Uint64:
				c.setUint64Slice(v, fieldName, vals)
			case reflect.Struct:
				c.setURLSlice(v, fieldName, vals)
			}
		}
	})
}

// Pick reads values for exported fields of a struct from either command-line flags, environment variables, or configuration files.
// You can also specify default values.
func Pick(config interface{}, opts ...Option) error {
	c := &controller{}

	// Applying options
	for _, opt := range opts {
		opt(c)
	}

	c.log(2, line)
	c.log(2, "Options: %s", c)
	c.log(2, line)

	v, err := validateStruct(config)
	if err != nil {
		return err
	}

	c.registerFlags(v)
	c.readConfig(v, false)

	return nil
}

// Watch first reads values for exported fields of a struct from either command-line flags, environment variables, or configuration files.
// It then watches any change to those fields that their values are read from configuration files and notifies subscribers on a channel.
func Watch(config sync.Locker, subscribers []chan Update, opts ...Option) (func(), error) {
	c := &controller{
		watchInterval: defaultInterval,
		subscribers:   subscribers,
	}

	// Applying options
	for _, opt := range opts {
		opt(c)
	}

	c.log(2, line)
	c.log(2, "Options: %s", c)
	c.log(2, line)

	v, err := validateStruct(config)
	if err != nil {
		return nil, err
	}

	c.registerFlags(v)
	c.readConfig(v, false)

	ticker := time.NewTicker(c.watchInterval)

	go func() {
		for range ticker.C {
			config.Lock()
			c.readConfig(v, true)
			config.Unlock()
		}
	}()

	stop := func() {
		ticker.Stop()
		for _, sub := range c.subscribers {
			close(sub)
		}
	}

	return stop, nil
}
