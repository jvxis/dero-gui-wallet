package page_node

import (
	"fmt"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_db"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/confirm_modal"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/node_manager"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageEditNodeForm struct {
	isActive            bool
	headerPageAnimation *prefabs.PageHeaderAnimation

	buttonEdit   *components.Button
	buttonDelete *components.Button
	txtEndpoint  *prefabs.TextField
	txtName      *prefabs.TextField
	nodeConn     app_db.NodeConnection

	list *widget.List
}

var _ router.Page = &PageEditNodeForm{}

func NewPageEditNodeForm() *PageEditNodeForm {
	list := new(widget.List)
	list.Axis = layout.Vertical

	saveIcon, _ := widget.NewIcon(icons.ContentSave)
	loadingIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	buttonEdit := components.NewButton(components.ButtonStyle{
		Rounded:     components.UniformRounded(unit.Dp(5)),
		Icon:        saveIcon,
		TextSize:    unit.Sp(14),
		IconGap:     unit.Dp(10),
		Inset:       layout.UniformInset(unit.Dp(10)),
		Animation:   components.NewButtonAnimationDefault(),
		LoadingIcon: loadingIcon,
	})
	buttonEdit.Label.Alignment = text.Middle
	buttonEdit.Style.Font.Weight = font.Bold

	txtName := prefabs.NewTextField()
	txtEndpoint := prefabs.NewTextField()

	deleteIcon, _ := widget.NewIcon(icons.ActionDelete)
	buttonDelete := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      deleteIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonDelete.Label.Alignment = text.Middle
	buttonDelete.Style.Font.Weight = font.Bold

	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_EDIT_NODE_FORM)
	return &PageEditNodeForm{
		headerPageAnimation: headerPageAnimation,

		buttonEdit:   buttonEdit,
		buttonDelete: buttonDelete,
		txtName:      txtName,
		txtEndpoint:  txtEndpoint,

		list: list,
	}
}

func (p *PageEditNodeForm) IsActive() bool {
	return p.isActive
}

func (p *PageEditNodeForm) Enter() {
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)
	page_instance.header.Title = func() string { return lang.Translate("Edit Node") }

	p.txtEndpoint.SetValue(p.nodeConn.Endpoint)
	p.txtName.SetValue(p.nodeConn.Name)
}

func (p *PageEditNodeForm) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)
}

func (p *PageEditNodeForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	defer p.headerPageAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

	if p.buttonEdit.Clicked(gtx) {
		p.submitForm(gtx)
	}

	if p.buttonDelete.Clicked(gtx) {
		go func() {
			yesChan := confirm_modal.Instance.Open(confirm_modal.ConfirmText{})

			if <-yesChan {
				err := p.removeNode()
				if err != nil {
					notification_modal.Open(notification_modal.Params{
						Type:  notification_modal.ERROR,
						Title: lang.Translate("Error"),
						Text:  err.Error(),
					})
				} else {
					notification_modal.Open(notification_modal.Params{
						Type:       notification_modal.SUCCESS,
						Title:      lang.Translate("Success"),
						Text:       lang.Translate("Node deleted."),
						CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
					})
					page_instance.header.GoBack()
				}
			}
		}()
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.txtName.Layout(gtx, th, lang.Translate("Name"), "Dero NFTs")
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtEndpoint.Layout(gtx, th, lang.Translate("Endpoint"), "wss://node.deronfts.com/ws")
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonEdit.Text = lang.Translate("SAVE")
			p.buttonEdit.Style.Colors = theme.Current.ButtonPrimaryColors
			return p.buttonEdit.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return prefabs.Divider(gtx, unit.Dp(5))
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonDelete.Text = lang.Translate("DELETE NODE")
			p.buttonDelete.Style.Colors = theme.Current.ButtonDangerColors
			return p.buttonDelete.Layout(gtx, th)
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	if p.txtName.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(0)
	}

	if p.txtEndpoint.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(1)
	}

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: theme.PagePadding, Right: theme.PagePadding,
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageEditNodeForm) removeNode() error {
	err := app_db.DelNodeConnection(p.nodeConn.ID)
	if err != nil {
		return err
	}

	if node_manager.CurrentNode != nil {
		if node_manager.CurrentNode.ID == p.nodeConn.ID {
			node_manager.Set(nil, true)
		}
	}

	return nil
}

func (p *PageEditNodeForm) submitForm(gtx layout.Context) {
	p.buttonEdit.SetLoading(true)
	go func() {
		setError := func(err error) {
			p.buttonEdit.SetLoading(false)
			notification_modal.Open(notification_modal.Params{
				Type:  notification_modal.ERROR,
				Title: lang.Translate("Error"),
				Text:  err.Error(),
			})
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

		_, err := TestConnect(txtEndpoint.Text())
		if err != nil {
			setError(err)
			return
		}

		node := app_db.NodeConnection{
			ID:          p.nodeConn.ID,
			Name:        txtName.Text(),
			Endpoint:    txtEndpoint.Text(),
			OrderNumber: p.nodeConn.OrderNumber,
		}

		err = app_db.UpdateNodeConnection(node)
		if err != nil {
			setError(err)
			return
		}

		p.buttonEdit.SetLoading(false)
		notification_modal.Open(notification_modal.Params{
			Type:       notification_modal.SUCCESS,
			Title:      lang.Translate("Success"),
			Text:       lang.Translate("Data saved."),
			CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
		})
		page_instance.header.GoBack()
	}()
}
