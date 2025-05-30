package cast_test

import (
	"testing"

	"github.com/spf13/cast"
)

type numberFuncs struct {
	specific func(any) (any, error)
	generic  func(any) (any, error)
}

func toAny[T cast.Number](fn func(any) (T, error)) func(i any) (any, error) {
	return func(i any) (any, error) { return fn(i) }
}

var numberCasters = map[string]numberFuncs{
	"int":     {toAny(cast.ToIntE), toAny(cast.ToNumberE[int])},
	"int8":    {toAny(cast.ToInt8E), toAny(cast.ToNumberE[int8])},
	"int16":   {toAny(cast.ToInt16E), toAny(cast.ToNumberE[int16])},
	"int32":   {toAny(cast.ToInt32E), toAny(cast.ToNumberE[int32])},
	"int64":   {toAny(cast.ToInt64E), toAny(cast.ToNumberE[int64])},
	"uint":    {toAny(cast.ToUintE), toAny(cast.ToNumberE[uint])},
	"uint8":   {toAny(cast.ToUint8E), toAny(cast.ToNumberE[uint8])},
	"uint16":  {toAny(cast.ToUint16E), toAny(cast.ToNumberE[uint16])},
	"uint32":  {toAny(cast.ToUint32E), toAny(cast.ToNumberE[uint32])},
	"uint64":  {toAny(cast.ToUint64E), toAny(cast.ToNumberE[uint64])},
	"float32": {toAny(cast.ToFloat32E), toAny(cast.ToNumberE[float32])},
	"float64": {toAny(cast.ToFloat64E), toAny(cast.ToNumberE[float64])},
}

func BenchmarkNumber(b *testing.B) {
	type testCase struct {
		name     string
		input    any
		specific func(any) (any, error)
		generic  func(any) (any, error)
	}

	var cases []testCase

	// TODO: sort keys before iterating (once Go version is updated)
	for typeName, caster := range numberCasters {
		cases = append(
			cases,
			testCase{
				name:     typeName,
				input:    "123",
				specific: caster.specific,
				generic:  caster.generic,
			},
			testCase{
				name:     typeName,
				input:    "1234567890123",
				specific: caster.specific,
				generic:  caster.generic,
			},
			testCase{
				name:     typeName,
				input:    "-123",
				specific: caster.specific,
				generic:  caster.generic,
			},
			testCase{
				name:     typeName,
				input:    "-1234567890123",
				specific: caster.specific,
				generic:  caster.generic,
			},
			testCase{
				name:     typeName,
				input:    "0000000000123",
				specific: caster.specific,
				generic:  caster.generic,
			},
			testCase{
				name:     typeName,
				input:    "00000000001234567890123",
				specific: caster.specific,
				generic:  caster.generic,
			},
		)
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			b.Run("Specific", func(b *testing.B) {
				// TODO: use b.Loop() once updated to Go 1.24
				for i := 0; i < b.N; i++ {
					_, _ = c.specific(c.input)
				}
			})

			b.Run("Generic", func(b *testing.B) {
				// TODO: use b.Loop() once updated to Go 1.24
				for i := 0; i < b.N; i++ {
					_, _ = c.generic(c.input)
				}
			})
		})
	}
}
