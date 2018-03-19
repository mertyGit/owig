package main

import (
  "github.com/lxn/walk"
  . "github.com/lxn/walk/declarative"
  "fmt"
  "os"
  "time"
  "strings"
  "strconv"
)

//Version of program
const VERSION = "Version 0.81"

//Windows window struct
type MyMainWindow struct {
  *walk.MainWindow
  paintWidget *walk.CustomWidget
}

var mainCanvas *walk.Canvas
var mainRect walk.Rectangle
var mainWindow *walk.MainWindow

var heroStats map[string][]string


type Sign struct {
  name string
  low Pixel
  high Pixel
}

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

// ----------------------------------------------------------------------------
// Constants for game state
const GS_NONE       = 1  // No game started or searching for game
const GS_START      = 2  // Game Started (waiting / assemble)
const GS_RUN        = 3  // Game is ongoing, stats might change later
const GS_END        = 4  // Game Ended, showing results, stats are final

// ----------------------------------------------------------------------------
// Global variable to hold every info about the game thats that is being played
type GameInfo struct {
  screen       int   // Found screen type (tab screen, overview screen etc.etc.)
  pscreen      int   // Previous screen type 
  mapname   string   // Name of map we are playing, like "ILIOS" 
  gametype  string   // Game type, like "MYSTERY HEROES"  or "QUICK PLAY"
  hero      string   // What the hero is playing at the moment = own.heroes[0]
  enemy   TeamComp   // Enemy team composition (see below)
  own     TeamComp   // Enemy team composition (see below)
  stats   OwnStats   // Own statistics
  currentSR    int
  highestSR    int
  result    string   // End result (won,lost,draw)
  side      string   // attack or defend
  dmsg   [4]string   // debug messages
  state        int   // state
  image       bool   // are we using images instead of screenshots ?
  medals [6]string   // Medals for common statistics
  lstats [6]string   // Common statistics (left bottom on TAB screen)
  rstats [6]string   // Special statistics (right bottom on TAB screen)
  snames [6]string   // Name of statistics field
  time      string   // Timeindicator during game
  ts        int64   // Timestamp since start
}

type OwnStats struct {
  eleminations   int
  objectiveKills int
  objectiveTime  string
  objectiveSecs  int     //objectiveTime converted to seconds
  damage         int
  deaths         int
  medals         [6]string  // "G","S" or "B" 0=eleminations, 1=objective .etc.
  stats          [6]string
  statsText      [6]string  //meaning of stat 1..6, based on hero choice
}

type TeamComp struct {
  hero      [6]string
  groupid   [6]int    // which group they belong to (groupid=0 => none)
  switches  [6]int    // How many switches has been detected in teamcomposition
  isChanged bool      // Flag if composition is changed during last write to it
}

var game GameInfo

var owig *OWImg

func ts(id string) {
  if config.dbg_time {
    fmt.Println("TIME:",time.Now().UnixNano()/1000000 - game.ts,id)
  }
}

// ----------------------------------------------------------------------------
// Try to figure out what screenshot is
//

func guessScreen() int {
  var s1 Pixel
  var s2 Pixel
  var crc1 string
  var crc2 string
  var same bool
  var pos int

  ret:=0

  if config.dbg_screen {
    fmt.Println("== guessScreen ==")
  }
  ts("gstart")

  // Tab statistics ?
  // grey bar + black left upper corner, right under title / time
  l1:=false
  l2:=false
  owig.All()
  switch owig.res {
    case SIZE_4K:
      l1=owig.At(90,145).isLike(Pixel{60,50,60},10)
      l2=owig.At(115,145).isLike(Pixel{5,5,5},5)
    case SIZE_1080:
      l1=owig.At(42,55).isLike(Pixel{60,50,60},10)
      l2=owig.At(70,75).isLike(Pixel{5,5,5},5)
  }
  if l1 && l2 {
    if config.dbg_screen {
      fmt.Println(" TAB statistics screen")
    }
    return SC_TAB
  }


  // Victory or Defeat Screen ? 
  // Vertical image line in middle with same value
  switch owig.res {
    case SIZE_4K:
      same=owig.Box(1820,1050,0,50).isLine()
      s1=owig.All().At(1820,1050).RGB()
    case SIZE_1080:
      same=owig.Box(910,500,0,20).isLine()
      s1=owig.All().At(910,500).RGB()
  }

  if same {
    // which one, Victory or Defeat, yellow or red ?
    if s1.R>220 && s1.G<5 {
      if config.dbg_screen {
        fmt.Println(" Defeat")
      }
      return SC_DEFEAT
    } else if s1.R>220 && s1.G>190 && int(s1.B/10)==6 {
      return SC_VICTORY
      if config.dbg_screen {
        fmt.Println(" Victory")
      }
    }
  }
  // Victory/Defeat, but voting & medals part
  // search for blue "Leave Game" button
  switch owig.res {
    case SIZE_4K:
      crc1=owig.From(3448,100).To(3720,172).Cs()
    case SIZE_1080:
      crc1=owig.From(1724,50).To(1860,86).Cs()
  }
  if crc1=="13232D" || crc1=="12232E" || crc1=="12232D" {
    if config.dbg_screen {
      fmt.Println(" Ending game")
    }
    return SC_ENDING
  }






  // Career Screen ? 
  // White Pixel border for player icon on left top
  switch owig.res {
    case SIZE_4K:
      same=owig.Box(100,400,0,200).isLine()
      s1=owig.All().At(100,400).RGB()
      pos=550
    case SIZE_1080:
      same=owig.Box(50,200,0,100).isLine()
      s1=owig.All().At(50,200).RGB()
      pos=280
  }
  owig.All()

  // + blue stripe behind icon, middle
  if same && s1.R>220 && s1.G>220 && s1.B>220 {
    s1=owig.At(500,pos).RGB()
    for x:=500;x<600 && same==true;x++ {
      s2=owig.All().At(x,pos).RGB()
      if same==true && int(s1.R/10) == int(s2.R/10) && int(s1.G/10) == int(s2.G/10) && int(s1.B/10) == int(s2.B/10) {
        if int(s2.R/10) == 3 && int(s2.G/10) == 4 && int(s2.B/10) == 7 {
          same=true
        } else {
          same=false
        }
      } else {
        same=false
      }
      s1=s2
    }
    // + same icon right top and left
    if same {
      switch owig.res {
        case SIZE_4K:
          crc1=owig.From(3040,60).To(3160,180).Cs()
          crc2=owig.From(114,336).To(399,621).Cs()
          same=owig.From(3040,60).To(3160,180).SameBase("B")
        case SIZE_1080:
          crc1=owig.From(1520,30).To(1579,89).Cs()
          crc2=owig.From(57,168).To(198,309).Cs()
          same=owig.From(1520,30).To(1579,89).SameBase("B")
      }
      if crc1==crc2 || same {
        if config.dbg_screen {
          fmt.Println(" Player Overview screen")
        }
        return SC_OVERVIEW
      }
    }
  }

  // is it an comp gain/loss screen ?
  owig.All().Th(190)
  switch owig.res {
    case SIZE_4K:
      // middle white bar (CTF) "Season High" & "Career high" white bars 
      same=owig.At(1700,2060).isAbove() && owig.At(2100,2060).isAbove()
      owig.At(1928,1306).RGB() // purple dot, above comp points count
    case SIZE_1080:
      same=owig.At(851,1030).isAbove() && owig.At(1049,1030).isAbove()
      owig.At(964,653).RGB()
  }
  owig.Th(-1)

//  if same && int(s1.R/10) == 17 && int(s1.G/10) == 0 && int(s1.B/10) == 22 {
  if same && owig.isLike(Pixel{170,0,220},9) {
      if config.dbg_screen {
        fmt.Println(" SR Gain/Loss screen")
      }
      return SC_SRGAIN
  }

  // is it "Assemble your team" screen ? (beginning of match)
  // find stripe and pull down menu for skin
  same=false
  owig.All()
  switch owig.res {
    case SIZE_4K:
      if owig.Box(3586,376,8,0).isLine() && owig.isWhite() && owig.Box(3684,450,8,0).isLine() && owig.isWhite() && owig.Box(3620,240,10,60).SameBase("C") && (owig.Box(1730,2000,38,60).SameBase("R")||owig.Box(1730,2000,38,60).SameBase("B")) {
        same=true
      }
    case SIZE_1080:
      if owig.Box(1793,188,4,0).isLine() && owig.isWhite() && owig.Box(1842,225,3,0).isLine() && owig.isWhite() && owig.Box(1810,120,5,30).SameBase("C") && (owig.Box(865,1000,19,30).SameBase("R")||owig.Box(865,1000,19,30).SameBase("B")) {
        same=true
      }
  }
  if same {
    if config.dbg_screen {
      fmt.Println(" Assemble Screen")
    }
    return SC_ASSEMBLE
  }

  //is it a "Main screen" (screen you see after logging in)
  same=false
  switch owig.res {
    case SIZE_4K:
      if owig.From(80,154).To(1000,154).Th(224).Holes()==15 {
        same=true
      }
    case SIZE_1080:
      if owig.From(40,77).To(500,77).Th(112).Holes()==13 {
        same=true
      }
  }
  owig.Th(-1)
  if same {
    if config.dbg_screen {
      fmt.Println(" Main Screen")
    }
    return SC_MAIN
  }
  return ret
}

// ----------------------------------------------------------------------------
// Wrapper for guessScreen, to set the right game
func getScreen() {
  game.pscreen=game.screen
  game.screen=guessScreen()
  ts("gstop")
  if game.pscreen != game.screen {
    if config.dbg_screen {
      fmt.Println("Screen change from",game.pscreen," to",game.screen)
    }
  }
}

// ----------------------------------------------------------------------------
// get received medals, returns gold,silver,bronze or "-"

func getMedal(pos int) string {
  var mxpos []int
  var mypos []int
  var medal ="( )"

  if config.dbg_screen {
    fmt.Println("== getMedal ==")
  }
  switch owig.res {
    case SIZE_4K:
      mxpos = []int{191,691,1191}
      mypos = []int{1801,1927}
    case SIZE_1080:
      mxpos = []int{95,345,595}
      mypos = []int{900,963}
  }
  x:=mxpos[pos%3]
  y:=mypos[0]
  if (pos>2) {
    y=mypos[1]
  }
  owig.All().At(x,y)
  if config.dbg_screen {
    fmt.Println(" Got Pix:",owig.RGB())
  }
  if (owig.Red()>89) {
    if (owig.Green()>89) {
      if (owig.Blue()>89) {
        medal = "S"
      } else {
        medal = "G"
      }
    } else {
      medal = "B"
    }
  }
  return medal
}

// ----------------------------------------------------------------------------
// guess teamcomposition by taking two owig.Atels from each player
// row (0=enemy team, 1=own team) and position 
// (col=0=always the player)
// returns "unknown" if not recognized or dead (red cross)

func guessHero(col int, row int) string {
  var heros   []Sign
  var xpos    []int
  var ypown   []int
  var ypenemy []int

  if config.dbg_screen {
    fmt.Println("== guessHero ==")
  }
  switch owig.res {
    case SIZE_4K:
      heros = []Sign{
        {name:"Ana"       ,low:Pixel{139,142,149},high:Pixel{101,95,94}},
        {name:"Bastion"   ,low:Pixel{176,164,132},high:Pixel{171,206,188}},
        {name:"Brigitte"  ,low:Pixel{113,35,22},high:Pixel{208,133,114}},
        {name:"Doomfist"  ,low:Pixel{84,61,46},high:Pixel{166,137,112}},
        {name:"D.Va"      ,low:Pixel{124,81,100},high:Pixel{197,155,125}},
        {name:"Genji"     ,low:Pixel{80,98,115},high:Pixel{154,148,134}},
        {name:"Hanzo"     ,low:Pixel{107,68,44},high:Pixel{59,54,58}},
        {name:"Junkrat"   ,low:Pixel{69,72,84},high:Pixel{199,190,160}},
        {name:"Lucio"     ,low:Pixel{79,48,36},high:Pixel{65,88,12}},
        {name:"McCree"    ,low:Pixel{99,23,26},high:Pixel{23,18,22}},
        {name:"Mei"       ,low:Pixel{171,168,164},high:Pixel{23,23,23}},
        {name:"Mercy"     ,low:Pixel{50,54,62},high:Pixel{223,191,156}},
        {name:"Moira"     ,low:Pixel{103,107,121},high:Pixel{188,125,104}},
        {name:"Orisa"     ,low:Pixel{43,38,35},high:Pixel{188,181,22}},
        {name:"Pharah"    ,low:Pixel{59,59,67},high:Pixel{130,94,77}},
        {name:"Reaper"    ,low:Pixel{26,25,26},high:Pixel{42,35,23}},
        {name:"Reinhardt" ,low:Pixel{46,42,39},high:Pixel{83,77,74}},
        {name:"Roadhog"   ,low:Pixel{171,138,122},high:Pixel{57,60,66}},
        {name:"Soldier 76",low:Pixel{172,188,194},high:Pixel{99,0,0}},
        {name:"Sombra"    ,low:Pixel{122,64,43},high:Pixel{195,132,87}},
        {name:"Symmetra"  ,low:Pixel{106,99,89},high:Pixel{67,73,67}},
        {name:"Torbjorn"  ,low:Pixel{203,185,135},high:Pixel{42,40,42}},
        {name:"Tracer"    ,low:Pixel{165,99,77},high:Pixel{132,32,0}},
        {name:"Widowmaker",low:Pixel{99,119,179},high:Pixel{23,24,23}},
        {name:"Winston"   ,low:Pixel{31,31,39},high:Pixel{95,100,109}},
        {name:"Zarya"     ,low:Pixel{55,97,107},high:Pixel{172,104,78}},
        {name:"Zenyatta"  ,low:Pixel{30,32,30},high:Pixel{34,28,21}},
      }

      xpos   = []int{959,1343,1727,2111,2495,2879}
      ypown  = []int{1280,1200}
      ypenemy= []int{670,590}

    case SIZE_1080:
      heros = []Sign{
        {name:"Ana"       ,low:Pixel{103,103,109},high:Pixel{153,143,138}},
        {name:"Bastion"   ,low:Pixel{132,126,100},high:Pixel{171,206,186}},
        {name:"Brigitte"  ,low:Pixel{166,143,138},high:Pixel{203,128,111}},
        {name:"Doomfist"  ,low:Pixel{57,44,33},high:Pixel{143,114,91}},
        {name:"D.Va"      ,low:Pixel{178,169,156},high:Pixel{138,101,79}},
        {name:"Genji"     ,low:Pixel{44,35,36},high:Pixel{144,141,128}},
        {name:"Hanzo"     ,low:Pixel{121,79,50},high:Pixel{22,17,22}},
        {name:"Junkrat"   ,low:Pixel{133,108,92},high:Pixel{180,172,153}},
        {name:"Lucio"     ,low:Pixel{55,46,36},high:Pixel{71,94,11}},
        {name:"McCree"    ,low:Pixel{90,21,22},high:Pixel{46,34,35}},
        {name:"Mei"       ,low:Pixel{157,158,157},high:Pixel{23,24,22}},
        {name:"Mercy"     ,low:Pixel{66,76,88},high:Pixel{136,109,87}},
        {name:"Moira"     ,low:Pixel{83,74,64},high:Pixel{182,117,99}},
        {name:"Orisa"     ,low:Pixel{20,19,20},high:Pixel{195,186,29}},
        {name:"Pharah"    ,low:Pixel{59,69,82},high:Pixel{120,83,68}},
        {name:"Reaper"    ,low:Pixel{13,13,13},high:Pixel{46,39,23}},
        {name:"Reinhardt" ,low:Pixel{74,71,64},high:Pixel{92,83,77}},
        {name:"Roadhog"   ,low:Pixel{84,43,36},high:Pixel{50,48,52}},
        {name:"Soldier 76",low:Pixel{78,87,92},high:Pixel{84,0,0}},
        {name:"Sombra"    ,low:Pixel{172,101,66},high:Pixel{195,132,86}},
        {name:"Symmetra"  ,low:Pixel{89,46,36},high:Pixel{37,51,75}},
        {name:"Torbjorn"  ,low:Pixel{196,186,151},high:Pixel{44,43,44}},
        {name:"Tracer"    ,low:Pixel{120,55,44},high:Pixel{29,6,0}},
        {name:"Widowmaker",low:Pixel{68,84,147},high:Pixel{23,21,22}},
        {name:"Winston"   ,low:Pixel{44,43,44},high:Pixel{105,111,105}},
        {name:"Zarya"     ,low:Pixel{52,47,44},high:Pixel{190,122,94}},
        {name:"Zenyatta"  ,low:Pixel{102,111,111},high:Pixel{8,7,6}},
      }

      xpos   = []int{480,672,864,1056,1248,1440}
      ypown  = []int{636,600}
      ypenemy= []int{331,295}

  }

  var dev=0
  var tot=0
  var score=1000
  var dif=0
  var found="unknown"
  var inl Pixel
  var inh Pixel

  if row>0 {
    inl=owig.All().At(xpos[col],ypown[0]).RGB()
    inh=owig.All().At(xpos[col],ypown[1]).RGB()
  } else {
    inl=owig.All().At(xpos[col],ypenemy[0]).RGB()
    inh=owig.All().At(xpos[col],ypenemy[1]).RGB()
  }
  if config.dbg_screen {
    fmt.Println(" Got pixels ",inl,inh)
  }

  for _, el := range heros {
      tot=0
      dev=0
      dif=getDif(inl.R,el.low.R);
      if (dif>dev) { dev=dif }
      tot+=dif
      dif=getDif(inl.G,el.low.G);
      if (dif>dev) { dev=dif }
      tot+=dif
      dif=getDif(inl.B,el.low.B);
      if (dif>dev) { dev=dif }
      tot+=dif
      dif=getDif(inh.R,el.high.R);
      if (dif>dev) { dev=dif }
      tot+=dif
      dif=getDif(inh.G,el.high.G);
      if (dif>dev) { dev=dif }
      tot+=dif
      dif=getDif(inh.B,el.high.B);
      if (dif>dev) { dev=dif }
      tot+=dif
      if ((dev<2) && (score>tot)) {
        score=tot
        found=el.name
      }
  }
  if (score>2) {
    found="unknown"
  }
  if config.dbg_screen {
    fmt.Println(" Returning ",found," for",col,row)
  }
  return found
}
// ----------------------------------------------------------------------------
// Determine group per player 
func getGroups() {
  var xpos    []int
  var ypos    []int
  var gcnt=0
  var fg=false

  if config.dbg_screen {
    fmt.Println("== getGroups ==")
  }
  switch owig.res {
    case SIZE_4K:
      xpos  = []int{1140,1524,1908,2292,2676}
      ypos  = []int{620,1240}
    case SIZE_1080:
      xpos  = []int{570,762,954,1146,1338}
      ypos  = []int{310,620}
  }
  for x:=0;x<6;x++ {
    game.enemy.groupid[x]=0
    game.own.groupid[x]=0
  }
  for y:=0;y<2;y++ {
    gcnt=0
    fg=false
    for x:=0;x<5;x++ {
      p:=owig.All().At(xpos[x],ypos[y]).RGB()
      if (int(p.R/10)==22 && int(p.G/10)==22 && int(p.B/10)==22) || (int(p.R/10)==13 && int(p.G/10)==22 && p.B==0) {
        if !fg {
          gcnt++
        }
        fg=true
      } else {
        fg=false
      }
      if fg {
        if y==0 {
          game.enemy.groupid[x]=gcnt
          game.enemy.groupid[x+1]=gcnt
        } else {
          game.own.groupid[x]=gcnt
          game.own.groupid[x+1]=gcnt
        }
      }
    }
  }
}

// ----------------------------------------------------------------------------
// get SR numbers from "overview"  screen
func getCurrentSR() {
  if config.dbg_screen {
    fmt.Println("== getCurrentSR ==")
  }
  c:=owig.SRCurrent()
  if (c>0) {
    game.currentSR=c
  }
  return
}

func getHighSR() {
  if config.dbg_screen {
    fmt.Println("== getHighSR ==")
  }
  c:=owig.SRHigh()
  if (c>0) {
    game.highestSR=c
  }
  return
}

// ----------------------------------------------------------------------------
// get SR numbers from "SR / comp. points"  screen
func getCompSR() {
  if config.dbg_screen {
    fmt.Println("== getCompSR ==")
  }
  gain:=owig.SRGain()
  if gain>0 {
    game.currentSR=gain
  }
  return
}

// ----------------------------------------------------------------------------
// get statistic string
func getStats(col int, row int) string {
  if config.dbg_screen {
    fmt.Println("== getStats ==")
  }
  return owig.TStat(col,row)
}

func guessCompObjective() string {
  ret:=""
  hc:=0

  switch owig.res {
    case SIZE_4K:
      hc=owig.From(1748,152).To(2080,152).Th(36).Holes()
    case SIZE_1080:
      hc=owig.From(877,75).To(1040,75).Th(36).Holes()
  }
  owig.Th(-1)
  hc++
  return ret

}

// ----------------------------------------------------------------------------
// Get all relevant information from TAB statistics screen

func parseTabStats() {

  if config.dbg_screen {
    fmt.Println("== parseTabStats ==")
  }
  ts("parsestart")

  // Get Title and Game type
  line:=owig.Title()
  if strings.Contains(line,"|") {
    game.mapname=strings.Split(line,"|")[0]
    game.gametype=strings.Split(line,"|")[1]
  }
  ts("parse1")

  // Get Time
  game.time=owig.TTime()

  // Figure out objective
  if game.gametype=="COMPETITIVE PLAY" && game.time != "0:00" {
    //fmt.Println("Objective=",guessCompObjective())
  }


  // Get hero composition
  game.enemy.isChanged=false
  for x := 0; x < 6; x++ {
    h:=guessHero(x,0)
    if h!="unknown" {
      if game.enemy.hero[x]!=h {
        game.enemy.hero[x]=h
        game.enemy.isChanged=true
        game.enemy.switches[x]++
      }
    }
  }
  game.enemy.isChanged=true
  for x := 0; x < 6; x++ {
    h:=guessHero(x,1)
    if h!="unknown" {
      if game.own.hero[x]!=h {
        game.own.hero[x]=h
        game.own.isChanged=true
        game.own.switches[x]++
      }
      if x==0 {
        game.hero=h
      }
    }
  }
  ts("parse2")

  // Get group composition
  getGroups()
  ts("parse3")

  // Get statistics & medals
  for i:=0;i<6;i++ {
    row:=0
    if (i>2) {
      row=1
    }
    game.lstats[i]=getStats(i%3,row)
    game.rstats[i]=getStats(i%3+3,row)
    game.medals[i]=getMedal(i)
  }
  ts("parsestop")
}

// ----------------------------------------------------------------------------
// Get statistics of Assemble screen
func parseAssembleScreen() {

  if config.dbg_screen {
    fmt.Println("== parseAssembleScreen ==")
  }

  var P Pixel
  switch owig.res {
    case SIZE_4K:
      P=owig.All().At(194,152).RGB()
    case SIZE_1080:
      P=owig.All().At(87,76).RGB()
  }
  if (P.R>P.B) {
    // Red color dominates, so attack
    game.side="attack"
  }
  if (P.B>P.R) {
    // blue color dominates, so defense
    game.side="defense"
  }
}

// ----------------------------------------------------------------------------
// Get statistics of End Screen
func parseEndScreen() {
  // Defeat , Victory or Draw ?
  var crc string

  if config.dbg_screen {
    fmt.Println("== parseEndScreen ==")
  }
  switch owig.res {
    case SIZE_4K:
      owig.From(100,94).To(468,178).Th(145).Filter()
      crc=owig.From(100,94).To(468,178).Cs()
    case SIZE_1080:
      owig.From(50,47).To(234,89).Th(145).Filter()
      crc=owig.From(50,47).To(234,89).Cs()
  }
  if (crc== "5F04") {
    game.result="lost"
  }
  if (crc== "352E0") {
    game.result="draw"
  }
  if (crc== "72735") {
    game.result="won"
  }
}

// ----------------------------------------------------------------------------
// Initialize all game related information

func initGameInfo() {

  if config.dbg_screen {
    fmt.Println("== initGameInfo ==")
  }
  game.mapname=""
  game.gametype=""
  game.hero=""
  game.result=""
  game.enemy.isChanged=false
  game.own.isChanged=false
  game.stats.eleminations=0
  game.stats.objectiveKills=0
  game.stats.objectiveTime=""
  game.stats.objectiveSecs=0
  game.stats.damage=0
  game.stats.deaths=0
  game.time=""
  game.ts=time.Now().UnixNano()/1000000
  for i:=0;i<6;i++ {
    game.stats.medals[i]=""
    game.stats.stats[i]=""
    game.stats.statsText[i]=""
    game.enemy.hero[i]=""
    game.enemy.groupid[i]=0
    game.enemy.switches[i]=0
    game.lstats[i]=""
    game.rstats[i]=""
    game.medals[i]=""
  }
}

// ----------------------------------------------------------------------------
// Main loop (or one shot, if debugging screenshots )
func interpret() {

  getScreen()
  switch game.screen {
    case SC_UNKNOWN:
      // just ignore

    case SC_GAME: {
        dbgWindow("Game screen")
        if game.state==GS_NONE {
          game.state=GS_START
        }
      }

    case SC_MAIN:
      if game.pscreen!=game.screen {
        initGameInfo()
        dbgWindow("Main screen")
        game.state=GS_NONE
        game.side=""
      }
    case SC_ASSEMBLE:
      if game.state!=GS_END||game.image {
        if game.pscreen!=game.screen {
          game.state=GS_START
        }
        dbgWindow("Assemble team detected "+game.side)
        parseAssembleScreen()
      }
    case SC_TAB:
      parseTabStats()
      if game.time=="0:00" {
        game.state=GS_START
      } else {
        game.state=GS_RUN
      }
      dbgWindow("Tab statistics read")
    case SC_VICTORY:
      if game.state==GS_RUN||game.image {
        dbgWindow("Victory!")
        game.state=GS_END
      }
    case SC_DEFEAT:
      if game.state==GS_RUN||game.image {
        dbgWindow("Defeat!")
        game.state=GS_END
      }
    case SC_ENDING:
      game.state=GS_END
      parseEndScreen()
      if game.pscreen!=game.screen {
        dbgWindow("End result: "+game.result)
      }
    case SC_OVERVIEW:
      getCurrentSR()
      getHighSR()
      dbgWindow("SR Current    : "+strconv.Itoa(game.currentSR))
      dbgWindow("SR Season High: "+strconv.Itoa(game.highestSR))
      if config.dbg_screen {
        fmt.Println("SR Current : ",game.currentSR)
        fmt.Println("SR Season High: ",game.highestSR)
      }
    case SC_SRGAIN:
      getCompSR()
      if game.currentSR>0 {
        dbgWindow("SR Current : "+strconv.Itoa(game.currentSR))
        if config.dbg_screen {
          fmt.Println("SR Current : ",game.currentSR)
        }
      }
    default:
      dbgWindow("Detected unknown screen type: "+strconv.Itoa(game.screen))
  }
}

// ----------------------------------------------------------------------------
// Mainloop with own thread

func mainLoop() {
  if (len(os.Args)>1) {
    game.image=true
    // testing, debug with screenshots
    for a:=1;a<len(os.Args);a++ {
      ts("open")
      owig.Open(os.Args[a])
      ts("interpret")
      interpret()
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
      ts("draw")
      mainWindow.Invalidate()
      ts("sleep")
      time.Sleep(time.Duration(config.sleep) * time.Millisecond)
    }
  }
}

func initOwig() {
  if config.dbg_screen {
    fmt.Println("== initOwig ==")
  }
  getIni()
  initGameInfo()
  loadIcons()
  heroStats=map[string][]string{
    "Ana":{
       "Unscoped Accuracy",
       "Scoped Accuracy",
       "Defensive Assists",
       "Nano Boost Assists",
       "Enemies Slept",
       "",
    },
    "Bastion":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Recon Kills",
       "Sentry Kills",
       "Tank Kills",
       "Self Healing",
    },
    "Brigitte":{
       "Offensive Assists",
       "Defensive Assists",
       "Damage Blocked",
       "Armor Provided",
       "",
       "",
       "",
    },
    "Doomfist":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Ability Damage Done",
       "Meteor Strike Kill",
       "Shields Created",
    },
    "D.Va":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Damage Blocked",
       "Self-Destruct Kills",
       "Mechs Called",
       "",
    },
    "Genji":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Damage Reflected",
       "Dragonblade Kills",
       "",
    },
    "Hanzo":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Damage Reflected",
       "Critical Hits",
       "Recon Assists",
       "Dragonstrike Kills",
    },
    "Junkrat":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Enemies Trapped",
       "Rip-Tire Kills",
       "",
    },
    "Lucio":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Sound Barriers Provided",
       "Offensive Assists",
       "Defensive Assists",
       "",
    },
    "McCree":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Critical Hits",
       "Deadeye Kills",
       "Fan The Hammer Kill",
    },
    "Mei":{
       "Damage Blocked",
       "Kill Streak - Best",
       "Enemies Frozen",
       "Blizzard Kills",
       "Self Healing",
       "",
    },
    "Mercy":{
       "Offensive Assists",
       "Defense Assists",
       "Player Resurrected",
       "Blaster Kills",
       "",
       "",
    },
    "Moira":{
       "Secondary Fire Accuracy",
       "Kill Streak - Best",
       "Defensive Assists",
       "Coalescence Kills",
       "Coalescence Healing",
       "Self Healing",
    },
    "Orisa":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Damage Blocked",
       "Offensive Assists",
       "Damage Amplified",
       "",
    },
    "Pharah":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Barrage Kill",
       "Rocket Direct Hits",
       "",
    },
    "Reaper":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Death Blossom Kills",
       "Self Healing",
       "",
    },
    "Reinhardt":{
       "Damage Blocked",
       "Kill Streak - Best",
       "Charge Kills",
       "Fire Strike Kills",
       "Earthshatter Kills",
       "",
    },
    "Roadhog":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Enemies Hooked",
       "Hook Accuracy",
       "Self Healing",
       "Whole Hog Kills",
    },
    "Soldier 76":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Helix Rockets Kills",
       "Tactical Visor Kills",
       "",
    },
    "Sombra":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Offensive Assists",
       "Enemies Hacked",
       "Enemies EMP'D",
       "",
    },
    "Symmetra":{
       "Sentry Turret Kills",
       "Kill Streak - Best",
       "Damage Blocked",
       "Players Teleported",
       "Teleporter Uptime - Average",
       "",
    },
    "Torbjorn":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Torbjorn Kills",
       "Turret Kills",
       "Molten Core Kills",
       "Armor Packs Created",
    },
    "Tracer":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Pulse Bomb kills",
       "Pulse Bombs Attached",
       "",
    },
    "Widowmaker":{
       "Scoped Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Scoped Critical Hits",
       "Recon Assists",
       "",
    },
    "Winston":{
       "Damage Blocked",
       "Kill Streak - Best",
       "Melee Kill",
       "Players Knocked Back",
       "",
       "",
    },
    "Zarya":{
       "Damage Blocked",
       "Kill Streak - Best",
       "High Energy Kills",
       "Average Energy",
       "Graviton Surge Kills",
       "",
    },
    "Zenyatta":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Offensive Assists",
       "Defense Assists",
       "Best Transcendence Heal",
       "",
    },
  }
}

func getStatsline(hero string, i int) string {
  lines,ok:=heroStats[hero]
  if ok {
    return lines[i]
  } else {
    return ""
  }
}


func main() {
  owig=new(OWImg)
  initOwig()
  go mainLoop()

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

