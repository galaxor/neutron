package textproto

import (
	"io"
	"io/ioutil"
	"mime/quotedprintable"

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
		if child.Type == "text" && i == 1 {
			continue
		}

		attachment := &backend.Attachment{
			ID: child.ID,
			Name: child.Params["name"],
			MIMEType: child.Type + "/" + child.SubType,
			Size: child.Size,
		}
		msg.Attachments = append(msg.Attachments, attachment)
	}
}

func GetMessagePreferredPart(structure *BodyStructure) (preferred *BodyStructure) {
	preferred = structure
	for _, child := range structure.Children {
		if child.Type == "multipart" && child.SubType == "alternative" {
			return GetMessagePreferredPart(child)
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

func ParseMessagePartContent(msg *backend.Message, structure *BodyStructure, r io.Reader) error {
	if structure.ContentEncoding == "quoted-printable" {
		r = quotedprintable.NewReader(r)
	}

	charset := structure.Params["charset"]
	if charset != "" {
		r = decoder(r, charset)
	}

	slurp, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	msg.Body = string(slurp)
	return nil
}
