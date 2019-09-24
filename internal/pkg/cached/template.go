package cached

import (
	"bytes"
	"sync"
	"text/template"
)

// Template caches compiled templates by content
type Template struct {
	cache sync.Map
}

// Get returns cached parsed template
func (t *Template) Get(contents string) (*template.Template, error) {
	if tmpl, ok := t.cache.Load(contents); ok {
		return tmpl.(*template.Template), nil
	}

	tmpl, err := template.New("").Parse(contents)
	if err != nil {
		return nil, err
	}

	t.cache.Store(contents, tmpl)

	return tmpl, nil
}

// StringTemplate renders cached template to string
type StringTemplate struct {
	Template
}

// Execute cached template and return result as string
func (t *StringTemplate) Execute(contents string, context interface{}) (string, error) {
	tmpl, err := t.Get(contents)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, context); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// MustExecute panics on error in Execute
func (t *StringTemplate) MustExecute(contents string, context interface{}) string {
	result, err := t.Execute(contents, context)
	if err != nil {
		panic(err)
	}

	return result
}
