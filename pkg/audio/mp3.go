package audio

import (
	"io"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
	"github.com/hajimehoshi/go-mp3"
)

type MP3 struct {
	filename string
	mime     string
	lenght   float64
	metadata tag.Metadata
}

func (m *MP3) FullPath() string {
	return m.filename
}

func (m *MP3) Filename() string {
	return filepath.Base(m.filename)
}

func (m *MP3) Length() float64 {
	return m.lenght
}

func (m *MP3) Metadata() tag.Metadata {
	return m.metadata
}

func (m *MP3) Title() string {
	if m.metadata.Title() != "" {
		return m.metadata.Title()
	} else {
		return ""
	}
}

func (m *MP3) Mime() string {
	return m.mime
}

func NewMP3(filename string) (Audiofile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d, err := mp3.NewDecoder(file)
	if err != nil {
		return nil, err
	}
	const sampleSize = 4
	samples := d.Length() / sampleSize
	audioLength := samples / int64(d.SampleRate())
	duration := float64(audioLength)

	// rewind the file
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	meta, err := tag.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	return &MP3{
		filename: filename,
		mime:     "audio/mpeg",
		lenght:   duration,
		metadata: meta}, nil
}
