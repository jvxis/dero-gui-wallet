package page_wallet_select

import (
	"fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageCreateWalletHexSeedForm struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	listStyle material.ListStyle

	txtHexSeed         *components.TextField
	txtWalletName      *components.TextField
	txtPassword        *components.TextField
	txtConfirmPassword *components.TextField
	buttonCreate       *components.Button
}

var _ router.Container = &PageCreateWalletHexSeedForm{}

func NewPageCreateWalletHexSeedForm() *PageCreateWalletHexSeedForm {
	th := app_instance.Current.Theme
	list := new(widget.List)
	list.Axis = layout.Vertical
	listStyle := material.List(th, list)
	listStyle.AnchorStrategy = material.Overlay

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .5, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .5, ease.OutCubic),
	))

	txtWalletName := components.NewTextField(th, "Wallet Name", "")
	txtPassword := components.NewTextField(th, "Password", "")
	txtPassword.EditorStyle.Editor.Mask = rune(42)
	txtConfirmPassword := components.NewTextField(th, "Confirm Password", "")
	txtConfirmPassword.EditorStyle.Editor.Mask = rune(42)

	txtHexSeed := components.NewTextField(th, "Hex Seed", "Enter hex seed of 64 chars")

	iconCreate, _ := widget.NewIcon(icons.ContentAddBox)
	buttonCreate := components.NewButton(components.ButtonStyle{
		Rounded:         unit.Dp(5),
		Text:            "RECOVER WALLET",
		Icon:            iconCreate,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})

	return &PageCreateWalletHexSeedForm{
		listStyle:      listStyle,
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		txtHexSeed:         txtHexSeed,
		txtWalletName:      txtWalletName,
		txtPassword:        txtPassword,
		txtConfirmPassword: txtConfirmPassword,
		buttonCreate:       buttonCreate,
	}
}

func (p *PageCreateWalletHexSeedForm) Enter() {
	page_instance.header.LabelTitle.Text = "Recover from Hex Seed"
	p.isActive = true
	p.animationEnter.Start()
	p.animationLeave.Reset()
}

func (p *PageCreateWalletHexSeedForm) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageCreateWalletHexSeedForm) IsActive() bool {
	return p.isActive
}

func (p *PageCreateWalletHexSeedForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	{
		state := p.animationEnter.Update(gtx)
		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	{
		state := p.animationLeave.Update(gtx)
		if state.Finished {
			p.isActive = false
			op.InvalidateOp{}.Add(gtx.Ops)
		}

		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	if p.buttonCreate.Clickable.Clicked() {
		name := p.txtWalletName.EditorStyle.Editor.Text()
		password := p.txtPassword.EditorStyle.Editor.Text()
		//confirmPassword := p.txtConfirmPassword.EditorStyle.Editor.Text()
		err := wallet_manager.Instance.CreateWalletFromSeed(name, password, "")
		fmt.Println(err)
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.txtHexSeed.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtWalletName.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtPassword.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtConfirmPassword.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.buttonCreate.Layout(gtx, th)
		},
	}

	return p.listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}