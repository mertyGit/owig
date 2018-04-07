package main

import (
  "github.com/lxn/walk"
  . "github.com/lxn/walk/declarative"
  "fmt"
  "os"
  "time"
  "strings"
)

//Version of program
const VERSION = "Version 0.92"

// ----------------------------------------------------------------------------
// Constants for screen resolutions
const SIZE_NONE     = 666
const SIZE_4K       = 0
const SIZE_1080     = 1

// ----------------------------------------------------------------------------
// Constants for screen types
const SC_UNKNOWN    = 0
const SC_TAB        = 1
const SC_VICTORY    = 2
const SC_DEFEAT     = 3
const SC_OVERVIEW   = 4
const SC_SRGAIN     = 5
const SC_ASSEMBLE   = 6
const SC_MAIN       = 7
const SC_ENDING     = 8
const SC_GAME       = 9
const SC_POTG       = 10
const SC_RESPAWN    = 11

// ----------------------------------------------------------------------------
// Constants for game state
const GS_NONE       = 1  // No game started or searching for game
const GS_START      = 2  // Game Started (waiting / assemble)
const GS_RUN        = 3  // Game is ongoing, stats might change later
const GS_END        = 4  // Game Ended, showing results, stats are final

// ----------------------------------------------------------------------------
// Global variable to hold every info about the game thats that is being played
type GameInfo struct {
  ts        int64    // Timestamp since start of program
  time      string   // Timeindicator during game
  screen       int   // Found screen type (tab screen, overview screen etc.etc.)
  pscreen      int   // Previous screen type 
  state        int   // state
  mapname   string   // Name of map we are playing, like "ILIOS" 
  gametype  string   // Game type, like "MYSTERY HEROES"  or "QUICK PLAY"
  side      string   // attack or defend
  hero      string   // What the hero is playing at the moment 
  currentSR    int
  highestSR    int
  medals [6]string   // Medals for common statistics
  lstats [6]string   // Common statistics (left bottom on TAB screen)
  rstats [6]string   // Special statistics (right bottom on TAB screen)
  group        int   // Group ID of player
  enemy   TeamComp   // Enemy team composition (see below)
  own     TeamComp   // Enemy team composition (see below)
  result    string   // End result (won,lost,draw)
  dmsg   [4]string   // debug messages
  forceM      bool   // Force found maptitle for rest of game
  forceT      bool   // Force found gametype for rest of game
  image       bool   // are we using images instead of screenshots ?
}

type TeamComp struct {
  hero      [6]string
  groupid   [6]int    // which group they belong to (groupid=0 => none)
  switches  [6]int    // How many switches has been detected in teamcomposition
  isChanged bool      // Flag if composition is changed during last write to it
}

var game GameInfo


func mainLoop() {

  // you can provide ini file as first argument and or list of png files 
  ic:=0
  if len(os.Args)>1 {
    if strings.HasSuffix(os.Args[1],"png") {
      ic=1
    } else if len(os.Args)>2 {
      if strings.HasSuffix(os.Args[2],"png") {
        ic=2
      }
    }
  }

  if ic>0 {
    game.image=true
    // testing, debug with screenshots
    for a:=ic;a<len(os.Args);a++ {
      ts("open")
      owig.Open(os.Args[a])
      ts("interpret")
      interpret()
      wrt.Send()
      ts("sleep")
      mainWindow.Invalidate()
      dbgWindow("Reading File: "+os.Args[a])
      if (a+1<len(os.Args)) {
        time.Sleep(time.Duration(config.dbg_pause) * time.Millisecond)
      }
      ts("afsleep")
    }
    if config.dbg_screen {
      fmt.Println(" Waiting for end")
    }
    for {
      time.Sleep(time.Duration(config.dbg_pause) * time.Millisecond)
    }
  } else {
    game.image=false
    for {
      ts("capture")
      owig.Capture()
      ts("interpret")
      interpret()
      wrt.Send()
      ts("draw")
      mainWindow.Invalidate()
      ts("sleep")
      time.Sleep(time.Duration(config.sleep) * time.Millisecond)
    }
  }
}

func main() {
  owig=new(OWImg)
  wrt=new(Oww)
  initOwig()
  go mainLoop()
  go wrt.Run()

  mw:= new(MyMainWindow)

  icon,_ := walk.NewIconFromFile("owig256.ico")

  MainWindow{
    AssignTo: &mainWindow,
    Title:   "OWIG",
    Icon: icon,
    MinSize:    Size{600, 500},
    Layout:  VBox{MarginsZero:true},
    Children: []Widget{
      CustomWidget{
        AssignTo:            &mw.paintWidget,
        ClearsBackground:    true,
        InvalidatesOnResize: true,
        Paint:               mw.drawWindow,
      },
    },
  }.Run()
}
