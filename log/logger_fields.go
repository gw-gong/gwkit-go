package log

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type field = zap.Field
type ObjectMarshaler = zapcore.ObjectMarshaler

func Err(err error) field {
	return zap.Error(err)
}

func Skip() field {
	return zap.Skip()
}

func Binary(key string, val []byte) field {
	return zap.Binary(key, val)
}

func Bool(key string, val bool) field {
	return zap.Bool(key, val)
}

func Boolp(key string, val *bool) field {
	return zap.Boolp(key, val)
}

func ByteString(key string, val []byte) field {
	return zap.ByteString(key, val)
}

func Complex128(key string, val complex128) field {
	return zap.Complex128(key, val)
}

func Complex128p(key string, val *complex128) field {
	return zap.Complex128p(key, val)
}

func Complex64(key string, val complex64) field {
	return zap.Complex64(key, val)
}

func Complex64p(key string, val *complex64) field {
	return zap.Complex64p(key, val)
}

func Float64(key string, val float64) field {
	return zap.Float64(key, val)
}

func Float64p(key string, val *float64) field {
	return zap.Float64p(key, val)
}

func Float32(key string, val float32) field {
	return zap.Float32(key, val)
}

func Float32p(key string, val *float32) field {
	return zap.Float32p(key, val)
}

func Int(key string, val int) field {
	return zap.Int(key, val)
}

func Intp(key string, val *int) field {
	return zap.Intp(key, val)
}

func Int64(key string, val int64) field {
	return zap.Int64(key, val)
}

func Int64p(key string, val *int64) field {
	return zap.Int64p(key, val)
}

func Int32(key string, val int32) field {
	return zap.Int32(key, val)
}

func Int32p(key string, val *int32) field {
	return zap.Int32p(key, val)
}

func Int16(key string, val int16) field {
	return zap.Int16(key, val)
}

func Int16p(key string, val *int16) field {
	return zap.Int16p(key, val)
}

func Int8(key string, val int8) field {
	return zap.Int8(key, val)
}

func Int8p(key string, val *int8) field {
	return zap.Int8p(key, val)
}

func Str(key string, val string) field {
	return zap.String(key, val)
}

func Strp(key string, val *string) field {
	return zap.Stringp(key, val)
}

func Uint(key string, val uint) field {
	return zap.Uint(key, val)
}

func Uintp(key string, val *uint) field {
	return zap.Uintp(key, val)
}

func Uint64(key string, val uint64) field {
	return zap.Uint64(key, val)
}

func Uint64p(key string, val *uint64) field {
	return zap.Uint64p(key, val)
}

func Uint32(key string, val uint32) field {
	return zap.Uint32(key, val)
}

func Uint32p(key string, val *uint32) field {
	return zap.Uint32p(key, val)
}

func Uint16(key string, val uint16) field {
	return zap.Uint16(key, val)
}

func Uint16p(key string, val *uint16) field {
	return zap.Uint16p(key, val)
}

func Uint8(key string, val uint8) field {
	return zap.Uint8(key, val)
}

func Uint8p(key string, val *uint8) field {
	return zap.Uint8p(key, val)
}

func Uintptr(key string, val uintptr) field {
	return zap.Uintptr(key, val)
}

func Uintptrp(key string, val *uintptr) field {
	return zap.Uintptrp(key, val)
}

func Reflect(key string, val interface{}) field {
	return zap.Reflect(key, val)
}

func NameSpace(key string) field {
	return zap.Namespace(key)
}

func Stringer(key string, val fmt.Stringer) field {
	return zap.Stringer(key, val)
}

func Time(key string, val time.Time) field {
	return zap.Time(key, val)
}

func Timep(key string, val *time.Time) field {
	return zap.Timep(key, val)
}

func Stack(key string) field {
	return zap.Stack(key)
}

func StackSkip(key string, skip int) field {
	return zap.StackSkip(key, skip)
}

func Duration(key string, val time.Duration) field {
	return zap.Duration(key, val)
}

func Durationp(key string, val *time.Duration) field {
	return zap.Durationp(key, val)
}

func Object(key string, val ObjectMarshaler) field {
	return zap.Object(key, val)
}

func Inline(val ObjectMarshaler) field {
	return zap.Inline(val)
}

func Dict(key string, val ...field) field {
	return zap.Dict(key, val...)
}

func Any(key string, val interface{}) field {
	return zap.Any(key, val)
}
