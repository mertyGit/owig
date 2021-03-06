package main

import (
  "fmt"
  "strings"
  "strconv"
)

type Sign struct {
  name string
  low Pixel
  high Pixel
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

  fnd:=""
  score:=0

  if config.dbg_screen {
    fmt.Println("== guessScreen ==")
  }
  ts("gstart")

  // Is there an image to process ? 
  if !owig.gotimg {
    return SC_UNKNOWN
  }


  // Any Chat icon messing up ?
  checkChat()

  // Tab statistics ?
  // grey bar + black left upper corner, right under title / time
  l1:=false
  l2:=false
  owig.All()
  switch owig.res {
    case SIZE_4K:
      l1=owig.At(90,145).isLike(Pixel{60,50,60},10)
      l2=owig.At(115,145).isLike(Pixel{5,5,5},5)
    case SIZE_WQHD:
      l1=owig.At(56,73).isLike(Pixel{59,50,60},14)
      l2=owig.At(93,100).isLike(Pixel{5,5,5},6)
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
    case SIZE_WQHD:
      same=owig.Box(1217,667,0,35).isLine()
      s1=owig.All().At(1215,667).RGB()
    case SIZE_1080:
      same=owig.Box(910,500,0,20).isLine()
      s1=owig.All().At(910,500).RGB()
  }
  if same {
    // which one, Victory or Defeat, yellow or red ?
    if s1.R>220 && s1.G<5 {
      // red line, but check "T" line in DEFEAT for sure ...
      same=false
      switch owig.res {
        case SIZE_4K:
          same=owig.Box(2330,960,0,120).isLine()
        case SIZE_WQHD:
          same=owig.Box(1554,640,0,60).isLine()
        case SIZE_1080:
          same=owig.Box(1165,480,0,60).isLine()
      }
      if same {
        if config.dbg_screen {
          fmt.Println(" Defeat")
        }
        return SC_DEFEAT
      }
    } else if s1.R>220 && s1.G>190 && int(s1.B/10)==6 {
      // yellow line, but check "I" line in VICTORY! for sure ...
      same=false
      switch owig.res {
        case SIZE_4K:
          same=owig.Box(1500,960,0,120).isLine()
        case SIZE_WQHD:
          same=owig.Box(1000,640,0,60).isLine()
        case SIZE_1080:
          same=owig.Box(750,480,0,60).isLine()
      }
      if same {
        if config.dbg_screen {
          fmt.Println(" Victory")
        }
        return SC_VICTORY
      }
    }
  }

  // Victory/Defeat, but voting & medals part
  // search for blue "Leave Game" button
  switch owig.res {
    case SIZE_4K:
      crc1=owig.From(3448,100).To(3720,172).Cs()
    case SIZE_WQHD:
      crc1=owig.From(2299,67).To(2480,115).Cs()
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
    case SIZE_WQHD:
      same=owig.Box(70,250,0,100).isLine()
      s1=owig.All().At(70,250).RGB()
      pos=380
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
        case SIZE_1080:
          crc1=owig.From(1521,30).To(1579,89).Cs()
          crc2=owig.From(58,168).To(198,309).Cs()
        case SIZE_WQHD:
          crc1=owig.From(2028,40).To(2105,119).Cs()
          crc2=owig.From(77,223).To(265,412).Cs()
      }
      //if crc1==crc2 {
      if crc1==crc1 {
        if config.dbg_screen {
          fmt.Println(" Player Overview screen")
        }
        return SC_OVERVIEW
      } else {
        if config.dbg_screen {
          fmt.Println(" Overview , but not player screen; ",crc1,crc2)
        }
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
    case SIZE_WQHD:
      same=owig.At(1150,1365).isAbove() && owig.At(1640,1365).isAbove()
      owig.At(1280,872).RGB()
    case SIZE_1080:
      same=owig.At(851,1030).isAbove() && owig.At(1049,1030).isAbove()
      owig.At(964,653).RGB()
  }
  owig.Th(-1)

  if same && owig.isLike(Pixel{170,0,220},9) {
      if config.dbg_screen {
        fmt.Println(" SR Gain/Loss screen")
      }
      return SC_SRGAIN
  }

  // is it "Assemble your team" screen ? (beginning of match)
  owig.All().Th(224)
  if config.dbg_screen {
    fmt.Println(" Assemble Check")
  }
  switch owig.res {
    case SIZE_4K:
      fnd,score=owig.At(1498,994).getPattern()
    case SIZE_WQHD:
      fnd,score=owig.At(999,663).getPattern()
    case SIZE_1080:
      fnd,score=owig.At(749,497).getPattern()
  }
  owig.Th(-1)
  if fnd == "ASSEMBLE" && score>850{
    if config.dbg_screen {
      fmt.Println(" Assemble Screen")
    }
    return SC_ASSEMBLE
  }

  // Play Of The Game screen ?
  if config.dbg_screen {
    fmt.Println(" POTG Check")
  }
  owig.All()
  switch owig.res {
    case SIZE_4K:
      fnd,score=owig.At(64,152).getPattern()
    case SIZE_WQHD:
      fnd,score=owig.At(43,101).getPattern()
    case SIZE_1080:
      fnd,score=owig.At(32,76).getPattern()
  }
  //fmt.Println("POTG: ",fnd,score)
  if fnd == "POTG" && score>900{
    if config.dbg_screen {
      fmt.Println(" POTG Screen")
    }
    return SC_POTG
  }

  // Waiting for Respawn ?

  // First, check for respawn message at "spectating" screen
  if config.dbg_screen {
    fmt.Println(" Spectating Check")
  }
  switch owig.res {
    case SIZE_4K:
      owig.Box(3338,134,50,50).Y2W()
      owig.All().Th(224)
      fnd,score=owig.At(3338,134).getPattern()
    case SIZE_WQHD:
      owig.Box(2225,89,30,30).Y2W()
      owig.All().Th(224)
      fnd,score=owig.At(2225,89).getPattern()
    case SIZE_1080:
      owig.Box(1669,67,25,25).Y2W()
      owig.All().Th(224)
      fnd,score=owig.At(1669,67).getPattern()
  }
  owig.Th(-1)
  //fmt.Println("RESPAWN 1: ",fnd,score)
  if fnd == "RESPAWN" && score>930 {
    return SC_RESPAWN
  }

  // Second, check for respawn message at "kill cam" screen (higher)
  //fmt.Println("searching respawn2: ")
  if config.dbg_screen {
    fmt.Println(" Kill cam Check")
  }
  switch owig.res {
    case SIZE_4K:
      owig.Box(3323,45,50,50).Y2W()
      owig.All().Th(224)
      fnd,score=owig.At(3323,45).getPattern()
    case SIZE_WQHD:
      owig.Box(2216,31,30,30).Y2W()
      owig.All().Th(224)
      fnd,score=owig.At(2216,31).getPattern()
    case SIZE_1080:
      owig.Box(1662,23,25,25).Y2W()
      owig.All().Th(224)
      fnd,score=owig.At(1662,23).getPattern()
  }
  owig.Th(-1)
  //fmt.Println("RESPAWN 2: ",fnd,score)
  if fnd == "RESPAWN" && score>930 {
    return SC_RESPAWN
  }
  //owig.Save("ocr.png")

  // Game screen ?
  if onFireIcon() {
    if config.dbg_screen {
      fmt.Println(" Game Screen")
    }
    return SC_GAME
  }

  //is it a "Main screen" (screen you see after logging in)
  owig.All().Th(224)
  switch owig.res {
    case SIZE_4K:
      fnd,score=owig.At(414,134).getPattern()
    case SIZE_WQHD:
      fnd,score=owig.At(271,90).getPattern()
    case SIZE_1080:
      fnd,score=owig.At(207,67).getPattern()
  }
  owig.Th(-1)
  //fmt.Println("MAIN: ",fnd,score)
  if fnd == "OVERWATCH" && score>900{
    if config.dbg_screen {
      fmt.Println(" Main Screen")
    }
    return SC_MAIN
  }
  return SC_UNKNOWN
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
  var medal =""

  if config.dbg_screen {
    fmt.Println("== getMedal ==")
  }
  switch owig.res {
    case SIZE_4K:
      mxpos = []int{191,691,1191}
      mypos = []int{1801,1927}
    case SIZE_WQHD:
      mxpos = []int{137,461,794}
      mypos = []int{1200,1283}
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
    fmt.Println(" Got Pix:",owig.RGB(),"(",pos,")")
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
        {"Ana"          ,Pixel{140,143,150},Pixel{99,92,92}},
        {"Bastion"      ,Pixel{176,164,132},Pixel{171,206,188}},
        {"Brigitte"     ,Pixel{113,35,22}  ,Pixel{208,133,114}},
        {"Doomfist"     ,Pixel{84,61,46}   ,Pixel{167,138,114}},
        {"D.Va"         ,Pixel{124,80,100} ,Pixel{199,156,127}},
        {"Genji"        ,Pixel{80,98,115},Pixel{155,150,134}},
        {"Hanzo"        ,Pixel{107,68,44}  ,Pixel{60,55,59}},
        {"Junkrat"      ,Pixel{69,72,84}   ,Pixel{199,190,160}},
        {"Lucio"        ,Pixel{79,48,36}   ,Pixel{65,88,12}},
        {"McCree"       ,Pixel{99,23,26}   ,Pixel{23,18,22}},
        {"Mei"          ,Pixel{171,168,164},Pixel{23,23,23}},
        {"Mercy"        ,Pixel{46,51,60}   ,Pixel{225,192,157}},
        {"Moira"        ,Pixel{102,106,120},Pixel{188,125,104}},
        {"Orisa"        ,Pixel{44,39,36}   ,Pixel{188,181,22}},
        {"Pharah"       ,Pixel{59,58,66}   ,Pixel{130,94,78}},
        {"Reaper"       ,Pixel{26,25,26}   ,Pixel{42,35,23}},
        {"Reinhardt"    ,Pixel{46,42,39}   ,Pixel{83,77,74}},
        {"Roadhog"      ,Pixel{171,138,122},Pixel{57,60,66}},
        {"Soldier 76"   ,Pixel{172,188,194},Pixel{99,0,0}},
        {"Sombra"       ,Pixel{123,65,44}  ,Pixel{195,132,87}},
        {"Symmetra"     ,Pixel{105,98,87}  ,Pixel{68,74,66}},
        {"Torbjörn"     ,Pixel{203,185,135},Pixel{42,40,42}},
        {"Tracer"       ,Pixel{165,99,77}  ,Pixel{132,32,0}},
        {"Widowmaker"   ,Pixel{99,119,179} ,Pixel{23,24,23}},
        {"Winston"      ,Pixel{32,31,40}   ,Pixel{95,101,111}},
        {"Wrecking Ball",Pixel{57,54,73}   ,Pixel{173,157,148}},
        {"Zarya"        ,Pixel{55,97,107}  ,Pixel{172,104,78}},
        {"Zenyatta"     ,Pixel{29,32,29}   ,Pixel{36,30,22}},
      }
      xpos   = []int{959,1343,1727,2111,2495,2879}
      ypown  = []int{1260,1180}
      ypenemy= []int{690,610}

    case SIZE_WQHD:
      heros = []Sign{
        {"Ana"          ,Pixel{106,106,111},Pixel{139,130,123}},
        {"Bastion"      ,Pixel{142,134,109},Pixel{171,207,187}},
        {"Brigitte"     ,Pixel{181,162,159},Pixel{203,129,112}},
        {"Doomfist"     ,Pixel{58,43,33}   ,Pixel{138,108,85}},
        {"D.Va"         ,Pixel{179,171,159},Pixel{166,124,97}},
        {"Genji"        ,Pixel{48,40,42}   ,Pixel{183,177,162}},
        {"Hanzo"        ,Pixel{127,84,52}  ,Pixel{24,19,23}},
        {"Junkrat"      ,Pixel{126,101,87} ,Pixel{189,182,160}},
        {"Lucio"        ,Pixel{55,45,36}   ,Pixel{66,91,6}},
        {"McCree"       ,Pixel{91,21,22}   ,Pixel{68,49,47}},
        {"Mei"          ,Pixel{156,157,156},Pixel{23,24,22}},
        {"Mercy"        ,Pixel{66,78,89}   ,Pixel{133,108,85}},
        {"Moira"        ,Pixel{84,77,67}   ,Pixel{182,117,99}},
        {"Orisa"        ,Pixel{20,19,20}   ,Pixel{195,185,29}},
        {"Pharah"       ,Pixel{58,68,83}   ,Pixel{132,92,77}},
        {"Reaper"       ,Pixel{13,13,13}   ,Pixel{45,38,23}},
        {"Reinhardt"    ,Pixel{67,65,61}   ,Pixel{91,82,76}},
        {"Roadhog"      ,Pixel{88,44,36}   ,Pixel{49,48,51}},
        {"Soldier 76"   ,Pixel{86,95,99}   ,Pixel{86,0,0}},
        {"Sombra"       ,Pixel{173,102,66} ,Pixel{195,132,87}},
        {"Symmetra"     ,Pixel{89,46,37}   ,Pixel{48,58,76}},
        {"Torbjörn"     ,Pixel{197,185,147},Pixel{44,43,44}},
        {"Tracer"       ,Pixel{124,59,48}  ,Pixel{39,9,0}},
        {"Widowmaker"   ,Pixel{70,84,147}  ,Pixel{24,22,22}},
        {"Winston"      ,Pixel{44,43,44}   ,Pixel{99,103,99}},
        {"Wrecking Ball",Pixel{66,52,57}   ,Pixel{173,158,143}},
        {"Zarya"        ,Pixel{52,50,47}   ,Pixel{195,125,97}},
        {"Zenyatta"     ,Pixel{100,110,109},Pixel{6,6,6}},
      }
      xpos   = []int{640,896,1152,1408,1664,1920}
      ypown  = []int{835,787}
      ypenemy= []int{455,407}

    case SIZE_1080:
      heros = []Sign{
        {"Ana"          ,Pixel{103,103,109},Pixel{153,143,138}},
        {"Bastion"      ,Pixel{132,126,100},Pixel{171,206,186}},
        {"Brigitte"     ,Pixel{175,155,149},Pixel{203,129,112}},
        {"Doomfist"     ,Pixel{55,42,32}   ,Pixel{144,115,91}},
        {"D.Va"         ,Pixel{179,169,156},Pixel{135,99,78}},
        {"Genji"        ,Pixel{44,35,36}   ,Pixel{154,150,136}},
        {"Hanzo"        ,Pixel{126,83,52}  ,Pixel{22,17,22}},
        {"Junkrat"      ,Pixel{133,108,92} ,Pixel{181,174,154}},
        {"Lucio"        ,Pixel{55,46,36}   ,Pixel{67,92,7}},
        {"McCree"       ,Pixel{90,21,22}   ,Pixel{46,34,35}},
        {"Mei"          ,Pixel{157,158,157},Pixel{23,24,22}},
        {"Mercy"        ,Pixel{66,77,89}   ,Pixel{139,111,89}},
        {"Moira"        ,Pixel{83,74,64}   ,Pixel{182,117,99}},
        {"Orisa"        ,Pixel{19,18,19}   ,Pixel{195,186,29}},
        {"Pharah"       ,Pixel{59,69,82}   ,Pixel{120,83,68}},
        {"Reaper"       ,Pixel{13,13,13}   ,Pixel{46,39,23}},
        {"Reinhardt"    ,Pixel{74,71,64}   ,Pixel{92,83,77}},
        {"Roadhog"      ,Pixel{84,43,36}   ,Pixel{51,49,53}},
        {"Soldier 76"   ,Pixel{82,91,97}   ,Pixel{84,0,0}},
        {"Sombra"       ,Pixel{173,102,66} ,Pixel{195,132,87}},
        {"Symmetra"     ,Pixel{89,46,36}   ,Pixel{34,50,75}},
        {"Torbjörn"     ,Pixel{196,186,151},Pixel{44,43,44}},
        {"Tracer"       ,Pixel{120,55,44}  ,Pixel{29,6,0}},
        {"Widowmaker"   ,Pixel{68,84,147}  ,Pixel{23,21,22}},
        {"Winston"      ,Pixel{44,43,44}   ,Pixel{107,112,107}},
        {"Wrecking Ball",Pixel{60,46,52}   ,Pixel{172,156,141}},
        {"Zarya"        ,Pixel{52,47,44}   ,Pixel{194,124,97}},
        {"Zenyatta"     ,Pixel{103,112,111},Pixel{6,6,6}},
      }
      xpos   = []int{480,672,864,1056,1248,1440}
      ypown  = []int{626,590}
      ypenemy= []int{341,305}
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
    fmt.Printf("Got %d,%d Pixel{%d,%d,%d},Pixel{%d,%d,%d}\n",col,row,inl.R,inl.G,inl.B,inh.R,inh.G,inh.B)
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
      ypos  = []int{644,1218}
    case SIZE_WQHD:
      xpos  = []int{760,1016,1272,1528,1784}
      ypos  = []int{429,812}
    case SIZE_1080:
      xpos  = []int{570,762,954,1146,1338}
      ypos  = []int{322,609}
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
  // Figure out own group id my name of choosen hero
  for x:=5;x>-1;x-- {
    if game.own.hero[x]==game.hero {
      game.group=game.own.groupid[x]
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

// ----------------------------------------------------------------------------
// Get all relevant information from TAB statistics screen

func parseTabStats() {

  if config.dbg_screen {
    fmt.Println("== parseTabStats ==")
  }
  ts("parsestart")

  // Dont bother to intepret title, game type or game time if
  // chat icon(s) blocking your view

  if !game.chat {
    // Get Title and Game type
    line:=owig.Title()
    if strings.Contains(line,"|") {
      game.mapname=cleanMapname(strings.Split(line,"|")[0])
      game.gametype=cleanGametype(strings.Split(line,"|")[1])
    }
    ts("parse1")

    // Get Time
    game.time=owig.TTime()
  }


  // Get own played hero, (not always the one on left bottom of composition)
  h:=owig.MyHero()
  if !(h=="") {
    game.hero=h
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
// Get info from game screen

func parseGameScreen() {

  var y int
  var y2 int
  var x int
  var xl int
  var xr int
  var wx int
  var wy int
  var cnt int
  var line string
  var font map[string][][]int

  if config.dbg_screen {
    fmt.Println("== parseGameScreen ==")
  }

  // Determine if we are attacking or defending
  //if game.gametype=="COMPETITIVE PLAY" && game.time != "0:00" {
  if true {
    switch owig.res {
      case SIZE_4K:
        y=200
        y2=190
        wx=90
        wy=70
        font=FontScore4K
      case SIZE_WQHD:
        y=133
        y2=123
        wx=60
        wy=47
        font=FontScoreWQHD
      case SIZE_1080:
        y=100
        y2=95
        wx=45
        wy=35
        font=FontScore1080
    }
    fnd :=false
    fndl:=false
    fndr:=false
    owig.All()
    for s:=0;s<400 && !fndl;s++ {
      xl=owig.width/2-s
      if owig.At(xl,y).Blue()-owig.Red() > 150 && owig.Green() < owig.Blue() {
        fndl=true
      }
    }
    for s:=0;s<400 && !fndr;s++ {
      xr=owig.width/2+s
      if owig.At(xr,y).Red()-owig.Blue() > 150 && owig.Green() < 50 {
        fndr=true
      }
    }
    if fndl && fndr {
      fnd=true
    }
    if fnd {
      // determine if we are attacking or defining, based on colored box
      if owig.At(xr,y2).Red()>150 && owig.At(xl,y2).Blue()<150 {
        game.side="defend"
      } else if owig.At(xr,y2).Red()<150 && owig.At(xl,y2).Blue()>150 {
        game.side="attack"
      }
      // now we have found the boxes, get the points out of it
      if config.dbg_screen {
        fmt.Println(" Getting scores: ")
      }
      owig.Box(xl-wx,y2-wy,wx,wy).Th(200)
      lscore:=-1
      for x:=0; x<wx && lscore<0; x++  {
        for y:=0; y<wy && lscore<0 ; y++ {
          if owig.At(x,y).isAbove() {
            ch,_,p:=owig.getChar(font)
            if p>700 {
              lscore,_=strconv.Atoi(ch)
              if config.dbg_screen {
                fmt.Println("  Score Left  :",lscore,"(",p,")")
              }
            }
          }
        }
      }
      owig.Box(xr,y2-wy,wx,wy)
      rscore:=-1
      for x:=0; x<wx && rscore<0; x++  {
        for y:=0; y<wy && rscore<0 ; y++ {
          if owig.At(x,y).isAbove() {
            ch,_,p:=owig.getChar(font)
            if p>700 {
              rscore,_=strconv.Atoi(ch)
              if config.dbg_screen {
                fmt.Println("  Score Right :",rscore,"(",p,")")
              }
            }
          }
        }
      }
      if lscore > -1 && rscore > -1 {
        game.compdef=rscore
        game.compatt=lscore
      }
      //owig.Save("score.png")
      owig.Th(-1)
    }
  }

  // Get opjective text

  // First, figure out bounderies of text to intepret (without counting 
  // time on left or any other things on the sides
  switch owig.res {
    case SIZE_4K:
      y=150
      y2=30
      owig.All().Th(220)
      font=FontObjective4K
    case SIZE_WQHD:
      y=100
      y2=20
      owig.All().Th(200)
      font=FontObjectiveWQHD
    case SIZE_1080:
      y=75
      y2=15
      owig.All().Th(200)
      font=FontObjective1080
  }

  // Find right border of text
  fnd:=false
  for s:=0;s<400 && !fnd;s++ {
    xr=owig.width/2+s
    if owig.At(xr,y).isAbove() {
      cnt=0
    } else {
      cnt++
    }
    if cnt>y2{
      fnd=true
      xr-=y2/2
    }
  }
  // Find left border of text
  if fnd {
    cnt=0
    fnd=false
    for s:=0;s<400 && !fnd;s++ {
      xl=owig.width/2-s
      if owig.At(xl,y).isAbove() {
        cnt=0
      } else {
        cnt++
      }
      if cnt>y2{
        fnd=true
        xl+=y2/2
      }
    }
  }

  // if both are found, ocr the text within
  if fnd {
    if config.dbg_ocr {
      fmt.Println(" Searching for Objective text")
    }
    cnt=1
    owig.Box(xl,y-(2*y2/3),xr-xl,4*y2/3)
    for cx:=0; cx<xr-xl; cx++ {
      for cy:=0; cy<4*y2/3; cy++ {
        if owig.At(cx,cy).isAbove() {
          if config.dbg_ocr {
            fmt.Println(" Character # ",cnt," at",cx,cy)
          }
          ch,_,_:=owig.getChar(font)
          line+=ch
          cnt++
        }
      }
    }
    if config.dbg_ocr {
      fmt.Println(" got line:",line)
    }
    if line!="" {
      // Check objective
      if strings.HasSuffix(line,"VEA") {
        // .....OBJECTIVE A"
        game.objective="A"
        game.state=GS_RUN
      } else if strings.HasSuffix(line,"VEB") {
        // .....OBJECTIVE B"
        game.objective="B"
        game.state=GS_RUN
      } else if strings.Contains(line,"PAY")||strings.HasSuffix(line,"OAD")||strings.HasPrefix(line,"ESC")||strings.HasPrefix(line,"STOP") {
        // ...PAYLOAD
        game.objective="PAYLOAD"
        game.state=GS_RUN
      }
      if strings.HasPrefix(line,"PREP") {
        // PREPARTE TO ...
        game.objective="WAITING"
        game.state=GS_START
        // Clear payload stats
        game.plpoint=0
        game.pltrack=0
        game.pltotal=0

      }


      if strings.HasPrefix(line,"DEF") || strings.HasPrefix(line,"BEF") || strings.HasPrefix(line,"STOP") || strings.HasSuffix(line,"NSES") {
        // DEFEND OBJECTIVE ...  or STOP THE ... or PREPARE YOUR DEFENSES...
        //(BEF=DEF, but sometimes wrong intepreted..)
        game.side="defend"
        if config.dbg_ocr {
          fmt.Println(" Defending")
        }
      } else if strings.HasPrefix(line,"AT") || strings.HasPrefix(line,"ESC") || strings.HasSuffix(line,"ACK") {
        // ATTACK OBJECTIVE ...  or ESCORD THE ... or PREPARE TO ATTACK...
        game.side="attack"
        if config.dbg_ocr {
          fmt.Println(" Attacking")
        }
      }

      if game.objective=="PAYLOAD" {
        // Figure out what points we have captured
        // scan for pointindicators
        switch owig.res {
          case SIZE_4K:
            y=286
            y2=320
            x=1571
            xr=2370
          case SIZE_WQHD:
            y=191
            y2=213
            x=1048
            xr=1580
          case SIZE_1080:
            y=143
            y2=160
            x=786
            xr=1185
        }
        ep:=0
        owig.All()
        last:="_"
        px:=make([]int,5,5)
        pt:=make([]string,5,5)
        pcnt:=0
        for cx:=x;cx<xr && ep==0;cx++ {

          owig.At(cx,y2)
          if owig.Red()>210 && owig.Green()<20 && owig.Blue()<owig.Red() {
            // Red pointer found ?
            if last != "R" {
              xl=cx
              pt[pcnt]="R"
              last="R"
            }
          } else if owig.Blue()>220 && owig.Red()<100 && owig.Green()<owig.Blue() {
            // Blue pointer found ?
            if last != "B" {
              xl=cx
              pt[pcnt]="B"
              last="B"
            }
          } else if owig.Red()>210 && owig.Green()>150 && owig.Blue()<100 {
            // Yellow pointer found ?
            if last != "Y" {
              xl=cx
              pt[pcnt]="Y"
              last="Y"
            }
          } else if owig.Red()>220 && owig.Green()>220 && owig.Blue()>220 {
            // White pointer found ?
            if last != "W" {
              xl=cx
              pt[pcnt]="W"
              last="W"
            }
          } else {
            if last != "_" {
              px[pcnt]=xl+((cx-xl)/2)
              if last == "W" || last == "Y" {
                ep=px[pcnt]
              }
              pcnt++
              last="_"
            }
          }
        }
        if config.dbg_screen {
          fmt.Print(" Got payload: ")
          for cx:=0;cx<pcnt;cx++ {
            fmt.Print(px[cx],":",pt[cx]," ")
          }
          fmt.Println();
          fmt.Println(" Endpoint",ep);
        }


        // Second, get start & lenght of line
        cnt:=0
        ch:=""
        for cx:=x;cx<ep;cx++ {
          owig.At(cx,y)
          if (ch=="" || ch=="B") && owig.Blue()>220 && owig.Red()<100 && owig.Green()<owig.Blue() {
            ch="B"
          } else if (ch=="" || ch=="R") && owig.Red()>210 && owig.Green()<20 && owig.Blue()<owig.Red() {
            ch="R"
          } else if cnt>0 {
            // end of line reached
            cx=ep
          }
          if !(ch=="") {
            if cnt==0 {
              x=cx
            }
            cnt++
          }
        }
        end:=x+cnt
        perc:=100*cnt/(ep-x)
        xl=x
        xr=0
        pp:=0
        for cx:=0;cx<pcnt && xr==0;cx++ {
          if px[cx]<end {
            xl=px[cx]
            pp++
          } else if px[cx]>end {
            xr=px[cx]
          }
        }
        tperc:=100*(end-xl)/(xr-xl)

        // At this point, we have all information of the progress line 
        game.plpoint=pp
        game.plamount=pcnt
        game.pltrack=tperc
        game.pltotal=perc

        if game.pltotal==100 {
          game.plpoint=pcnt
          game.pltrack=0
        }

        // weird bugfix
        if game.pltrack<0 {
          game.pltrack=0
        }

        if config.dbg_screen {
          fmt.Println(" Color: ",ch," Start: ",x," End: ",end," Percentage: ",perc)
          fmt.Println(" On traject: ",xl," to: ",xr," Percentage: ",tperc)
          fmt.Println(" Position: ",pp," of",pcnt)
        }
      } // PAYLOAD

      // If objective, delete any payload info, just to be sure
      if game.objective=="A" || game.objective=="B" {
        game.plpoint=0
        game.plamount=0
        game.pltrack=0
        game.pltotal=0
      }
    } // Objective text found
  } // Objective found
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
    case SIZE_WQHD:
      P=owig.All().At(116,101).RGB()
    case SIZE_1080:
      P=owig.All().At(87,76).RGB()
  }
  if (P.R>P.B) {
    // Red color dominates, so attack
    game.side="attack"
  }
  if (P.B>P.R) {
    // blue color dominates, so defense
    game.side="defend"
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
    case SIZE_WQHD:
      owig.From(67,63).To(312,119).Th(145).Filter()
      crc=owig.From(67,63).To(312,119).Cs()
    case SIZE_1080:
      owig.From(50,47).To(234,89).Th(145).Filter()
      crc=owig.From(50,47).To(234,89).Cs()
  }
  if config.dbg_screen {
    fmt.Println("Endscreen crc =",crc)
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
// Try to figure out if someone is chatting (and f*cking up stats with overlay)

func checkChat() {
  var x1,y,w,x2 int
  switch owig.res {
    case SIZE_4K:
      x1=122
      x2=202
      y=130
      w=10
    case SIZE_WQHD:
      x1=81
      x2=133
      y=88
      w=5
    case SIZE_1080:
      x1=61
      x2=101
      y=65
      w=5
  }
  if owig.Box(x1,y,w,w).SameColor(Pixel{134,225,0},5) || owig.Box(x1,y,w,w).SameColor(Pixel{0,220,225},5) || owig.Box(x2,y,w,w).SameColor(Pixel{134,225,0},5) || owig.Box(x2,y,w,w).SameColor(Pixel{0,220,225},5) {
    if config.dbg_screen {
      fmt.Println("chat icon present")
    }
    game.chat=true
  } else {
    game.chat=false
  }
}

// ----------------------------------------------------------------------------
// Try to figure out if we have a "on fire" icon, so we are in a game
func onFireIcon() bool {

  var ystart int
  var yend   int
  var xstart int
  var x,dx   int
  var y,dy   int
  var score  int
  var found  string

  fnd:=false


  if config.dbg_screen {
    fmt.Println("== onFireIcon ==")
  }
  switch owig.res {
    case SIZE_4K:
      ystart=2040
      yend=1980
      xstart=528
      dx=24
      dy=38
    case SIZE_WQHD:
      ystart=1360
      yend=1320
      xstart=352
      dx=16
      dy=25
    case SIZE_1080:
      ystart=1020
      yend=990
      xstart=264
      dx=12
      dy=19
  }
  owig.All()
  x=xstart
  for y=ystart;y>yend&&!fnd;y-- {
    if owig.At(x,y).isAbove() {
      fnd=true
    }
  }
  if !fnd {
    return false
  }
  x=x-dx
  y=y-dy
  found,score=owig.At(x,y).getPattern()
  //fmt.Println("FIREICON: ",found,score)
  if (score < 800) {
    // Try again, but this time, only for "real white", in case background
    // is already light
    fnd=false
    owig.All()
    x=xstart
    for y=ystart;y>yend&&!fnd;y-- {
      if owig.At(x,y).isAbove() {
        // Make sure it is a "true" white color and
        // filter out all other whites below
        if owig.Red()==owig.Blue() && owig.Blue()==owig.Green() {
          owig.Th(owig.Red()-1)
          fnd=true
        }
      }
    }
    if !fnd {
      return false
    }
    x=x-dx
    y=y-dy
    found,score=owig.At(x,y).getPattern()
    //fmt.Println("FIREICON2: ",found,score)
  }
  owig.Th(-1)
  if found=="FIRE" && score>800 {
    return true
  }
  return false
}


// ----------------------------------------------------------------------------
// Main loop (or one shot, if debugging screenshots )
func interpret() {

  getScreen()
  switch game.screen {
    case SC_UNKNOWN:
      // just ignore

    case SC_GAME:
      if game.pscreen!=game.screen {
        dbgWindow("Game screen")
      }
      if game.state==GS_NONE {
        game.state=GS_START
      }
      parseGameScreen()
    case SC_MAIN:
      if game.pscreen!=game.screen {
        initGameInfo()
        dbgWindow("Main screen")
        game.state=GS_NONE
        game.side=""
      }
    case SC_ASSEMBLE:
      if game.state!=GS_END||game.image {
        if game.pscreen != game.screen {
          game.state=GS_START
          game.forceM=false
          game.forceT=false
          game.plpoint=0
          game.plamount=0
          game.pltrack=0
          game.pltotal=0
        }
        dbgWindow("Assemble team detected "+game.side)
      }
      parseAssembleScreen()
    case SC_TAB:
      parseTabStats()
      if game.time=="0:00" {
        game.state=GS_START
        game.plpoint=0
        game.plamount=0
        game.pltrack=0
        game.pltotal=0
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
    case SC_POTG:
      if game.pscreen!=game.screen {
        dbgWindow("POTG playing")
      }
      if game.state!=GS_END {
        game.state=GS_END
      }
    case SC_RESPAWN:
      if game.pscreen!=game.screen {
        dbgWindow("Waiting for Respawn")
      }
    default:
      dbgWindow("Detected unknown screen type: "+strconv.Itoa(game.screen))
  }
}
