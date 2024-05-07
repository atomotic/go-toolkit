package audio

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/araddon/dateparse"
	"github.com/readium/go-toolkit/pkg/manifest"
)

type Audiobook struct {
	SourceDirectory string
	DestinationFile string
	ID              string
	Metadata        Metadata
	Duration        float64
	Manifest        manifest.Manifest
	Files           Audiofiles
	Cover           string
}

type Metadata struct {
	Title     string
	Author    string
	Narrator  string
	Publisher string
	Year      int
}

// MetadatafromID3 extracts metadata from the ID3 tags of the audiobook's first file.
// It returns a pointer to a Metadata struct containing the extracted metadata.
func (audiobook *Audiobook) MetadatafromID3() *Metadata {
	metadata := Metadata{}
	m := audiobook.Files[0].Metadata()
	metadata.Title, metadata.Author, metadata.Narrator, metadata.Year =
		m.Album(), m.Artist(), m.AlbumArtist(), m.Year()
	return &metadata
}

// AddTrack adds a track to the audiobook at the specified position.
// It takes an `Audiofile` object representing the track to be added and the position at which it should be inserted.
// If the track has a title, it will be used as the table of contents (TOC) title for the track.
// Otherwise, a default title in the format "Track <position>" will be used.
// The track's filename, MIME type, and duration will be added to the audiobook's manifest.
// The reading order of the audiobook's manifest will be updated to include the newly added track.
// The duration of the audiobook will be updated to include the duration of the added track.
// Returns an error if there was a problem adding the track.
func (audiobook *Audiobook) AddTrack(track Audiofile, pos int) error {
	var toctitle string
	if track.Title() != "" {
		toctitle = track.Title()
	} else {
		toctitle = fmt.Sprintf("Track %d", pos+1)
	}

	audiobook.Manifest.Links = append(audiobook.Manifest.Links, manifest.Link{
		Title:    toctitle,
		Href:     track.Filename(),
		Type:     track.Mime(),
		Duration: track.Length()})

	audiobook.Manifest.ReadingOrder = audiobook.Manifest.Links
	audiobook.Duration = audiobook.Duration + track.Length()
	audiobook.Manifest.Metadata.Duration = &audiobook.Duration
	return nil
}

func (audiobook *Audiobook) AddCover(cover string) error {

	_, err := os.Stat(cover)
	if err != nil {
		return err
	}

	audiobook.Cover = cover

	var resources []manifest.Link
	resources = append(resources, manifest.Link{
		Title: "Cover",
		Href:  "cover.jpg",
		Type:  "image/jpg",
		Rels:  []string{"cover"},
	})
	audiobook.Manifest.Resources = resources

	return nil
}

// WriteManifest writes the manifest of the audiobook to a file.
// It marshals the manifest into JSON format and saves it as "manifest.json" in the source directory of the audiobook.
// Returns an error if there was a problem writing the manifest.
func (audiobook *Audiobook) WriteManifest() error {
	manifest, err := audiobook.Manifest.MarshalJSON()
	if err != nil {
		return err
	}
	err = os.WriteFile(audiobook.SourceDirectory+"/manifest.json", manifest, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (audiobook *Audiobook) Package() error {
	out, err := os.Create(audiobook.DestinationFile)
	if err != nil {
		return err
	}
	defer out.Close()

	zipwriter := zip.NewWriter(out)
	defer zipwriter.Close()

	zmanifest, err := zipwriter.Create("manifest.json")
	if err != nil {
		return err
	}
	manifest, _ := os.Open(audiobook.SourceDirectory + "/manifest.json")
	defer manifest.Close()
	_, err = io.Copy(zmanifest, manifest)
	if err != nil {
		return err
	}

	// audio tracks
	for _, file := range audiobook.Files {
		audiofile, _ := os.Open(file.FullPath())
		defer audiofile.Close()

		mp3info, _ := audiofile.Stat()
		header, err := zip.FileInfoHeader(mp3info)
		if err != nil {
			return err
		}
		header.Method = zip.Store
		zmp3, err := zipwriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(zmp3, audiofile)
		if err != nil {
			return err
		}
	}

	// cover image
	zcover, err := zipwriter.Create("cover.jpg")
	if err != nil {
		return err
	}
	cover, err := os.Open(audiobook.Cover)
	if err != nil {
		return err
	}
	_, err = io.Copy(zcover, cover)
	if err != nil {
		return err
	}

	return nil
}

// NewAudiobook creates a new Audiobook instance with the provided source, destination, identifier, and metafromid3 parameters.
// It returns the created Audiobook instance and an error, if any.
// The source parameter specifies the source directory of the audiobook files.
// The dest parameter specifies the destination file path for the generated audiobook.
// The identifier parameter is a unique identifier for the audiobook.
// The metafromid3 parameter indicates whether to extract metadata from ID3 tags.
// If metafromid3 is true, the function will attempt to extract metadata from ID3 tags.
func NewAudiobook(source, dest string, identifier string, cover string, metafromid3 bool) (Audiobook, error) {
	metadata := &Metadata{}

	audiobook := Audiobook{
		SourceDirectory: source,
		DestinationFile: dest,
		ID:              identifier,
		Metadata:        *metadata}

	files, err := NewAudiofilesFromDirectory(source)
	if err != nil {
		return Audiobook{}, err
	}
	audiobook.Files = files

	for pos, file := range audiobook.Files {
		audiobook.AddTrack(file, pos)
	}

	audiobook.AddCover(cover)

	if metafromid3 {
		metadata = audiobook.MetadatafromID3()
	}

	audiobook.Manifest.Metadata.Type = "https://schema.org/Audiobook"
	audiobook.Manifest.Metadata.ConformsTo = append(
		audiobook.Manifest.Metadata.ConformsTo,
		manifest.ProfileAudiobook)

	audiobook.Manifest.Metadata.Identifier = audiobook.ID
	audiobook.Manifest.Metadata.LocalizedTitle = manifest.NewLocalizedStringFromString(metadata.Title)
	audiobook.Manifest.Metadata.Languages = append(audiobook.Manifest.Metadata.Languages, "ita")
	audiobook.Manifest.Metadata.Authors = append(audiobook.Manifest.Metadata.Authors,
		manifest.Contributor{LocalizedName: manifest.NewLocalizedStringFromString(metadata.Author)})
	audiobook.Manifest.Metadata.Narrators = append(audiobook.Manifest.Metadata.Narrators,
		manifest.Contributor{LocalizedName: manifest.NewLocalizedStringFromString(metadata.Narrator)})
	audiobook.Manifest.Metadata.Publishers = append(audiobook.Manifest.Metadata.Publishers,
		manifest.Contributor{LocalizedName: manifest.NewLocalizedStringFromString(metadata.Publisher)})

	t := time.Now()
	tp, _ := dateparse.ParseAny(fmt.Sprint(metadata.Year))
	audiobook.Manifest.Metadata.Published = &tp
	audiobook.Manifest.Metadata.Modified = &t

	return audiobook, nil
}
