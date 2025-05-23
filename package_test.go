// package_test.go - Tests for package.go
//
// Copyright (c) 2013 The go-alpm Authors
//
// MIT Licensed. See LICENSE for details.

package alpm_test

import (
	"bytes"
	"testing"
	"text/template"
	"time"

	alpm "github.com/Jguer/go-alpm/v2"
	"github.com/stretchr/testify/assert"
)

// Auxiliary formatting
const pkginfoTemplate = `
Name         : {{ .Name }}
Version      : {{ .Version }}
Architecture : {{ .Architecture }}
Description  : {{ .Description }}
URL          : {{ .URL }}
Groups       : {{ .Groups.Slice }}
Licenses     : {{ .Licenses.Slice }}
Dependencies : {{ range .Depends.Slice }}{{ . }} {{ end }}
Provides     : {{ range .Provides.Slice }}{{ . }} {{ end }}
Replaces     : {{ range .Replaces.Slice }}{{ . }} {{ end }}
Conflicts    : {{ range .Conflicts.Slice }}{{ . }} {{ end }}
Packager     : {{ .Packager }}
Build Date   : {{ .PrettyBuildDate }}
Install Date : {{ .PrettyInstallDate }}
Package Size : {{ .Size }} bytes
Install Size : {{ .ISize }} bytes
MD5 Sum      : {{ .MD5Sum }}
SHA256 Sum   : {{ .SHA256Sum }}
Reason       : {{ .Reason }}

Required By  : {{ .ComputeRequiredBy }}
Files        : {{ range .Files }}
               {{ .Name }} {{ .Size }}{{ end }}
`

type PrettyPackage struct {
	*alpm.Package
}

func (p PrettyPackage) PrettyBuildDate() string {
	return p.Package.BuildDate().Format(time.RFC1123)
}

func (p PrettyPackage) PrettyInstallDate() string {
	return p.Package.InstallDate().Format(time.RFC1123)
}

// Tests package attribute getters.
func TestPkginfo(t *testing.T) {
	t.Parallel()
	pkginfoTemp, er := template.New("info").Parse(pkginfoTemplate)
	assert.NoError(t, er, "couldn't compile template")

	h, er := alpm.Initialize(root, dbpath)
	defer h.Release()
	if er != nil {
		t.Errorf("Failed at alpm initialization: %s", er)
	}

	db, _ := h.LocalDB()

	pkg := db.Pkg("glibc")
	buf := bytes.NewBuffer(nil)
	pkginfoTemp.Execute(buf, PrettyPackage{pkg.(*alpm.Package)})

	pkg = db.Pkg("linux")
	if pkg != nil {
		buf = bytes.NewBuffer(nil)
		pkginfoTemp.Execute(buf, PrettyPackage{pkg.(*alpm.Package)})
	}
}

func TestPkgNoExist(t *testing.T) {
	t.Parallel()
	h, er := alpm.Initialize(root, dbpath)
	defer h.Release()
	if er != nil {
		t.Errorf("Failed at alpm initialization: %s", er)
	}

	db, _ := h.LocalDB()

	pkg := db.Pkg("non-existing-package-fa93f4af")
	if pkg != nil {
		t.Errorf("pkg should be nil but got %v", pkg)
	}
}

func TestPkgFiles(t *testing.T) {
	t.Parallel()
	h, er := alpm.Initialize(root, dbpath)
	defer h.Release()
	if er != nil {
		t.Errorf("Failed at alpm initialization: %s", er)
	}

	db, _ := h.LocalDB()

	pkg := db.Pkg("glibc")
	_, err := pkg.ContainsFile("etc/locale.gen")
	if err != nil {
		t.Errorf("File should not be nil but got %v", err)
	}
	_, err = pkg.ContainsFile("etc/does-not-exist")
	if err == nil {
		t.Errorf("File should be nil but got %v", err)
	}
}
