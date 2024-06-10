package main

import (
	"net/url"
	"strconv"
	"strings"
)

// The readString() helper returns a string value from the query string, or the provided
// default value if no matching key could be found.
func (app *application) readString(queryString url.Values, key, defaultValue string) string {
	value := queryString.Get(key)
	if value == "" {
		return defaultValue
	} else {
		return value
	}
}

// The readCSV() helper reads a string value from the query string and then splits it
// into a slice on the comma character. If no matching key could be found, it returns
// the provided default value.
func (app *application) readCSV(queryString url.Values, key string, defaultValue []string) []string {
	csv := queryString.Get(key)
	if csv == "" {
		return defaultValue
	} else {
		return strings.Split(csv, ",")
	}
}

// The readInt() helper reads a string value from the query string and converts it to an
// integer before returning. If no matching key could be found it returns the provided
// default value. If the value couldn't be converted to an integer, then we record an
// error messag
func (app *application) readInt(queryString url.Values, key string, defaultValue int) int {
	val := queryString.Get(key)
	if val == "" {
		return defaultValue
	} else {
		num, err := strconv.Atoi(val)
		if err != nil {
			app.errorLogger.Println("query string: must have an interger value", err)
			return defaultValue
		} else {
			return num
		}
	}
}
