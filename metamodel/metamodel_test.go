package metamodel_test

import (
	"encoding/json"
	"testing"

	"github.com/pflow-dev/go-metamodel/metamodel"
)

// sampleUrl is a base64 encoded zip file containing a json file
const sampleUrl = "https://pflow.dev/p/?z=UEsDBAoAAAAAACSCnFfPFUHSdwIAAHcCAAAKAAAAbW9kZWwuanNvbnsKICJtb2RlbFR5cGUiOiAicGV0cmlOZXQiLAogInZlcnNpb24iOiAidjAiLAogInBsYWNlcyI6IHsKICAiZm9vIjogewogICAib2Zmc2V0IjogMCwKICAgIngiOiAzNjQsCiAgICJ5IjogMzI2LAogICAiaW5pdGlhbCI6IDEKICB9CiB9LAogInRyYW5zaXRpb25zIjogewogICJhZGQiOiB7CiAgICJ4IjogMjQ2LAogICAieSI6IDIwNQogIH0sCiAgInN1YiI6IHsKICAgIngiOiA0NzUsCiAgICJ5IjogMjA2CiAgfSwKICAiYmFyIjogewogICAieCI6IDI0MywKICAgInkiOiA0MzkKICB9LAogICJiYXoiOiB7CiAgICJ4IjogNDc2LAogICAieSI6IDQ0MwogIH0KIH0sCiAiYXJjcyI6IFsKICB7CiAgICJzb3VyY2UiOiAiYWRkIiwKICAgInRhcmdldCI6ICJmb28iLAogICAid2VpZ2h0IjogMQogIH0sCiAgewogICAic291cmNlIjogImZvbyIsCiAgICJ0YXJnZXQiOiAic3ViIiwKICAgIndlaWdodCI6IDEKICB9LAogIHsKICAgInNvdXJjZSI6ICJiYXIiLAogICAidGFyZ2V0IjogImZvbyIsCiAgICJ3ZWlnaHQiOiAzLAogICAiaW5oaWJpdCI6IHRydWUKICB9LAogIHsKICAgInNvdXJjZSI6ICJmb28iLAogICAidGFyZ2V0IjogImJheiIsCiAgICJ3ZWlnaHQiOiAxLAogICAiaW5oaWJpdCI6IHRydWUKICB9CiBdCn1QSwECFAAKAAAAAAAkgpxXzxVB0ncCAAB3AgAACgAAAAAAAAAAAAAAAAAAAAAAbW9kZWwuanNvblBLBQYAAAAAAQABADgAAACfAgAAAAA="

func testModelDeclaration(m metamodel.Declaration) {
	cell, fn := m.Cell, m.Fn

	// block/callback style assignment
	foo := cell(func(p *metamodel.Place) {
		p.Label = "foo"
		p.Initial = 1
		p.X = 170
		p.Y = 230
	})
	baz := cell().Label("baz").Position(330, 110).Capacity(0).Initial(0)

	// chaining assignment
	bar := fn().Label("bar").Position(170, 110, 0)
	qux := fn().Label("qux").Position(330, 230)
	quux := fn().Label("quux").Position(50, 230)
	plugh := fn().Label("plugh").Role("test2").Position(460, 110)

	foo.Tx(1, bar)

	baz.Tx(1, qux)
	bar.Tx(1, baz)

	baz.Guard(1, plugh)
	foo.Guard(1, quux)
}

func TestModel_Define(t *testing.T) {
	var mm = metamodel.New().Define(testModelDeclaration)

	if !mm.Node("foo").IsPlace() {
		t.Fatal("failed to retrieve element")
	}
	if len(mm.Net().Places) != 2 {
		t.Fatalf("expected 2 places")
	}
	if mm.Net().Places["foo"] == nil {
		t.Fatalf("expected foo")
	}
	data, err := json.Marshal(mm)
	t.Logf("%s", data)
	if err != nil {
		t.Fatalf("failed to marshal")
	}

	mm.Edit().Graph()
}

func TestModel_GetState(t *testing.T) {
	var mm = metamodel.New().Define(testModelDeclaration)
	p := mm.Execute()
	s := p.GetState()
	s[0] = 7 // try to alter state
	s2 := p.GetState()
	if s2[0] == 7 {
		t.Fatalf("state should be immutable %v", s2)
	}
}

type testCmd struct {
	metamodel.Process
	call       func(op metamodel.Op) (bool, string, metamodel.Vector)
	action     string
	role       string
	expectPass bool
	expectFail bool
}

func (c testCmd) tx(t *testing.T) metamodel.Vector {
	flag, msg, v := c.call(metamodel.Op{Action: c.action, Role: c.role})
	if c.expectFail && flag {
		t.Fatalf("expected Failure - %s", msg)
	}
	if c.expectPass && !flag {
		t.Fatalf("expected OK - %s", msg)
	}
	t.Logf("%v, %v, %v", flag, msg, v)
	return v
}

func (c testCmd) assertInhibited(t *testing.T) {
	isInhibited, label := c.Process.Inhibited(metamodel.Op{Action: c.action})
	if isInhibited {
		if c.expectFail {
			t.Fatalf("expected %s not to be inhibited by %s", c.action, label)
		}
	} else {
		if c.expectPass {
			t.Fatalf("expected %s to be inhibited, got msg: %s", c.action, label)
		}
	}
}

func TestModel_Execute(t *testing.T) {

	testEditAndExecute := func(rebuild bool) {
		var mm = metamodel.New().Define(testModelDeclaration)
		if rebuild {
			mm.Edit().Graph().Index()
		}
		p := mm.Execute()
		testCmd{Process: p, action: "quux", expectPass: true}.assertInhibited(t)
		testCmd{Process: p, action: "bar", expectFail: true}.assertInhibited(t)
		testCmd{Process: p, action: "plugh", expectFail: true}.assertInhibited(t)
		testCmd{call: p.TestFire, action: "bar", expectPass: true}.tx(t)
		testCmd{call: p.Fire, action: "bar", role: "badRole", expectFail: true}.tx(t)
		testCmd{call: p.Fire, action: "bar", role: "default", expectPass: true}.tx(t)
		testCmd{Process: p, action: "plugh", expectPass: true}.assertInhibited(t)
		testCmd{Process: p, action: "quux", expectFail: true}.assertInhibited(t)
		testCmd{call: p.TestFire, action: "bar", expectFail: true}.tx(t)
		testCmd{call: p.Fire, action: "bar", expectFail: true}.tx(t)
	}

	testEditAndExecute(false)
	testEditAndExecute(true)
}

func TestVectorFromBytes(t *testing.T) {
	v := metamodel.Vector{0, 1, -127}
	vb := metamodel.VectorToBytes(v)
	t.Logf("vector: %v\n", v)
	t.Logf("bytes: %v\n", vb)
	v2 := metamodel.VectorFromBytes(vb)
	t.Logf("vector: %v\n", v2)
	for i, v := range v {
		if v != v2[i] {
			t.Fatalf("mismatch %v <=> %v", v, v2[i])
		}
	}
}

func TestZipAndUnzipUrl(t *testing.T) {
	mm := metamodel.New()
	_, ok := mm.UnzipUrl(sampleUrl, "model.json")
	if !ok {
		t.Fatalf("failed to unzip")
	}

	urlOut, zipOk := mm.ZipUrl("https://pflow.dev/p/")
	if !zipOk {
		t.Fatalf("failed to zip")
	}

	m := metamodel.New()
	_, unzipOk := m.UnzipUrl(urlOut, "model.json")
	if !unzipOk {
		t.Fatalf("failed to unzip")
	}
	t.Logf("generated url: %s", urlOut)
}
func TestUnzipUrl(t *testing.T) {
	mm := metamodel.New()
	_, ok := mm.UnzipUrl(sampleUrl, "model.json")
	if !ok {
		t.Fatalf("failed to unzip")
	}
	p := mm.Execute()
	testCmd{Process: p, action: "bar", expectPass: true}.assertInhibited(t)
	testCmd{call: p.Fire, action: "add", expectPass: true}.tx(t)
	testCmd{call: p.Fire, action: "add", expectPass: true}.tx(t)
	testCmd{Process: p, action: "bar", expectFail: true}.assertInhibited(t)

	testCmd{Process: p, action: "baz", expectPass: true}.assertInhibited(t)
	testCmd{call: p.Fire, action: "sub", expectPass: true}.tx(t)
	testCmd{call: p.Fire, action: "sub", expectPass: true}.tx(t)
	testCmd{call: p.Fire, action: "sub", expectPass: true}.tx(t)
	testCmd{Process: p, action: "baz", expectFail: true}.assertInhibited(t)
}
