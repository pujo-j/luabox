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

package luabox

type Log interface {
	WithFields(context map[string]interface{}) Log
	Debug(msg string, context map[string]interface{})
	Info(msg string, context map[string]interface{})
	Warn(msg string, context map[string]interface{})
	Error(msg string, context map[string]interface{})
	Fatal(msg string, context map[string]interface{})
}
