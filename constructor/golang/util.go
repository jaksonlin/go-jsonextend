package golang

import (
	"html"
	"reflect"
	"strconv"
	"unicode/utf8"
)

// string to number receiver
func convertStringToNumericReceiver(receiver reflect.Value, value string) (reflect.Value, error) {
	// reflect.Value is struct not pointer return a new one
	switch receiver.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		val, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(val).Convert(receiver.Type()), err
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(value, 10, 64)
		return reflect.ValueOf(val).Convert(receiver.Type()), err
	default:
		return reflect.Value{}, ErrorUnsupportedDataKind
	}
}

// number reflect.value to string reflect.value
func convertNumericToString(value reflect.Value) (string, error) {
	switch value.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		return strconv.FormatInt(value.Int(), 10), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10), nil
	case reflect.Float32:
		return strconv.FormatFloat(value.Float(), 'f', -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64), nil
	default:
		return "", ErrorUnsupportedDataKind
	}
}

func getMemoryAddress(v reflect.Value) uintptr {
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return v.Pointer()
	case reflect.Interface:
		if !v.IsNil() && v.Elem().CanAddr() {
			return v.Elem().UnsafeAddr()
		}
		return v.UnsafeAddr()
	default:
		return v.UnsafeAddr()
	}
}

// repair string as json standard request
func repairUTF8(s string) string {
	if utf8.ValidString(s) {
		return s // Already valid UTF-8.
	}

	var repaired []rune
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		repaired = append(repaired, r)
		s = s[size:]
	}

	return string(repaired)
}

func htmlEscape(s string) string {
	return html.EscapeString(s)
}

// json input value is always float64, convert to different numeric value based on out element kind
func convertNumberBaseOnKind(val reflect.Value) (float64, error) {

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(val.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(val.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return val.Float(), nil
	default:
		return 0.0, ErrNotNumericValueField
	}
}
