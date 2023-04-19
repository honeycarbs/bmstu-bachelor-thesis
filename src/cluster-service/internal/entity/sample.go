package entity

type Sample struct {
	ID        string `db:"uuid"`
	AudioPath string `db:"audio_path"`
	Emotion   string `db:"emotion"`
	Frames    []Frame
}
