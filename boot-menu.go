package main

import (
	tt "github.com/electricface/grub-theme-viewer/themetxt"

	"github.com/fogleman/gg"
)

type menuItem struct {
	icon string
	text string
}

var menuItems = []*menuItem{
	{
		icon: "deepin",
		text: "Deepin GNU/Linux",
	},
	{
		text: "Advanced options for Deepin GNU/Linux",
	},
	{
		icon: "windows",
		text: "Window XP",
	},
	{
		text: "System setup",
	},
}

type BootMenu struct {
	CompCommon

	itemFont                string
	itemColor               string
	itemPixmapStyle         string
	selectedItemFont        string
	selectedItemColor       string
	selectedItemPixmapStyle string

	itemHeight  tt.Length
	itemPadding tt.Length
	itemSpacing tt.Length

	iconWidth     tt.Length
	iconHeight    tt.Length
	itemIconSpace tt.Length
}

func getLengthExprHorizontal(l tt.Length, n *Node) Expr {
	return getLengthExpr(l, n.getWidth())
}

func getLengthExprVertical(l tt.Length, n *Node) Expr {
	return getLengthExpr(l, n.getHeight())
}

func (bm *BootMenu) getItemHeight() Expr {
	return getLengthExprVertical(bm.itemHeight, bm.node)
}
func (bm *BootMenu) getItemPadding() Expr {
	return getLengthExprVertical(bm.itemPadding, bm.node)
}
func (bm *BootMenu) getItemSpacing() Expr {
	return getLengthExprVertical(bm.itemSpacing, bm.node)
}

func (bm *BootMenu) getIconWidth() Expr {
	return getLengthExprHorizontal(bm.iconWidth, bm.node)
}

func (bm *BootMenu) getIconHeight() Expr {
	return getLengthExprVertical(bm.iconHeight, bm.node)
}

func (bm *BootMenu) getItemIconSpace() Expr {
	return getLengthExprHorizontal(bm.itemIconSpace, bm.node)
}

func (cc *CompCommon) fillCommonOptions(comp *tt.Component) {
	var ok bool
	cc.left, ok = comp.GetPropLength("left")
	if !ok {
		cc.left = tt.AbsNum(0)
	}
	cc.node.left = cc.left

	cc.top, ok = comp.GetPropLength("top")
	if !ok {
		cc.top = tt.AbsNum(0)
	}
	cc.node.top = cc.top

	cc.width, ok = comp.GetPropLength("width")
	if !ok {
		cc.width = tt.AbsNum(0)
	}
	cc.node.width = cc.width

	cc.height, ok = comp.GetPropLength("top")
	if !ok {
		cc.height = tt.AbsNum(0)
	}
	cc.node.height = cc.height
}

func newBootMenu(comp *tt.Component, parent *Node) *BootMenu {
	bm := &BootMenu{}
	bm.node = &Node{
		parent: parent,
	}

	bm.fillCommonOptions(comp)
	var ok bool

	bm.itemHeight, ok = comp.GetPropLength("item_height")
	if !ok {
		// set default value
		bm.itemHeight = tt.AbsNum(42)
	}

	bm.itemPadding, ok = comp.GetPropLength("item_padding")
	if !ok {
		bm.itemPadding = tt.AbsNum(14)
	}

	bm.itemSpacing, ok = comp.GetPropLength("item_spacing")
	if !ok {
		bm.itemSpacing = tt.AbsNum(16)
	}

	bm.iconWidth, ok = comp.GetPropLength("icon_width")
	if !ok {
		bm.iconWidth = tt.AbsNum(32)
	}

	bm.iconHeight, ok = comp.GetPropLength("icon_height")
	if !ok {
		bm.iconHeight = tt.AbsNum(32)
	}

	bm.itemIconSpace, ok = comp.GetPropLength("item_icon_space")
	if !ok {
		bm.itemIconSpace = tt.AbsNum(4)
	}

	return bm
}

func compBootMenuToNode(comp *tt.Component, parent *Node) *Node {
	bm := newBootMenu(comp, parent)
	textFontSize := 32
	bmNode := bm.node

	y := bm.getItemPadding()

	//itemWidth := bmNode.getWidth() - (2 * float64(itemPadding)) - 2
	itemWidthExpr := sub(bmNode.getWidth(),
		mul(AbsNum(2), bm.getItemPadding()))
	itemWidthExpr = sub(itemWidthExpr, AbsNum(2))

	for i := 0; i < 4; i++ {
		// add item
		item := &Node{
			left:      bm.itemPadding,
			topExpr:   y,
			widthExpr: itemWidthExpr,
			height:    bm.itemHeight,
		}

		//iconTop := float64(itemHeight-iconHeight) / 2
		iconTopExpr := div(sub(bm.getItemHeight(), bm.getIconHeight()), AbsNum(2))

		icon := &Node{
			left:    tt.AbsNum(0),
			topExpr: iconTopExpr,

			width:  bm.iconWidth,
			height: bm.iconHeight,
		}
		idx := i
		icon.draw = func(n *Node, ctx *gg.Context, ec *EvalContext) {
			iconName := menuItems[idx].icon
			n.drawImage(ctx, ec, "icons/"+iconName+".png")
		}

		//textTop := float64(itemHeight-textFontSize) / 2
		textTopExpr := div(sub(bm.getItemHeight(), AbsNum(textFontSize)), AbsNum(2))

		//textWidth = tt.AbsNum(int(itemWidth) - iconWidth - itemIconSpace),
		textWidthExpr := sub(sub(itemWidthExpr, bm.getIconWidth()),
			bm.getItemIconSpace())

		// textLeft = iconWidth + itemIconSpace
		textLeftExpr := add(bm.getIconWidth(), bm.getItemIconSpace())
		text := &Node{
			leftExpr:  textLeftExpr,
			topExpr:   textTopExpr,
			widthExpr: textWidthExpr,
			height:    tt.AbsNum(textFontSize),
		}
		text.draw = func(n *Node, ctx *gg.Context, ec *EvalContext) {
			textStr := menuItems[idx].text
			n.drawText(ctx, ec, textStr)
		}

		item.addChild(icon)
		item.addChild(text)

		//y += itemHeight + itemSpacing
		y = add(y, add(bm.getItemHeight(), bm.getItemSpacing()))

		bmNode.addChild(item)

	}
	return bmNode
}
