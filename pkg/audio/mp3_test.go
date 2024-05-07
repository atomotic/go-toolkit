package audio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMP3(t *testing.T) {
	filename := "../../test/readium-audio-test/nature_crime_ml_2302_librivox/natureofacrime_00_conrad_64kb.mp3"
	audio, err := NewMP3(filename)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, float64(659), audio.Length())
	assert.Equal(t, "00 - Preface", audio.Title())
	assert.Equal(t, "natureofacrime_00_conrad_64kb.mp3", audio.Filename())
	assert.NotNil(t, audio.Metadata())
}
