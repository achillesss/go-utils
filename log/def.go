package log

const (
	// spilts
	slash = "/"
	point = "."
	// print type
	logInformation = "I"
	logWarning     = "W"
	logError       = "E"
	// log is always on
	logLog = "L"

	newline = "\n"
	inline  = "\t"
)

type logAgent struct {
	// slash with package name while point not when calls by funcName
	symble string
	// prints with a skip. 0 means the function that first calls a print
	skip int
	// inline, newline
	end string
	// logInfomation, logWarning, logError
	printType string
}
