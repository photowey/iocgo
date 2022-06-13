/*
 * Copyright © 2022 photowey (photowey@gmail.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package environment

var _ Environment = (*environment)(nil)

type Environment interface {

	// ---------------------------------------------------------------- Get

	GetProperty(key string, standBy any) (any, error)

	// ---------------------------------------------------------------- Set

	SetProperty(key string, standBy any)
}

type StandardEnvironment interface {
	Environment
	// ---------------------------------------------------------------- Get

	GetString(key string) (string, error)
	GetInt64(key string) (int64, error)
	GetInt32(key string) (int32, error)
	GetInt16(key string) (int16, error)
	GetInt8(key string) (int8, error)
	GetUInt64(key string) (uint64, error)
	GetUInt32(key string) (uint32, error)
	GetUInt16(key string) (uint16, error)
	GetUInt8(key string) (uint8, error)
	GetFloat64(key string) (float64, error)
	GetFloat32(key string) (float32, error)

	// ---------------------------------------------------------------- Set

	SetString(key string, value string)
	SetInt64(key string, value int64)
	SetInt32(key string, value int32)
	SetInt16(key string, value int16)
	SetInt8(key string, value int8)
	SetUInt64(key string, value uint64)
	SetUInt32(key string, value uint32)
	SetUInt16(key string, value uint16)
	SetUInt8(key string, value uint8)
	SetFloat64(key string, value float64)
	SetFloat32(key string, value float32)
}

type environment struct {
	ctx map[string]any
}

// ---------------------------------------------------------------- Get

// ---------------------------------------------------------------- any

func (evn *environment) GetProperty(key string, standBy any) (any, error) {
	return standBy, nil
}

// ---------------------------------------------------------------- string

func (evn *environment) GetString(key string) (string, error) {
	return "", nil
}

// ---------------------------------------------------------------- int

func (evn *environment) GetInt64(key string) (int64, error) {
	return 0, nil
}

func (evn *environment) GetInt32(key string) (int32, error) {
	return 0, nil
}

func (evn *environment) GetInt16(key string) (int16, error) {
	return 0, nil
}

func (evn *environment) GetInt8(key string) (int8, error) {
	return 0, nil
}

// ---------------------------------------------------------------- uint

func (evn *environment) GetUInt64(key string) (uint64, error) {
	return 0, nil
}

func (evn *environment) GetUInt32(key string) (uint32, error) {
	return 0, nil
}

func (evn *environment) GetUInt16(key string) (uint16, error) {
	return 0, nil
}

func (evn *environment) GetUInt8(key string) (uint8, error) {
	return 0, nil
}

// ---------------------------------------------------------------- uint

func (evn *environment) GetFloat64(key string) (float64, error) {
	return 0.0, nil
}

func (evn *environment) GetFloat32(key string) (float32, error) {
	return 0.0, nil
}

// ---------------------------------------------------------------- Set

func (evn *environment) SetProperty(key string, value any) {

}

func (evn *environment) SetString(key string, value string) {

}

func (evn *environment) SetInt64(key string, value int64) {

}

func (evn *environment) SetInt32(key string, value int32) {

}

func (evn *environment) SetInt16(key string, value int16) {

}

func (evn *environment) SetInt8(key string, value int8) {

}

func (evn *environment) SetUInt64(key string, value uint64) {

}

func (evn *environment) SetUInt32(key string, value uint32) {

}

func (evn *environment) SetUInt16(key string, value uint16) {

}

func (evn *environment) SetUInt8(key string, value uint8) {

}

func (evn *environment) SetFloat64(key string, value float64) {

}

func (evn *environment) SetFloat32(key string, value float32) {

}
