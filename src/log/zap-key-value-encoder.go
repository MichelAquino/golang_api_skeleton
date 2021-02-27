package log

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"sync"
	"time"
	"unicode/utf8"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

const _hex = "0123456789abcdef"

var bufferPool = buffer.NewPool()

var _kvPool = sync.Pool{New: func() interface{} {
	return &keyValueEncoder{}
}}

func getKeyValueEncoder() *keyValueEncoder {
	return _kvPool.Get().(*keyValueEncoder)
}

func putKeyValueEncoder(enc *keyValueEncoder) {
	enc.EncoderConfig = nil
	enc.buf = nil
	_kvPool.Put(enc)
}

type keyValueEncoder struct {
	*zapcore.EncoderConfig
	buf *buffer.Buffer
}

// NewKeyValueEncoder creates a key=value encoder
func NewKeyValueEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return &keyValueEncoder{
		EncoderConfig: &cfg,
		buf:           bufferPool.Get(),
	}
}

func (enc *keyValueEncoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	enc.addKey(key)
	return enc.AppendArray(arr)
}

func (enc *keyValueEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	enc.addKey(key)
	return enc.AppendObject(obj)
}

func (enc *keyValueEncoder) AddBinary(key string, val []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (enc *keyValueEncoder) AddByteString(key string, val []byte) {
	enc.addKey(key)
	enc.AppendByteString(val)
}

func (enc *keyValueEncoder) AddBool(key string, val bool) {
	enc.addKey(key)
	enc.AppendBool(val)
}

func (enc *keyValueEncoder) AddComplex128(key string, val complex128) {
	enc.addKey(key)
	enc.AppendComplex128(val)
}

func (enc *keyValueEncoder) AddDuration(key string, val time.Duration) {
	enc.addKey(key)
	enc.AppendDuration(val)
}

func (enc *keyValueEncoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.AppendFloat64(val)
}

func (enc *keyValueEncoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.AppendInt64(val)
}

func (enc *keyValueEncoder) AddReflected(key string, obj interface{}) error {
	marshaled, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	enc.addKey(key)
	_, err = enc.buf.Write(marshaled)
	return err
}

func (enc *keyValueEncoder) OpenNamespace(key string) {
}

func (enc *keyValueEncoder) AddString(key, val string) {
	enc.addKey(key)
	enc.AppendString(val)
}

func (enc *keyValueEncoder) AddTime(key string, val time.Time) {
	enc.addKey(key)
	enc.AppendTime(val)
}

func (enc *keyValueEncoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.AppendUint64(val)
}

func (enc *keyValueEncoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	return arr.MarshalLogArray(enc)
}

func (enc *keyValueEncoder) AppendObject(obj zapcore.ObjectMarshaler) error {
	return obj.MarshalLogObject(enc)
}

func (enc *keyValueEncoder) AppendBool(val bool) {
	enc.buf.AppendBool(val)
}

func (enc *keyValueEncoder) AppendByteString(val []byte) {
	enc.safeAddByteString(val)
}

func (enc *keyValueEncoder) AppendComplex128(val complex128) {
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(val)), float64(imag(val))
	enc.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}

func (enc *keyValueEncoder) AppendDuration(val time.Duration) {
	cur := enc.buf.Len()
	enc.EncodeDuration(val, enc)
	if cur == enc.buf.Len() {
		// User-supplied EncodeDuration is a no-op. Fall back to nanoseconds to keep
		// JSON valid.
		enc.AppendInt64(int64(val))
	}
}

func (enc *keyValueEncoder) AppendInt64(val int64) {
	enc.buf.AppendInt(val)
}

func (enc *keyValueEncoder) AppendReflected(val interface{}) error {
	marshaled, err := json.Marshal(val)
	if err != nil {
		return err
	}
	_, err = enc.buf.Write(marshaled)
	return err
}

func (enc *keyValueEncoder) AppendString(val string) {
	enc.safeAddString(val)
}

func (enc *keyValueEncoder) AppendTime(val time.Time) {
	cur := enc.buf.Len()
	enc.EncodeTime(val, enc)
	if cur == enc.buf.Len() {
		// User-supplied EncodeTime is a no-op. Fall back to nanos since epoch to keep
		// output JSON valid.
		enc.AppendInt64(val.UnixNano())
	}
}

func (enc *keyValueEncoder) AppendUint64(val uint64) {
	enc.buf.AppendUint(val)
}

func (enc *keyValueEncoder) AddComplex64(k string, v complex64) { enc.AddComplex128(k, complex128(v)) }
func (enc *keyValueEncoder) AddFloat32(k string, v float32)     { enc.AddFloat64(k, float64(v)) }
func (enc *keyValueEncoder) AddInt(k string, v int)             { enc.AddInt64(k, int64(v)) }
func (enc *keyValueEncoder) AddInt32(k string, v int32)         { enc.AddInt64(k, int64(v)) }
func (enc *keyValueEncoder) AddInt16(k string, v int16)         { enc.AddInt64(k, int64(v)) }
func (enc *keyValueEncoder) AddInt8(k string, v int8)           { enc.AddInt64(k, int64(v)) }
func (enc *keyValueEncoder) AddUint(k string, v uint)           { enc.AddUint64(k, uint64(v)) }
func (enc *keyValueEncoder) AddUint32(k string, v uint32)       { enc.AddUint64(k, uint64(v)) }
func (enc *keyValueEncoder) AddUint16(k string, v uint16)       { enc.AddUint64(k, uint64(v)) }
func (enc *keyValueEncoder) AddUint8(k string, v uint8)         { enc.AddUint64(k, uint64(v)) }
func (enc *keyValueEncoder) AddUintptr(k string, v uintptr)     { enc.AddUint64(k, uint64(v)) }
func (enc *keyValueEncoder) AppendComplex64(v complex64)        { enc.AppendComplex128(complex128(v)) }
func (enc *keyValueEncoder) AppendFloat64(v float64)            { enc.appendFloat(v, 64) }
func (enc *keyValueEncoder) AppendFloat32(v float32)            { enc.appendFloat(float64(v), 32) }
func (enc *keyValueEncoder) AppendInt(v int)                    { enc.AppendInt64(int64(v)) }
func (enc *keyValueEncoder) AppendInt32(v int32)                { enc.AppendInt64(int64(v)) }
func (enc *keyValueEncoder) AppendInt16(v int16)                { enc.AppendInt64(int64(v)) }
func (enc *keyValueEncoder) AppendInt8(v int8)                  { enc.AppendInt64(int64(v)) }
func (enc *keyValueEncoder) AppendUint(v uint)                  { enc.AppendUint64(uint64(v)) }
func (enc *keyValueEncoder) AppendUint32(v uint32)              { enc.AppendUint64(uint64(v)) }
func (enc *keyValueEncoder) AppendUint16(v uint16)              { enc.AppendUint64(uint64(v)) }
func (enc *keyValueEncoder) AppendUint8(v uint8)                { enc.AppendUint64(uint64(v)) }
func (enc *keyValueEncoder) AppendUintptr(v uintptr)            { enc.AppendUint64(uint64(v)) }

func (enc *keyValueEncoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	clone.buf.Write(enc.buf.Bytes())
	return clone
}

func (enc *keyValueEncoder) clone() *keyValueEncoder {
	clone := getKeyValueEncoder()
	clone.EncoderConfig = enc.EncoderConfig
	clone.buf = bufferPool.Get()
	return clone
}

func (enc *keyValueEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	final := enc.clone()

	if final.LevelKey != "" {
		final.addKey(final.LevelKey)
		cur := final.buf.Len()
		final.EncodeLevel(ent.Level, final)
		if cur == final.buf.Len() {
			final.AppendString(ent.Level.String())
		}
		final.addElementSeparator()
	}
	if final.TimeKey != "" {
		final.AddTime(final.TimeKey, ent.Time)
		final.addElementSeparator()
	}
	if ent.LoggerName != "" && final.NameKey != "" {
		final.addKey(final.NameKey)
		cur := final.buf.Len()
		nameEncoder := final.EncodeName

		// if no name encoder provided, fall back to FullNameEncoder for backwards
		// compatibility
		if nameEncoder == nil {
			nameEncoder = zapcore.FullNameEncoder
		}

		nameEncoder(ent.LoggerName, final)
		if cur == final.buf.Len() {
			// User-supplied EncodeName was a no-op. Fall back to strings to
			// keep output valid.
			final.AppendString(ent.LoggerName)
		}
		final.addElementSeparator()
	}
	if ent.Caller.Defined && final.CallerKey != "" {
		final.addKey(final.CallerKey)
		cur := final.buf.Len()
		final.EncodeCaller(ent.Caller, final)
		if cur == final.buf.Len() {
			// User-supplied EncodeCaller was a no-op. Fall back to strings to
			// keep JSON valid.
			final.AppendString(ent.Caller.String())
		}
		final.addElementSeparator()
	}
	if final.MessageKey != "" {
		final.addKey(enc.MessageKey)
		final.buf.AppendByte('"')
		final.AppendString(ent.Message)
		final.buf.AppendByte('"')
		final.addElementSeparator()
	}
	if enc.buf.Len() > 0 {
		final.buf.Write(enc.buf.Bytes())
	}
	addFields(final, final, fields)
	final.addElementSeparator()
	if ent.Stack != "" && final.StacktraceKey != "" {
		final.AddString(final.StacktraceKey, ent.Stack)
		final.addElementSeparator()
	}
	if final.LineEnding != "" {
		final.buf.AppendString(final.LineEnding)
	} else {
		final.buf.AppendString(zapcore.DefaultLineEnding)
	}

	ret := final.buf
	putKeyValueEncoder(final)
	return ret, nil
}

func (enc *keyValueEncoder) addKey(key string) {
	enc.buf.AppendString(key)
	enc.buf.AppendByte('=')
}

func (enc *keyValueEncoder) addElementSeparator() {
	enc.buf.AppendByte(' ')
}

func (enc *keyValueEncoder) appendFloat(val float64, bitSize int) {
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}

// safeAddString JSON-escapes a string and appends it to the internal buffer.
// Unlike the standard library's encoder, it doesn't attempt to protect the
// user from browser vulnerabilities or JSONP-related problems.
func (enc *keyValueEncoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.AppendString(s[i : i+size])
		i += size
	}
}

// safeAddByteString is no-alloc equivalent of safeAddString(string(s)) for s []byte.
func (enc *keyValueEncoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.Write(s[i : i+size])
		i += size
	}
}

// tryAddRuneSelf appends b if it is valid UTF-8 character represented in a single byte.
func (enc *keyValueEncoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.AppendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte(b)
	case '\n':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('n')
	case '\r':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('r')
	case '\t':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf.AppendString(`\u00`)
		enc.buf.AppendByte(_hex[b>>4])
		enc.buf.AppendByte(_hex[b&0xF])
	}
	return true
}

func (enc *keyValueEncoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf.AppendString(`\ufffd`)
		return true
	}
	return false
}

func addFields(kvEnc *keyValueEncoder, enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
		kvEnc.buf.AppendByte(' ')
	}
}
