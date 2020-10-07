package constants

//ServerMode 用于指定当前为生产环境还是开发环境
type ServerMode uint

const (
	_ ServerMode = iota
	Prod
	Dev
)

func (mode ServerMode) String() string {
	if mode == Prod {
		return "prod"
	} else if mode == Dev {
		return "dev"
	} else {
		panic("unknown")
	}
}
