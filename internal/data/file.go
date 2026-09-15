package data

import "github.com/devsebastianops/x/parser"

type FileData struct {
	Type string `yaml:"type"`
	Path string `yaml:"path"`
}

func (f FileData) Load() (map[string]any, error) {
	return f.LoadAt(f.Path)
}

func (f FileData) SourcePath() string {
	return f.Path
}

func (f FileData) LoadAt(path string) (map[string]any, error) {
	parser, err := parser.NewParser(path)
	if err != nil {
		return nil, err
	}

	data, err := parser.Parse(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}
