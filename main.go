package main

import (
  "github.com/mertyGit/owig/screenshot"
  "github.com/mertyGit/owig/owocr"
  "github.com/lxn/walk"
  . "github.com/lxn/walk/declarative"
  "fmt"
  "image"
  "os"
  "time"
  "strings"
  "strconv"
)

//Windows window struct
type MyMainWindow struct {
  *walk.MainWindow
  paintWidget *walk.CustomWidget
}
var mainCanvas *walk.Canvas
var mainRect walk.Rectangle
var mainWindow *walk.MainWindow


// Pixel struct 
type Pixel struct {
  R int
  G int
  B int
}

type Sign struct {
  name string
  low Pixel
  high Pixel
}

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

// ----------------------------------------------------------------------------
// Constants for game state
const GS_NONE       = 1  // No game started or searching for game
const GS_START      = 2  // Game Started (waiting / assemble)
const GS_RUN        = 3  // Game is ongoing, stats might change later
const GS_END        = 4  // Game Ended, showing results, stats are final

// ----------------------------------------------------------------------------
// Global variable to hold every information about the game thats is being played
type GameInfo struct {
  screen      int    // Found screen type (tab screen, overview screen etc.etc.)
  pscreen     int    // Previous screen type 
  mapname  string    // Name of map we are playing, like "ILIOS" 
  gametype string    // Game type, like "MYSTERY HEROES"  or "QUICK PLAY"
  hero     string    // What the hero is playing at the moment = own.heroes[0]
  enemy  TeamComp    // Enemy team composition (see below)
  own    TeamComp    // Enemy team composition (see below)
  stats  OwnStats    // Own statistics
  currentSR   int
  highestSR   int
  result   string    // End result (won,lost,draw)
  side     string    // attack or defend
  dmsg     [4]string // debug messages
  state       int    // state
  image      bool    // are we using images instead of screenshots ?
}

type OwnStats struct {
  eleminations   int
  objectiveKills int
  objectiveTime  string
  objectiveSecs  int     //objectiveTime converted to seconds
  damage         int
  deaths         int
  medals         [6]string  // "G","S" or "B" 0=medal for eleminations, 1=for objectiveKills ... etc.
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


// ----------------------------------------------------------------------------
// Configuration struct

type Ini struct {
  sleep int
  stats string
  divider string
  header bool
  dbg_screen bool
  dbg_window bool
  dbg_ocr bool
  dbg_pause int
}

var config Ini

// ----------------------------------------------------------------------------
// captured or loaded image with screen information
var img *image.RGBA

// Resolution type
// 0=3840x2160 (4K)
// 1=1920x1080 (1080p)
var res int

// ----------------------------------------------------------------------------
// Determine resolution type
func setRes() {
  x:=img.Bounds().Max.X;
  y:=img.Bounds().Max.Y;

  res=666 // unknown or unsupported

  if x==3840 && y==2160 {
    res=0
  } else if x==1920 && y==1080 {
    res=1
  }
  if (res==666) {
    fmt.Println("Error: unsupported screen format (",x,",",y,")")
    os.Exit(2)
  }
}


// ----------------------------------------------------------------------------
// Capture Screen 

func grabScreen() {
  var err error
  img, err = screenshot.CaptureScreen()
  if err != nil {
    fmt.Println("Error: Can't make screenshot")
    os.Exit(3)
  }
  setRes()
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
  // Tab statistics ?
  // grey bar + black left upper corner, right under title / time
  switch res {
    case 0:
      s1=pix(90,145)
      s2=pix(115,145)
    case 1:
      s1=pix(42,55)
      s2=pix(70,75)
  }
  if s2.R+s2.G+s2.B < 15 && (int(s1.R/10) == 5 || (int(s1.R/10) == 6)) && (int(s1.G/10) == 5 || (int(s1.G/10) == 6)) && (int(s1.B/10) == 5 || (int(s1.B/10) == 6)) {
    return SC_TAB
  }
  // Victory or Defeat Screen ? 
  // Vertical pixel line in middle with same value
  switch res {
    case 0:
      same=isLine(1820,1050,1820,1100)
      s1=pix(1820,1050)
    case 1:
      same=isLine(910,500,910,520)
      s1=pix(910,500)
  }

  if same {
    // which one, Victory or Defeat, yellow or red ?
    if s1.R>220 && s1.G<5 {
      return SC_DEFEAT
    } else if s1.R>220 && s1.G>190 && int(s1.B/10)==6 {
      return SC_VICTORY
    }
  }
  // Victory/Defeat, but voting & medals part
  // search for blue "Leave Game" button
  switch res {
    case 0:
      crc1=areaAverage(3448,100,3720,172)
    case 1:
      crc1=areaAverage(1724,50,1860,86)
  }
  if crc1=="13232D" || crc1=="12232E" || crc1=="12232D" {
    return SC_ENDING
  }






  // Career Screen ? 
  // White Pixel border for player icon on left top
  switch res {
    case 0:
      same=isLine(100,400,100,600)
      s1=pix(100,400)
      pos=550
    case 1:
      same=isLine(50,200,50,300)
      s1=pix(50,200)
      pos=280
  }

  // + blue stripe behind icon, middle
  if same && s1.R>220 && s1.G>220 && s1.B>220 {
    s1=pix(500,pos)
    for x:=500;x<600 && same==true;x++ {
      s2=pix(x,pos)
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
    switch res {
      case 0:
        crc1=areaAverage(3040,60,3160,180)
        crc2=areaAverage(114,336,399,621)
      case 1:
        crc1=areaAverage(1520,30,1579,89)
        crc2=areaAverage(57,168,198,309)
    }
    if same && crc1==crc2 {
      return SC_OVERVIEW
    }
  }

  // is it an comp gain/loss screen ?
  switch res {
    case 0:
      s1=pix(1928,1306) // purple dot, above comp points count
      // middle white bar (CTF) "Season High" & "Career high" white bars 
      same=isAbove(1700,2060,190) && isAbove(2100,2060,190)
    case 1:
      s1=pix(964,653)
      same=isAbove(851,1030,190) && isAbove(1049,1030,190)
  }
//  if same && int(s1.R/10) == 17 && int(s1.G/10) == 0 && int(s1.B/10) == 22 {
  if same && like(s1,17,0,22,10) {
      return SC_SRGAIN
  }

  // is it "Assemble your team" screen ? (beginning of match)
  // find stripe and pull down menu for skin
  same=false
  switch res {
    case 0:
      if isLine(3586,376,3594,376) && isWhite(3586,376) && isLine(3684,450,3690,450) && isWhite(3684,450) && hasColor(3620,240,2630,300,"C") && (hasColor(1730,2000,1768,2060,"R")||hasColor(1730,2000,1768,2060,"B")) {
        same=true
      }
    case 1:
      if isLine(1793,188,1797,188) && isWhite(1793,188) && isLine(1842,225,1845,225) && isWhite(1842,225) && hasColor(1810,120,1315,150,"C") && (hasColor(865,1000,884,1030,"R")||hasColor(865,1000,884,1030,"B")) {
        same=true
      }
  }
  if same {
    return SC_ASSEMBLE
  }

  //is it a "Main screen" (screen you see after logging in)
  switch res {
    case 0:
      if holes(80,154,1000,224)==15 {
        same=true
      }
    case 1:
      if holes(40,77,500,112)==13 {
        same=true
      }
  }
  if same {
    return SC_MAIN
  }
  return ret
}

// ----------------------------------------------------------------------------
// Wrapper for guessScreen, to set the right game
func getScreen() {
  game.pscreen=game.screen
  game.screen=guessScreen()
}

// ----------------------------------------------------------------------------
// get received medals, returns gold,silver,bronze or "-"

func getMedal(pos int) string {
  var mxpos []int
  var mypos []int
  var medal ="( )"

  switch res {
    case 0:
      mxpos = []int{191,691,1191}
      mypos = []int{1801,1927}
    case 1:
      mxpos = []int{95,345,595}
      mypos = []int{900,963}
  }
  x:=mxpos[pos%3]
  y:=mypos[0]
  if (pos>2) {
    y=mypos[1]
  }
  if config.dbg_screen {
    fmt.Println("GETMEDAL Got Pix:",pix(x,y).R,pix(x,y).G,pix(x,y).B)
  }
  if (pix(x,y).R>89) {
    if (pix(x,y).G>89) {
      if (pix(x,y).B>89) {
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
// guess teamcomposition by taking two pixels from each player
// row (0=enemy team, 1=own team) and position 
// (col=0=always the player)
// returns "unknown" if not recognized or dead (red cross)

func guessHero(col int, row int) string {
  var heros   []Sign
  var xpos    []int
  var ypown   []int
  var ypenemy []int

  switch res {
    case 0:
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

    case 1:
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
    inl=pix(xpos[col],ypown[0])
    inh=pix(xpos[col],ypown[1])
  } else {
    inl=pix(xpos[col],ypenemy[0])
    inh=pix(xpos[col],ypenemy[1])
  }
  if config.dbg_screen {
    fmt.Println("GUESSHERO: Got pix ",inl,inh)
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
      if config.dbg_screen {
        fmt.Println("GUESSHERO: Checked ",el.name," score=",tot," dev=",dev)
      }
  }
  if (score>2) {
    found="unknown"
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

  switch res {
    case 0:
      xpos  = []int{1140,1524,1908,2292,2676}
      ypos  = []int{620,1240}
    case 1:
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
      p:=pix(xpos[x],ypos[y])
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
  game.currentSR=owocr.Img2CurrentSR(img,res)
  return
}

func getHighSR() {
  game.highestSR=owocr.Img2HighSR(img,res)
  return
}

// ----------------------------------------------------------------------------
// get SR numbers from "SR / comp. points"  screen
func getCompSR() {
  game.currentSR=owocr.Img2CompSR(img,res)
  return
}

// ----------------------------------------------------------------------------
// get statistic string
func getStats(col int, row int) string {
  return owocr.GetStats(img,col,row,res)
}

// ----------------------------------------------------------------------------
// Get all relevant information from TAB statistics screen

func parseTabStats() {


  // Get Title and Game type
  line:=owocr.Img2Title(img,res)
  if strings.Contains(line,"|") {
    game.mapname=strings.Split(line,"|")[0]
    game.gametype=strings.Split(line,"|")[1]
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

  // Get group composition
  getGroups()
}

// ----------------------------------------------------------------------------
// Get statistics of Assemble screen
func parseAssembleScreen() {
  var P Pixel
  switch res {
    case 0:
      P=pix(194,152)
    case 1:
      P=pix(87,76)
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

  switch res {
    case 0:
      filter(100,94,468,178,145)
      crc=areaAverage(100,94,468,178)
    case 1:
      filter(50,47,234,89,145)
      crc=areaAverage(50,47,234,89)
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
  for i:=0;i<6;i++ {
    game.stats.medals[i]=""
    game.stats.stats[i]=""
    game.stats.statsText[i]=""
    game.enemy.hero[i]=""
    game.enemy.groupid[i]=0
    game.enemy.switches[i]=0
  }
}

// ----------------------------------------------------------------------------
// Main loop (or one shot, if debugging screenshots )
func interpret() {

  getScreen()
  switch game.screen {
    case SC_UNKNOWN:
      // just ignore

    case SC_MAIN:
      if game.pscreen!=game.screen {
        initGameInfo()
        dbgWindow("Main screen")
        game.state=GS_NONE
        game.side=""
      }
    case SC_ASSEMBLE:
      if game.state!=GS_END||game.image {
        parseAssembleScreen()
        if game.pscreen!=game.screen {
          dbgWindow("Assemble team, we are on "+game.side)
          game.state=GS_START
        }
      }
    case SC_TAB:
      game.state=GS_RUN
      parseTabStats()
      if config.dbg_screen {
        dumpTabStats()
      }
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
      if game.pscreen!=game.screen {
        getCurrentSR()
        getHighSR()
        dbgWindow("SR Current    : "+strconv.Itoa(game.currentSR))
        dbgWindow("SR Season High: "+strconv.Itoa(game.highestSR))
      }
    case SC_SRGAIN:
      getCompSR()
      dbgWindow("SR Current : "+strconv.Itoa(game.currentSR))
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
      loadFile(os.Args[a])
      interpret()
      //fmt.Println("File: ",os.Args[a])
      dbgWindow("Reading File: "+os.Args[a])
      if (a+1<len(os.Args)) {
        time.Sleep(time.Duration(config.dbg_pause) * time.Millisecond)
      }
    }
    for {
      // wait till window is closed
    }
  } else {
    game.image=false
    for {
      grabScreen()
      interpret()
      time.Sleep(time.Duration(config.sleep) * time.Millisecond)
    }
  }
}

func main() {
  getIni()
  initGameInfo()

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
