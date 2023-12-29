package metamodel_test

import (
	"encoding/json"
	"testing"

	"github.com/pflow-dev/go-metamodel/metamodel"
)

// sampleUrl is a base64 encoded zip file containing a json file
const sampleUrl = "https://pflow.dev/p/?z=UEsDBAoAAAAAACSCnFfPFUHSdwIAAHcCAAAKAAAAbW9kZWwuanNvbnsKICJtb2RlbFR5cGUiOiAicGV0cmlOZXQiLAogInZlcnNpb24iOiAidjAiLAogInBsYWNlcyI6IHsKICAiZm9vIjogewogICAib2Zmc2V0IjogMCwKICAgIngiOiAzNjQsCiAgICJ5IjogMzI2LAogICAiaW5pdGlhbCI6IDEKICB9CiB9LAogInRyYW5zaXRpb25zIjogewogICJhZGQiOiB7CiAgICJ4IjogMjQ2LAogICAieSI6IDIwNQogIH0sCiAgInN1YiI6IHsKICAgIngiOiA0NzUsCiAgICJ5IjogMjA2CiAgfSwKICAiYmFyIjogewogICAieCI6IDI0MywKICAgInkiOiA0MzkKICB9LAogICJiYXoiOiB7CiAgICJ4IjogNDc2LAogICAieSI6IDQ0MwogIH0KIH0sCiAiYXJjcyI6IFsKICB7CiAgICJzb3VyY2UiOiAiYWRkIiwKICAgInRhcmdldCI6ICJmb28iLAogICAid2VpZ2h0IjogMQogIH0sCiAgewogICAic291cmNlIjogImZvbyIsCiAgICJ0YXJnZXQiOiAic3ViIiwKICAgIndlaWdodCI6IDEKICB9LAogIHsKICAgInNvdXJjZSI6ICJiYXIiLAogICAidGFyZ2V0IjogImZvbyIsCiAgICJ3ZWlnaHQiOiAzLAogICAiaW5oaWJpdCI6IHRydWUKICB9LAogIHsKICAgInNvdXJjZSI6ICJmb28iLAogICAidGFyZ2V0IjogImJheiIsCiAgICJ3ZWlnaHQiOiAxLAogICAiaW5oaWJpdCI6IHRydWUKICB9CiBdCn1QSwECFAAKAAAAAAAkgpxXzxVB0ncCAAB3AgAACgAAAAAAAAAAAAAAAAAAAAAAbW9kZWwuanNvblBLBQYAAAAAAQABADgAAACfAgAAAAA="

const badUrl = "https://pflow.dev/p.svg?z=UEsDBAoAAAAAAI6hnVeykZqF2gsAANoLAAAKAAAAbW9kZWwuanNvbnsKICAibW9kZWxUeXBlIjogInBldHJpTmV0IiwKICAidmVyc2lvbiI6ICJ2MCIsCiAgInBsYWNlcyI6IHsKICAgICJ1cC51cCI6IHsKICAgICAgIm9mZnNldCI6IDAsCiAgICAgICJ4IjogMTYwLAogICAgICAieSI6IDE2MCwKICAgICAgImluaXRpYWwiOiAyLAogICAgICAiY2FwYWNpdHkiOiAyCiAgICB9LAogICAgImRvd24uZG93biI6IHsKICAgICAgIm9mZnNldCI6IDEsCiAgICAgICJ4IjogODAsCiAgICAgICJ5IjogMjQwLAogICAgICAiaW5pdGlhbCI6IDIsCiAgICAgICJjYXBhY2l0eSI6IDIKICAgIH0sCiAgICAidHdvLmxlZnRzIjogewogICAgICAib2Zmc2V0IjogMiwKICAgICAgIngiOiA4MCwKICAgICAgInkiOiA4MCwKICAgICAgImluaXRpYWwiOiAyLAogICAgICAiY2FwYWNpdHkiOiAyCiAgICB9LAogICAgInRoZW4ucmlnaHQiOiB7CiAgICAgICJvZmZzZXQiOiAzLAogICAgICAieCI6IDI0MCwKICAgICAgInkiOiAyNDAsCiAgICAgICJjYXBhY2l0eSI6IDEKICAgIH0sCiAgICAidHdvLnJpZ2h0cyI6IHsKICAgICAgIm9mZnNldCI6IDQsCiAgICAgICJ4IjogMjQwLAogICAgICAieSI6IDgwCiAgICB9LAogICAgInRoZW4uQSI6IHsKICAgICAgIm9mZnNldCI6IDUsCiAgICAgICJ4IjogNDgwLAogICAgICAieSI6IDgwLAogICAgICAiY2FwYWNpdHkiOiAxCiAgICB9LAogICAgInRoZW4uU2VsZWN0IjogewogICAgICAib2Zmc2V0IjogNiwKICAgICAgIngiOiAzMjAsCiAgICAgICJ5IjogMjQwLAogICAgICAiY2FwYWNpdHkiOiAxCiAgICB9LAogICAgInRoZW4uU3RhcnQiOiB7CiAgICAgICJvZmZzZXQiOiA3LAogICAgICAieCI6IDQwMCwKICAgICAgInkiOiAyNDAsCiAgICAgICJjYXBhY2l0eSI6IDEKICAgIH0KICB9LAogICJ0cmFuc2l0aW9ucyI6IHsKICAgICJ1cCI6IHsKICAgICAgInJvbGUiOiB7CiAgICAgICAgImxhYmVsIjogImRlZmF1bHQiCiAgICAgIH0sCiAgICAgICJ4IjogMTYwLAogICAgICAieSI6IDgwCiAgICB9LAogICAgImRvd24iOiB7CiAgICAgICJyb2xlIjogewogICAgICAgICJsYWJlbCI6ICJkZWZhdWx0IgogICAgICB9LAogICAgICAieCI6IDE2MCwKICAgICAgInkiOiAyNDAKICAgIH0sCiAgICAibGVmdCI6IHsKICAgICAgInJvbGUiOiB7CiAgICAgICAgImxhYmVsIjogImRlZmF1bHQiCiAgICAgIH0sCiAgICAgICJ4IjogODAsCiAgICAgICJ5IjogMTYwCiAgICB9LAogICAgInJpZ2h0IjogewogICAgICAicm9sZSI6IHsKICAgICAgICAibGFiZWwiOiAiZGVmYXVsdCIKICAgICAgfSwKICAgICAgIngiOiAyNDAsCiAgICAgICJ5IjogMTYwCiAgICB9LAogICAgInNlbGVjdCI6IHsKICAgICAgInJvbGUiOiB7CiAgICAgICAgImxhYmVsIjogImRlZmF1bHQiCiAgICAgIH0sCiAgICAgICJ4IjogMzIwLAogICAgICAieSI6IDE2MAogICAgfSwKICAgICJzdGFydCI6IHsKICAgICAgInJvbGUiOiB7CiAgICAgICAgImxhYmVsIjogImRlZmF1bHQiCiAgICAgIH0sCiAgICAgICJ4IjogNDAwLAogICAgICAieSI6IDE2MAogICAgfSwKICAgICJiIjogewogICAgICAicm9sZSI6IHsKICAgICAgICAibGFiZWwiOiAiZGVmYXVsdCIKICAgICAgfSwKICAgICAgIngiOiA0ODAsCiAgICAgICJ5IjogMTYwCiAgICB9LAogICAgImEiOiB7CiAgICAgICJyb2xlIjogewogICAgICAgICJsYWJlbCI6ICJkZWZhdWx0IgogICAgICB9LAogICAgICAieCI6IDU2MCwKICAgICAgInkiOiAxNjAKICAgIH0KICB9LAogICJhcmNzIjogWwogICAgewogICAgICAic291cmNlIjogInVwLnVwIiwKICAgICAgInRhcmdldCI6ICJ1cCIsCiAgICAgICJ3ZWlnaHQiOiAxCiAgICB9LAogICAgewogICAgICAic291cmNlIjogInVwLnVwIiwKICAgICAgInRhcmdldCI6ICJkb3duIiwKICAgICAgIndlaWdodCI6IDEsCiAgICAgICJpbmhpYml0IjogdHJ1ZQogICAgfSwKICAgIHsKICAgICAgInNvdXJjZSI6ICJkb3duLmRvd24iLAogICAgICAidGFyZ2V0IjogImRvd24iLAogICAgICAid2VpZ2h0IjogMQogICAgfSwKICAgIHsKICAgICAgInNvdXJjZSI6ICJkb3duLmRvd24iLAogICAgICAidGFyZ2V0IjogImxlZnQiLAogICAgICAid2VpZ2h0IjogMSwKICAgICAgImluaGliaXQiOiB0cnVlCiAgICB9LAogICAgewogICAgICAic291cmNlIjogInR3by5sZWZ0cyIsCiAgICAgICJ0YXJnZXQiOiAibGVmdCIsCiAgICAgICJ3ZWlnaHQiOiAxCiAgICB9LAogICAgewogICAgICAic291cmNlIjogInRoZW4ucmlnaHQiLAogICAgICAidGFyZ2V0IjogInJpZ2h0IiwKICAgICAgIndlaWdodCI6IDEKICAgIH0sCiAgICB7CiAgICAgICJzb3VyY2UiOiAibGVmdCIsCiAgICAgICJ0YXJnZXQiOiAidGhlbi5yaWdodCIsCiAgICAgICJ3ZWlnaHQiOiAxCiAgICB9LAogICAgewogICAgICAic291cmNlIjogInR3by5yaWdodHMiLAogICAgICAidGFyZ2V0IjogImIiLAogICAgICAid2VpZ2h0IjogMgogICAgfSwKICAgIHsKICAgICAgInNvdXJjZSI6ICJyaWdodCIsCiAgICAgICJ0YXJnZXQiOiAidHdvLnJpZ2h0cyIsCiAgICAgICJ3ZWlnaHQiOiAxCiAgICB9LAogICAgewogICAgICAic291cmNlIjogInRoZW4uQSIsCiAgICAgICJ0YXJnZXQiOiAiYSIsCiAgICAgICJ3ZWlnaHQiOiAxCiAgICB9LAogICAgewogICAgICAic291cmNlIjogImIiLAogICAgICAidGFyZ2V0IjogInRoZW4uQSIsCiAgICAgICJ3ZWlnaHQiOiAxCiAgICB9LAogICAgewogICAgICAic291cmNlIjogInRoZW4uU2VsZWN0IiwKICAgICAgInRhcmdldCI6ICJzZWxlY3QiLAogICAgICAid2VpZ2h0IjogMQogICAgfSwKICAgIHsKICAgICAgInNvdXJjZSI6ICJhIiwKICAgICAgInRhcmdldCI6ICJ0aGVuLlNlbGVjdCIsCiAgICAgICJ3ZWlnaHQiOiAxCiAgICB9LAogICAgewogICAgICAic291cmNlIjogInRoZW4uU3RhcnQiLAogICAgICAidGFyZ2V0IjogInN0YXJ0IiwKICAgICAgIndlaWdodCI6IDEKICAgIH0sCiAgICB7CiAgICAgICJzb3VyY2UiOiAic2VsZWN0IiwKICAgICAgInRhcmdldCI6ICJ0aGVuLlN0YXJ0IiwKICAgICAgIndlaWdodCI6IDEKICAgIH0KICBdCn1QSwECFAAKAAAAAACOoZ1XspGahdoLAADaCwAACgAAAAAAAAAAAAAAAAAAAAAAbW9kZWwuanNvblBLBQYAAAAAAQABADgAAAACDAAAAAA="

const goodUrl = "https://pflow.dev/p.svg?z=UEsDBAoAAAAAAMKhnVdIN14wOgoAADoKAAAKAAAAbW9kZWwuanNvbnsKICAibW9kZWxUeXBlIjogInBldHJpTmV0IiwKICAidmVyc2lvbiI6ICJ2MCIsCiAgInBsYWNlcyI6IHsKICAgICJ1cC51cCI6IHsKICAgICAgIm9mZnNldCI6IDAsCiAgICAgICJ4IjogMTYwLAogICAgICAieSI6IDE2MCwKICAgICAgImluaXRpYWwiOiAyLAogICAgICAiY2FwYWNpdHkiOiAyCiAgICB9LAogICAgImRvd24uZG93biI6IHsKICAgICAgIm9mZnNldCI6IDEsCiAgICAgICJ4IjogODAsCiAgICAgICJ5IjogMjQwLAogICAgICAiaW5pdGlhbCI6IDIsCiAgICAgICJjYXBhY2l0eSI6IDIKICAgIH0sCiAgICAidHdvLmxlZnRzIjogewogICAgICAib2Zmc2V0IjogMiwKICAgICAgIngiOiA4MCwKICAgICAgInkiOiA4MCwKICAgICAgImluaXRpYWwiOiAyLAogICAgICAiY2FwYWNpdHkiOiAyCiAgICB9LAogICAgInRoZW4ucmlnaHQiOiB7CiAgICAgICJvZmZzZXQiOiAzLAogICAgICAieCI6IDI0MCwKICAgICAgInkiOiAyNDAsCiAgICAgICJjYXBhY2l0eSI6IDEKICAgIH0sCiAgICAidHdvLnJpZ2h0cyI6IHsKICAgICAgIm9mZnNldCI6IDQsCiAgICAgICJ4IjogMjQwLAogICAgICAieSI6IDgwCiAgICB9LAogICAgInRoZW4uQSI6IHsKICAgICAgIm9mZnNldCI6IDUsCiAgICAgICJ4IjogNDgwLAogICAgICAieSI6IDgwLAogICAgICAiY2FwYWNpdHkiOiAxCiAgICB9LAogICAgInRoZW4uU2VsZWN0IjogewogICAgICAib2Zmc2V0IjogNiwKICAgICAgIngiOiAzMjAsCiAgICAgICJ5IjogMjQwLAogICAgICAiY2FwYWNpdHkiOiAxCiAgICB9LAogICAgInRoZW4uU3RhcnQiOiB7CiAgICAgICJvZmZzZXQiOiA3LAogICAgICAieCI6IDQwMCwKICAgICAgInkiOiAyNDAsCiAgICAgICJjYXBhY2l0eSI6IDEKICAgIH0KICB9LAogICJ0cmFuc2l0aW9ucyI6IHsKICAgICJ1cCI6IHsKICAgICAgIngiOiAxNjAsCiAgICAgICJ5IjogODAKICAgIH0sCiAgICAiZG93biI6IHsKICAgICAgIngiOiAxNjAsCiAgICAgICJ5IjogMjQwCiAgICB9LAogICAgImxlZnQiOiB7CiAgICAgICJ4IjogODAsCiAgICAgICJ5IjogMTYwCiAgICB9LAogICAgInJpZ2h0IjogewogICAgICAieCI6IDI0MCwKICAgICAgInkiOiAxNjAKICAgIH0sCiAgICAic2VsZWN0IjogewogICAgICAieCI6IDMyMCwKICAgICAgInkiOiAxNjAKICAgIH0sCiAgICAic3RhcnQiOiB7CiAgICAgICJ4IjogNDAwLAogICAgICAieSI6IDE2MAogICAgfSwKICAgICJiIjogewogICAgICAieCI6IDQ4MCwKICAgICAgInkiOiAxNjAKICAgIH0sCiAgICAiYSI6IHsKICAgICAgIngiOiA1NjAsCiAgICAgICJ5IjogMTYwCiAgICB9CiAgfSwKICAiYXJjcyI6IFsKICAgIHsKICAgICAgInNvdXJjZSI6ICJ1cC51cCIsCiAgICAgICJ0YXJnZXQiOiAidXAiLAogICAgICAid2VpZ2h0IjogMQogICAgfSwKICAgIHsKICAgICAgInNvdXJjZSI6ICJ1cC51cCIsCiAgICAgICJ0YXJnZXQiOiAiZG93biIsCiAgICAgICJ3ZWlnaHQiOiAxLAogICAgICAiaW5oaWJpdCI6IHRydWUKICAgIH0sCiAgICB7CiAgICAgICJzb3VyY2UiOiAiZG93bi5kb3duIiwKICAgICAgInRhcmdldCI6ICJkb3duIiwKICAgICAgIndlaWdodCI6IDEKICAgIH0sCiAgICB7CiAgICAgICJzb3VyY2UiOiAiZG93bi5kb3duIiwKICAgICAgInRhcmdldCI6ICJsZWZ0IiwKICAgICAgIndlaWdodCI6IDEsCiAgICAgICJpbmhpYml0IjogdHJ1ZQogICAgfSwKICAgIHsKICAgICAgInNvdXJjZSI6ICJ0d28ubGVmdHMiLAogICAgICAidGFyZ2V0IjogImxlZnQiLAogICAgICAid2VpZ2h0IjogMQogICAgfSwKICAgIHsKICAgICAgInNvdXJjZSI6ICJ0aGVuLnJpZ2h0IiwKICAgICAgInRhcmdldCI6ICJyaWdodCIsCiAgICAgICJ3ZWlnaHQiOiAxCiAgICB9LAogICAgewogICAgICAic291cmNlIjogImxlZnQiLAogICAgICAidGFyZ2V0IjogInRoZW4ucmlnaHQiLAogICAgICAid2VpZ2h0IjogMQogICAgfSwKICAgIHsKICAgICAgInNvdXJjZSI6ICJ0d28ucmlnaHRzIiwKICAgICAgInRhcmdldCI6ICJiIiwKICAgICAgIndlaWdodCI6IDIKICAgIH0sCiAgICB7CiAgICAgICJzb3VyY2UiOiAicmlnaHQiLAogICAgICAidGFyZ2V0IjogInR3by5yaWdodHMiLAogICAgICAid2VpZ2h0IjogMQogICAgfSwKICAgIHsKICAgICAgInNvdXJjZSI6ICJ0aGVuLkEiLAogICAgICAidGFyZ2V0IjogImEiLAogICAgICAid2VpZ2h0IjogMQogICAgfSwKICAgIHsKICAgICAgInNvdXJjZSI6ICJiIiwKICAgICAgInRhcmdldCI6ICJ0aGVuLkEiLAogICAgICAid2VpZ2h0IjogMQogICAgfSwKICAgIHsKICAgICAgInNvdXJjZSI6ICJ0aGVuLlNlbGVjdCIsCiAgICAgICJ0YXJnZXQiOiAic2VsZWN0IiwKICAgICAgIndlaWdodCI6IDEKICAgIH0sCiAgICB7CiAgICAgICJzb3VyY2UiOiAiYSIsCiAgICAgICJ0YXJnZXQiOiAidGhlbi5TZWxlY3QiLAogICAgICAid2VpZ2h0IjogMQogICAgfSwKICAgIHsKICAgICAgInNvdXJjZSI6ICJ0aGVuLlN0YXJ0IiwKICAgICAgInRhcmdldCI6ICJzdGFydCIsCiAgICAgICJ3ZWlnaHQiOiAxCiAgICB9LAogICAgewogICAgICAic291cmNlIjogInNlbGVjdCIsCiAgICAgICJ0YXJnZXQiOiAidGhlbi5TdGFydCIsCiAgICAgICJ3ZWlnaHQiOiAxCiAgICB9CiAgXQp9UEsBAhQACgAAAAAAwqGdV0g3XjA6CgAAOgoAAAoAAAAAAAAAAAAAAAAAAAAAAG1vZGVsLmpzb25QSwUGAAAAAAEAAQA4AAAAYgoAAAAA"

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
	var mm = metamodel.New("test").Define(testModelDeclaration)

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
	var mm = metamodel.New("test").Define(testModelDeclaration)
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
		var mm = metamodel.New("test").Define(testModelDeclaration)
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

func TestUnzipGoodAndBadUrl(t *testing.T) {
	mm := metamodel.New("test")
	_, ok := mm.UnzipUrl(goodUrl)
	if !ok {
		t.Fatalf("failed to unzip")
	}
	_, ok = mm.UnzipUrl(badUrl)
	if ok {
		t.Fatalf("should have failed to unzip")
	}
}

func TestZipAndUnzipUrl(t *testing.T) {
	mm := metamodel.New("test")
	_, ok := mm.UnzipUrl(sampleUrl)
	if !ok {
		t.Fatalf("failed to unzip")
	}

	urlOut, zipOk := mm.ZipUrl("https://pflow.dev/p/")
	if !zipOk {
		t.Fatalf("failed to zip")
	}

	m := metamodel.New("test2")
	_, unzipOk := m.UnzipUrl(urlOut)
	if !unzipOk {
		t.Fatalf("failed to unzip")
	}
	t.Logf("generated url: %s", urlOut)
}
func TestUnzipUrl(t *testing.T) {
	mm := metamodel.New("test")
	_, ok := mm.UnzipUrl(sampleUrl)
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
