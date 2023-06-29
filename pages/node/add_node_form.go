package page_node

import (
	"fmt"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/node_manager"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageAddNodeForm struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonAddNode *components.Button
	txtEndpoint   *components.TextField
	txtName       *components.TextField
	submitting    bool

	list *widget.List
}

var _ router.Page = &PageAddNodeForm{}

func NewPageAddNodeForm() *PageAddNodeForm {
	th := app_instance.Theme

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .5, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .5, ease.OutCubic),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	addIcon, _ := widget.NewIcon(icons.ContentAdd)
	buttonAddNode := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            addIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonAddNode.Label.Alignment = text.Middle
	buttonAddNode.Style.Font.Weight = font.Bold

	txtName := components.NewTextField(th, lang.Translate("Name"), "Dero NFTs")
	txtEndpoint := components.NewTextField(th, lang.Translate("Endpoint"), "wss://node.deronfts.com/ws")

	return &PageAddNodeForm{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonAddNode: buttonAddNode,
		txtName:       txtName,
		txtEndpoint:   txtEndpoint,

		list: list,
	}
}

func (p *PageAddNodeForm) IsActive() bool {
	return p.isActive
}

func (p *PageAddNodeForm) Enter() {
	p.isActive = true
	page_instance.header.SetTitle(lang.Translate("Add Node"))
	p.animationEnter.Start()
	p.animationLeave.Reset()
}

func (p *PageAddNodeForm) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageAddNodeForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.buttonAddNode.Clickable.Clicked() {
		p.submitForm()
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.txtName.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtEndpoint.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonAddNode.Text = lang.Translate("ADD NODE")
			return p.buttonAddNode.Layout(gtx, th)
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	if p.txtName.Input.Clickable.Clicked() {
		p.list.ScrollTo(0)
	}

	if p.txtEndpoint.Input.Clickable.Clicked() {
		p.list.ScrollTo(0)
	}

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageAddNodeForm) submitForm() {
	if p.submitting {
		return
	}

	p.submitting = true

	go func() {
		setError := func(err error) {
			p.submitting = false
			notification_modals.ErrorInstance.SetText("Error", err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		}

		txtName := p.txtName.Editor()
		txtEndpoint := p.txtEndpoint.Editor()

		if txtName.Text() == "" {
			setError(fmt.Errorf("enter name"))
			return
		}

		if txtEndpoint.Text() == "" {
			setError(fmt.Errorf("enter endpoint"))
			return
		}

		_, err := walletapi.TestConnect(txtEndpoint.Text())
		if err != nil {
			setError(err)
			return
		}

		err = node_manager.AddNode(node_manager.NodeConnection{
			Name:     txtName.Text(),
			Endpoint: txtEndpoint.Text(),
		})
		if err != nil {
			setError(err)
			return
		}

		p.submitting = false
		notification_modals.SuccessInstance.SetText(lang.Translate("Success"), "new noded added")
		notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		page_instance.pageRouter.SetCurrent(PAGE_SELECT_NODE)
	}()
}
