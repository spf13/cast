// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package cast provides easy and safe casting in Go.
package cast

import "time"

// ToBool casts an interface to a bool type.
func ToBool(i interface{}) bool {
	v, _ := ToBoolE(i)
	return v
}

// ToBoolP casts an interface to a bool type.
func ToBoolP(fn func(i interface{}) (bool, error), i interface{}) bool {
	v, _ := ToBoolPE(fn, i)
	return v
}

// ToTime casts an interface to a time.Time type.
func ToTime(i interface{}) time.Time {
	v, _ := ToTimeE(i)
	return v
}

// ToTimeP casts an interface to a time.Time type.
func ToTimeP(fn func(i interface{}) (time.Time, error), i interface{}) time.Time {
	v, _ := ToTimePE(fn, i)
	return v
}

func ToTimeInDefaultLocation(i interface{}, location *time.Location) time.Time {
	v, _ := ToTimeInDefaultLocationE(i, location)
	return v
}

// ToTimeInDefaultLocationP casts an interface to a time.Time type.
func ToTimeInDefaultLocationP(fn func(i interface{}, location *time.Location) (time.Time, error), i interface{}, location *time.Location) time.Time {
	v, _ := ToTimeInDefaultLocationPE(fn, i, location)
	return v
}

// ToDuration casts an interface to a time.Duration type.
func ToDuration(i interface{}) time.Duration {
	v, _ := ToDurationE(i)
	return v
}

// ToDurationP casts an interface to a time.Duration type.
func ToDurationP(fn func(i interface{}) (time.Duration, error), i interface{}) time.Duration {
	v, _ := ToDurationPE(fn, i)
	return v
}

// ToFloat64 casts an interface to a float64 type.
func ToFloat64(i interface{}) float64 {
	v, _ := ToFloat64E(i)
	return v
}

// ToFloat64P casts an interface to a float64 type.
func ToFloat64P(fn func(i interface{}) (float64, error), i interface{}) float64 {
	v, _ := ToFloat64PE(fn, i)
	return v
}

// ToFloat32 casts an interface to a float32 type.
func ToFloat32(i interface{}) float32 {
	v, _ := ToFloat32E(i)
	return v
}

// ToFloat32P casts an interface to a float32 type.
func ToFloat32P(fn func(i interface{}) (float32, error), i interface{}) float32 {
	v, _ := ToFloat32PE(fn, i)
	return v
}

// ToInt64 casts an interface to an int64 type.
func ToInt64(i interface{}) int64 {
	v, _ := ToInt64E(i)
	return v
}

// ToInt64P casts an interface to an int64 type.
func ToInt64P(fn func(i interface{}) (int64, error), i interface{}) int64 {
	v, _ := ToInt64PE(fn, i)
	return v
}

// ToInt32 casts an interface to an int32 type.
func ToInt32(i interface{}) int32 {
	v, _ := ToInt32E(i)
	return v
}

// ToInt32P casts an interface to an int32 type.
func ToInt32P(fn func(i interface{}) (int32, error), i interface{}) int32 {
	v, _ := ToInt32PE(fn, i)
	return v
}

// ToInt16 casts an interface to an int16 type.
func ToInt16(i interface{}) int16 {
	v, _ := ToInt16E(i)
	return v
}

// ToInt16P casts an interface to an int16 type.
func ToInt16P(fn func(i interface{}) (int16, error), i interface{}) int16 {
	v, _ := ToInt16PE(fn, i)
	return v
}

// ToInt8 casts an interface to an int8 type.
func ToInt8(i interface{}) int8 {
	v, _ := ToInt8E(i)
	return v
}

// ToInt8P casts an interface to an int8 type.
func ToInt8P(fn func(i interface{}) (int8, error), i interface{}) int8 {
	v, _ := ToInt8PE(fn, i)
	return v
}

// ToInt casts an interface to an int type.
func ToInt(i interface{}) int {
	v, _ := ToIntE(i)
	return v
}

// ToIntP casts an interface to an int type.
func ToIntP(fn func(i interface{}) (int, error), i interface{}) int {
	v, _ := ToIntPE(fn, i)
	return v
}

// ToUint casts an interface to a uint type.
func ToUint(i interface{}) uint {
	v, _ := ToUintE(i)
	return v
}

// ToUintP casts an interface to a uint type.
func ToUintP(fn func(i interface{}) (uint, error), i interface{}) uint {
	v, _ := ToUintPE(fn, i)
	return v
}

// ToUint64 casts an interface to a uint64 type.
func ToUint64(i interface{}) uint64 {
	v, _ := ToUint64E(i)
	return v
}

// ToUint64P casts an interface to a uint64 type.
func ToUint64P(fn func(i interface{}) (uint64, error), i interface{}) uint64 {
	v, _ := ToUint64PE(fn, i)
	return v
}

// ToUint32 casts an interface to a uint32 type.
func ToUint32(i interface{}) uint32 {
	v, _ := ToUint32E(i)
	return v
}

// ToUint32P casts an interface to a uint32 type.
func ToUint32P(fn func(i interface{}) (uint32, error), i interface{}) uint32 {
	v, _ := ToUint32PE(fn, i)
	return v
}

// ToUint16 casts an interface to a uint16 type.
func ToUint16(i interface{}) uint16 {
	v, _ := ToUint16E(i)
	return v
}

// ToUint16P casts an interface to a uint16 type.
func ToUint16P(fn func(i interface{}) (uint16, error), i interface{}) uint16 {
	v, _ := ToUint16PE(fn, i)
	return v
}

// ToUint8 casts an interface to a uint8 type.
func ToUint8(i interface{}) uint8 {
	v, _ := ToUint8E(i)
	return v
}

// ToUint8P casts an interface to a uint8 type.
func ToUint8P(fn func(i interface{}) (uint8, error), i interface{}) uint8 {
	v, _ := ToUint8PE(fn, i)
	return v
}

// ToString casts an interface to a string type.
func ToString(i interface{}) string {
	v, _ := ToStringE(i)
	return v
}

// ToStringP casts an interface to a string type.
func ToStringP(fn func(i interface{}) (string, error), i interface{}) string {
	v, _ := ToStringPE(fn, i)
	return v
}

// ToStringMapString casts an interface to a map[string]string type.
func ToStringMapString(i interface{}) map[string]string {
	v, _ := ToStringMapStringE(i)
	return v
}

// ToStringMapStringP casts an interface to a map[string]string type.
func ToStringMapStringP(fn func(i interface{}) (map[string]string, error), i interface{}) map[string]string {
	v, _ := ToStringMapStringPE(fn, i)
	return v
}

// ToStringMapStringSlice casts an interface to a map[string][]string type.
func ToStringMapStringSlice(i interface{}) map[string][]string {
	v, _ := ToStringMapStringSliceE(i)
	return v
}

// ToStringMapStringSliceP casts an interface to a map[string][]string type.
func ToStringMapStringSliceP(fn func(i interface{}) (map[string][]string, error), i interface{}) map[string][]string {
	v, _ := ToStringMapStringSlicePE(fn, i)
	return v
}

// ToStringMapBool casts an interface to a map[string]bool type.
func ToStringMapBool(i interface{}) map[string]bool {
	v, _ := ToStringMapBoolE(i)
	return v
}

// ToStringMapBoolP casts an interface to a map[string]bool type.
func ToStringMapBoolP(fn func(i interface{}) (map[string]bool, error), i interface{}) map[string]bool {
	v, _ := ToStringMapBoolPE(fn, i)
	return v
}

// ToStringMapInt casts an interface to a map[string]int type.
func ToStringMapInt(i interface{}) map[string]int {
	v, _ := ToStringMapIntE(i)
	return v
}

// ToStringMapIntP casts an interface to a map[string]int type.
func ToStringMapIntP(fn func(i interface{}) (map[string]int, error), i interface{}) map[string]int {
	v, _ := ToStringMapIntPE(fn, i)
	return v
}

// ToStringMapInt64 casts an interface to a map[string]int64 type.
func ToStringMapInt64(i interface{}) map[string]int64 {
	v, _ := ToStringMapInt64E(i)
	return v
}

// ToStringMapInt64P casts an interface to a map[string]int64 type.
func ToStringMapInt64P(fn func(i interface{}) (map[string]int64, error), i interface{}) map[string]int64 {
	v, _ := ToStringMapInt64PE(fn, i)
	return v
}

// ToStringMap casts an interface to a map[string]interface{} type.
func ToStringMap(i interface{}) map[string]interface{} {
	v, _ := ToStringMapE(i)
	return v
}

// ToStringMapP casts an interface to a map[string]interface{} type.
func ToStringMapP(fn func(i interface{}) (map[string]interface{}, error), i interface{}) map[string]interface{} {
	v, _ := ToStringMapPE(fn, i)
	return v
}

// ToSlice casts an interface to a []interface{} type.
func ToSlice(i interface{}) []interface{} {
	v, _ := ToSliceE(i)
	return v
}

// ToSliceP casts an interface to a []interface{} type.
func ToSliceP(fn func(i interface{}) ([]interface{}, error), i interface{}) []interface{} {
	v, _ := ToSlicePE(fn, i)
	return v
}

// ToBoolSlice casts an interface to a []bool type.
func ToBoolSlice(i interface{}) []bool {
	v, _ := ToBoolSliceE(i)
	return v
}

// ToBoolSliceP casts an interface to a []bool type.
func ToBoolSliceP(fn func(i interface{}) ([]bool, error), i interface{}) []bool {
	v, _ := ToBoolSlicePE(fn, i)
	return v
}

// ToStringSlice casts an interface to a []string type.
func ToStringSlice(i interface{}) []string {
	v, _ := ToStringSliceE(i)
	return v
}

// ToStringSliceP casts an interface to a []string type.
func ToStringSliceP(fn func(i interface{}) ([]string, error), i interface{}) []string {
	v, _ := ToStringSlicePE(fn, i)
	return v
}

// ToIntSlice casts an interface to a []int type.
func ToIntSlice(i interface{}) []int {
	v, _ := ToIntSliceE(i)
	return v
}

// ToIntSliceP casts an interface to a []int type.
func ToIntSliceP(fn func(i interface{}) ([]int, error), i interface{}) []int {
	v, _ := ToIntSlicePE(fn, i)
	return v
}

// ToDurationSlice casts an interface to a []time.Duration type.
func ToDurationSlice(i interface{}) []time.Duration {
	v, _ := ToDurationSliceE(i)
	return v
}

// ToDurationSliceP casts an interface to a []time.Duration type.
func ToDurationSliceP(fn func(i interface{}) ([]time.Duration, error), i interface{}) []time.Duration {
	v, _ := ToDurationSlicePE(fn, i)
	return v
}
