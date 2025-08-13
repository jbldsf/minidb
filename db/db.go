package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"reflect"
	"time"
)

const extension string = ".json"

func foreignKey(models, model reflect.Value) bool {
	name := model.Type().Name() + extension
	data, err := os.ReadFile(name)
	if err == nil {
		err = json.Unmarshal(data, models.Addr().Interface())
		if err == nil {
			if !model.IsZero() {
				for f := range model.NumField() {
					for l := models.Field(0).Len() - 1; l >= 0; l-- {
						if model.Type().Field(f).Tag.Get("db") == "fk" {
							if models.Field(0).Index(l).Field(f).Interface() == model.Field(f).Interface() {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

func Delete(models, model reflect.Value) error {
	name := model.Type().Name() + extension
	data, err := os.ReadFile(name)
	if err == nil {
		err = json.Unmarshal(data, models.Addr().Interface())
		if err == nil {
			for f := range model.NumField() {
				for l := models.Field(0).Len() - 1; l >= 0; l-- {
					if model.Type().Field(f).Tag.Get("db") == "fk" {
						fkms := model.Field(f)
						fkm := reflect.New(fkms.Field(0).Type().Elem()).Elem()
						fkm.FieldByName(model.Type().Field(f).Tag.Get("fk")).Set(model.FieldByName(model.Type().Field(f).Tag.Get("refer")))
						if foreignKey(fkms, fkm) {
							return fmt.Errorf("foreign key constraint on %s field", model.Type().Field(f).Name)
						}
					}
					if model.Type().Field(f).Tag.Get("db") == "id" {
						if models.Field(0).Index(l).Field(f).Interface() == model.Field(f).Interface() {
							models.Field(0).Set(reflect.AppendSlice(models.Field(0).Slice(0, l), models.Field(0).Slice(l+1, models.Field(0).Len())))
							break
						}
					}
				}
			}
			data, err = json.Marshal(models.Interface())
			if err == nil {
				return os.WriteFile(name, data, 0644)
			}
		}
	}
	return err
}

func Insert(models, model reflect.Value) error {
	name := model.Type().Name() + extension
	_, err := os.Stat(name)
	if errors.Is(err, fs.ErrNotExist) {
		data, err := json.Marshal(models.Interface())
		if err != nil {
			return err
		}
		err = os.WriteFile(name, data, 0644)
		if err != nil {
			return err
		}
	}
	data, err := os.ReadFile(name)
	if err == nil {
		err = json.Unmarshal(data, models.Addr().Interface())
		if err == nil {
			for f := range model.NumField() {
				for l := models.Field(0).Len() - 1; l >= 0; l-- {
					if model.Type().Field(f).Tag.Get("db") == "pk" {
						if models.Field(0).Index(l).Field(f).Interface() == model.Field(f).Interface() {
							return fmt.Errorf("primary key constraint on %s field", model.Type().Field(f).Name)
						}
					}
				}
				if model.Type().Field(f).Name == "CreatedAt" {
					model.Field(f).Set(reflect.ValueOf(time.Now()))
				}
			}
			models.Field(0).Set(reflect.Append(models.Field(0), model))
			data, err = json.Marshal(models.Interface())
			if err == nil {
				return os.WriteFile(name, data, 0644)
			}
		}
	}
	return err
}

func Select(models, model reflect.Value) error {
	name := model.Type().Name() + extension
	data, err := os.ReadFile(name)
	if err == nil {
		err = json.Unmarshal(data, models.Addr().Interface())
		if err == nil {
			if !model.IsZero() {
				for f := range model.NumField() {
					for l := models.Field(0).Len() - 1; l >= 0; l-- {
						if model.Type().Field(f).Tag.Get("db") == "id" {
							if models.Field(0).Index(l).Field(f).Interface() != model.Field(f).Interface() {
								models.Field(0).Set(reflect.AppendSlice(models.Field(0).Slice(0, l), models.Field(0).Slice(l+1, models.Field(0).Len())))
							}
						}
					}
				}
			}
		}
	}
	return err
}

func Update(models, model reflect.Value) error {
	name := model.Type().Name() + extension
	data, err := os.ReadFile(name)
	if err == nil {
		err = json.Unmarshal(data, models.Addr().Interface())
		if err == nil {
			for f := range model.NumField() {
				for l := models.Field(0).Len() - 1; l >= 0; l-- {
					if model.Type().Field(f).Tag.Get("db") == "id" {
						if models.Field(0).Index(l).Field(f).Interface() == model.Field(f).Interface() {
							for fs := range model.NumField() {
								if model.Type().Field(fs).Name == "UpdatedAt" {
									model.Field(fs).Set(reflect.ValueOf(time.Now()))
								}
								if model.Type().Field(fs).Tag.Get("db") != "id" && model.Type().Field(fs).Tag.Get("db") != "pk" {
									models.Field(0).Index(l).Field(fs).Set(model.Field(fs))
								}
							}
						}
					}
				}
			}
			fmt.Println(models)
			data, err = json.Marshal(models.Interface())
			if err == nil {
				return os.WriteFile(name, data, 0644)
			}
		}
	}
	return err
}
