package audio

import "github.com/dhowden/tag"

type WAV struct {
	filename string
	mime     string
	lenght   float64
	metadata tag.Metadata
}

func (w *WAV) FullPath() string {
	return w.filename
}

func (w *WAV) Filename() string {
	return w.filename
}

func (w *WAV) Mime() string {
	return w.mime
}

func (w *WAV) Length() float64 {
	return 0
}

func (w *WAV) Metadata() tag.Metadata {
	return nil
}

func (w *WAV) Title() string {
	return w.metadata.Title()
}

func NewWAV(filename string) (Audiofile, error) {
	return &WAV{
		filename: filename,
		mime:     "audio/wav",
		lenght:   0,
		metadata: nil}, nil
}
