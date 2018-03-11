package main

import (
  "fmt"
  "image"
  "image/png"
  "os"
  "github.com/mertyGit/owig/owocr"
)

// ----------------------------------------------------------------------------
// determine positive difference between two ints

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
      if s2.R != s1.R || s2.G != s1.G || s2.B != s1.B {
        same=false
      }
    }
  } else {
    if yfrom==yto {
      for x:=xfrom;x<xto;x++ {
        s2=pix(x,yfrom)
        if s2.R != s1.R || s2.G != s1.G || s2.B != s1.B {
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
// Do pixel r,g,b matches given rgb (divided by div)

func like(s1 Pixel,r int,g int,b int, div int) bool {
  if div<1 {
    div=1
  }
  if int(s1.R/div) == r && int(s1.G/div) == g && int(s1.B/div) == b  {
    return true
  }
  return false
}

// ----------------------------------------------------------------------------
// How many "holes" do we encounter if we go from xt to xb, treshold tr

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

