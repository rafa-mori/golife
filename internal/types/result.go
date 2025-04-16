package types

type IResult interface {
	GetPackage() string
	GetStatus() string
	GetError() string

	SetPackage(string)
	SetStatus(string)
	SetError(string)

	ToJSON(outputTarget string) string
	ToXML(outputTarget string) string
	ToCSV(outputTarget string) string

	ToMap() map[string]interface{}

	DataTable() error
}
