package env_loader

import (
	"errors"
	"reflect"
	"sport_platform/internal/env_loader/utils"
)

type EnvLoader struct {
	env *utils.Env
}

func (loader *EnvLoader) LoadDataIntoStruct(v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Pointer {
		return errors.New("value should be a pointer")
	}

	if value.Elem().Kind() != reflect.Struct {
		return errors.New("value under pointer should be a struct")
	}

	return loader.loadDataIntoValue(value)
}

func (loader *EnvLoader) loadDataIntoValue(v reflect.Value) error {
	vType := v.Type().Elem()
	v = v.Elem()

	for fieldIdx := range v.NumField() {
		fieldValue := v.Field(fieldIdx)
		fieldType := vType.Field(fieldIdx)

		envKey, tagParseError := utils.ParseTag(vType, fieldType)
		if tagParseError != nil {
			return tagParseError
		}

		if envKey == nil {
			continue
		}

		if err := loader.env.GetIntoReflectValue(fieldValue.Addr(), *envKey); err != nil {
			return err
		}
	}

	return nil
}

func CreateLoaderFromFile(filePath string) (*EnvLoader, error) {
	env, err := utils.CreateEnvFromFile(filePath)
	if err != nil {
		return nil, err
	}

	return &EnvLoader{env: env}, nil
}

func CreateLoaderFromEnv() *EnvLoader {
	env := utils.CreateEnvFromEnvironment()
	return &EnvLoader{env: env}
}
