package prefabs

import (
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/ui/components"
)

type SelectModal struct {
	SelectedKey string

	list     *widget.List
	modal    *components.Modal
	selected bool
}

func NewSelectModal(w *app.Window) *SelectModal {
	list := new(widget.List)
	list.Axis = layout.Vertical

	modal := components.NewModal(w, components.ModalStyle{
		CloseOnOutsideClick: true,
		CloseOnInsideClick:  false,
		Direction:           layout.S,
		BgColor:             color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		Rounded:             components.UniformRounded(unit.Dp(10)),
		Inset:               layout.UniformInset(25),
		Animation:           components.NewModalAnimationUp(),
		Backdrop:            components.NewModalBackground(),
	})

	return &SelectModal{
		modal: modal,
		list:  list,
	}
}

func (l *SelectModal) Selected() bool {
	return l.selected
}

func (l *SelectModal) Layout(gtx layout.Context, th *material.Theme, items []*SelectListItem) layout.Dimensions {
	l.selected = false
	return l.modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(10), Bottom: unit.Dp(10),
			Left: unit.Dp(10), Right: unit.Dp(10),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.Y = gtx.Dp(200)
			listStyle := material.List(th, l.list)
			listStyle.AnchorStrategy = material.Overlay

			return listStyle.Layout(gtx, len(items), func(gtx layout.Context, index int) layout.Dimensions {
				if items[index].clickable.Clicked() {
					l.SelectedKey = items[index].key
					l.selected = true
					op.InvalidateOp{}.Add(gtx.Ops)
				}

				return items[index].Layout(gtx, th, index)
			})
		})
	})
}

type SelectListItem struct {
	key       string
	render    layout.ListElement
	clickable *widget.Clickable
}

func NewSelectListItem(key string, render layout.ListElement) *SelectListItem {
	return &SelectListItem{
		key:       key,
		render:    render,
		clickable: new(widget.Clickable),
	}
}

func (c *SelectListItem) Layout(gtx layout.Context, th *material.Theme, index int) layout.Dimensions {
	dims := c.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return c.render(gtx, index)
		})
	})

	if c.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)

		paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 100},
			clip.UniformRRect(
				image.Rectangle{Max: image.Pt(dims.Size.X, dims.Size.Y)},
				gtx.Dp(15),
			).Op(gtx.Ops),
		)
	}

	return dims
}
