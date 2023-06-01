package page_settings

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
)

type Page struct {
	isActive bool
}

var _ router.Container = &Page{}

func NewPage() *Page {
	return &Page{}
}

func (p *Page) IsActive() bool {
	return p.isActive
}

func (p *Page) Enter() {
	app_instance.Current.BottomBar.SetActive("settings")
	p.isActive = true
}

func (p *Page) Leave() {
	p.isActive = false
}

func (p *Page) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return SettingsItem{Title: "App Dir"}.Layout(gtx, th, func(gtx layout.Context) layout.Dimensions {
				dir := settings.Instance.AppDir
				label := material.Label(th, unit.Sp(16), dir)
				return label.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return SettingsItem{Title: "Node Dir"}.Layout(gtx, th, func(gtx layout.Context) layout.Dimensions {
				dir := settings.Instance.NodeDir
				label := material.Label(th, unit.Sp(16), dir)
				return label.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return SettingsItem{Title: "Wallets Dir"}.Layout(gtx, th, func(gtx layout.Context) layout.Dimensions {
				dir := settings.Instance.WalletsDir
				label := material.Label(th, unit.Sp(16), dir)
				return label.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return SettingsItem{Title: "Version"}.Layout(gtx, th, func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(16), settings.Version)
				return label.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return SettingsItem{Title: "Git Version"}.Layout(gtx, th, func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(16), settings.GitVersion)
				return label.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return SettingsItem{Title: "Build Time"}.Layout(gtx, th, func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(16), settings.BuildTime)
				return label.Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return app_instance.Current.BottomBar.Layout(gtx, th)
		}),
	)
}

type SettingsItem struct {
	Title string
}

func (s SettingsItem) Layout(gtx layout.Context, th *material.Theme, w layout.Widget) layout.Dimensions {
	dims := layout.Inset{
		Top: unit.Dp(10), Bottom: unit.Dp(10),
		Left: unit.Dp(10), Right: unit.Dp(10),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(18), s.Title)
				label.Font.Weight = font.Bold
				return label.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return w(gtx)
			}),
		)
	})

	cl := clip.Rect{Max: image.Pt(dims.Size.X, gtx.Dp(1))}.Push(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{A: 50}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	cl.Pop()

	return dims
}
