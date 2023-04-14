package entity

type Sample struct {
	Hash      string `db:"hash"`
	AudioPath string `db:"audio_path"`
	Emotion   string `db:"emotion"`
	Frames    []Frame
}
