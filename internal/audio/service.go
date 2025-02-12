package audio

type AudioProvider interface {
	SearchAudio(query string, limit int) ([]AudioResult, error)
}
