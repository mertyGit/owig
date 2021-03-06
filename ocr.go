package main

import (
  "fmt"
  "strings"
  "strconv"
)


//------------------------------------------------------------
// Try to match location of char to complete char for given font

func (g *OWImg) getChar(font map[string][][]int) (string,int,int) {
  var cleanCol=false
  var bitc uint8
  var bytec,oct int

  bx:=g.pos.ax
  by:=g.pos.y2-g.pos.y
  bx2:=0
  by2:=0
  found:="?"
  hitscore:=0
  highscore:=0
  pcnt:=0



  g.Flood(Pixel{255,255,160}) // Bright Yellow, used as "marker"
  if g.fc<5 && g.res==SIZE_4K { // just picked up noise, ignore
    return "",0,0
  }
  for x:=bx; x<g.pos.x2-g.pos.x && !cleanCol; x++ {
    cleanCol=true
    for y:=0; y<g.pos.y2-g.pos.y ; y++ {
      if g.At(x,y).isLike(Pixel{255,255,160},0) {
        cleanCol=false
        g.Plot(Pixel{160,160,255})
        if bx2<x {
          bx2=x
        }
        if by2<y {
          by2=y
        }
        if by>y {
          by=y
        }
      }
    }
  }
  // Find closed character in font
  w:=len(font["search"][0])
  h:=len(font["search"])

  //Clean Search Entry
  for x:=0;x<w;x++ {
    for y:=0;y<h;y++ {
      font["search"][y][x]=0
    }
  }
  if config.dbg_ocr {
    fmt.Println("border: ",g.pos.x+bx,g.pos.y+by,g.pos.x+bx2,g.pos.y+by2)
  }

  //Place character into "search" entry of font
  for y:=by;y<=by2;y++ {
    bitc=0
    bytec=0
    oct=0
    for x:=bx;x<=bx2;x++ {
      if g.At(x,y).isLike(Pixel{160,160,255},0) {
        oct|=1<<(7-bitc)
        // Color pixel far below thresshold, to prevent rescan of pixel
        g.Plot(Pixel{0,10,20})
      }
      bitc++
      if bitc>7 {
        if y-by > -1 && y-by < h && bytec < w {
          font["search"][y-by][bytec]=oct
        }
        bitc=0
        bytec++
        oct=0
      }
    }
    if bitc>0 {
      if y-by > -1 && y-by < h && bytec < w {
        font["search"][y-by][bytec]=oct
      }
    }
  }

  // Dump character (used for debug/retrieving font info)
  if config.dbg_ocr {
    fmt.Println("  \"?\": {")
    for y:=0;y<h;y++ {
      fmt.Print("    {")
      for x:=0;x<w;x++ {
        fmt.Print(font["search"][y][x])
        if x<w-1 {
          fmt.Print(",")
        }
      }
      fmt.Println("},")
    }
    fmt.Println("  },")
  }



  // Get score for each of the found character
  for k, v := range font {
    if k != "search" {
      hitscore=0
      pcnt=0
      for r:=0; r<h; r++ {
        bitc=0
        bytec=0
        oct=0
        for c:=0; c<w*8; c++ {
          pcnt++
          mask:=(1<<(7-bitc))
          if v[r][bytec] & mask == font["search"][r][bytec] & mask {
            hitscore++
          } else {
            hitscore--
          }
          bitc++
          if (bitc>7) {
            bitc=0
            bytec++
            oct=0
          }
        }
      }
      // make hitscore percentage of matched pixels
      hitscore=1000*hitscore/pcnt
      //fmt.Println(" key ",k," hitscore:",hitscore)
      if hitscore>highscore {
        highscore=hitscore
        // font can have double entries for same char, denoted by trailing _
        found=strings.Replace(k,"_","",-1)
      }
    }
  }


  // Draw box around for debugging purposes
  if config.dbg_ocr {
    for x:=bx;x<bx2;x++ {
      if !g.At(x,by).isAbove() {
        g.Plot(Pixel{20,20,20})
      }
      if !g.At(x,by2).isAbove() {
        g.Plot(Pixel{20,20,20})
      }
    }
    for y:=by;y<by2;y++ {
      if !g.At(bx,y).isAbove() {
        g.Plot(Pixel{20,20,20})
      }
      if !g.At(bx2,y).isAbove() {
        g.Plot(Pixel{20,20,20})
      }
    }
  }


  return found,bx2-bx,highscore
}

//------------------------------------------------------------
// Try to match pattern of screen in icon font list
// Will get square of 50x50 (4K) or 25x25 (1080)

func (g *OWImg) getPattern() (string,int) {
  var bitc uint8
  var bytec,oct int
  var font map[string][][]int

  bx,by:=0,0
  bx2,by2:=0,0

  if config.dbg_screen {
    fmt.Println("== getPattern ==")
  }

  bx=g.pos.ax
  by=g.pos.ay
  found:="?"
  hitscore:=0
  highscore:=0
  pcnt:=0

  switch g.res {
    case SIZE_4K:
      bx2=bx+50
      by2=by+50
      font=Icon4K
    case SIZE_WQHD:
      bx2=bx+40
      by2=by+40
      font=IconWQHD
    case SIZE_1080:
      bx2=bx+25
      by2=by+25
      font=Icon1080
  }

  w:=len(font["search"][0])
  h:=len(font["search"])

  //Place character into "search" entry of font
  for y:=by;y<=by2;y++ {
    bitc=0
    bytec=0
    oct=0
    for x:=bx;x<=bx2;x++ {
      if g.At(x,y).isAbove() {
        oct|=1<<(7-bitc)
      } else {
        g.Plot(Pixel{255,255,0})
      }
      bitc++
      if bitc>7 {
        if y-by > -1 && y-by < h && bytec < w {
          font["search"][y-by][bytec]=oct
        }
        bitc=0
        bytec++
        oct=0
      }
    }
    if bitc>0 {
      if y-by > -1 && y-by < h && bytec < w {
        font["search"][y-by][bytec]=oct
      }
    }
  }

  // Dump character (used for debug/retrieving font info)
  if config.dbg_ocr {
    fmt.Println("  \"?\": {")
    for y:=0;y<h;y++ {
      fmt.Print("    {")
      for x:=0;x<w;x++ {
        fmt.Print(font["search"][y][x])
        if x<w-1 {
          fmt.Print(",")
        }
      }
      fmt.Println("},")
    }
    fmt.Println("  },")
  }



  // Get score for each of the found character
  for k, v := range font {
    if k != "search" {
      hitscore=0
      pcnt=0
      for r:=0; r<h; r++ {
        bitc=0
        bytec=0
        oct=0
        for c:=0; c<w*8; c++ {
          pcnt++
          mask:=(1<<(7-bitc))
          if v[r][bytec] & mask == font["search"][r][bytec] & mask {
            hitscore++
          } else {
            hitscore--
          }
          bitc++
          if (bitc>7) {
            bitc=0
            bytec++
            oct=0
          }
        }
      }
      // make hitscore percentage of matched pixels
      hitscore=1000*hitscore/pcnt
      //fmt.Println(" key ",k," hitscore:",hitscore)
      if hitscore>highscore {
        highscore=hitscore
        // font can have double entries for same char, denoted by trailing _
        found=strings.Replace(k,"_","",-1)
      }
    }
  }


  // Draw box around for debugging purposes
  if config.dbg_ocr {
    for x:=bx;x<bx2;x++ {
      if !g.At(x,by).isAbove() {
        g.Plot(Pixel{20,20,20})
      }
      if !g.At(x,by2).isAbove() {
        g.Plot(Pixel{20,20,20})
      }
    }
    for y:=by;y<by2;y++ {
      if !g.At(bx,y).isAbove() {
        g.Plot(Pixel{20,20,20})
      }
      if !g.At(bx2,y).isAbove() {
        g.Plot(Pixel{20,20,20})
      }
    }
  }


  //g.Save("ocr.png")
  return found,highscore
}

//------------------------------------------------------------
// Read map name and game type from TAB information screen

func (g *OWImg) Title() string {
  var line=""
  var ccnt=1
  var px=0
  var pw=0
  var space=0
  var divider=0
  var font map[string][][]int
  var bx,by,bx2,by2 int

  if config.dbg_screen {
    fmt.Println("== Title ==")
  }

  switch g.res {
    case SIZE_4K:
      bx,by=127,70
      bx2,by2=1000,105
      font=FontTitle4K
      space=10
      divider=30
    case SIZE_WQHD:
      bx,by=85,46
      bx2,by2=667,70
      font=FontTitleWQHD
      space=7
      divider=20
    case SIZE_1080:
      bx,by=64,35
      bx2,by2=500,52
      font=FontTitle1080
      space=5
      divider=15

  }
  // Turn red & yellow into white, to get all text read
  g.All()
  for x:=bx; x<bx2; x++ {
    for y:=by; y<by2 ; y++ {
      if g.At(x,y).isRed() {
        g.Plot(Pixel{255,255,245})
      }
    }
  }


  for x:=0; x<bx2-bx; x++ {
    for y:=0; y<by2-by; y++ {
      if g.From(bx,by).To(bx2,by2).At(x,y).isAbove() {
        // Did we hit more then a space ? -> divider map name and game type
        if px>0 && x-px-pw > divider {
          line+="|"
        } else {
          // Did we hit a space ?
          if px>0 && x-px-pw>space {
            line+=" "
          }
        }
        if config.dbg_ocr {
          fmt.Println(" Character # ",ccnt," at",bx+x,by+y)
        }
        ch,w,_:=g.getChar(font)
        line+=ch
        ccnt++
        px=x
        pw=w
      }
    }
  }
  //g.Save("ocr.png")
  line=strings.Replace(line,"..",":",1) // : is intepreted as two single dots

  return line
}

//------------------------------------------------------------
// Read Season High Comp SR from overview

func (g *OWImg) SRHigh() int {
  var line=""
  var ccnt=1
  var font map[string][][]int
  var bx,by,bx2,by2 int
  var total=0

  if config.dbg_screen {
    fmt.Println("== SRHigh ==")
  }

  switch g.res {
    case SIZE_4K:
      bx,by=2906,405
      bx2,by2=3026,465
      font=FontSR4K
    case SIZE_1080:
      bx,by=1453,202
      bx2,by2=1513,233
      font=FontSR1080
    case SIZE_WQHD:
      bx,by=1934,270
      bx2,by2=2018,310
      font=FontSRWQHD
  }

  for x:=0; x<bx2-bx; x++ {
    for y:=0; y<by2-by; y++ {
      if g.From(bx,by).To(bx2,by2).At(x,y).isAbove() {
        if config.dbg_ocr {
          fmt.Println(" Character # ",ccnt," at",bx+x,by+y)
        }
        ch,_,score:=g.getChar(font)
        if config.dbg_ocr {
          fmt.Println(" Score: ",score)
        }
        total+=score
        line+=ch
        ccnt++
      }
    }
  }
  total=int(total/(10*ccnt))
  if (len(line)>4) {
    line=line[len(line)-4:]
  }
  //g.Save("ocr.png")
  g.Th(-1)
  if config.dbg_ocr {
    fmt.Println(" Total Score: ",total)
  }
  ret,_:=strconv.Atoi(line)

  return ret
}

//------------------------------------------------------------
// Read current Comp SR from overview

func (g *OWImg) SRCurrent() int {
  var line=""
  var ccnt=1
  var font map[string][][]int
  var bx,by,bx2,by2 int
  var total=0

  if config.dbg_screen {
    fmt.Println("== SRCurrent==")
  }

  switch g.res {
    case SIZE_4K:
      bx,by=2906,245
      bx2,by2=3026,305
      font=FontSR4K
    case SIZE_1080:
      bx,by=1453,122
      bx2,by2=1513,153
      font=FontSR1080
    case SIZE_WQHD:
      bx,by=1934,164
      bx2,by2=2018,203
      font=FontSRWQHD
  }

  for x:=0; x<bx2-bx; x++ {
    for y:=0; y<by2-by; y++ {
      if g.From(bx,by).To(bx2,by2).At(x,y).isAbove() {
        if config.dbg_ocr {
          fmt.Println(" Character # ",ccnt," at",bx+x,by+y)
        }
        ch,_,score:=g.getChar(font)
        if config.dbg_ocr {
          fmt.Println(" Score: ",score)
        }
        total+=score
        line+=ch
        ccnt++
      }
    }
  }
  if (len(line)>4) {
    line=line[len(line)-4:]
  }
  total=int(total/(10*ccnt))
  g.Th(-1)
  if config.dbg_ocr {
    fmt.Println(" Total Score: ",total)
  }
  ret,_:=strconv.Atoi(line)

  return ret
}

//------------------------------------------------------------
// Read Comp SR from Gain/Loss screen

func (g *OWImg) SRGain() int {
  var line=""
  var ccnt=1
  var font map[string][][]int
  var bx,by,bx2,by2 int
  var total=0

  if config.dbg_screen {
    fmt.Println("== SRGain ==")
  }

  switch g.res {
    case SIZE_4K:
      bx,by=1738,915
      bx2,by2=2030,1060
      font=FontBigSR4K
    case SIZE_WQHD:
      bx,by=1143,609
      bx2,by2=1357,706
      font=FontBigSRWQHD
    case SIZE_1080:
      bx,by=859,457
      bx2,by2=1018,530
      font=FontBigSR1080
  }

  for x:=0; x<bx2-bx; x++ {
    for y:=0; y<by2-by; y++ {
      if g.From(bx,by).To(bx2,by2).At(x,y).isAbove() {
        if config.dbg_ocr {
          fmt.Println(" Character # ",ccnt," at",bx+x,by+y)
        }
        ch,_,score:=g.getChar(font)
        if config.dbg_ocr {
          fmt.Println(" Score: ",score)
        }
        total+=score
        line+=ch
        ccnt++
      }
    }
  }
  total=int(total/(10*ccnt))
  // Since SR Gain starts with an animation with "growing" numbers, 
  // it might intepret this too early, so check hitscore 
  // Certainty less then 60% ? => unreliable....
  if (total<60) {
    return 0
  }
  //g.Save("ocr.png")
  line=strings.Replace(line,".","",-1)  // : noise & artifacs
  g.Th(-1)
  if config.dbg_ocr {
    fmt.Println(" Total Score: ",total)
  }
  ret,_:=strconv.Atoi(line)

  return ret
}

//------------------------------------------------------------
// Read time TAB information screen

func (g *OWImg) TTime() string {
  var line=""
  var ccnt=1
  var font map[string][][]int
  var bx,by,bx2,by2 int

  if config.dbg_screen {
    fmt.Println("== TTime ==")
  }

  switch g.res {
    case SIZE_4K:
      bx,by=200,110
      bx2,by2=320,136
      font=FontTime4K
    case SIZE_WQHD:
      bx,by=133,73
      bx2,by2=213,91
      font=FontTimeWQHD
    case SIZE_1080:
      bx,by=100,55
      bx2,by2=160,68
      font=FontTime1080
      g.Th(100)
  }
  // Turn red & yellow into white, to get all text read
  g.All()
  for x:=bx; x<bx2; x++ {
    for y:=by; y<by2 ; y++ {
      if g.At(x,y).isRed() {
        g.Plot(Pixel{255,255,245})
      }
    }
  }


  for x:=0; x<bx2-bx; x++ {
    for y:=0; y<by2-by; y++ {
      if g.From(bx,by).To(bx2,by2).At(x,y).isAbove() {
        if config.dbg_ocr {
          fmt.Println(" Character # ",ccnt," at",bx+x,by+y)
        }
        ch,_,_:=g.getChar(font)
        line+=ch
        ccnt++
      }
    }
  }
  //g.Save("ocr.png")
  line=strings.Replace(line,"..",":",1) // : is intepreted as two single dots
  line=strings.Replace(line,".","",-1)  // : noise
  if !strings.Contains(line,":") || len(line)<4 {
    // unreliable information, should have a ":" and at least x:xx format
    // blank out to prevent any other mishap
    line=""
  }

  g.Th(-1)

  return line
}

//------------------------------------------------------------
// Get statistic line from TAB screen, below, counting from upper left medal

func (g *OWImg) TStat(col,row int) string {
  var line=""
  var ccnt=1
  var font map[string][][]int
  var cols_l []int
  var cols_r []int
  var rows_u []int
  var rows_b []int
  var bx,by,bx2,by2 int

  if config.dbg_screen {
    fmt.Println("== TStat ==")
  }

  switch g.res {
    case SIZE_4K:
      font=FontStats4K
      cols_l= []int{255,755,1255,2060,2600,3140}
      cols_r= []int{500,1000,1500,2360,3000,3540}
      rows_u= []int{1780,1830}
      rows_b= []int{1910,1960}
    case SIZE_WQHD:
      font=FontStatsWQHD
      cols_l= []int{172,505,838,1373,1733,2093}
      cols_r= []int{333,665,997,1569,1995,2301}
      rows_u= []int{1193,1219}
      rows_b= []int{1277,1303}
    case SIZE_1080:
      font=FontStats1080
      cols_l= []int{128,378,627,1030,1300,1570}
      cols_r= []int{250,500,750,1180,1500,1770}
      rows_u= []int{890,920}
      rows_b= []int{955,985}
  }
  bx=cols_l[col]
  by=rows_u[row]
  bx2=cols_r[col]
  by2=rows_b[row]

  if config.dbg_ocr {
    fmt.Println(" Stats for ",col,row)
  }

  for x:=0; x<bx2-bx; x++ {
    for y:=0; y<by2-by; y++ {
      if g.From(bx,by).To(bx2,by2).At(x,y).isAbove() {
        if config.dbg_ocr {
          fmt.Println(" Character # ",ccnt," at",bx+x,by+y)
        }
        ch,_,_:=g.getChar(font)
        line+=ch
        ccnt++
      }
    }
  }
  //g.Save("ocr.png")
  line=strings.Replace(line,",,",":",1)  // : is intepreted as two single dots
  line=strings.Replace(line,",","",1)    // : 9,999 => 9999
  line=strings.Replace(line,"%%%","%",1) // : percentage is read as 3 chars
  line=strings.Replace(line,"..",":",1)  // : is intepreted as two single dots
  line=strings.Replace(line,".","",-1)   // : noise
  g.Th(-1)

  return line
}

//------------------------------------------------------------
// Read name of hero choosen by yourself

func (g *OWImg) MyHero() string {
  var line=""
  var font map[string][][]int
  var bx,by,bx2,by2 int
  var ccnt=0
  g.Th(224)

  if config.dbg_screen {
    fmt.Println("== MyHero ==")
  }

  switch g.res {
    case SIZE_4K:
      bx,by=1900,1660
      bx2,by2=2400,1750
      font=FontHero4K
    case SIZE_WQHD:
      bx,by=1270,1110
      bx2,by2=1570,1160
      font=FontHeroWQHD
      g.Th(190)
    case SIZE_1080:
      bx,by=750,830
      bx2,by2=1200,875
      font=FontHero1080
  }

  for x:=0; x<bx2-bx; x++ {
    for y:=0; y<by2-by; y++ {
      if g.From(bx,by).To(bx2,by2).At(x,y).isAbove() {
        if config.dbg_ocr {
          fmt.Println(" Character # ",ccnt," at",bx+x,by+y)
        }
        ch,_,_:=g.getChar(font)
        line+=ch
        ccnt++
      }
    }
  }
  g.Th(-1)

  // Quick fix to found name and name used internally
  switch line {
    case "ANA":
      line="Ana"

    case "BASIION":
      line="Bastion"

    case "BASTION":
      line="Bastion"

    case "BRIGITTE":
      line="Brigitte"

    case "BRIGITFE":
      line="Brigitte"

    case "BRIGIFFE":
      line="Brigitte"

    case "BRIGIFTE":
      line="Brigitte"

    case "DOOMFIST":
      line="Doomfist"

    case "D.VA":
      line="D.Va"

    case "GENJI":
      line="Genji"

    case "HANZO":
      line="Hanzo"

    case "JUNKRAT":
      line="Junkrat"

    case "JUNKRRAT":
      line="Junkrat"

    case "LU.CIO":
      line="Lucio"

    case "MGCREE":
      line="McCree"

    case "MCCREE":
      line="McCree"

    case "MEI":
      line="Mei"

    case "MERCY":
      line="Mercy"

    case "MOIRA":
      line="Moira"

    case "ORI6A":
      line="Orisa"

    case "ORISA":
      line="Orisa"

    case "PHARAH":
      line="Pharah"

    case "REAPER":
      line="Reaper"

    case "REKHHARDI":
      line="Reinhardt"

    case "REINHARDT":
      line="Reinhardt"

    case "ROAOHOG":
      line="Roadhog"

    case "ROADHOG":
      line="Roadhog"

    case "SOLDIER..76":
      line="Soldier 76"

    case "SOMBRA":
      line="Sombra"

    case "SYMMETRA":
      line="Symmetra"

    case "TOR6JO..RN":
      line="Torbjörn"

    case "TORBJO..RN":
      line="Torbjörn"

    case "TORBJORN":
      line="Torbjörn"

    case "TRACER":
      line="Tracer"

    case "WIOOWMAIRER":
      line="Widowmaker"

    case "WIDOWMAIIER":
      line="Widowmaker"

    case "WIDOWMAKER":
      line="Widowmaker"

    case "WIDUWMAKER":
      line="Widowmaker"

    case "WINSTON":
      line="Winston"

    case "WRECKINGBALL":
      line="Wrecking Ball"

    case "WRECKINGBAIL":
      line="Wrecking Ball"

    case "ZARYA":
      line="Zarya"

    case "ZENYAIIA":
      line="Zenyatta"

    case "ZENYATTA":
      line="Zenyatta"
  }

  if config.dbg_ocr {
    fmt.Println("got heroname:",line)
  }
  //g.Save("ocr.png")

  return line
}
