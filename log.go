package luabox

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type Log interface {
	GetLevel() Level
	WithFields(context map[string]interface{}) Log
	Debug(msg string, context map[string]interface{})
	Info(msg string, context map[string]interface{})
	Warn(msg string, context map[string]interface{})
	Error(msg string, context map[string]interface{})
	Fatal(msg string, context map[string]interface{})
}
