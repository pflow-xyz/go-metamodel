package model

import (
	"github.com/pflow-dev/go-metamodel/codec"
	"github.com/pflow-dev/go-metamodel/metamodel"
	"time"
)

type Model struct {
	ID                int       `json:"id"`
	IpfsCid           string    `json:"cid"`
	Base64GzippedJson string    `json:"data"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Keywords          string    `json:"keywords"`
	CreatedAt         time.Time `json:"created"`
	PublishedAt       time.Time `json:"published"`
}

type Document struct {
	ModelCid    string                      `json:"model_cid"`
	Title       string                      `json:"title"`
	Description string                      `json:"description"`
	Keywords    string                      `json:"keywords"`
	Declaration metamodel.DeclarationObject `json:"declaration"`
	PublishedAt time.Time                   `json:"published"`
}

func (d Document) Cid() string {
	return codec.ToOid(codec.Marshal(d)).String()
}

func (m *Model) ToDocument() Document {
	mm := metamodel.New()
	_, ok := mm.UnzipUrl("?z="+m.Base64GzippedJson, "model.json")
	if !ok {
		panic("Failed to unzip model")
	}
	return Document{
		ModelCid:    m.IpfsCid,
		Title:       m.Title,
		Description: m.Description,
		Keywords:    m.Keywords,
		Declaration: mm.ToDeclarationObject(),
		PublishedAt: m.PublishedAt,
	}
}

func (m *Model) MetaModel() (string, metamodel.MetaModel) {
	mm := metamodel.New()
	jsonData, ok := mm.UnzipUrl("?z="+m.Base64GzippedJson, "model.json")
	if !ok {
		panic("Failed to unzip model")
	}
	return jsonData, mm
}

func (m *Model) Declare(args ...func(metamodel.Declaration)) {
	mm := metamodel.New()
	mm.Define(args...)
	url, _ := mm.ZipUrl()
	m.Base64GzippedJson = url[3:]
	m.IpfsCid = codec.ToOid(codec.Marshal(m.Base64GzippedJson)).String()
}
