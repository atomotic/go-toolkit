package models

import "errors"

// Publication Main structure for a publication
type Publication struct {
	Context   []string `json:"@context,omitempty"`
	Metadata  Metadata `json:"metadata"`
	Links     []Link   `json:"links"`
	Spine     []Link   `json:"spine"`
	Resources []Link   `json:"resources,omitempty"` //Replaces the manifest but less redundant
	TOC       []Link   `json:"toc,omitempty"`
	PageList  []Link   `json:"page-list,omitempty"`
	Landmarks []Link   `json:"landmarks,omitempty"`
	LOI       []Link   `json:"loi,omitempty"` //List of illustrations
	LOA       []Link   `json:"loa,omitempty"` //List of audio files
	LOV       []Link   `json:"lov,omitempty"` //List of videos
	LOT       []Link   `json:"lot,omitempty"` //List of tables

	MediaOverlays    []MediaOverlayNode      `json:"-"`
	OtherLinks       []Link                  `json:"-"` //Extension point for links that shouldn't show up in the manifest
	OtherCollections []PublicationCollection `json:"-"` //Extension point for collections that shouldn't show up in the manifest
	Internal         []Internal              `json:"-"`
}

// Internal TODO
type Internal struct {
	Name  string
	Value interface{}
}

// Link object used in collections and links
type Link struct {
	Href       string   `json:"href"`
	TypeLink   string   `json:"type,omitempty"`
	Rel        []string `json:"rel,omitempty"`
	Height     int      `json:"height,omitempty"`
	Width      int      `json:"width,omitempty"`
	Title      string   `json:"title,omitempty"`
	Properties []string `json:"properties,omitempty"`
	Duration   string   `json:"duration,omitempty"`
	Templated  bool     `json:"templated,omitempty"`
}

// PublicationCollection is used as an extension points for other collections in a Publication
type PublicationCollection struct {
	Role     string
	Metadata []Meta
	Links    []Link
	Children []PublicationCollection
}

// GetCover return the link for the cover
func (publication *Publication) GetCover() (Link, error) {
	return publication.searchLinkByRel("cover")
}

// GetNavDoc return the link for the navigation document
func (publication *Publication) GetNavDoc() (Link, error) {
	return publication.searchLinkByRel("contents")
}

func (publication *Publication) searchLinkByRel(rel string) (Link, error) {
	for _, resource := range publication.Resources {
		for _, resRel := range resource.Rel {
			if resRel == rel {
				return resource, nil
			}
		}
	}

	for _, item := range publication.Spine {
		for _, spineRel := range item.Rel {
			if spineRel == rel {
				return item, nil
			}
		}
	}

	for _, link := range publication.Links {
		for _, linkRel := range link.Rel {
			if linkRel == rel {
				return link, nil
			}
		}
	}

	return Link{}, errors.New("Can't find " + rel + " in publication")
}
