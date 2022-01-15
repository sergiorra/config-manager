package config_manager

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// globalConfigPath defines the global config path
	globalConfigPath = "config"
	// globalConfigFile defines the global config file name
	globalConfigFile = "config"
	// globalConfigType defines the global config file type
	globalConfigType = "json"

	// envConfigFile defines the environment config file name structure
	envConfigFile = "config.%v"

	// reserved keys for host configuration
	envKey = "env"
)

// HostConfiguration Host configuration struct
type HostConfiguration struct {
	env string
}

// reservedHostKeys slice of reserved Host keys
var reservedHostKeys = []string{envKey}

// Manager representation of manager into data struct
type Manager struct {
	viper              *viper.Viper
	envVariablesPrefix string
}

// Option function structure for the optional functions
type Option func(*Manager)

// NewManager initializes a new Manager and applies all the optional functions
func NewManager(opts ...Option) *Manager {
	mgr := &Manager{
		viper.New(),
		"",
	}

	for _, opt := range opts {
		opt(mgr)
	}

	return mgr
}

// WithDefault optional function to define default values in the constructor
func WithDefault(key string, value interface{}) Option {
	return func(mgr *Manager) {
		mgr.viper.SetDefault(key, value)
	}
}

// WithEnvPrefix optional function to define the environment variables prefix in the constructor
func WithEnvPrefix(value string) Option {
	return func(mgr *Manager) {
		mgr.envVariablesPrefix = value + "_"
	}
}

// Load manages all the configuration providers for the Host and App configuration
func (m *Manager) Load(cfg interface{}) error {
	if err := m.checkConfigParam(cfg); err != nil {
		return err
	}

	if err := m.bindHost(); err != nil {
		return err
	}

	if err := m.bindApp(cfg); err != nil {
		return err
	}

	// parsing flags after all Host and App flags binding are created
	pflag.Parse()

	if err := m.validateHost(); err != nil {
		return err
	}

	if err := m.loadConfigFiles(cfg); err != nil {
		return err
	}

	return nil
}

// checkConfigParam checks input parameters
func (m *Manager) checkConfigParam(cfg interface{}) error {
	if cfg == nil || reflect.ValueOf(cfg).Kind() != reflect.Ptr {
		return NewErrInvalidConfig()
	}

	return nil
}

// bindHost binds a flag and a environment variable for each Host configuration field
func (m *Manager) bindHost() error {
	v := reflect.ValueOf(HostConfiguration{})
	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name
		fieldKind := v.Type().Field(i).Type.Kind()
		m.defineFlag(fieldName, fieldKind)

		if err := m.viper.BindPFlag(fieldName, pflag.Lookup(fieldName)); err != nil {
			return NewErrBindFlags(err, fieldName)
		}

		if err := m.viper.BindEnv(fieldName, strings.ToUpper(m.envVariablesPrefix+fieldName)); err != nil {
			return NewErrBindEnvVar(err, fieldName)
		}

	}

	return nil
}

// validateHost validates host configuration values are not zero values
func (m *Manager) validateHost() error {
	v := reflect.ValueOf(HostConfiguration{})
	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name
		value := m.viper.Get(fieldName)
		if isZeroValue(value) {
			return NewErrMandatoryField(fieldName)
		}
	}

	return nil
}

// isZeroValue checks if a value has its zero value
func isZeroValue(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}

// bindApp binds a flag and a environment variable for each App configuration field and inner field
func (m *Manager) bindApp(cfg interface{}) error {
	replacer := strings.NewReplacer(".", "_")
	m.viper.SetEnvKeyReplacer(replacer)

	v := reflect.Indirect(reflect.ValueOf(cfg))
	for i := 0; i < v.NumField(); i++ {
		fieldName := strings.ToLower(v.Type().Field(i).Name)
		if contains(reservedHostKeys, fieldName) {
			continue
		}

		fieldKind := v.Type().Field(i).Type.Kind()
		if fieldKind != reflect.Struct && fieldKind != reflect.String && fieldKind != reflect.Int && fieldKind != reflect.Bool {
			return NewErrTypeNotSupported(fieldKind.String())
		}

		if fieldKind == reflect.Struct {
			innerStruct := reflect.Indirect(reflect.ValueOf(v.Field(i).Interface()))
			for j := 0; j < innerStruct.NumField(); j++ {
				fieldKind = innerStruct.Type().Field(j).Type.Kind()
				if fieldKind != reflect.Struct && fieldKind != reflect.String && fieldKind != reflect.Int && fieldKind != reflect.Bool {
					return NewErrTypeNotSupported(fieldKind.String())
				}

				if fieldKind == reflect.Struct {
					return NewErrMultiLevelNestedStruct()
				}

				innerFieldName := strings.ToLower(innerStruct.Type().Field(j).Name)
				fieldNameKey := fieldName + "." + innerFieldName

				m.defineFlag(fieldNameKey, fieldKind)
				if err := m.viper.BindPFlag(fieldNameKey, pflag.Lookup(fieldNameKey)); err != nil {
					return NewErrBindFlags(err, fieldNameKey)
				}

				if err := m.viper.BindEnv(fieldNameKey, strings.ToUpper(m.envVariablesPrefix+fieldName+"_"+innerFieldName)); err != nil {
					return NewErrBindEnvVar(err, fieldNameKey)
				}
			}
		} else {
			m.defineFlag(fieldName, fieldKind)
			if err := m.viper.BindPFlag(fieldName, pflag.Lookup(fieldName)); err != nil {
				return NewErrBindFlags(err, fieldName)
			}

			if err := m.viper.BindEnv(fieldName, strings.ToUpper(m.envVariablesPrefix+fieldName)); err != nil {
				return NewErrBindEnvVar(err, fieldName)
			}
		}

	}

	return nil
}

// contains checks if an specific string is in a string slice
func contains(values []string, val string) bool {
	for _, v := range values {
		if v == val {
			return true
		}
	}
	return false
}

// defineFlag defines a flag for the fieldName and its type received
func (m *Manager) defineFlag(fieldName string, fieldKind reflect.Kind) {
	switch fieldKind {
	case reflect.String:
		pflag.String(fieldName, "", "")
	case reflect.Int:
		pflag.Int(fieldName, 0, "")
	case reflect.Bool:
		pflag.Bool(fieldName, false, "")
	default:
		pflag.String(fieldName, "", "")
	}
}

// loadConfigFiles loads config files
func (m *Manager) loadConfigFiles(cfg interface{}) error {
	m.viper.AddConfigPath(globalConfigPath)
	m.viper.SetConfigName(globalConfigFile)
	m.viper.SetConfigType(globalConfigType)

	// Global file
	fileName := fmt.Sprintf("%s/%s.%s", globalConfigPath, globalConfigFile, globalConfigType)
	if err := m.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return NewErrNotFound(fileName)
		} else {
			return NewErrReadFile(err, fileName)
		}
	}

	if err := m.viper.Unmarshal(cfg); err != nil {
		return NewErrUnmarshalConfig(err, fileName)
	}

	// Specific file for environment
	envFile := fmt.Sprintf(envConfigFile, m.viper.Get(envKey))
	m.viper.SetConfigName(envFile)
	fileName = fmt.Sprintf("%s/%s.%s", globalConfigPath, envFile, globalConfigType)
	if err := m.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return NewErrReadFile(err, fileName)
		}
	}

	if err := m.viper.Unmarshal(cfg); err != nil {
		return NewErrUnmarshalConfig(err, fileName)
	}

	return nil
}
