/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/DarlingGoose/bare/apps/dashboard"
	"github.com/spf13/cobra"
)

var guiFile string

// guiCmd represents the gui command
var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Launch the desktop tail viewer",
	Long:  "Launch the desktop tail viewer for live-following a file with selectable text.",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("starting gio app")

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
					dash.Layout(gtx)

					e.Frame(gtx.Ops)
				}
			}
		}()

		app.Main()
	},
}

func init() {
	rootCmd.AddCommand(guiCmd)
	guiCmd.Flags().StringVar(&guiFile, "file", "", "path to a file to open immediately")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// guiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// guiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
