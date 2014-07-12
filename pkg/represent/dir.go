// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package represent

import (
	"code.google.com/p/go.tools/present"
	"fmt"
	"io"
	"os"
	"bytes"
	"os/exec"
	"html/template"
	"path/filepath"
)

// extensions maps the presentable file extensions to the name of the
// template to be executed.
var extensions = map[string]string{
	".slide":   "slides.tmpl",
	".article": "article.tmpl",
}

// isDoc tests if the path is to a file with a Present format file extension.
func isDoc(path string) bool {
	_, ok := extensions[filepath.Ext(path)]
	return ok
}

// parse reads the given file path and parses into a Present document structure.
func parse(name string, mode present.ParseMode) (*present.Doc, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := EolConvert(f, LF)
	return present.Parse(r, name, 0)
}

func pygmentizeFile(content []byte, ext string) ([]byte, error) {
	cmd := exec.Command("pygmentize", "-f", "html", "-l", ext[1:])
	cmd.Stdin = bytes.NewReader(content)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.Bytes(), err
}

func pygmentizeDoc(doc *present.Doc) error {
	for _, section := range doc.Sections {
		for i, elem := range section.Elem {
			if code, ok := elem.(present.Code); ok {
				// we can not use code.FileName as it is not absolute
				text, err := pygmentizeFile(code.Raw, code.Ext)
				if err != nil {
					fmt.Println(err, string(text))
					continue
				}
				// fmt.Println("text", string(text), code.Text)
				code.Text = template.HTML(text)
				section.Elem[i] = code
			}
		}
	}
	return nil
}

// renderDoc reads the present file, builds its template representation,
// and executes the template, sending output to w.
func renderDoc(w io.Writer, base, docFile string) error {
	// Read the input and build the doc structure.
	doc, err := parse(docFile, 0)
	if err != nil {
		return err
	}

	err = pygmentizeDoc(doc)
	if err != nil {
		return err
	}

	// Find which template should be executed.
	ext := filepath.Ext(docFile)
	contentTmpl, ok := extensions[ext]
	if !ok {
		return fmt.Errorf("no template for extension %v", ext)
	}

	// Locate the template file.
	actionTmpl := filepath.Join(base, "templates/action.tmpl")
	contentTmpl = filepath.Join(base, "templates", contentTmpl)

	// Read and parse the input.
	tmpl := present.Template()
	if _, err := tmpl.ParseFiles(actionTmpl, contentTmpl); err != nil {
		return err
	}

	// Execute the template.
	return doc.Render(w, tmpl)
}
