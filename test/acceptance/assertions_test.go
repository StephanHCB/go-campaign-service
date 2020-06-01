package acceptance

import (
	"fmt"
	"github.com/smartystreets/assertions"
	"regexp"
)

// a better name for the GoConvey ShouldResemble assertion, which actually checks deep equality
var ShouldDeepEqual = assertions.ShouldResemble

// a custom assertion for GoConvey (can't believe they don't have this already?)
func shouldMatchRegex(actual interface{}, expected ...interface{}) string {
	match := matchRegexp(expected[0], actual)

	if match {
		return ""
	}
	return fmt.Sprintf("Assertion failed, actual value '%v' did not match regex '%v'", actual, expected[0])
}

// matchRegexp returns true if a specified regexp matches a string.
func matchRegexp(rx interface{}, str interface{}) bool {
	var r *regexp.Regexp
	var err error
	if rr, ok := rx.(*regexp.Regexp); ok {
		r = rr
	} else {
		r, err = regexp.Compile(fmt.Sprint(rx))
		if err != nil {
			return false
		}
	}
	return (r.FindStringIndex(fmt.Sprint(str)) != nil)
}
