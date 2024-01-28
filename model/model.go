package model

import (
	"github.com/pflow-dev/go-metamodel/v2/codec"
	"github.com/pflow-dev/go-metamodel/v2/metamodel"
	"time"
)

// Zblob is a data wrapper for encapsulating a model
type Zblob struct {
	ID           int64     `json:"-"`
	IpfsCid      string    `json:"cid"`
	Base64Zipped string    `json:"data"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Keywords     string    `json:"keywords"`
	Referer      string    `json:"-"`
	CreatedAt    time.Time `json:"created"`
}

type Model struct {
	*Zblob
}

type Document struct {
	ModelCid    string                      `json:"model_cid"`
	Title       string                      `json:"title"`
	Description string                      `json:"description"`
	Keywords    string                      `json:"keywords"`
	Declaration metamodel.DeclarationObject `json:"declaration"`
}

func getMetaModel(data string) metamodel.MetaModel {
	mm := metamodel.New()
	_, ok := mm.UnpackFromUrl("?z="+data, "model.json")
	if !ok {
		panic("Failed to unzip model")
	}
	return mm
}

func assertValid(data string) {
	_ = getMetaModel(data)
}

func (z *Zblob) ToModel() Model {
	assertValid(z.Base64Zipped)
	return Model{Zblob: z}
}
func (z *Zblob) ToDocument() Document {
	mm := getMetaModel(z.Base64Zipped)
	return Document{
		ModelCid:    z.IpfsCid,
		Title:       z.Title,
		Description: z.Description,
		Keywords:    z.Keywords,
		Declaration: mm.ToDeclarationObject(),
	}
}

func (d Document) Cid() string {
	return codec.ToOid(codec.Marshal(d)).String()
}

func (m *Model) MetaModel() (string, metamodel.MetaModel) {
	mm := metamodel.New()
	jsonData, ok := mm.UnpackFromUrl("?z="+m.Base64Zipped, "model.json")
	if !ok {
		panic("Failed to unzip model")
	}
	return jsonData, mm
}

func (m *Model) Declare(args ...func(metamodel.Declaration)) {
	mm := metamodel.New()
	mm.Define(args...)
	url, _ := mm.ZipUrl()
	m.Base64Zipped = url[3:]
	m.IpfsCid = codec.ToOid(codec.Marshal(m.Base64Zipped)).String()
}
