package bindme

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func ReadString(qs url.Values, key, defaultValue string) string {
	v := qs.Get(key)
	if v == "" {
		return defaultValue
	}
	return v
}

func ReadCSV(qs url.Values, key string, defaultValue []string) []string {
	v := qs.Get(key)
	if v == "" {
		return defaultValue
	}
	return strings.Split(v, ",")
}

func ReadInt(qs url.Values, key string, defaultValue int) (int, error) {
	v := qs.Get(key)
	if v == "" {
		return defaultValue, nil
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue, fmt.Errorf("key %s: must be an integer", key)
	}
	return i, nil
}
