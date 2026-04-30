/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log/slog"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/DarlingGoose/bare/apps/dashboard"
	"github.com/DarlingGoose/bare/cmd"
)

func main() {
	if cmd.IsGUICommand(os.Args[1:]) {
		go func() {
			w := new(app.Window)
			w.Option(
				app.Title("Gio Dashboard"),
				app.Size(unit.Dp(1200), unit.Dp(800)),
			)
			w.Invalidate()

			dash := dashboard.NewApp()
			var ops op.Ops

			for {
				e := w.Event()

				switch e := e.(type) {
				case app.DestroyEvent:
					slog.Info("gio destroyed", "err", e.Err)
					return

				case app.FrameEvent:
					ops.Reset()

					gtx := app.NewContext(&ops, e)

					paint.Fill(gtx.Ops, dash.Theme.Color.Background)

					dash.Layout(gtx)

					e.Frame(gtx.Ops)
				}
			}
		}()

		app.Main()
		return
	}

	cmd.Execute()

}
