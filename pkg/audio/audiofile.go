package audio

import (
	"fmt"
	"os"
	"path/filepath"

	"facette.io/natsort"
	"github.com/dhowden/tag"
	"github.com/gabriel-vasile/mimetype"
)

var ErrUnsupportedFormat = fmt.Errorf("unsupported audio format")

func Mime(filename string) (string, error) {
	audiofile, _ := os.Open(filename)
	defer audiofile.Close()
	mtype, err := mimetype.DetectReader(audiofile)
	if err != nil {
		return "", err
	}
	return mtype.String(), nil
}

type Audiofiles []Audiofile

type Audiofile interface {
	Filename() string
	FullPath() string
	Length() float64
	Metadata() tag.Metadata
	Title() string
	Mime() string
}

func NewAudiofile(filename string) (Audiofile, error) {
	mime, err := Mime(filename)
	if err != nil {
		return nil, err
	}

	// handle supportedFormats
	if mime == "audio/mpeg" {
		audio, err := NewMP3(filename)
		if err != nil {
			return nil, err
		}
		return audio, nil
	} else {
		return nil, ErrUnsupportedFormat
	}
}

func NewAudiofilesFromDirectory(directory string) (audiodirectory Audiofiles, err error) {
	dir, err := filepath.Glob(directory + "/*")
	if err != nil {
		return nil, err
	}

	// natural string sorting of filenames
	natsort.Sort(dir)

	for _, file := range dir {
		audio, err := NewAudiofile(file)
		if err != nil {
			if err == ErrUnsupportedFormat {
				continue
			}
			return nil, err
		}
		audiodirectory = append(audiodirectory, audio)
	}

	return audiodirectory, nil
}
