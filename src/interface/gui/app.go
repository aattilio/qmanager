package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"qmanager/src/backend/hypervisor"
	"qmanager/src/core"
	"qmanager/src/backend/filesystem"
)

type QManagerUI struct {
	Application fyne.App
	Window      fyne.Window
	
	Connector   *hypervisor.LibvirtHypervisorConnector
	Provisioner *core.AutomatedVirtualMachineProvisioner
	
	VmList      *widget.List
	VmNames     []string
}

func NewQManagerUI() *QManagerUI {
	instance := &QManagerUI{
		Application: app.NewWithID("com.aattilio.qmanager"),
	}
	
	instance.Window = instance.Application.NewWindow("QManager - Hypervisor Excellence")
	instance.Window.Resize(fyne.NewSize(1024, 768))
	
	instance.initializeBackend()
	instance.setupContent()
	
	return instance
}

func (ui *QManagerUI) initializeBackend() {
	connector, err := hypervisor.NewLibvirtHypervisorConnector("qemu:///system")
	if err == nil {
		ui.Connector = connector
	}
	
	diskManager, _ := filesystem.NewVirtualDiskManager("data/vms")
	ui.Provisioner = core.NewAutomatedVirtualMachineProvisioner(
		ui.Connector,
		diskManager,
		"data",
	)
}

func (ui *QManagerUI) setupContent() {
	ui.refreshVmNames()
	
	ui.VmList = widget.NewList(
		func() int {
			return len(ui.VmNames)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.ComputerIcon()),
				widget.NewLabel("Virtual Machine Name"),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(ui.VmNames[id])
		},
	)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), ui.showExpressWizard),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), ui.refreshVmNames),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {}),
	)

	content := container.NewBorder(
		toolbar,
		nil,
		nil,
		nil,
		ui.VmList,
	)

	ui.Window.SetContent(content)
}

func (ui *QManagerUI) refreshVmNames() {
	if ui.Connector != nil {
		names, err := ui.Connector.ListAllVirtualMachineNames()
		if err == nil {
			ui.VmNames = names
			if ui.VmList != nil {
				ui.VmList.Refresh()
			}
		}
	}
}

func (ui *QManagerUI) showExpressWizard() {
	// Implementazione del wizard modale in Go
}

func (ui *QManagerUI) Run() {
	ui.Window.ShowAndRun()
}
