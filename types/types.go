package types

const (
	DefaultPluginPath="/var/lib/"
)


type Plugin interface {
	Load()(string,int32)
	UnLoad()
	AddFunction(name string, function interface{}) error
	InvokeFunction(name string, args []interface{}) ([]interface{}, error)
}