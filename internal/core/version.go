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
{{- if .DirtyWorktree }}
            Dirty Worktree
{{- end }}
{{- end }}
Go version: {{ .GoVersion }}
OS/Arch:    {{ .Os }}/{{ .Arch }}
`

type Version struct {
	Tag           string `json:"tag,omitempty"`
	Commit        string `json:"commit,omitempty"`
	DirtyWorktree bool   `json:"dirtyWorktree,omitempty"`
	GoVersion     string `json:"goVersion,omitempty"`
	Os            string `json:"os,omitempty"`
	Arch          string `json:"arch,omitempty"`
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
				v.DirtyWorktree = true
			}
		}
	}

	return v
}

func (v Version) Print() {
	t := template.Must(template.New("version").Parse(versionTemplate[1:]))
	t.Execute(os.Stdout, v)
}
