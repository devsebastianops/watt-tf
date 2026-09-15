package data

const (
	DataTypeFile = "file"
)

type Data interface {
	Load() (map[string]any, error)
	SourcePath() string
	LoadAt(path string) (map[string]any, error)
}

func NewData(dataType, path string) Data {
	switch dataType {
	case DataTypeFile:
		return FileData{
			Type: DataTypeFile,
			Path: path,
		}
	default:
		return nil
	}
}
