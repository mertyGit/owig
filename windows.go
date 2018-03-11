package main

import (
  "github.com/lxn/walk"
)

// ----------------------------------------------------------------------------
// Draws window with information
func (mw *MyMainWindow) drawWindow(canvas *walk.Canvas, updateBounds walk.Rectangle) error {
  bounds := mw.paintWidget.ClientBounds()
  rbounds := bounds
  mainCanvas=canvas
  mainRect=updateBounds

  // Color background black
  blackBrush,err:=walk.NewSolidColorBrush(walk.RGB(0,0,0))
  if err != nil {
    return err
  }
  defer blackBrush.Dispose()
  canvas.FillRectangle(blackBrush,bounds)

  // Display Game map and Gametype
  tfont, err := walk.NewFont("Abadi MT Condensed Light", 12, 0)
  if err != nil {
    return err
  }
  defer tfont.Dispose()
  tifont, err := walk.NewFont("Abadi MT Condensed Light", 10, walk.FontItalic)
  if err != nil {
    return err
  }
  defer tifont.Dispose()
  rbounds.X=10
  rbounds.X=10
  rbounds.Y=10
  rbounds.Width=bounds.Width-10
  rbounds.Height=20
  if game.state==GS_RUN||game.state==GS_END {
    canvas.DrawText(game.mapname,tfont,walk.RGB(212,212,212),rbounds, 0)
    rbounds.Width=200
    fbounds,_,_:=canvas.MeasureText(game.mapname,tfont,rbounds, 0)
    rbounds.X=10+fbounds.Width+10
    canvas.DrawText(game.gametype,tfont,walk.RGB(212,212,0),rbounds, 0)
    fbounds,_,_=canvas.MeasureText(game.gametype,tfont,rbounds, 0)
    rbounds.X=fbounds.X+fbounds.Width+10
    rbounds.Y+=2
    if game.side=="attack" {
      canvas.DrawText("Attack",tifont,walk.RGB(212,100,100),rbounds, 0)
    } else {
      if game.side=="attack" {
        canvas.DrawText("Defense",tifont,walk.RGB(100,100,212),rbounds, 0)
      }
    }
  }

  // Display debug info
  if config.dbg_window {
    dfont, err := walk.NewFont("Consolas", 8, 0)
    if err != nil {
      return err
    }
    defer dfont.Dispose()
    rbounds.X=10
    rbounds.Y=bounds.Height-43
    rbounds.Width=bounds.Width-10
    rbounds.Height=13
    canvas.DrawText(game.dmsg[0],dfont,walk.RGB(90,90,60),rbounds, walk.TextWordbreak)
    rbounds.Y=bounds.Height-33
    canvas.DrawText(game.dmsg[1],dfont,walk.RGB(90,90,60),rbounds, walk.TextWordbreak)
    rbounds.Y=bounds.Height-23
    canvas.DrawText(game.dmsg[2],dfont,walk.RGB(90,90,60),rbounds, walk.TextWordbreak)
    rbounds.Y=bounds.Height-13
    canvas.DrawText(game.dmsg[3],dfont,walk.RGB(90,90,60),rbounds, walk.TextWordbreak)
  }



  return nil
}

// ----------------------------------------------------------------------------
// Used to display 4 lines of debug text in Window

func dbgWindow(msg string) {
  if config.dbg_window {
    game.dmsg[0]=game.dmsg[1]
    game.dmsg[1]=game.dmsg[2]
    game.dmsg[2]=game.dmsg[3]
    game.dmsg[3]=msg
    mainWindow.Invalidate()
  }
}
