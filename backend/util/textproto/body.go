package textproto

import (
	"io"
	"strings"
	"strconv"

	"github.com/emersion/neutron/backend"
)

type BodyStructure struct {
	ID string
	Type string
	SubType string
	Params map[string]string
	ContentId string
	ContentDescription string
	ContentEncoding string
	Content io.Reader
	Size int
	Children []*BodyStructure
}

func (s *BodyStructure) Get(id string) *BodyStructure {
	if id == "" {
		return s
	}

	parts := strings.SplitN(id, ".", 2)
	index, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil
	}

	var childId string
	if len(parts) == 2 {
		childId = parts[1]
	}

	for i, child := range s.Children {
		if i == index - 1 {
			return child.Get(childId)
		}
	}

	return nil
}

func (s *BodyStructure) GetPreferredPart() (preferred *BodyStructure) {
	preferred = s
	for _, child := range s.Children {
		if child.Type == "multipart" && child.SubType == "alternative" {
			return child.GetPreferredPart()
		}
		if child.Type != "text" {
			continue
		}
		if preferred.Type == "multipart" || child.SubType == "html" {
			preferred = child
		}
	}
	return
}

func (s *BodyStructure) DecodeContent(r io.Reader) io.Reader {
	contentEncoding := s.ContentEncoding
	if contentEncoding != "" {
		r = decodeContentEncoding(r, contentEncoding)
	}

	charset := s.Params["charset"]
	if charset != "" {
		r = decodeCharset(r, charset)
	}

	s.Content = r
	return r
}

func (s *BodyStructure) Attachment() *backend.Attachment {
	name := s.Params["name"]
	if name == "" {
		name = s.Params["filename"]
	}

	return &backend.Attachment{
		ID: s.ID,
		Name: name,
		MIMEType: s.Type + "/" + s.SubType,
		Size: s.Size,
	}
}

func ParseMessageStructure(msg *backend.Message, structure *BodyStructure) {
	if structure.Type != "multipart" || structure.SubType == "alternative" {
		return
	}

	for i, child := range structure.Children {
		if child.Type == "multipart" {
			ParseMessageStructure(msg, child)
			continue
		}

		// AppleMail doesn't format well headers
		// First child is message content
		if child.Type == "text" && i == 0 {
			continue
		}

		msg.Attachments = append(msg.Attachments, child.Attachment())
	}
}
