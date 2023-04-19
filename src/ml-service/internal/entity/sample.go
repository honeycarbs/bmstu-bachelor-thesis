package entity

type Sample struct {
	ID        string `db:"uuid"`
	AudioPath string `db:"audio_path"`
	Emotion   Label  `db:"emotion"`
	Frames    []Frame
}
