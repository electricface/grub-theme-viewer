package main

import (
	"image/color"
	"strings"

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

	visible         bool
	menuPixmapStyle string

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

	scrollbar      bool
	scrollbarWidth tt.Length
	scrollbarFrame string
	scrollbarThumb string
}

func (bm *BootMenu) getItemHeight() Expr {
	return AbsNum(bm.itemHeight.(tt.AbsNum))
}

func (bm *BootMenu) getItemPadding() Expr {
	return AbsNum(bm.itemPadding.(tt.AbsNum))
}

func (bm *BootMenu) getItemSpacing() Expr {
	return AbsNum(bm.itemSpacing.(tt.AbsNum))
}

func (bm *BootMenu) getIconWidth() Expr {
	return AbsNum(bm.iconWidth.(tt.AbsNum))
}

func (bm *BootMenu) getIconHeight() Expr {
	return AbsNum(bm.iconHeight.(tt.AbsNum))
}

func (bm *BootMenu) getItemIconSpace() Expr {
	return AbsNum(bm.itemIconSpace.(tt.AbsNum))
}

func (bm *BootMenu) getItemColor() color.Color {
	return parseColor(bm.itemColor)
}

func (bm *BootMenu) getSelectedItemColor() color.Color {
	return parseColor(bm.selectedItemColor)
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
	bm.visible, ok = comp.GetPropBool("visible")
	if !ok {
		bm.visible = true
	}

	bm.menuPixmapStyle, _ = comp.GetPropString("menu_pixmap_style")

	bm.itemFont, ok = comp.GetPropString("item_font")
	if !ok {
		bm.itemFont = "Unknown Regular 16"
	}

	bm.itemColor, ok = comp.GetPropString("item_color")
	if !ok {
		bm.itemColor = "black"
	}

	bm.itemPixmapStyle, _ = comp.GetPropString("item_pixmap_style")

	bm.selectedItemFont, ok = comp.GetPropString("selected_item_font")
	if !ok {
		bm.selectedItemFont = bm.itemFont
	}

	bm.selectedItemColor, ok = comp.GetPropString("selected_item_color")
	if !ok {
		bm.selectedItemColor = bm.itemColor
	}

	bm.selectedItemPixmapStyle, _ = comp.GetPropString("selected_item_pixmap_style")

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

	bm.scrollbar, ok = comp.GetPropBool("scrollbar")
	if !ok {
		bm.scrollbar = true
	}

	bm.scrollbarWidth, ok = comp.GetPropLength("scrollbar_width")
	if !ok {
		bm.scrollbarWidth = tt.AbsNum(16)
	}

	bm.scrollbarFrame, _ = comp.GetPropString("scrollbar_frame")
	bm.scrollbarThumb, _ = comp.GetPropString("scrollbar_thumb")

	return bm
}

func compBootMenuToNode(comp *tt.Component, parent *Node) *Node {
	bm := newBootMenu(comp, parent)
	textFontSize := 32
	bmNode := bm.node

	y := bm.getItemPadding()

	// itemWidth = bootMenu.width - (2 * itemPadding) - 2
	itemWidthExpr := sub(sub(bmNode.getWidth(),
		mul(AbsNum(2), bm.getItemPadding())), AbsNum(2))

	for i := 0; i < 4; i++ {
		// add item
		item := &Node{
			left:      bm.itemPadding,
			topExpr:   y,
			widthExpr: itemWidthExpr,
			height:    bm.itemHeight,
		}

		// select first item
		if i == 0 {
			item.draw = func(n *Node, ctx *gg.Context, ec *EvalContext) {
				n.drawStyleBox(ctx, ec, bm.selectedItemPixmapStyle)
			}
		} else {
			item.draw = func(n *Node, ctx *gg.Context, ec *EvalContext) {
				n.drawStyleBox(ctx, ec, bm.itemPixmapStyle)
			}
		}

		// iconTop = (itemHeight-iconHeight) / 2
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

		// textTop = (itemHeight-textFontSize) / 2
		textTopExpr := div(sub(bm.getItemHeight(), AbsNum(textFontSize)), AbsNum(2))

		// textWidth = itemWidth - iconWidth - itemIconSpace
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

		var textColor color.Color
		if i == 0 {
			textColor = bm.getSelectedItemColor()
		} else {
			textColor = bm.getItemColor()
		}
		text.draw = func(n *Node, ctx *gg.Context, ec *EvalContext) {
			textStr := menuItems[idx].text
			n.drawText(ctx, ec, textStr, textColor)
		}

		item.addChild(icon)
		item.addChild(text)

		// y += itemHeight + itemSpacing
		y = add(y, add(bm.getItemHeight(), bm.getItemSpacing()))

		bmNode.addChild(item)

	}

	bmNode.draw = func(n *Node, ctx *gg.Context, ec *EvalContext) {
		n.drawStyleBox(ctx, ec, bm.menuPixmapStyle)
	}

	return bmNode
}

const (
	styleBoxNorthwest = iota
	styleBoxNorth
	styleBoxNortheast
	styleBoxWest
	styleBoxCenter
	styleBoxEast
	styleBoxSouthwest
	styleBoxSouth
	styleBoxSoutheast
)

func getPixmapName(name string, part int) string {
	var partStr string
	switch part {
	case styleBoxNorthwest:
		partStr = "nw"
	case styleBoxNorth:
		partStr = "n"
	case styleBoxNortheast:
		partStr = "ne"
	case styleBoxWest:
		partStr = "w"
	case styleBoxCenter:
		partStr = "c"
	case styleBoxEast:
		partStr = "e"
	case styleBoxSouthwest:
		partStr = "sw"
	case styleBoxSouth:
		partStr = "s"
	case styleBoxSoutheast:
		partStr = "se"
	}
	return strings.Replace(name, "*", partStr, 1)
}
