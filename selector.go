package warpcache

import (
	"fmt"
	"net/url"
	"strings"
)

// Selector is a warp10 selector
type Selector struct {
	Classname string
	Labels    map[string]string
}

// String is returning the string version of a Warp10 Selector
func (s *Selector) String() string {
	return fmt.Sprintf("%s{%s}", s.Classname, s.getLabels())
}

func (s *Selector) getLabels() string {

	var tmp string
	for key, value := range s.Labels {

		tmp = tmp + url.QueryEscape(key) + "=" + url.QueryEscape(value) + ","
	}
	// Removing last comma
	tmp = strings.TrimSuffix(tmp, ",")
	return tmp
}
