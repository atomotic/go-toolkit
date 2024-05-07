package audio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMime(t *testing.T) {
	filename := "../../test/readium-audio-test/nature_crime_ml_2302_librivox/natureofacrime_00_conrad_64kb.mp3"
	expectedMime := "audio/mpeg"
	mime, err := Mime(filename)
	assert.NoError(t, err)
	assert.Equal(t, expectedMime, mime)
}
func TestNewAudiodirectory(t *testing.T) {
	directory := "../../test/readium-audio-test/nature_crime_ml_2302_librivox"

	audiodirectory, err := NewAudiofilesFromDirectory(directory)
	assert.NoError(t, err)
	assert.Equal(t, 9, len(audiodirectory))
}
