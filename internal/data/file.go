package data

import "github.com/devsebastianops/x/parser"

type FileData struct {
	Type string `yaml:"type"`
	Path string `yaml:"path"`
}

func (f FileData) Load() (map[string]any, error) {
	parser, err := parser.NewParser(f.Path)
	if err != nil {
		return nil, err
	}

	data, err := parser.Parse(f.Path)
	if err != nil {
		return nil, err
	}
	return data, nil
}
