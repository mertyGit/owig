package main

import (
  "github.com/lxn/walk"
  "strings"
  "fmt"
)

type Icons struct {
  name string
  file string
  img walk.Image
}

//Windows window struct
type MyMainWindow struct {
  *walk.MainWindow
  paintWidget *walk.CustomWidget
}

var mainCanvas *walk.Canvas
var mainRect walk.Rectangle
var mainWindow *walk.MainWindow



var heroIcons []Icons

var redflair walk.Image
var blueflair walk.Image
var owigsign walk.Image





// ----------------------------------------------------------------------------
// Load all icon information

func loadIcons() {
  heroIcons = []Icons{
    {name:"Ana"       ,file:"Ana_icon.png"       ,img:nil},
    {name:"Bastion"   ,file:"Bastion_icon.png"   ,img:nil},
    {name:"Brigitte"  ,file:"Brigitte_icon.png"  ,img:nil},
    {name:"Doomfist"  ,file:"Doomfist_icon.png"  ,img:nil},
    {name:"D.Va"      ,file:"DVa_icon.png"       ,img:nil},
    {name:"Genji"     ,file:"Genji_icon.png"     ,img:nil},
    {name:"Hanzo"     ,file:"Hanzo_icon.png"     ,img:nil},
    {name:"Junkrat"   ,file:"Junkrat_icon.png"   ,img:nil},
    {name:"Lucio"     ,file:"Lucio_icon.png"     ,img:nil},
    {name:"McCree"    ,file:"McCree_icon.png"    ,img:nil},
    {name:"Mei"       ,file:"Mei_icon.png"       ,img:nil},
    {name:"Mercy"     ,file:"Mercy_icon.png"     ,img:nil},
    {name:"Moira"     ,file:"Moira_icon.png"     ,img:nil},
    {name:"Orisa"     ,file:"Orisa_icon.png"     ,img:nil},
    {name:"Pharah"    ,file:"Pharah_icon.png"    ,img:nil},
    {name:"Reaper"    ,file:"Reaper_icon.png"    ,img:nil},
    {name:"Reinhardt" ,file:"Reinhardt_icon.png" ,img:nil},
    {name:"Roadhog"   ,file:"Roadhog_icon.png"   ,img:nil},
    {name:"Soldier 76",file:"Soldier76_icon.png" ,img:nil},
    {name:"Sombra"    ,file:"Sombra_icon.png"    ,img:nil},
    {name:"Symmetra"  ,file:"Symmetra_icon.png"  ,img:nil},
    {name:"Torbjorn"  ,file:"Torbjorn_icon.png"  ,img:nil},
    {name:"Tracer"    ,file:"Tracer_icon.png"    ,img:nil},
    {name:"Widowmaker",file:"Widowmaker_icon.png",img:nil},
    {name:"Winston"   ,file:"Winston_icon.png"   ,img:nil},
    {name:"Zarya"     ,file:"zarya_icon.png"     ,img:nil},
    {name:"Zenyatta"  ,file:"Zenyatta_icon.png"  ,img:nil},
  }
  for i,v := range heroIcons {
    //fmt.Println("Got i,v",i,v.file)
    v.img,_=walk.NewImageFromFile("images/"+v.file)
    heroIcons[i]=v
  }
  redflair,_=walk.NewImageFromFile("images/redflair.png")
  blueflair,_=walk.NewImageFromFile("images/blueflair.png")
  owigsign,_=walk.NewImageFromFile("images/owigsign.png")
}

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

  hfont, err := walk.NewFont("Abadi MT Condensed Light", 12, walk.FontBold)
  if err != nil {
    return err
  }
  defer hfont.Dispose()

  tifont, err := walk.NewFont("Abadi MT Condensed Light", 10, walk.FontItalic)
  if err != nil {
    return err
  }
  defer tifont.Dispose()

  // Initialize Dark Grey and Light Grey pens
  drawLG, err:= walk.NewSolidColorBrush(walk.RGB(200,200,200))
  if err != nil {
    return err
  }
  defer drawLG.Dispose()

  LGPen, err := walk.NewGeometricPen(walk.PenSolid, 1 , drawLG)
  if err != nil {
    return err
  }
  defer LGPen.Dispose()

  drawDG, err:= walk.NewSolidColorBrush(walk.RGB(50,50,50))
  if err != nil {
    return err
  }
  defer drawDG.Dispose()

  DGPen, err := walk.NewGeometricPen(walk.PenSolid, 1 , drawDG)
  if err != nil {
    return err
  }
  defer DGPen.Dispose()

  xoff:=(bounds.Width-600)/2
  if xoff < 0 {
    xoff=0
  }


  // Color background
  for x:=0;x<bounds.Width;x+=1000 {
    canvas.DrawImage(redflair, walk.Point{x,30})
    canvas.DrawImage(blueflair, walk.Point{x,120})
  }

  // Headlines on top
  rbounds.X=10
  rbounds.Y=10
  rbounds.Width=bounds.Width-10
  rbounds.Height=20
  canvas.DrawText(game.mapname,tfont,walk.RGB(212,212,212),rbounds, 0)
  rbounds.Width=200

  // BugFix: Walk.TextCalcRec added, otherwise width will not be returned
  fbounds,_,_:=canvas.MeasureText(game.mapname,tfont,rbounds, walk.TextCalcRect)
  rbounds.X=10+fbounds.Width+10
  canvas.DrawText(game.gametype,tfont,walk.RGB(212,212,0),rbounds, 0)
  fbounds,_,_=canvas.MeasureText(game.gametype,tfont,rbounds, walk.TextCalcRect)

  rbounds.X=bounds.Width-80
  canvas.DrawText(game.time,tfont,walk.RGB(212,212,212),rbounds, 0)

  rbounds.X=fbounds.X+fbounds.Width+10
  rbounds.Y+=2
  if game.side=="attack" {
    canvas.DrawText("Attack",tifont,walk.RGB(212,100,100),rbounds, 0)
  } else {
    if game.side=="defend" {
      canvas.DrawText("Defense",tifont,walk.RGB(100,100,212),rbounds, 0)
    }
  }



  // Display icons, borders and grouping of known players
  gfont, err := walk.NewFont("Arial Black", 25, walk.FontBold)
  if err != nil {
    return err
  }
  defer gfont.Dispose()
  for i:=0;i<6;i++ {
    for s,v := range heroIcons {
      if game.enemy.hero[i]==v.name {
        canvas.DrawImage(heroIcons[s].img, walk.Point{xoff+10+100*i,50})
        if i<5 && game.enemy.groupid[i] >0 && game.enemy.groupid[i]==game.enemy.groupid[i+1] {
          rbounds.X=xoff+83+100*i
          rbounds.Y=50
          rbounds.Width=40
          rbounds.Height=40
          canvas.DrawText("-",gfont,walk.RGB(212,212,212),rbounds, 0)
        }
      }
      if game.own.hero[i]==v.name {
        canvas.DrawImage(heroIcons[s].img, walk.Point{xoff+10+100*i,150})
        if i<5 && game.own.groupid[i] >0 && game.own.groupid[i]==game.own.groupid[i+1] {
          rbounds.X=xoff+83+100*i
          rbounds.Y=150
          rbounds.Width=40
          rbounds.Height=40
          if game.group==game.own.groupid[i] {
            canvas.DrawText("-",gfont,walk.RGB(50,212,50),rbounds, 0)
          } else {
            canvas.DrawText("-",gfont,walk.RGB(212,212,212),rbounds, 0)
          }
        }
      }
    }
    // Draw border
    canvas.DrawLine(DGPen,walk.Point{xoff+10+100*i,50},walk.Point{xoff+66+100*i,50})
    canvas.DrawLine(DGPen,walk.Point{xoff+10+100*i,49},walk.Point{xoff+66+100*i,49})
    canvas.DrawLine(DGPen,walk.Point{xoff+10+100*i,50},walk.Point{xoff+10+100*i,100})
    canvas.DrawLine(DGPen,walk.Point{xoff+9+100*i,50},walk.Point{xoff+9+100*i,100})
    canvas.DrawLine(LGPen,walk.Point{xoff+10+100*i,100},walk.Point{xoff+66+100*i,100})
    canvas.DrawLine(LGPen,walk.Point{xoff+66+100*i,50},walk.Point{xoff+66+100*i,100})

    canvas.DrawLine(DGPen,walk.Point{xoff+10+100*i,150},walk.Point{xoff+66+100*i,150})
    canvas.DrawLine(DGPen,walk.Point{xoff+10+100*i,149},walk.Point{xoff+66+100*i,149})
    canvas.DrawLine(DGPen,walk.Point{xoff+10+100*i,150},walk.Point{xoff+10+100*i,200})
    canvas.DrawLine(DGPen,walk.Point{xoff+9+100*i,150},walk.Point{xoff+9+100*i,200})
    canvas.DrawLine(LGPen,walk.Point{xoff+10+100*i,200},walk.Point{xoff+66+100*i,200})
    canvas.DrawLine(LGPen,walk.Point{xoff+66+100*i,150},walk.Point{xoff+66+100*i,200})


    }

  // Dividers
  canvas.DrawLine(LGPen,walk.Point{0,30},walk.Point{bounds.Width,30})
  canvas.DrawLine(DGPen,walk.Point{0,31},walk.Point{bounds.Width,31})
  canvas.DrawLine(DGPen,walk.Point{0,32},walk.Point{bounds.Width,32})

  canvas.DrawLine(LGPen,walk.Point{0,220},walk.Point{bounds.Width,220})
  canvas.DrawLine(DGPen,walk.Point{0,221},walk.Point{bounds.Width,221})
  canvas.DrawLine(DGPen,walk.Point{0,222},walk.Point{bounds.Width,222})

  canvas.DrawLine(LGPen,walk.Point{0,bounds.Height-52},walk.Point{bounds.Width,bounds.Height-52})
  canvas.DrawLine(DGPen,walk.Point{0,bounds.Height-51},walk.Point{bounds.Width,bounds.Height-51})
  canvas.DrawLine(DGPen,walk.Point{0,bounds.Height-50},walk.Point{bounds.Width,bounds.Height-50})

  // Display common statistics 
  sfont, err := walk.NewFont("Abadi MT Condensed Light", 10, 0)
  if err != nil {
    return err
  }
  defer sfont.Dispose()
  rbounds.X=30+xoff
  rbounds.Y=235
  rbounds.Width=bounds.Width-10
  rbounds.Height=20
  canvas.DrawText("Eliminations",sfont,walk.RGB(180,180,180),rbounds, 0)
  rbounds.Y+=20
  canvas.DrawText("Objective Kills",sfont,walk.RGB(180,180,180),rbounds, 0)
  rbounds.Y+=20
  canvas.DrawText("Objective Time",sfont,walk.RGB(180,180,180),rbounds, 0)
  rbounds.Y+=20
  canvas.DrawText("Hero Damage Done",sfont,walk.RGB(180,180,180),rbounds, 0)
  rbounds.Y+=20
  canvas.DrawText("Healing Done",sfont,walk.RGB(180,180,180),rbounds, 0)
  rbounds.Y+=20
  canvas.DrawText("Deaths",sfont,walk.RGB(180,180,180),rbounds, 0)
  for y:=0;y<6;y++ {
    rbounds.X=150+xoff
    rbounds.Y=235+y*20
    rbounds.Width=20
    rbounds.Height=20
    canvas.DrawText(":",sfont,walk.RGB(180,180,180),rbounds, 0)
    rbounds.X=160+xoff
    rbounds.Width=50
    canvas.DrawText(game.lstats[y],sfont,walk.RGB(250,250,250),rbounds, walk.TextRight)
  }
  GoldBrush,_:=walk.NewSolidColorBrush(walk.RGB(176,134,46))
  SilverBrush,_:=walk.NewSolidColorBrush(walk.RGB(124,129,130))
  BronzeBrush,_:=walk.NewSolidColorBrush(walk.RGB(99,50,46))
  defer GoldBrush.Dispose()
  defer SilverBrush.Dispose()
  defer BronzeBrush.Dispose()

  // Display Medals
  for y:=0;y<6;y++ {
    rbounds.X=10+xoff
    rbounds.Y=235+y*20
    rbounds.Width=15
    rbounds.Height=15
    switch game.medals[y] {
      case "G":
        canvas.FillEllipse(GoldBrush,rbounds)
      case "S":
        canvas.FillEllipse(SilverBrush,rbounds)
      case "B":
        canvas.FillEllipse(BronzeBrush,rbounds)
    }
  }

  // Display special Statistics
  hero:=game.own.hero[0]
  for y:=0;y<6;y++ {
    rbounds.X=250+xoff
    rbounds.Y=235+y*20
    rbounds.Height=20
    rbounds.Width=200
    canvas.DrawText(getStatsline(hero,y),sfont,walk.RGB(180,180,180),rbounds, 0)
    rbounds.X=400+xoff
    rbounds.Width=20
    canvas.DrawText(":",sfont,walk.RGB(180,180,180),rbounds, 0)

    rbounds.X=410+xoff
    rbounds.Width=50
    if strings.Contains(game.rstats[y],"%") {
      rbounds.X=422+xoff
    }
    canvas.DrawText(game.rstats[y],sfont,walk.RGB(250,250,250),rbounds, walk.TextRight)
  }


  // Choosen Hero
  rbounds.X=10
  rbounds.Y=380
  rbounds.Width=bounds.Width-10
  rbounds.Height=20
  canvas.DrawText(game.hero,hfont,walk.RGB(212,212,212),rbounds, 0)
  rbounds.Width=200


  // Display SR
  rbounds.X=bounds.Width-100
  rbounds.Y=380
  rbounds.Width=bounds.Width-10
  rbounds.Height=20
  sr:=fmt.Sprintf("%d",game.currentSR)
  canvas.DrawText(sr,hfont,walk.RGB(212,212,212),rbounds, 0)
  rbounds.Width=200


  // Bottom part of screen
  rbounds.X=10
  rbounds.Y=bounds.Height-43
  rbounds.Width=bounds.Width-10
  rbounds.Height=13
  dfont, err := walk.NewFont("Consolas", 8, 0)
  if err != nil {
    return err
  }
  defer dfont.Dispose()

  if config.dbg_window {
    // Display debug info
    canvas.DrawText(game.dmsg[0],dfont,walk.RGB(90,90,60),rbounds, walk.TextWordbreak)
    rbounds.Y=bounds.Height-33
    canvas.DrawText(game.dmsg[1],dfont,walk.RGB(90,90,60),rbounds, walk.TextWordbreak)
    rbounds.Y=bounds.Height-23
    canvas.DrawText(game.dmsg[2],dfont,walk.RGB(90,90,60),rbounds, walk.TextWordbreak)
    rbounds.Y=bounds.Height-13
    canvas.DrawText(game.dmsg[3],dfont,walk.RGB(90,90,60),rbounds, walk.TextWordbreak)
  } else {
    // Or display logo with version information
    canvas.DrawImage(owigsign, walk.Point{bounds.Width/2-200,bounds.Height-43})
    rbounds.Y=bounds.Height-13
    rbounds.X=bounds.Width/2-50
    canvas.DrawText(VERSION,dfont,walk.RGB(90,90,60),rbounds, walk.TextWordbreak)
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
  }
  mainWindow.Invalidate()
}
