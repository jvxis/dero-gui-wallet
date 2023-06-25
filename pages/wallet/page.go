package page_wallet

import (
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/containers/bottom_bar"
	"github.com/g45t345rt/g45w/containers/node_status_bar"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type Page struct {
	isActive       bool
	animationEnter *animation.Animation
	animationLeave *animation.Animation

	header *prefabs.Header

	pageBalanceTokens *PageBalanceTokens
	pageSendForm      *PageSendForm
	pageSCToken       *PageSCToken

	pageRouter *router.Router
}

var _ router.Page = &Page{}

var page_instance *Page

var (
	PAGE_SETTINGS        = "page_settings"
	PAGE_SEND_FORM       = "page_send_form"
	PAGE_RECEIVE_FORM    = "page_receive_form"
	PAGE_BALANCE_TOKENS  = "page_balance_tokens"
	PAGE_ADD_SC_FORM     = "page_add_sc_form"
	PAGE_TXS             = "page_txs"
	PAGE_SC_TOKEN        = "page_sc_token"
	PAGE_REGISTER_WALLET = "page_register_wallet"
	PAGE_CONTACTS        = "page_contacts"
)

func New() *Page {
	th := app_instance.Theme

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .5, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .5, ease.OutCubic),
	))

	pageRouter := router.NewRouter()
	pageBalanceTokens := NewPageBalanceTokens()
	pageRouter.Add(PAGE_BALANCE_TOKENS, pageBalanceTokens)

	pageSendForm := NewPageSendForm()
	pageRouter.Add(PAGE_SEND_FORM, pageSendForm)

	pageReceiveForm := NewPageReceiveForm()
	pageRouter.Add(PAGE_RECEIVE_FORM, pageReceiveForm)

	pageSettings := NewPageSettings()
	pageRouter.Add(PAGE_SETTINGS, pageSettings)

	pageAddSCForm := NewPageAddSCForm()
	pageRouter.Add(PAGE_ADD_SC_FORM, pageAddSCForm)

	pageTxs := NewPageTxs()
	pageRouter.Add(PAGE_TXS, pageTxs)

	pageSCToken := NewPageSCToken()
	pageRouter.Add(PAGE_SC_TOKEN, pageSCToken)

	pageRegisterWallet := NewPageRegisterWallet()
	pageRouter.Add(PAGE_REGISTER_WALLET, pageRegisterWallet)

	pageContacts := NewPageContacts()
	pageRouter.Add(PAGE_CONTACTS, pageContacts)

	labelHeaderStyle := material.Label(th, unit.Sp(22), "")
	labelHeaderStyle.Font.Weight = font.Bold

	header := prefabs.NewHeader(labelHeaderStyle, pageRouter)

	page := &Page{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		header: header,

		pageBalanceTokens: pageBalanceTokens,
		pageSendForm:      pageSendForm,
		pageSCToken:       pageSCToken,

		pageRouter: pageRouter,
	}
	page_instance = page
	return page
}

func (p *Page) IsActive() bool {
	return p.isActive
}

func (p *Page) Enter() {
	openedWallet := wallet_manager.Instance.OpenedWallet
	if openedWallet != nil {
		p.isActive = true
		if p.pageRouter.Current == PAGE_TXS {
			bottom_bar.Instance.SetButtonActive(bottom_bar.BUTTON_TXS)
		} else {
			bottom_bar.Instance.SetButtonActive(bottom_bar.BUTTON_WALLET)
		}

		w := app_instance.Window
		w.Option(app.StatusColor(color.NRGBA{A: 255}))

		p.animationLeave.Reset()
		p.animationEnter.Start()
		p.pageRouter.SetCurrent(PAGE_BALANCE_TOKENS)
	} else {
		app_instance.Router.SetCurrent(app_instance.PAGE_WALLET_SELECT)
	}
}

func (p *Page) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *Page) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	openedWallet := wallet_manager.Instance.OpenedWallet
	if openedWallet == nil {
		return layout.Dimensions{Size: gtx.Constraints.Max}
	}

	if p.pageBalanceTokens.displayBalance.buttonSend.Clickable.Clicked() {
		p.pageRouter.SetCurrent(PAGE_SEND_FORM)
	}

	if p.pageBalanceTokens.displayBalance.buttonReceive.Clickable.Clicked() {
		p.pageRouter.SetCurrent(PAGE_RECEIVE_FORM)
	}

	if p.pageBalanceTokens.tokenBar.buttonAddToken.Clickable.Clicked() {
		p.pageRouter.SetCurrent(PAGE_ADD_SC_FORM)
	}

	//if bottom_bar.Instance.ButtonTxs.Button.Clickable.Clicked() {
	//	p.pageRouter.SetCurrent(PAGE_TXS)
	//}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			{
				state := p.animationEnter.Update(gtx)
				if state.Active {
					defer animation.TransformY(gtx, state.Value).Push(gtx.Ops).Pop()
				}
			}

			{
				state := p.animationLeave.Update(gtx)

				if state.Active {
					defer animation.TransformY(gtx, state.Value).Push(gtx.Ops).Pop()
				}

				if state.Finished {
					p.isActive = false
					op.InvalidateOp{}.Add(gtx.Ops)
				}
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return node_status_bar.Instance.Layout(gtx, th)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, clip.RRect{
						Rect: image.Rectangle{Max: gtx.Constraints.Max},
						NW:   gtx.Dp(15),
						NE:   gtx.Dp(15),
						SE:   gtx.Dp(0),
						SW:   gtx.Dp(0),
					}.Op(gtx.Ops))

					dr := image.Rectangle{Max: gtx.Constraints.Max}
					paint.LinearGradientOp{
						Stop1:  f32.Pt(0, float32(dr.Min.Y)),
						Stop2:  f32.Pt(0, float32(dr.Max.Y)),
						Color1: color.NRGBA{R: 0, G: 0, B: 0, A: 5},
						Color2: color.NRGBA{R: 0, G: 0, B: 0, A: 50},
					}.Add(gtx.Ops)
					defer clip.Rect(dr).Push(gtx.Ops).Pop()
					paint.PaintOp{}.Add(gtx.Ops)

					return layout.Inset{
						Left: unit.Dp(30), Right: unit.Dp(30),
						Top: unit.Dp(30), Bottom: unit.Dp(20),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return p.header.Layout(gtx, th)
					})
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return p.pageRouter.Layout(gtx, th)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return bottom_bar.Instance.Layout(gtx, th)
		}),
	)
}
