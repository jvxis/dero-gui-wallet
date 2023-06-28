package page_wallet

import (
	"bytes"
	"image"
	"log"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/g45t345rt/g45w/utils"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type PageReceiveForm struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	list      *widget.List
	labelAddr material.EditorStyle
	addrImage *components.Image
}

var _ router.Page = &PageReceiveForm{}

func NewPageReceiveForm() *PageReceiveForm {
	th := app_instance.Theme

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	editor := new(widget.Editor)
	labelAddr := material.Editor(th, editor, "")
	labelAddr.TextSize = unit.Sp(16)
	labelAddr.Editor.Alignment = text.Middle
	labelAddr.Font.Weight = font.Bold
	labelAddr.Editor.ReadOnly = true

	return &PageReceiveForm{
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		list:           list,
		labelAddr:      labelAddr,
	}
}

func (p *PageReceiveForm) IsActive() bool {
	return p.isActive
}

func (p *PageReceiveForm) Enter() {
	p.isActive = true
	p.animationEnter.Start()
	p.animationLeave.Reset()

	// gio ui does not implement character break yet https://todo.sr.ht/~eliasnaur/gio/467
	addr := "dero1qyvzwypmgqrqpsr8xlz209mwr6sz8fu9a4fphkpnesg29du40zw22qqpm2nkv"
	splitAddr := utils.SplitString(addr, 22)
	addr = strings.Join(splitAddr, "\n")

	imgBytes, err := qrcode.Encode(addr, qrcode.Medium, 256)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(bytes.NewBuffer(imgBytes))
	if err != nil {
		log.Fatal(err)
	}

	p.addrImage = &components.Image{
		Src: paint.NewImageOp(img),
		Fit: components.Contain,
	}

	p.labelAddr.Editor.SetText(addr)
}

func (p *PageReceiveForm) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageReceiveForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	{
		state := p.animationEnter.Update(gtx)
		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	{
		state := p.animationLeave.Update(gtx)
		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}

		if state.Finished {
			p.isActive = false
			op.InvalidateOp{}.Add(gtx.Ops)
		}
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.labelAddr.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.Y = gtx.Dp(250)
				return p.addrImage.Layout(gtx)
			})
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}
