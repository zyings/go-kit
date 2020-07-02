
package config

import (
	"net"
	"sync/atomic"
	"time"
)

type ReadOnlyConfig interface {
	Get(key string) (interface{}, error)
	Unmarshal(v interface{}) error
	Sub(key string) ReadOnlyConfig
	SubArr(key string) ([]ReadOnlyConfig, error)
	SubMap(key string) (map[string]ReadOnlyConfig, error)
	Bind(key string, v interface{}, opts ...BindOption) *atomic.Value
	BindVar(key string, v interface{}, av *atomic.Value, opts ...BindOption)
	
	Bool(key string, opts ...BindOption) *AtomicBool
	GetBoolD(key string, dftVal bool) bool
	GetBoolE(key string) (bool, error)
	GetBool(key string) bool
	GetBoolP(key string) bool

	Int(key string, opts ...BindOption) *AtomicInt
	GetIntD(key string, dftVal int) int
	GetIntE(key string) (int, error)
	GetInt(key string) int
	GetIntP(key string) int

	Uint(key string, opts ...BindOption) *AtomicUint
	GetUintD(key string, dftVal uint) uint
	GetUintE(key string) (uint, error)
	GetUint(key string) uint
	GetUintP(key string) uint

	Int64(key string, opts ...BindOption) *AtomicInt64
	GetInt64D(key string, dftVal int64) int64
	GetInt64E(key string) (int64, error)
	GetInt64(key string) int64
	GetInt64P(key string) int64

	Int32(key string, opts ...BindOption) *AtomicInt32
	GetInt32D(key string, dftVal int32) int32
	GetInt32E(key string) (int32, error)
	GetInt32(key string) int32
	GetInt32P(key string) int32

	Int16(key string, opts ...BindOption) *AtomicInt16
	GetInt16D(key string, dftVal int16) int16
	GetInt16E(key string) (int16, error)
	GetInt16(key string) int16
	GetInt16P(key string) int16

	Int8(key string, opts ...BindOption) *AtomicInt8
	GetInt8D(key string, dftVal int8) int8
	GetInt8E(key string) (int8, error)
	GetInt8(key string) int8
	GetInt8P(key string) int8

	Uint64(key string, opts ...BindOption) *AtomicUint64
	GetUint64D(key string, dftVal uint64) uint64
	GetUint64E(key string) (uint64, error)
	GetUint64(key string) uint64
	GetUint64P(key string) uint64

	Uint32(key string, opts ...BindOption) *AtomicUint32
	GetUint32D(key string, dftVal uint32) uint32
	GetUint32E(key string) (uint32, error)
	GetUint32(key string) uint32
	GetUint32P(key string) uint32

	Uint16(key string, opts ...BindOption) *AtomicUint16
	GetUint16D(key string, dftVal uint16) uint16
	GetUint16E(key string) (uint16, error)
	GetUint16(key string) uint16
	GetUint16P(key string) uint16

	Uint8(key string, opts ...BindOption) *AtomicUint8
	GetUint8D(key string, dftVal uint8) uint8
	GetUint8E(key string) (uint8, error)
	GetUint8(key string) uint8
	GetUint8P(key string) uint8

	Float64(key string, opts ...BindOption) *AtomicFloat64
	GetFloat64D(key string, dftVal float64) float64
	GetFloat64E(key string) (float64, error)
	GetFloat64(key string) float64
	GetFloat64P(key string) float64

	Float32(key string, opts ...BindOption) *AtomicFloat32
	GetFloat32D(key string, dftVal float32) float32
	GetFloat32E(key string) (float32, error)
	GetFloat32(key string) float32
	GetFloat32P(key string) float32

	String(key string, opts ...BindOption) *AtomicString
	GetStringD(key string, dftVal string) string
	GetStringE(key string) (string, error)
	GetString(key string) string
	GetStringP(key string) string

	Duration(key string, opts ...BindOption) *AtomicDuration
	GetDurationD(key string, dftVal time.Duration) time.Duration
	GetDurationE(key string) (time.Duration, error)
	GetDuration(key string) time.Duration
	GetDurationP(key string) time.Duration

	Time(key string, opts ...BindOption) *AtomicTime
	GetTimeD(key string, dftVal time.Time) time.Time
	GetTimeE(key string) (time.Time, error)
	GetTime(key string) time.Time
	GetTimeP(key string) time.Time

	IP(key string, opts ...BindOption) *AtomicIP
	GetIPD(key string, dftVal net.IP) net.IP
	GetIPE(key string) (net.IP, error)
	GetIP(key string) net.IP
	GetIPP(key string) net.IP

}
