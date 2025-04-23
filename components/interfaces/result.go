package interfaces

type IResult interface {
	GetPackage() string
	GetStatus() string
	GetError() string
	GetAstFile() string
	GetInfo() interface{}

	SetPackage(string)
	SetStatus(string)
	SetError(string)

	ToJSON(outputTarget string) string
	ToXML(outputTarget string) string
	ToCSV(outputTarget string) string

	ToMap() map[string]interface{}

	DataTable() error
	GetStatusCode() int
	GetStatusText() string
	GetAst() interface{}
}
