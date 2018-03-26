package main

import (
  "fmt"
  "image"
  "image/png"
  "image/color"
  "os"
  "github.com/mertyGit/owig/screenshot"
)


type Pixel struct {
  R int
  G int
  B int
}

type Pos struct {
  x int
  y int
  x2 int
  y2 int
  ax int
  ay int
}

type OWImg struct {
  img *image.RGBA // Image itself
  pix Pixel       // RGB pixel values of latest At or From (if no At present)
  pos Pos         // Positions of From (x,y), To (x2,y2) and At (ax,ay)
  width int       // Width of image
  height int      // Height of image
  res int         // Resolution type of image (4K or 1080 or nothing)
  th int          // Thresshold for recognizing pixels as "on" or "off"
  fc int          // Fill Count, amount of pixels in flood fill
  gotimg bool     // Flag to indicate we have loaded image
  gotfrom bool    // Flag to indicate we have "From" coordinates
  gotto bool      // Flag to indicate we have "To"   coordinates
  gotat bool      // Flag to indicate we have "At"   coordinates
}

var owig *OWImg

// ----------------------------------------------------------------------------
// Initialize parameters
func (g *OWImg) Init() {
  g.img=nil
  g.pix=Pixel{0,0,0}
  g.pos.x=0
  g.pos.y=0
  g.pos.x2=0
  g.pos.y2=0
  g.pos.ax=0
  g.pos.ay=0
  g.width=0
  g.height=0
  g.th=150
  g.fc=0
  g.res=SIZE_NONE
  g.gotimg=false
  g.gotfrom=false
  g.gotto=false
  g.gotat=false
}


// ----------------------------------------------------------------------------
// Set resolution (called by FromFile or FromScreen)

func (g *OWImg) setRes() *OWImg {
  if !g.gotimg {
    return g
  }
  if g.width==3840 && g.height==2160 {
    g.res=SIZE_4K
  } else if g.width==1920 && g.height==1080 {
    g.res=SIZE_1080
  } else {
    if config.dbg_screen {
      fmt.Println("WARNING: Unknown screen size ",g.width," x",g.height)
    }
  }
  return g
}


// ----------------------------------------------------------------------------
// Load image using given filename (must be .png file)

func (g *OWImg) Open(name string) *OWImg {

  g.Init()

  file, err := os.Open(name)
  if err != nil {
    fmt.Println("Warning: File could not be opened")
    return g
  }

  defer file.Close()
  myimg, err := png.Decode(file)

  if err != nil {
    fmt.Println("Warning: PNG could not be decoded")
    return g
  }
  g.img=myimg.(*image.RGBA)
  g.width=g.img.Bounds().Max.X
  g.height=g.img.Bounds().Max.Y
  g.gotimg=true
  g.setRes()
  if config.dbg_screen {
    fmt.Println("Image loaded from file :",name)
    fmt.Println("Size :",g.width,"x ",g.height,"(",g.res,")")
  }
  return g
}

// ----------------------------------------------------------------------------
// Load image from screen capture
func (g *OWImg) Capture() *OWImg {
  var err error

  g.Init()

  g.img, err = screenshot.CaptureScreen()
  if err != nil {
    fmt.Println("Error: Can't make screenshot")
    return g
  }
  g.width=g.img.Bounds().Max.X
  g.height=g.img.Bounds().Max.Y
  g.gotimg=true
  g.setRes()

  return g
}

// ----------------------------------------------------------------------------
// Create new image given width and height
func (g *OWImg) Create(w,h int) *OWImg {

  g.Init()

  g.img=image.NewRGBA(image.Rect(0,0,w,h))
  g.width=w
  g.height=h
  g.gotimg=true

  g.Fill(Pixel{0,0,0})

  return g
}



// ----------------------------------------------------------------------------
// Save image
func (g *OWImg) Save(name string) *OWImg {
  df,err := os.OpenFile(name,os.O_WRONLY|os.O_CREATE, 0777)
  if err!=nil {
    fmt.Println("Error: Image could not be saved")
    os.Exit(2)
  }
  defer df.Close()
  err=png.Encode(df,g.img)
  if err!=nil {
    fmt.Println("Error: Image could not be encoded:",err)
    os.Exit(3)
  }
  return g
}



// ----------------------------------------------------------------------------
// Sets coordinates x,y at given point and get RGB info
// Used for point requests or "from" or "upper-left box" coordinates

func (g *OWImg) From(x,y int) *OWImg {

  if !g.gotimg {
    return g
  }

  if x<0 {
    x=0
  }
  if y<0 {
    y=0
  }
  g.pos.x=x
  g.pos.y=y
  g.gotfrom=true
  if g.pos.x>g.width {
    g.pos.x=g.width
  }
  if g.pos.y>g.height {
    g.pos.y=g.height
  }
  R,G,B,_ := g.img.At(x,y).RGBA()
  g.pix.R=int(R/257)
  g.pix.G=int(G/257)
  g.pix.B=int(B/257)
  g.gotat=false
  return g
}

// ----------------------------------------------------------------------------
// Sets coordinates x2,y2 at given point and get RGB info
// Used for actions that do require second point (lines, boxes, areas)

func (g *OWImg) To(x,y int) *OWImg {

  if !g.gotimg {
    return g
  }

  if x<0 {
    x=0
  }
  if y<0 {
    y=0
  }
  g.pos.x2=x
  g.pos.y2=y
  g.gotto=true
  if g.gotfrom {
    if g.pos.x2<g.pos.x {
      g.pos.x2=g.pos.x
    }
    if g.pos.y2<g.pos.y {
      g.pos.y2=g.pos.y
    }
  } else {
    g.pos.x=0
    g.pos.y=0
  }
  if g.pos.x2>g.width {
    g.pos.x2=g.width
  }
  if g.pos.y2>g.height {
    g.pos.y2=g.height
  }
  g.gotat=false
  return g
}

// ----------------------------------------------------------------------------
// Shortcut to define Bounding Box
func (g *OWImg) Box(x,y,w,h int) *OWImg {

  if !g.gotimg {
    return g
  }

  g.From(x,y).To(x+w,y+h)
  return g
}

// ----------------------------------------------------------------------------
// Shortcut to define whole area
func (g *OWImg) All() *OWImg {

  if !g.gotimg {
    return g
  }

  g.Box(0,0,g.width,g.height)
  return g
}

// ----------------------------------------------------------------------------
// Make following function realative to "From" (and not further then "To")
// If set... otherwise relative to left upper corner

func (g *OWImg) At(x,y int) *OWImg{

  if !g.gotimg {
    return g
  }

  if !g.gotfrom {
    g.pos.x=0
    g.pos.y=0
  }
  if g.gotto {
    if x+g.pos.x>g.pos.x2 {
      x=g.pos.x2-g.pos.x
    }
    if y+g.pos.y>g.pos.y2 {
      y=g.pos.y2-g.pos.y
    }
  } else {
    if x+g.pos.x>g.width {
      x=g.pos.x+g.width
    }
    if y+g.pos.y>g.height {
      y=g.pos.y+g.height
    }
  }
  R,G,B,_ := g.img.At(g.pos.x+x,g.pos.y+y).RGBA()
  g.pix.R=int(R/257)
  g.pix.G=int(G/257)
  g.pix.B=int(B/257)
  g.pos.ax=x
  g.pos.ay=y
  g.gotat=true
  return g
}


// ----------------------------------------------------------------------------
// Set thresshold, negative number resets to default (150)
func (g *OWImg) Th(t int) *OWImg {

  if t>0 && t<256 {
    g.th=t
  } else {
    g.th=150
  }

  return g
}

// ----------------------------------------------------------------------------
// Get Pixel at given point

func (g *OWImg) RGB() Pixel {

  if !g.gotimg {
    return Pixel{0,0,0}
  }

  return g.pix
}

// ----------------------------------------------------------------------------
// Get Red value of pixel at given point 

func (g *OWImg) Red() int {

  if !g.gotimg {
    return 0
  }

  return g.pix.R
}

// ----------------------------------------------------------------------------
// Check if Red value is above thresshold
func (g *OWImg) isRed() bool {

  if !g.gotimg || !(g.gotfrom || g.gotat) {
    return false
  }

  if g.pix.R>g.th {
    return true
  }
  return false
}

// ----------------------------------------------------------------------------
// Get Green value of pixel at given point 

func (g *OWImg) Green() int {

  if !g.gotimg {
    return 0
  }

  return g.pix.G
}

// ----------------------------------------------------------------------------
// Check if Green value is above thresshold
func (g *OWImg) isGreen() bool {

  if !g.gotimg || !(g.gotfrom || g.gotat) {
    return false
  }

  if g.pix.G>g.th {
    return true
  }
  return false
}


// ----------------------------------------------------------------------------
// Get Blue value of pixel at given point 

func (g *OWImg) Blue() int {

  if !g.gotimg || !(g.gotfrom || g.gotat) {
    return 0
  }

  return g.pix.B
}


// ----------------------------------------------------------------------------
// Check if Blue value is above thresshold
func (g *OWImg) isBlue() bool {

  if !g.gotimg || !(g.gotfrom || g.gotat) {
    return false
  }

  if g.pix.B>g.th {
    return true
  }
  return false
}

// ----------------------------------------------------------------------------
// Check if Yellow value is above thresshold
func (g *OWImg) isYellow() bool {

  if !g.gotimg || !(g.gotfrom || g.gotat) {
    return false
  }

  if g.isRed() && g.isGreen() {
    return true
  }
  return false
}


// ----------------------------------------------------------------------------
// Check if R,G,B are above thresshold

func (g *OWImg) isAbove() bool {

  if !g.gotimg || !(g.gotfrom || g.gotat) {
    return false
  }

  if g.pix.R>g.th && g.pix.G>g.th && g.pix.B>g.th {
    return true
  }
  return false
}

// ----------------------------------------------------------------------------
// Check if R,G,B are below thresshold

func (g *OWImg) isBelow() bool {

  if !g.gotimg || !(g.gotfrom || g.gotat) {
    return false
  }

  if g.pix.R<g.th && g.pix.G<g.th && g.pix.B<g.th {
    return true
  }
  return false
}


// ----------------------------------------------------------------------------
// Check to see if point is white

func (g *OWImg) isWhite() bool {

  if !g.gotimg || !(g.gotfrom || g.gotat) {
    return false
  }

  if g.pix.R>244 && g.pix.G>244 && g.pix.B>244 {
    return true
  }
  return false
}


// ----------------------------------------------------------------------------
// Check to see if there is a straight line with same color 
// (horizontal or vertical)

func (g *OWImg) isLine() bool {
  var same bool


  if !g.gotimg {
    return false
  }

  same=true
  if !g.gotto || !g.gotfrom {
    // line without start or end ? So just a point...
    return true
  }
  if g.pos.x==g.pos.x2 {
    for y:=0;y+g.pos.y<g.pos.y2;y++ {
      px := g.At(0,y).RGB()
      if px.R != g.pix.R || px.G != g.pix.G || px.B != g.pix.B {
        same=false
      }
    }
  } else {
    if g.pos.y==g.pos.y2 {
      for x:=0;x+g.pos.x<g.pos.x2;x++ {
        px := g.At(x,0).RGB()
        if px.R != g.pix.R || px.G != g.pix.G || px.B != g.pix.B {
          same=false
        }
      }
    } else {
      // not supported, diagonal lines
      return false
    }
  }
  return same
}


// ----------------------------------------------------------------------------
// Check to see if pixel matches given pixels values , +/- given range

func (g *OWImg) isLike(p Pixel,dev int) bool {

  if !g.gotimg || !(g.gotfrom || g.gotat) {
    return false
  }

  if p.R > g.pix.R+dev || p.R < g.pix.R-dev || p.G > g.pix.G+dev || p.G < g.pix.G-dev || p.B > g.pix.B+dev || p.B < g.pix.B-dev {
    return false
  }
  return true
}

// ----------------------------------------------------------------------------
// Count "holes" (amount wher line goes from below to above thresshold) 
// in given horizontal line

func (g *OWImg) Holes() int {

  if !g.gotimg || !g.gotfrom || !g.gotto  {
    return 0
  }

  hole:=false
  cnt:=0
  for x:=0;x+g.pos.x<g.pos.x2;x++ {
    if g.At(x,0).isAbove() {
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
// generate checksum based on average R,G,B color values in whole area
// Used to check or validate if some menu items are on screen or not

func (g *OWImg) Cs() string {
  var R,G,B uint64

  if !g.gotimg || !g.gotfrom || !g.gotto  {
    return ""
  }

  R=0
  G=0
  B=0
  for x:=0;x+g.pos.x<g.pos.x2;x++ {
    for y:=0;y+g.pos.y<g.pos.y2;y++ {
      R+=uint64(g.At(x,y).pix.R)
      G+=uint64(g.At(x,y).pix.G)
      B+=uint64(g.At(x,y).pix.B)
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
// Set pixel value
func (g *OWImg) Plot(pix Pixel) *OWImg {
  if !g.gotimg || !(g.gotfrom || g.gotat){
    return g
  }
  if g.gotat {
    g.img.Set(g.pos.x+g.pos.ax,g.pos.y+g.pos.ay,color.RGBA{uint8(pix.R),uint8(pix.G),uint8(pix.B),255})
  } else {
    g.img.Set(g.pos.x,g.pos.y,color.RGBA{uint8(pix.R),uint8(pix.G),uint8(pix.B),255})
  }
  g.pix=pix
  return g
}

// ----------------------------------------------------------------------------
// Draw Horizontal or Vertical Line with Pixel color

func (g *OWImg) DrawLine(pix Pixel) *OWImg {


  if !g.gotimg {
    return g
  }

  if !g.gotto || !g.gotfrom {
    // line without start or end ? So just a point...
    return g
  }
  if g.pos.x==g.pos.x2 {
    for y:=0;y+g.pos.y<=g.pos.y2;y++ {
      g.At(0,y).Plot(pix)
    }
  } else {
    if g.pos.y==g.pos.y2 {
      for x:=0;x+g.pos.x<=g.pos.x2;x++ {
        g.At(x,0).Plot(pix)
      }
    }
  }
  return g
}

// ----------------------------------------------------------------------------
// Draw Box with Pixel color

func (g *OWImg) DrawBox(pix Pixel) *OWImg {


  if !g.gotimg {
    return g
  }

  if !g.gotto || !g.gotfrom {
    // line without start or end ? So just a point...
    return g
  }
  for y:=0;y+g.pos.y<=g.pos.y2;y++ {
    g.At(0,y).Plot(pix)
    g.At(g.pos.x2-g.pos.x,y).Plot(pix)
  }
  for x:=0;x+g.pos.x<=g.pos.x2;x++ {
    g.At(x,0).Plot(pix)
    g.At(x,g.pos.y2-g.pos.y).Plot(pix)
  }
  return g
}


// ----------------------------------------------------------------------------
// Fill area with color
func (g *OWImg) Fill(pix Pixel) *OWImg {

  var w=0
  var h=0


  if !g.gotimg {
    return g
  }

  if g.gotto {
    w=g.pos.x2-g.pos.x
    h=g.pos.y2-g.pos.y
  } else {
    w=g.width
    h=g.height
  }
  for x:=0;x<w;x++ {
    for y:=0;y<h;y++ {
      g.At(x,y).Plot(pix)
    }
  }
  return g
}

// ----------------------------------------------------------------------------
// Flood Fill area with color (fills all adjecent pixels > thresshold)
func (g *OWImg) Flood(pix Pixel) *OWImg {
  if !g.gotimg {
    return g
  }
  g.fc=0
  if g.gotat {
    g.ff(g.pos.ax,g.pos.ay,pix)
  } else {
    g.ff(0,0,pix)
  }
  return g
}

// ----------------------------------------------------------------------------
// recursive function to do floodfill
func (g *OWImg) ff(x,y int,p Pixel) {
  if x<0 || y<0 || x>g.width+g.pos.x || y>g.height+g.pos.y {
    return
  }
  if g.At(x,y).isLike(p,0) || !g.At(x,y).isAbove() {
    return
  }
  //fmt.Println("color at ",x,y,g.RGB())
  g.At(x,y).Plot(p)
  g.fc++
  g.ff(x-1,y-1,p)
  g.ff(x,y-1,p)
  g.ff(x+1,y-1,p)
  g.ff(x+1,y,p)
  g.ff(x+1,y+1,p)
  g.ff(x,y+1,p)
  g.ff(x-1,y+1,p)
  g.ff(x-1,y,p)
}


// ----------------------------------------------------------------------------
// Black out artifacts and colors below thresshold
func (g *OWImg) Filter() *OWImg {

  if !g.gotimg || !g.gotfrom || !g.gotto  {
    return g
  }

  for x:=0;x+g.pos.x<=g.pos.x2;x++ {
    for y:=0;y+g.pos.y<=g.pos.y2;y++ {
      if g.At(x,y).isBelow() {
        g.Plot(Pixel{0,0,0})
      }
    }
  }
  return g
}

// ----------------------------------------------------------------------------
// Check to see if area has same color value (within max deviation)
func (g *OWImg) SameColor(p Pixel,dev int) bool {
  var ret=true

  if !g.gotimg || !g.gotfrom || !g.gotto  {
    return false
  }

  for x:=0;x+g.pos.x<=g.pos.x2;x++ {
    for y:=0;y+g.pos.y<=g.pos.y2;y++ {
      if g.At(x,y).isLike(p,dev) {
        ret=false
      }
    }
  }
  return ret
}

// ----------------------------------------------------------------------------
// Check to see if area has same color, based on which RGB values are the
// highest R=Red, G=Green, Y=yellow, B=blue, C=cyan, M=magenta
func (g *OWImg) SameBase(c string) bool {
  same:=true

  if !g.gotimg || !g.gotfrom || !g.gotto  {
    return false
  }

  for x:=0;x+g.pos.x<g.pos.x2;x++ {
    for y:=0;y+g.pos.y<g.pos.y2;y++ {
      P:=g.At(x,y).RGB()
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
