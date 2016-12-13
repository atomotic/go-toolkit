package parser

import (
	"archive/zip"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/feedbooks/webpub-streamer/models"
)

func init() {
	parserList = append(parserList, List{fileExt: "cbz", parser: CbzParser})
}

// CbzParser TODO add doc
func CbzParser(filePath string, selfURL string) models.Publication {
	var publication models.Publication

	publication.Metadata.Title = filePathToTitle(filePath)
	publication.Metadata.Identifier = filePath
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		fmt.Println("failed to open zip " + filePath)
		fmt.Println(err)
		return publication
	}

	publication.Internal = append(publication.Internal, models.Internal{Name: "type", Value: "cbz"})
	publication.Internal = append(publication.Internal, models.Internal{Name: "cbz", Value: zipReader})

	for _, f := range zipReader.File {
		linkItem := models.Link{}
		linkItem.TypeLink = getMediaTypeByName(f.Name)
		linkItem.Href = f.Name
		if linkItem.TypeLink != "" {
			publication.Spine = append(publication.Spine, linkItem)
		}
	}

	return publication
}

func filePathToTitle(filePath string) string {
	_, filename := filepath.Split(filePath)
	filename = strings.Split(filename, ".")[0]
	title := strings.Replace(filename, "_", " ", -1)

	return title
}

func getMediaTypeByName(filePath string) string {
	ext := filepath.Ext(filePath)

	switch strings.ToLower(ext) {
	case ".jpg":
		return "image/jpeg"
	case ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	default:
		return ""
	}
}
