package shared

const (
	NotStarted = iota
	Running
	Finished
)

const (
	Downloading = iota
	Searching
)

const (
	PinkTheme   = "pink"
	BlueTheme   = "blue"
	PurpleTheme = "purple"
)

type PState uint

const (
	Playing PState = iota
	Paused
	Stopped
)

const (
	HashPrefixLength = 5
)
