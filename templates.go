package warpcache

import (
	"bytes"
	"text/template"
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

var templateFetchSingle = `
'{{.ReadToken}}' 'token' STORE
'{{.Selector}}' PARSESELECTOR 'labels' STORE 'classname' STORE
[ $token $classname $labels NOW -1 ] FETCH 0 GET VALUES 0 GET`

type templateFetchSingleData struct {
	ReadToken string
	Selector  string
}

func generateFetchSingleWarpScript(token string, selector string) (string, error) {
	tmpl, err := template.New("FetchSingle").Parse(templateFetchSingle)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, templateFetchSingleData{ReadToken: token, Selector: selector})

	return tpl.String(), nil
}

var templateFetchMultiple = `
'{{.ReadToken}}' 'token' STORE
'{{.Pivot}}' 'pivot' STORE
'{{.Selector}}' PARSESELECTOR 'labels' STORE 'classname' STORE
[ $token $classname $labels NOW -1 ] FETCH 
<%
    DUP LABELS $pivot GET 
    SWAP VALUES 0 GET
%> FOREACH
DEPTH ->MAP
`

type templateFetchMultipleData struct {
	ReadToken string
	Selector  string
	Pivot     string
}

func generateFetchMultipleWarpScript(token string, selector string, pivot string) (string, error) {
	tmpl, err := template.New("FetchMultiple").Parse(templateFetchMultiple)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, templateFetchMultipleData{ReadToken: token, Selector: selector, Pivot: pivot})

	return tpl.String(), nil
}
