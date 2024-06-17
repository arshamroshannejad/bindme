package bindme

import (
	"net/url"
	"strconv"
	"strings"
)

func (v *Validator) ReadString(qs url.Values, key, defaultValue string) string {
	k := qs.Get(key)
	if k == "" {
		return defaultValue
	}
	return k
}

func (v *Validator) ReadCSV(qs url.Values, key string, defaultValue []string) []string {
	k := qs.Get(key)
	if k == "" {
		return defaultValue
	}
	return strings.Split(k, ",")
}

func (v *Validator) ReadInt(qs url.Values, key string, defaultValue int) int {
	k := qs.Get(key)
	if k == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(k)
	if err != nil {
		v.Add(key, "must be an integer value")
		return defaultValue
	}
	return i
}
