package bm

type Logger interface {
	Panicln(args ...interface{})
	Fatalln(args ...interface{})
	Errorln(args ...interface{})
	Warnln(args ...interface{})
	Infoln(args ...interface{})
	Debugln(args ...interface{})
}
