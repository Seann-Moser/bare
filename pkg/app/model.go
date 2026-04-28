package app

type App struct {
	Title       string
	Icon        string
	DefaultSize struct {
		Width  int
		Height int
	}
	Resizable bool
	TargetFPS int
}
