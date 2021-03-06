package gui

import (
	"github.com/jamesrr39/taskrunner-app/taskrunner"
	"github.com/jamesrr39/taskrunner-app/taskrunnerdal"
	"github.com/jamesrr39/taskrunner-app/triggers"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

type TaskrunnerGUI struct {
	mainFrame   gtk.IBox
	PaneContent Scene
	paneWidget  gtk.IWidget
	Window      *gtk.Window
	*taskrunnerdal.TaskrunnerDAL
	JobStatusChangeChan chan *taskrunner.JobRun // job runs
	titleLabel          *gtk.Label
	options             TaskrunnerGUIOptions
	udevRulesDAL        *triggers.UdevRulesDAL
}

func NewTaskrunnerGUI(taskrunnerDAL *taskrunnerdal.TaskrunnerDAL, udevRulesDAL *triggers.UdevRulesDAL, options TaskrunnerGUIOptions) *TaskrunnerGUI {

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetSizeRequest(800, 600)
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		gtk.MainQuit()
	})
	window.SetTitle("(Alpha) :: Taskrunner (" + taskrunnerDAL.String() + ")")
	window.ModifyBG(gtk.STATE_NORMAL, gdk.NewColor("white"))

	mainFrame := gtk.NewVBox(false, 10)

	titleLabel := gtk.NewLabel("")
	titleLabel.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("white"))

	taskrunnerGUI := &TaskrunnerGUI{
		mainFrame:           gtk.IBox(mainFrame),
		Window:              window,
		TaskrunnerDAL:       taskrunnerDAL,
		JobStatusChangeChan: make(chan *taskrunner.JobRun),
		titleLabel:          titleLabel,
		options:             options,
		udevRulesDAL:        udevRulesDAL,
	}

	go func() {
		for {
			jobRun := <-taskrunnerGUI.JobStatusChangeChan
			taskrunnerGUI.PaneContent.OnJobRunStatusChange(jobRun)
		}
	}()

	mainFrame.PackStart(buildToolbar(taskrunnerGUI), false, false, 0)
	window.Add(mainFrame)

	return taskrunnerGUI
}

func (taskrunnerGUI *TaskrunnerGUI) RenderScene(scene Scene) {
	if nil != taskrunnerGUI.paneWidget {
		taskrunnerGUI.paneWidget.Destroy()
	}
	taskrunnerGUI.PaneContent = scene

	taskrunnerGUI.titleLabel.SetText(scene.Title())

	taskrunnerGUI.paneWidget = scene.Content()

	taskrunnerGUI.mainFrame.Add(taskrunnerGUI.paneWidget)

	taskrunnerGUI.Window.ShowAll()
}
