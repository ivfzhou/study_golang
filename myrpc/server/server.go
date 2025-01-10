/*
 * Copyright (c) 2023 ivfzhou
 * myrpc is licensed under Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *          http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
 * MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 */

package server

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

var (
	ctx     context.Context
	ctxType = reflect.TypeOf(&ctx).Elem()
)

type Server interface {
	Register(svr Service) error
	ListenAndServe() error
	Close(chan struct{}) error
}

type Service interface {
	Name() string
}

type server map[string]*reflect.Value

func New(opts ...Option) Server {
	return server{}
}

func (s server) Register(svr Service) error {
	if err := checkServiceName(svr.Name()); err != nil {
		return err
	}

	methodMap, err := parseServiceMethods(svr)
	if err != nil {
		return err
	}

	for k, v := range methodMap {
		s[svr.Name()+"/"+k] = v
	}

	return nil
}

func (s server) ListenAndServe() error {
	fmt.Println("服务注册列表...")
	for k, v := range s {
		fmt.Printf("服务名: %s, 方法: %s\n", k, v.Type().String())
	}
	return nil
}

func (s server) Close(chan struct{}) error {
	return nil
}

func checkServiceName(name string) error {
	for _, b := range name {
		if !(b >= 48 && b <= 57 || b >= 65 && b <= 90 || b >= 97 && b <= 122) {
			return errors.New("服务名只能包含数字字母")
		}
	}
	return nil
}

func parseServiceMethods(svr interface{}) (map[string]*reflect.Value, error) {
	val := reflect.ValueOf(svr)
	if val.IsZero() {
		return nil, errors.New("无法注册nil服务")
	}
	res := make(map[string]*reflect.Value, val.NumMethod())
	for i := 0; i < val.NumMethod(); i++ {
		method := val.Type().Method(i)
		if !method.IsExported() {
			continue
		}
		if method.Func.Type().NumIn() != 3 {
			continue
		}
		if method.Func.Type().NumOut() != 2 {
			continue
		}
		if !method.Func.Type().In(1).AssignableTo(ctxType) {
			continue
		}
		value := val.Method(i)
		res[method.Name] = &value
	}
	return res, nil
}
