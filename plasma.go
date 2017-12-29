package warpcache

import (
	"strconv"
	"strings"
)

func parseInputFormat(message string) (string, map[string]string, float64, error) {

	message = strings.Replace(message, "\n", "", -1)

	var classname string
	labels := make(map[string]string)
	var value float64
	var err error

	elts := strings.Split(message, " ")
	value, err = strconv.ParseFloat(elts[2], 64)
	if err != nil {
		return classname, labels, value, err
	}
	classname = strings.Split(elts[1], "{")[0]

	selector := strings.Split(elts[1], "{")[1]
	selector = selector[0 : len(selector)-1]

	l := strings.Split(selector, ",")
	for _, la := range l {
		lab := strings.Split(la, "=")
		if lab[0] != ".app" {
			labels[lab[0]] = lab[1]
		}
	}

	return classname, labels, value, err
}
