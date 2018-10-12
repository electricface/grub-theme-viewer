package main

import (
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

func (n *Node) drawText(ctx *gg.Context, ec *EvalContext, str string) {
	x := n.getLeft().Eval(ec)
	y := n.getTop().Eval(ec)
	ctx.SetRGB(1, 0, 0)
	ctx.DrawStringAnchored(str, x, y, 0, 1)
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
