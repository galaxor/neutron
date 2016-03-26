package textproto

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"mime/quotedprintable"
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
	switch s.ContentEncoding {
	case "quoted-printable":
		r = quotedprintable.NewReader(r)
	case "base64":
		r = base64.NewDecoder(base64.StdEncoding, r)
	}

	charset := s.Params["charset"]
	if charset != "" {
		r = decoder(r, charset)
	}

	return r
}

func (s *BodyStructure) Attachment() *backend.Attachment {
	return &backend.Attachment{
		ID: s.ID,
		Name: s.Params["name"],
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

func ParseMessagePartContent(msg *backend.Message, structure *BodyStructure, r io.Reader) error {
	r = structure.DecodeContent(r)

	slurp, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	msg.Body = string(slurp)
	return nil
}
