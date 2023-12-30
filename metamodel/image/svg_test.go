package image_test

import (
	"github.com/pflow-dev/go-metamodel/metamodel"
	"github.com/pflow-dev/go-metamodel/metamodel/image"
	"testing"
)

const sampleUrl = "https://pflow.dev/p/?z=UEsDBAoAAAAAACSCnFfPFUHSdwIAAHcCAAAKAAAAbW9kZWwuanNvbnsKICJtb2RlbFR5cGUiOiAicGV0cmlOZXQiLAogInZlcnNpb24iOiAidjAiLAogInBsYWNlcyI6IHsKICAiZm9vIjogewogICAib2Zmc2V0IjogMCwKICAgIngiOiAzNjQsCiAgICJ5IjogMzI2LAogICAiaW5pdGlhbCI6IDEKICB9CiB9LAogInRyYW5zaXRpb25zIjogewogICJhZGQiOiB7CiAgICJ4IjogMjQ2LAogICAieSI6IDIwNQogIH0sCiAgInN1YiI6IHsKICAgIngiOiA0NzUsCiAgICJ5IjogMjA2CiAgfSwKICAiYmFyIjogewogICAieCI6IDI0MywKICAgInkiOiA0MzkKICB9LAogICJiYXoiOiB7CiAgICJ4IjogNDc2LAogICAieSI6IDQ0MwogIH0KIH0sCiAiYXJjcyI6IFsKICB7CiAgICJzb3VyY2UiOiAiYWRkIiwKICAgInRhcmdldCI6ICJmb28iLAogICAid2VpZ2h0IjogMQogIH0sCiAgewogICAic291cmNlIjogImZvbyIsCiAgICJ0YXJnZXQiOiAic3ViIiwKICAgIndlaWdodCI6IDEKICB9LAogIHsKICAgInNvdXJjZSI6ICJiYXIiLAogICAidGFyZ2V0IjogImZvbyIsCiAgICJ3ZWlnaHQiOiAzLAogICAiaW5oaWJpdCI6IHRydWUKICB9LAogIHsKICAgInNvdXJjZSI6ICJmb28iLAogICAidGFyZ2V0IjogImJheiIsCiAgICJ3ZWlnaHQiOiAxLAogICAiaW5oaWJpdCI6IHRydWUKICB9CiBdCn1QSwECFAAKAAAAAAAkgpxXzxVB0ncCAAB3AgAACgAAAAAAAAAAAAAAAAAAAAAAbW9kZWwuanNvblBLBQYAAAAAAQABADgAAACfAgAAAAA="

func TestNewSvg(t *testing.T) {

	// Load a model from a URL
	mm := metamodel.New("test")
	mm.UnzipUrl(sampleUrl)
	w, h := mm.GetSize()
	x := 110
	y := 80
	marginX := 40
	marginY := 40
	width := w - marginX
	height := h - marginY
	i := image.NewSvgFile("/tmp/test.svg", width, height, x, y, width, height)
	i.Rect(x, y, width, height, "fill: #fff; stroke: #000; stroke-width: 1px;")
	i.Render(mm)
}
