package main

import (
  "github.com/fatih/color"
  "github.com/go-ini/ini"
  "github.com/mertyGit/owig/screenshot"
  "github.com/mertyGit/owig/owocr"
  "fmt"
  "image"
  "image/png"
  "os"
  "os/exec"
  "time"
  "strings"
  "path/filepath"
)


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
// Global variable to hold every information about the game thats is being played
type GameInfo struct {
  screen      int  // Found screen type (tab screen, overview screen etc.etc.)
  pscreen     int  // Previous screen type (used to determine start/end of games, updates)
  mapname  string  // Name of map we are playing, like "ILIOS" 
  gametype string  // Game type, like "MYSTERY HEROES"  or "QUICK PLAY"
  hero     string  // What the hero is playing at the moment  (equals to own.heroes[0])
  enemy  TeamComp  // Enemy team composition (see below)
  own    TeamComp  // Enemy team composition (see below)
  stats  OwnStats  // Own statistics
  currentSR   int
  highestSR   int
  result   string  // End result (won,lost,draw)
  side     string  // attack or defend
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
  screen bool
  ocr bool
  pause int
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
// return pixel value at given coordinates

func pix(x int, y int) Pixel {
  if x>img.Bounds().Max.X {
    x=img.Bounds().Max.X
  }
  if y>img.Bounds().Max.Y {
    y=img.Bounds().Max.Y
  }
  r,g,b,_ := img.At(x, y).RGBA()
  return Pixel{int(r/257),int(g/257),int(b/257)}
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
// generate checksum based on average R,G,B color values in whole area
// Used to check or validate if some menu items are on screen or not

func areaAverage(xt int,yt int,xr int,yb int) string {
  var R uint64
  var G uint64
  var B uint64

  R=0
  G=0
  B=0
  for x:=xt;x<xr;x++ {
    for y:=yt;y<yb;y++ {
      R+=uint64(pix(x,y).R)
      G+=uint64(pix(x,y).G)
      B+=uint64(pix(x,y).B)
      //owocr.Plot(img,x,y,255,0,0)
    }
  }
  if R+G+B<1 {
    B=1
  }
  rc:=100*R/(R+G+B)
  gc:=100*G/(R+G+B)
  bc:=100*B/(R+G+B)
  line:=fmt.Sprintf("%X%X%X",rc,gc,bc)
  return line
}

// ----------------------------------------------------------------------------
// Filter out artifacts and colors below thresshold
func filter(xt int,yt int,xb int,yb int,tr int) {
  for x:=xt;x<xb;x++ {
    for y:=yt;y<yb;y++ {
      P:=pix(x,y)
      if P.R<tr && P.G<tr && P.B<tr {
        owocr.Plot(img,x,y,0,0,0)
      }
    }
  }
  //owocr.SaveImg(img,"test.png")
}

// ----------------------------------------------------------------------------
// Check to see if area has same color value (within max deviation)
func sameColor(xt int,yt int,xb int,yb int,r int,g int,b int,dev int) bool {
  var ret=true

  for x:=xt;x<xb;x++ {
    for y:=yt;y<yb;y++ {
      P:=pix(x,y)
      if ((P.R>r+dev)||(P.R<r-dev)||(P.G>g+dev)||(P.G<g-dev)||(P.B>b+dev)||(P.B<b-dev)) {
        ret=false
      }
    }
  }
  return ret
}

// ----------------------------------------------------------------------------
// Check to see if area has same color, based on which RGB values are the 
// highest R=Red, G=Green, Y=yellow, B=blue, C=cyan, M=magenta
func hasColor(xt int,yt int,xb int,yb int,c string) bool {
  same:=true
  for x:=xt;x<xb;x++ {
    for y:=yt;y<yb;y++ {
      P:=pix(x,y)
      switch c {
        case "R":
          if !(P.R>P.G && P.R>P.B) {
            same=false
          }
        case "G":
          if !(P.G>P.R && P.G>P.B) {
            same=false
          }
        case "Y":
          if !(P.G>P.B && P.R>P.B) {
            same=false
          }
        case "B":
          if !(P.B>P.R && P.B>P.G) {
            same=false
          }
        case "C":
          if !(P.G>P.R && P.B>P.R) {
            same=false
          }
        case "M":
          if !(P.R>P.G && P.B>P.G) {
            same=false
          }
        default:
          same=false
      }
    }
  }
  return same
}



// Functions to intepret screen
//

// ----------------------------------------------------------------------------
// Get positive difference between two ints (used for OCR stuff)

func getDif(l int,r int) int {
  var d=0;
  if (l>r) {
    d=l-r
  } else {
    d=r-l
  }
  return d
}

// ----------------------------------------------------------------------------
// check if pixel is white (R,G,B > 224)

func isWhite(x int,y int) bool {
  return isAbove(x,y,224)
}

// ----------------------------------------------------------------------------
// check if pixel all R,G and B do have value higher then given value

func isAbove(x int,y int,tr int) bool {
  if pix(x,y).R > tr && pix(x,y).G > tr && pix(x,y).B > tr {
    return true
  }
  return false
}


// ----------------------------------------------------------------------------
// Do we have a line with same r,g,b values ? 
func isLine(xfrom int,yfrom int,xto int,yto int) bool {
  var s1 Pixel
  var s2 Pixel
  var same bool

  s1=pix(xfrom,yfrom)
  same=true
  if xfrom==xto {
    for y:=yfrom;y<yto;y++ {
      s2=pix(xfrom,y)
      //fmt.Println("GOT",xfrom,y," ",s1.R,s2.R," ",s1.G,s2.G," ",s1.B,s2.B)
      if s2.R != s1.R || s2.G != s1.G || s2.B != s1.B {
        same=false
      }
    }
  } else {
    if yfrom==yto {
      for x:=xfrom;x<xto;x++ {
        s2=pix(x,yfrom)
        //fmt.Println("GOT",x,yfrom," ",s1.R,s2.R," ",s1.G,s2.G," ",s1.B,s2.B)
        if s2.R != s1.R || s2.G != s1.G || s2.B != s1.B {
          same=false
        }
      }
    } else {
      // not supported, diagonal lines
      fmt.Println("Warning: matchine diagonal lines")
      return false
    }
  }
  return same
}

func like(s1 Pixel,r int,g int,b int, div int) bool {
  if div<1 {
    div=1
  }
  if int(s1.R/div) == r && int(s1.G/div) == g && int(s1.B/div) == b  {
    return true
  }
  return false
}

func holes(xt int,yt int,xb int,tr int) int {
  hole:=false
  cnt:=0
  for x:=xt;x<xb;x++ {
    if isAbove(x,yt,tr) {
      if hole {
        cnt++
      }
      hole=false
    } else {
      hole=true
    }
  }
  return cnt
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



// ============================================================================
// In game statistics screen (when pressing TAB)
// ============================================================================

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
  //fmt.Println("Got Pix:",pix(x,y).R,pix(x,y).G,pix(x,y).B)
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
  //fmt.Println("Got pix ",inl,inh)

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
      //fmt.Println("Checked ",el.name," score=",tot," dev=",dev)
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
// Load image using given filename (must be .png file)
func loadFile(name string) {
    file, err := os.Open(name)

    if err != nil {
        fmt.Println("Error: File could not be opened")
        os.Exit(1)
    }

    defer file.Close()
    myimg, err := png.Decode(file)

    if err != nil {
      fmt.Println("Error: PNG could not be decoded")
      os.Exit(4)
    }
    img=myimg.(*image.RGBA)
    setRes()
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
// Clear screen for console output

func cls() {
  cmd := exec.Command("cmd", "/c", "clear")
  cmd.Stdout = os.Stdout
  cmd.Run()
}


// ----------------------------------------------------------------------------
// Print medal in appropiate color

func printMedal(m string) {
  switch m {
    case "G":
      color.Set(color.FgYellow)
      fmt.Print("(G)")
      color.Set(color.FgWhite)
    case "S":
      color.Set(color.FgHiWhite)
      fmt.Print("(S)")
      color.Set(color.FgWhite)
    case "B":
      color.Set(color.FgRed)
      fmt.Print("(B)")
      color.Set(color.FgWhite)
    default:
      fmt.Print("( )")
  }
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

func dumpTabStats() {
    cls()
    fmt.Println()
    fmt.Print("             ")
    color.Set(color.FgWhite,color.Bold)
    fmt.Print(game.mapname)
    color.Set(color.FgWhite)
    fmt.Print(" | ")
    color.Set(color.FgYellow,color.Bold)
    fmt.Println(game.gametype)
    color.Unset()
    color.Set(color.FgHiBlue,color.Underline)
    fmt.Println("                                                                                ")
    color.Unset()
    fmt.Println()
    for x := 0; x < 6; x++ {
      if game.enemy.groupid[x]>0 {
        color.Set(color.FgHiWhite)
      }
      fmt.Printf(" %-10s",guessHero(x,0))
      if x<5 && game.enemy.groupid[x] >0 && game.enemy.groupid[x] == game.enemy.groupid[x+1] {
        fmt.Print(" - ");
      } else {
        fmt.Print("   ");
      }
      color.Unset()
    }
    fmt.Println()
    fmt.Println()
    color.HiBlue("-------------------------------------= V S =------------------------------------")
    fmt.Println()
    for x := 0; x < 6; x++ {
      if game.own.groupid[0]==1 && game.own.groupid[x]==1 {
        color.Set(color.FgHiGreen)
      } else {
        if game.own.groupid[x]>0 {
          color.Set(color.FgHiWhite)
        } else {
          color.Set(color.FgWhite)
        }
      }
      fmt.Printf(" %-10s",guessHero(x,1))
      if x<5 && game.own.groupid[x] >0 && game.own.groupid[x] == game.own.groupid[x+1] {
        fmt.Print(" - ")
      } else {
        fmt.Print("   ")
      }
      color.Unset()
    }
    fmt.Println()
    color.Set(color.FgHiBlue,color.Underline)
    fmt.Println("                                                                                ")
    color.Unset()
    fmt.Println()
    fmt.Print(" Eliminations     ")
    printMedal(getMedal(0))
    fmt.Printf(":%8s\n",getStats(0,0))
    fmt.Print(" Objective kills  ")
    printMedal(getMedal(1))
    fmt.Printf(":%8s\n",getStats(1,0))
    fmt.Print(" Objective time   ")
    printMedal(getMedal(2))
    fmt.Printf(":%8s\n",getStats(2,0))
    fmt.Print(" Hero Damage Done ")
    printMedal(getMedal(3))
    fmt.Printf(":%8s\n",getStats(0,1))
    fmt.Print(" Healing Done     ")
    printMedal(getMedal(4))
    fmt.Printf(":%8s\n",getStats(1,1))
    fmt.Print(" Deaths              ")
    fmt.Printf(":%8s\n",getStats(2,1))
    color.Set(color.FgHiBlue,color.Underline)
    fmt.Println("                                                                                ")
    fmt.Println()
    color.Unset()
    fmt.Println(" Stat 1: ",getStats(3,0))
    fmt.Println(" Stat 2: ",getStats(4,0))
    fmt.Println(" Stat 3: ",getStats(5,0))
    fmt.Println(" stat 4: ",getStats(3,1))
    fmt.Println(" stat 5: ",getStats(4,1))
    fmt.Println(" stat 6: ",getStats(5,1))
    color.Set(color.FgHiBlue,color.Underline)
    fmt.Println("                                                                                ")
    color.Unset()
    fmt.Println()
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
  game.side=""
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
        fmt.Println("Main screen",game.side)
      }
    case SC_ASSEMBLE:
      parseAssembleScreen()
      if game.pscreen!=game.screen {
        fmt.Println("Assemble team, we are on",game.side)
      }
    case SC_TAB:
      parseTabStats()
      dumpTabStats()
    case SC_VICTORY:
      if game.pscreen!=game.screen {
        fmt.Println("Victory!")
      }
    case SC_DEFEAT:
      if game.pscreen!=game.screen {
        fmt.Println("Defeat!")
      }
    case SC_ENDING:
      parseEndScreen()
      if game.pscreen!=game.screen {
        fmt.Println("Ending:",game.result)
      }
    case SC_OVERVIEW:
      if game.pscreen!=game.screen {
        getCurrentSR()
        getHighSR()
        fmt.Println("SR Current    :",game.currentSR)
        fmt.Println("SR Season High:",game.highestSR)
      }
    case SC_SRGAIN:
      getCompSR()
      fmt.Println("SR Current    :",game.currentSR)
    default:
      fmt.Println("Detected unknown screen type:",game.screen)
  }
}

// ----------------------------------------------------------------------------
// Read "owig.ini" file

func getIni() {
  config.sleep=1000
  config.stats="owig_stats.csv"
  config.divider=","
  config.header=true
  config.screen=false
  config.ocr=false
  config.pause=2000

  wd,_:=os.Getwd();
  // Try working directory
  inifile:=wd+"\\owig.ini"
  cfg,err := ini.InsensitiveLoad(inifile)
  if err != nil {
    // Try directory .exe is located
    bd:=filepath.Dir(os.Args[0])
    inifile2:=bd+"\\owig.ini"
    cfg,err = ini.InsensitiveLoad(inifile2)
    fmt.Println("Warning: can't read inifile ",inifile," or ",inifile2)
    return
  }
  // Got INI file, so lets read it //
  if cfg.Section("main").HasKey("sleep") {
    config.sleep,_=cfg.Section("main").Key("sleep").Int()
  }
  if cfg.Section("output").HasKey("stats") {
    config.stats=cfg.Section("output").Key("stats").String()
  }
  if cfg.Section("output").HasKey("divider") {
    config.divider=cfg.Section("output").Key("divider").String()
  }
  if cfg.Section("output").HasKey("header") {
    config.header,_=cfg.Section("output").Key("header").Bool()
  }
  if cfg.Section("debug").HasKey("screen") {
    config.screen,_=cfg.Section("debug").Key("screen").Bool()
  }
  if cfg.Section("debug").HasKey("ocr") {
    config.ocr,_=cfg.Section("debug").Key("ocr").Bool()
  }
  if cfg.Section("debug").HasKey("pause") {
    config.pause,_=cfg.Section("debug").Key("pause").Int()
  }
  // Set appropiate values, if needed
  if config.ocr {
    owocr.Debug=true
  }
}

func main() {

  getIni()
  initGameInfo()
  if (len(os.Args)>1) {
    // testing, debug with screenshots
    for a:=1;a<len(os.Args);a++ {
      initGameInfo()
      loadFile(os.Args[a])
      interpret()
      fmt.Println("File: ",os.Args[a])
      if (a+1<len(os.Args)) {
        time.Sleep(time.Duration(config.pause) * time.Millisecond)
      }
    }
  } else {
    for {
      grabScreen()
      interpret()
      time.Sleep(time.Duration(config.sleep) * time.Millisecond)
    }
  }
}
