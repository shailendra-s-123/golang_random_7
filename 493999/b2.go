package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Params struct {
	data map[string]interface{}
}

func NewParams() *Params {
	return &Params{
		data: make(map[string]interface{}),
	}
}
func (p *Params) Set(key string, value interface{}) {
	p.data[key] = value
}
func (p *Params) Get(key string) (interface{}, bool) {
	v, ok := p.data[key]
	return v, ok
}
func (p *Params) MustGet(key string) interface{} {
	v, ok := p.Get(key)
	if !ok {
		panic(fmt.Errorf("parameter '%s' not found", key))
	}
	return v
}
func (p *Params) String(key string) string {
	v, ok := p.Get(key)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v", v)
}
func (p *Params) Int(key string) (int, error) {
	s := p.String(key)
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}
func (p *Params) Bool(key string) (bool, error) {
	s := p.String(key)
	if s == "" {
		return false, nil
	}
	return strconv.ParseBool(s)
}
func (p *Params) Float64(key string) (float64, error) {
	s := p.String(key)
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}
func (p *Params) Slice(key string) ([]string, error) {
	v, ok := p.Get(key)
	if !ok {
		return nil, nil
	}
	switch v := v.(type) {
	case []string:
		return v, nil
	default:
		return nil, errors.New("invalid type for slice parameter")
	}
}
func (p *Params) Map(key string) (map[string]string, error) {
	v, ok := p.Get(key)
	if !ok {
		return nil, nil
	}
	switch v := v.(type) {
	case map[string]string:
		return v, nil
	default:
		return nil, errors.New("invalid type for map parameter")
	}
}
func (p *Params) ParseQuery(query string) error {
	u, err := url.ParseQuery(query)
	if err != nil {
		return err
	}
	for key, values := range u {
		for _, value := range values {
			if err := p.parseValue(key, value); err != nil {
				return err
			}
		}
	}
	return nil
}
func (p *Params) parseValue(key, value string) error {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		p.Set(key, value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		p.Set(key, i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		p.Set(key, f)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err