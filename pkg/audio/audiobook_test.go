package audio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAudiobook(t *testing.T) {
	in := "../../test/readium-audio-test/nature_crime_ml_2302_librivox/"
	out := ""
	identifier := "1234"
	metafromid3 := true
	cover := "../../test/readium-audio-test/nature_crime_ml_2302_librivox/cover.jpg"

	audiobook, err := NewAudiobook(in, out, identifier, cover, metafromid3)
	audiobook.WriteManifest()
	assert.NoError(t, err)

	assert.NotNil(t, audiobook.Files)
	assert.NotNil(t, audiobook.Manifest)
}

func TestPackageAudiobook(t *testing.T) {
	in := "../../test/readium-audio-test/nature_crime_ml_2302_librivox/"
	out := "../../test/readium-audio-test/naturecrime.audiobook"
	identifier := "W9RAWSi1TvoAWtWvmrbjC"
	metafromid3 := true
	cover := "../../test/readium-audio-test/nature_crime_ml_2302_librivox/cover.jpg"

	audiobook, err := NewAudiobook(in, out, identifier, cover, metafromid3)
	assert.NoError(t, err)

	audiobook.WriteManifest()

	err = audiobook.Package()
	assert.NoError(t, err)
}
