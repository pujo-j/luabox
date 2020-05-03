/*
 *    Copyright 2020 Josselin Pujo
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 *
 */

package test

import "errors"

//go:generate luaboxgen func.go

func Test(greeting string, name string) (string, error) {
	if greeting == "" {
		return "", errors.New("greeting is mandatory")
	}
	if name == "" {
		return "", errors.New("name is mandatory")
	}
	return greeting + " " + name, nil
}
