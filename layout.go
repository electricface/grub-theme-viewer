package main

import (
	"image/color"
	"log"

	tt "github.com/electricface/grub-theme-viewer/themetxt"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

type Node struct {
	parent   *Node
	Children []*Node

	left   tt.Length
	top    tt.Length
	width  tt.Length
	height tt.Length

	leftExpr   Expr
	topExpr    Expr
	widthExpr  Expr
	heightExpr Expr

	draw func(n *Node, ctx *gg.Context, ec *EvalContext)
}

func getLengthExpr(l tt.Length, val Expr) Expr {
	switch ll := l.(type) {
	case tt.AbsNum:
		return AbsNum(int(ll))
	case tt.RelNum:
		// (val * (ll  / 100))
		return mul(val, div(AbsNum(int(ll)), AbsNum(100)))

	case tt.CombinedNum:
		a := mul(val, div(AbsNum(int(ll.Rel)), AbsNum(100)))
		// TODO add or sub
		return sub(a, AbsNum(ll.Abs))
	}
	panic("not expect")
	return nil
}

func (n *Node) getLeft() Expr {
	if n.parent == nil {
		// root
		return AbsNum(0)
	}

	pl := n.parent.getLeft()
	if n.leftExpr != nil {
		return add(pl, n.leftExpr)
	}

	pw := n.parent.getWidth()
	return add(pl, getLengthExpr(n.left, pw))
}

func (n *Node) getTop() Expr {
	if n.parent == nil {
		// root
		return AbsNum(0)
	}

	pt := n.parent.getTop()
	if n.topExpr != nil {
		return add(pt, n.topExpr)
	}

	ph := n.parent.getHeight()
	return add(pt, getLengthExpr(n.top, ph))
}

func (n *Node) getWidth() Expr {
	if n.widthExpr != nil {
		return n.widthExpr
	}

	if n.parent == nil {
		// root
		return &Unknown{name: "screen-width"}
	}
	pw := n.parent.getWidth()
	return getLengthExpr(n.width, pw)
}

func (n *Node) getHeight() Expr {
	if n.heightExpr != nil {
		return n.heightExpr
	}

	if n.parent == nil {
		// root
		return &Unknown{name: "screen-height"}
	}
	ph := n.parent.getHeight()
	return getLengthExpr(n.height, ph)
}

//func (n *Node) getExtent() (x, y, w, h float64) {
//	x = n.getTop()
//	y = n.getLeft()
//	w = n.getWidth()
//	h = n.getHeight()
//	return
//}

func (n *Node) addChild(child *Node) {
	child.parent = n
	n.Children = append(n.Children, child)
}

func (n *Node) drawImage(ctx *gg.Context, ec *EvalContext, name string) error {
	img, err := gg.LoadImage(getResourceFile(name))
	if err != nil {
		return err
	}

	x := n.getLeft().Eval(ec)
	y := n.getTop().Eval(ec)
	width := n.getWidth().Eval(ec)
	height := n.getHeight().Eval(ec)

	img = resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	ctx.DrawImage(img, int(x), int(y))
	return nil
}

func (n *Node) drawStyleBox(ctx *gg.Context, ec *EvalContext, name string) {
	if name == "" {
		return
	}

	x := int(n.getLeft().Eval(ec))
	y := int(n.getTop().Eval(ec))
	width := int(n.getWidth().Eval(ec))
	height := int(n.getHeight().Eval(ec))

	color1 := "#f9f806"
	color2 := "#f97306"

	// nw
	imgNW, err := gg.LoadImage(getResourceFile(getPixmapName(name, styleBoxNorthwest)))
	if err != nil {
		log.Println(err)
	} else {
		imgNWWidth := imgNW.Bounds().Dx()
		imgNWHeight := imgNW.Bounds().Dy()
		ctx.DrawImage(imgNW, x-imgNWWidth, y-imgNWHeight)

		if optDrawOutline {
			ctx.SetHexColor(color1)
			ctx.DrawRectangle(float64(x-imgNWWidth), float64(y-imgNWHeight),
				float64(imgNWWidth), float64(imgNWHeight))
			ctx.Stroke()
		}
	}

	// n
	imgN, err := gg.LoadImage(getResourceFile(getPixmapName(name, styleBoxNorth)))
	if err != nil {
		log.Println(err)
	} else {
		imgNHeight := imgN.Bounds().Dy()
		imgN = resize.Resize(uint(width), uint(imgNHeight), imgN, resize.Lanczos3)
		ctx.DrawImage(imgN, x, y-imgNHeight)

		if optDrawOutline {
			ctx.SetHexColor(color2)
			ctx.DrawRectangle(float64(x), float64(y-imgNHeight),
				float64(width), float64(imgNHeight))
			ctx.Stroke()
		}
	}

	// ne
	imgNE, err := gg.LoadImage(getResourceFile(getPixmapName(name, styleBoxNortheast)))
	if err != nil {
		log.Println(err)
	} else {
		imgNEWidth := imgNE.Bounds().Dx()
		imgNEHeight := imgNE.Bounds().Dy()
		ctx.DrawImage(imgNE, x+width, y-imgNEHeight)

		if optDrawOutline {
			ctx.SetHexColor(color1)
			ctx.DrawRectangle(float64(x+width), float64(y-imgNEHeight),
				float64(imgNEWidth), float64(imgNEHeight))
			ctx.Stroke()
		}
	}

	// w
	imgW, err := gg.LoadImage(getResourceFile(getPixmapName(name, styleBoxWest)))
	if err != nil {
		log.Println(err)
	} else {
		imgWWidth := imgW.Bounds().Dx()
		imgW = resize.Resize(uint(imgWWidth), uint(height), imgW, resize.Lanczos3)
		ctx.DrawImage(imgW, x-imgWWidth, y)

		if optDrawOutline {
			ctx.SetHexColor(color2)
			ctx.DrawRectangle(float64(x-imgWWidth), float64(y),
				float64(uint(imgWWidth)), float64(uint(height)))
			ctx.Stroke()
		}
	}

	// c
	imgC, err := gg.LoadImage(getResourceFile(getPixmapName(name, styleBoxCenter)))
	if err != nil {
		log.Println(err)
	} else {
		imgC = resize.Resize(uint(width), uint(height), imgC, resize.Lanczos3)
		ctx.DrawImage(imgC, x, y)
	}

	// e
	imgE, err := gg.LoadImage(getResourceFile(getPixmapName(name, styleBoxEast)))
	if err != nil {
		log.Println(err)
	} else {
		imgEWidth := imgE.Bounds().Dx()
		imgE = resize.Resize(uint(imgEWidth), uint(height), imgE, resize.Lanczos3)
		ctx.DrawImage(imgE, x+width, y)

		if optDrawOutline {
			ctx.SetHexColor(color2)
			ctx.DrawRectangle(float64(x+width), float64(y),
				float64(uint(imgEWidth)), float64(uint(height)))
			ctx.Stroke()
		}
	}

	// sw
	imgSW, err := gg.LoadImage(getResourceFile(getPixmapName(name, styleBoxSouthwest)))
	if err != nil {
		log.Println(err)
	} else {
		imgSWWidth := imgSW.Bounds().Dx()
		imgSWHeight := imgSW.Bounds().Dy()
		ctx.DrawImage(imgSW, x-imgSWWidth, y+height)

		if optDrawOutline {
			ctx.SetHexColor(color1)
			ctx.DrawRectangle(float64(x-imgSWWidth), float64(y+height),
				float64(imgSWWidth), float64(imgSWHeight))
			ctx.Stroke()
		}
	}

	// s
	imgS, err := gg.LoadImage(getResourceFile(getPixmapName(name, styleBoxSouth)))
	if err != nil {
		log.Println(err)
	} else {
		imgSHeight := imgS.Bounds().Dy()
		imgS = resize.Resize(uint(width), uint(imgSHeight), imgS, resize.Lanczos3)
		ctx.DrawImage(imgS, x, y+height)

		if optDrawOutline {
			ctx.SetHexColor(color2)
			ctx.DrawRectangle(float64(x), float64(y+height),
				float64(width), float64(imgSHeight))
			ctx.Stroke()
		}
	}

	// se
	imgSE, err := gg.LoadImage(getResourceFile(getPixmapName(name, styleBoxSoutheast)))
	if err != nil {
		log.Println(err)
	} else {
		ctx.DrawImage(imgSE, x+width, y+height)

		if optDrawOutline {
			imgSEWidth := imgSE.Bounds().Dx()
			imgSEHeight := imgSE.Bounds().Dy()
			ctx.SetHexColor(color1)
			ctx.DrawRectangle(float64(x+width), float64(y+height),
				float64(imgSEWidth), float64(imgSEHeight))
			ctx.Stroke()
		}
	}
}

func (n *Node) drawText(ctx *gg.Context, ec *EvalContext, str string, color color.Color, fontSize int) {
	x := n.getLeft().Eval(ec)
	y := n.getTop().Eval(ec)
	ctx.SetColor(color)

	ctx.LoadFontFace(globalFontFile, float64(fontSize))
	ctx.DrawStringAnchored(str, x, y, 0, 1)
}

func (n *Node) drawText1(ctx *gg.Context, ec *EvalContext, str string, color color.Color, fontSize int, width float64, align gg.Align) {
	x := n.getLeft().Eval(ec)
	y := n.getTop().Eval(ec)
	ctx.SetColor(color)

	ctx.LoadFontFace(globalFontFile, float64(fontSize))
	//ctx.DrawStringAnchored(str, x, y, 0, 1)
	ctx.DrawStringWrapped(str, x, y, 0, 1, width, 0, align)
}

func (n *Node) DrawTo(ctx *gg.Context, ec *EvalContext) {
	if optDrawOutline {
		x := n.getLeft().Eval(ec)
		y := n.getTop().Eval(ec)
		w := n.getWidth().Eval(ec)
		h := n.getHeight().Eval(ec)
		ctx.DrawRectangle(x, y, w, h)
		ctx.SetRGB(1, 0, 0)
		ctx.Stroke()
	}

	if n.draw != nil {
		n.draw(n, ctx, ec)
	}

	for _, c := range n.Children {
		c.DrawTo(ctx, ec)
	}
}
