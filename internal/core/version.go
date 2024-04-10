package core

import (
	"os"
	"runtime"
	"runtime/debug"
	"text/template"
)

var VERSION = "dev"

const versionTemplate = `
Tag:        {{ .Tag }}
{{- if .Commit }}
Commit:     {{ .Commit }}
{{- end }}
Go version: {{ .GoVersion }}
OS/Arch:    {{ .Os }}/{{ .Arch }}
`

type Version struct {
	Tag       string `json:"tag,omitempty"`
	Commit    string `json:"commit,omitempty"`
	GoVersion string `json:"goVersion,omitempty"`
	Os        string `json:"os,omitempty"`
	Arch      string `json:"arch,omitempty"`
}

func GetVersion() Version {
	v := Version{
		Tag:       VERSION,
		GoVersion: runtime.Version(),
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

	info, ok := debug.ReadBuildInfo()
	if ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				v.Commit = setting.Value
			case "vcs.modified":
				if setting.Value == "true" && VERSION != "dev" {
					v.Tag += "-dirty"
				}
			}
		}
	}

	return v
}

func (v Version) Print() {
	t := template.Must(template.New("version").Parse(versionTemplate[1:]))
	t.Execute(os.Stdout, v)
}
