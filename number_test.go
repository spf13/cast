package cast_test

import (
	"math"
	"testing"

	"github.com/spf13/cast"
)

type numberContext struct {
	specific func(any) (any, error)
	generic  func(any) (any, error)

	// Order of samples:
	// zero, one, 8, -8, 8.3, -8.3, min, max
	samples []any
}

func toAny[T cast.Number](fn func(any) (T, error)) func(i any) (any, error) {
	return func(i any) (any, error) { return fn(i) }
}

var numberContexts = map[string]numberContext{
	"int": {
		specific: toAny(cast.ToIntE),
		generic:  toAny(cast.ToNumberE[int]),
		samples:  []any{int(0), int(1), int(8), int(-8), int(8), int(-8), math.MinInt, math.MaxInt},
	},
	"int8": {
		specific: toAny(cast.ToInt8E),
		generic:  toAny(cast.ToNumberE[int8]),
		samples:  []any{int8(0), int8(1), int8(8), int8(-8), int8(8), int8(-8), math.MinInt8, math.MaxInt8},
	},
	"int16": {
		specific: toAny(cast.ToInt16E),
		generic:  toAny(cast.ToNumberE[int16]),
		samples:  []any{int16(0), int16(1), int16(8), int16(-8), int16(8), int16(-8), math.MinInt16, math.MaxInt16},
	},
	"int32": {
		specific: toAny(cast.ToInt32E),
		generic:  toAny(cast.ToNumberE[int32]),
		samples:  []any{int32(0), int32(1), int32(8), int32(-8), int32(8), int32(-8), math.MinInt32, math.MaxInt32},
	},
	"int64": {
		specific: toAny(cast.ToInt64E),
		generic:  toAny(cast.ToNumberE[int64]),
		samples:  []any{int64(0), int64(1), int64(8), int64(-8), int64(8), int64(-8), math.MinInt64, math.MaxInt64},
	},
	"uint": {
		specific: toAny(cast.ToUintE),
		generic:  toAny(cast.ToNumberE[uint]),
		samples:  []any{uint(0), uint(1), uint(8), uint(0), uint(8), uint(0), uint(0), uint(math.MaxUint)},
	},
	"uint8": {
		specific: toAny(cast.ToUint8E),
		generic:  toAny(cast.ToNumberE[uint8]),
		samples:  []any{uint8(0), uint8(1), uint8(8), uint8(0), uint8(8), uint8(0), uint8(0), uint8(math.MaxUint8)},
	},
	"uint16": {
		specific: toAny(cast.ToUint16E),
		generic:  toAny(cast.ToNumberE[uint16]),
		samples:  []any{uint16(0), uint16(1), uint16(8), uint16(0), uint16(8), uint16(0), uint16(0), uint16(math.MaxUint16)},
	},
	"uint32": {
		specific: toAny(cast.ToUint32E),
		generic:  toAny(cast.ToNumberE[uint32]),
		samples:  []any{uint32(0), uint32(1), uint32(8), uint32(0), uint32(8), uint32(0), uint32(0), uint32(math.MaxUint32)},
	},
	"uint64": {
		specific: toAny(cast.ToUint64E),
		generic:  toAny(cast.ToNumberE[uint64]),
		samples:  []any{uint64(0), uint64(1), uint64(8), uint64(0), uint64(8), uint64(0), uint64(0), uint64(math.MaxUint64)},
	},
	"float32": {
		specific: toAny(cast.ToFloat32E),
		generic:  toAny(cast.ToNumberE[float32]),
		samples:  []any{float32(0), float32(1), float32(8), float32(0), float32(8.3), float32(-8.3), float32(-math.MaxFloat32), float32(math.MaxFloat32)},
	},
	"float64": {
		specific: toAny(cast.ToFloat64E),
		generic:  toAny(cast.ToNumberE[float64]),
		samples:  []any{float64(0), float64(1), float64(8), float64(0), float64(8.3), float64(-8.3), float64(-math.MaxFloat64), float64(math.MaxFloat64)},
	},
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
	for typeName, ctx := range numberContexts {
		cases = append(
			cases,
			testCase{
				name:     typeName,
				input:    "123",
				specific: ctx.specific,
				generic:  ctx.generic,
			},
			testCase{
				name:     typeName,
				input:    "1234567890123",
				specific: ctx.specific,
				generic:  ctx.generic,
			},
			testCase{
				name:     typeName,
				input:    "-123",
				specific: ctx.specific,
				generic:  ctx.generic,
			},
			testCase{
				name:     typeName,
				input:    "-1234567890123",
				specific: ctx.specific,
				generic:  ctx.generic,
			},
			testCase{
				name:     typeName,
				input:    "0000000000123",
				specific: ctx.specific,
				generic:  ctx.generic,
			},
			testCase{
				name:     typeName,
				input:    "00000000001234567890123",
				specific: ctx.specific,
				generic:  ctx.generic,
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
