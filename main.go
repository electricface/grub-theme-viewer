package main

import (
	tt "github.com/electricface/grub-theme-viewer/themetxt"

	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
)

var optThemeFile string
var optDraw bool
var optDump bool

var optScreenWidth int
var optScreenHeight int

var globalThemeDir string

func init() {
	flag.StringVar(&optThemeFile, "theme", "", "theme file")
	flag.BoolVar(&optDraw, "draw", false, "draw out.png")
	flag.BoolVar(&optDump, "dump", false, "dump theme")

	flag.IntVar(&optScreenWidth, "width", 1366, "screen width (px)")
	flag.IntVar(&optScreenHeight, "height", 768, "screen height (px)")
}

func testMain() {
	ec := newEvalContent()
	ec.setUnknown("screen-width", 500)
	ec.setUnknown("screen-height", 600)

	root := newRootNode(500, 600)
	dc := gg.NewContext(500, 600)

	c1 := &Node{
		left:   tt.AbsNum(0),
		top:    tt.AbsNum(0),
		width:  tt.RelNum(50),
		height: tt.RelNum(50),
	}

	c11 := &Node{
		left:   tt.RelNum(50),
		top:    tt.RelNum(50),
		width:  tt.RelNum(50),
		height: tt.CombinedNum{50, 10},
	}
	c1.addChild(c11)

	root.addChild(c1)

	root.DrawTo(dc, ec)
	dc.SavePNG("./test.png")
	os.Exit(0)
}

func main() {
	//testMain()
	flag.Parse()

	globalThemeDir = filepath.Dir(optThemeFile)

	theme, err := tt.ParseThemeFile(optThemeFile)
	if err != nil {
		log.Fatal(err)
	}

	if optDump {
		theme.Dump()
	}

	if optDraw {
		// draw
		draw(theme)
	}
}

func draw(theme *tt.Theme) {

	//screenWidth := 1366
	//screenHeight := 768
	//screenWidth = 3000
	//screenHeight = 2000

	ec := newEvalContent()
	ec.setUnknown("screen-width", float64(optScreenWidth))
	ec.setUnknown("screen-height", float64(optScreenHeight))

	root := themeToNodeTree(theme, optScreenWidth, optScreenHeight)
	//textFontSize := 32
	ctx := gg.NewContext(optScreenWidth, optScreenHeight)
	fontFile := "/usr/share/fonts/truetype/noto/NotoSans-Regular.ttf"
	err := ctx.LoadFontFace(fontFile, 32)
	if err != nil {
		log.Fatal(err)
	}
	// 画背景
	root.draw = func(n *Node, ctx *gg.Context, ec *EvalContext) {
		n.drawImage(ctx, ec, "background.png")
	}

	root.DrawTo(ctx, ec)
	ctx.SavePNG("./out.png")
}

func getResourceFile(name string) string {
	// TODO
	//dir := "/boot/grub/themes/deepin-green"
	dir := globalThemeDir
	return filepath.Join(dir, name)
}

func themeToNodeTree(theme *tt.Theme, w, h int) *Node {
	root := newRootNode(w, h)
	for _, comp := range theme.Components {
		if comp.Id == "boot_menu" {
			root.addChild(compBootMenuToNode(comp, root))
		}
	}
	return root
}

type CompCommon struct {
	left   tt.Length
	top    tt.Length
	width  tt.Length
	height tt.Length
	id     string

	node *Node
}
