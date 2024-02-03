package shared

const (
	NotStarted = iota
	Running
	Finished
)

const (
	Download = iota
	Search
)

const (
	Playing = iota
	Paused
	Stopped
)

var Separator = "_#__#_"

var SystemCacheDir = "/var/cache/goplay"

var Pathytldpl = "yt-dlp"
var Pathffmpeg = "ffmpeg"
var Pathffprobe = "ffprobe"
