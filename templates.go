package warpcache

import (
	"bytes"
	"html/template"
)

var templateFindSize = `
'{{.ReadToken}}' 'token' STORE
'{{.Selector}}' PARSESELECTOR 'labels' STORE 'classname' STORE
[ $token $classname $labels ] FIND SIZE`

type templateFindSizeData struct {
	ReadToken string
	Selector  string
}

func generateFindSizeWarpScript(token string, selector string) (string, error) {
	tmpl, err := template.New("findSize").Parse(templateFindSize)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, templateFindSizeData{ReadToken: token, Selector: selector})

	return tpl.String(), nil
}
