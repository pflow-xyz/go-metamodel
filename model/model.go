package model

import (
	"github.com/pflow-xyz/go-metamodel/metamodel"
	"github.com/pflow-xyz/go-metamodel/oid"
	. "github.com/pflow-xyz/go-metamodel/zblob"
)

type Model struct {
	*Zblob
}

func FromZblob(z *Zblob) *Model {
	return &Model{z}
}

func (m *Model) MetaModel() (string, metamodel.MetaModel) {
	mm := metamodel.New()
	jsonData, ok := mm.UnpackFromUrl("?z=" + m.Base64Zipped)
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
	m.IpfsCid = oid.ToOid(oid.Marshal(m.Base64Zipped)).String()
}
