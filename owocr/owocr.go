package owocr

import (
  "github.com/mertyGit/owig/owfonts"
  "fmt"
  "image"
  "image/png"
  "image/color"
  "os"
  "strings"
  "strconv"
)


type Pixel struct {
  R int
  G int
  B int
}

var imgIn *image.RGBA
var imgOut *image.RGBA
var width,height int
var thresshold=150
var filled=0

func pixIn(x int, y int) Pixel {
  if x>width {
    x=width
  }
  if y>height {
    y=height
  }
  r,g,b,_ := imgIn.At(x, y).RGBA()
  return Pixel{int(r/257),int(g/257),int(b/257)}
}

func isPixIn(x int,y int) bool {
  if pixIn(x,y).R>thresshold {
    return true
  }
  return false
}

func pixOut(x int, y int) Pixel {
  if x>width {
    x=width
  }
  if y>height {
    y=height
  }
  r,g,b,_ := imgOut.At(x, y).RGBA()
  return Pixel{int(r/257),int(g/257),int(b/257)}
}

func isPixOut(x int,y int) bool {
  if pixOut(x,y).G>thresshold || pixOut(x,y).B>thresshold {
    return true
  }
  return false
}

func isBlue(x int,y int) bool {
  if pixOut(x,y).B>thresshold  {
    return true
  }
  return false
}

func isYellow(x int,y int) bool {
  if pixOut(x,y).G>thresshold && pixOut(x,y).R>thresshold  {
    return true
  }
  return false
}

// Copy - part of screenshot- to imgIn
func setImg(myimg *image.RGBA,xt int,yt int,xb int,yb int) {
  width=xb-xt
  height=yb-yt
  imgIn=image.NewRGBA(image.Rect(0,0,width,height))
  imgOut=image.NewRGBA(image.Rect(0,0,width,height))
  for x := xt; x < xb; x++ {
    for y := yt; y < yb; y++ {
      r,g,b,_ := myimg.At(x, y).RGBA()
      imgIn.Set(x-xt,y-yt,color.RGBA{uint8(r/257),uint8(g/257),uint8(b/257),255})
    }
  }
}

// Helper / Debug functions to check if right pixels are choosen 
func Plot(img *image.RGBA,x int,y int,r uint8, g uint8, b uint8) {
  img.Set(x,y,color.RGBA{r,g,b,255})
}


// Load imgIn from file
func LoadFile(name string) {
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
  imgIn=myimg.(*image.RGBA)
  imgOut=image.NewRGBA(image.Rect(0,0,imgIn.Bounds().Max.X,imgIn.Bounds().Max.Y))
  width=imgIn.Bounds().Max.X
  height=imgIn.Bounds().Max.Y
}

// Save imgOut to file
func SaveFile(name string) {
  SaveImg(imgOut,name)
}

func SaveImg(img *image.RGBA,name string) {
  df,err := os.OpenFile(name,os.O_WRONLY|os.O_CREATE, 0777)
  if err!=nil {
    fmt.Println("Error: Image could not be saved")
    os.Exit(2)
  }
  defer df.Close()
  err=png.Encode(df,img)
  if err!=nil {
    fmt.Println("Error: Image could not be encoded:",err)
    os.Exit(3)
  }
}


func blackOut() {
  for x := 0; x < width; x++ {
    for y := 0; y < height; y++ {
      imgOut.Set(x,y,color.RGBA{0,0,0,255})
    }
  }
}

func floodFill(x int,y int) {
  if x >= width || y >= height || x < 0 || y < 0 {
    return
  }
  if isPixOut(x,y) || !isPixIn(x,y) {
    return
  }
  imgOut.Set(x,y,color.RGBA{255,255,0,255})
  filled++
  floodFill(x-1,y-1)
  floodFill(x,y-1)
  floodFill(x+1,y-1)
  floodFill(x+1,y)
  floodFill(x+1,y+1)
  floodFill(x,y+1)
  floodFill(x-1,y+1)
  floodFill(x-1,y)
}

// dump character information to golang-readable code
// input rectangle (top-left and bottom-right coordinates) + width x height

func dumpChar(xl int,yt int,xr int,yb int,font map[string][][]int ) {
  var bytec=0
  var bitc uint8
  var oct=0

  w:=len(font["search"][0])*8
  h:=len(font["search"])

  fmt.Println("  \"?\": {")
  for yi:=0; yi<=h; yi++ {
    bitc=0
    bytec=0
    oct=0
    fmt.Print("    {")
    for xi:=0; xi<w; xi++ {
      if xi+xl<=xr && yi+yt<=yb {
        if isPixOut(xi+xl,yi+yt)&&isBlue(xi+xl,yi+yt) {
          oct|=1<<(7-bitc)
        }
      }
      bitc++
      if (bitc>7) {
        bitc=0
        bytec++
        fmt.Print(oct)
        if (bytec<w/8) {
          fmt.Print(",")
        }
        oct=0
      }
    }
    fmt.Println("},")
  }
  fmt.Println("  },")
}


func searchChar(xl int,yt int,xr int,yb int,font map[string][][]int) string {
  var bytec=0
  var bitc uint8
  var oct=0
  var found="?"
  var hitscore=0
  var highscore=0

  w:=len(font["search"][0])*8
  h:=len(font["search"])

  // use font "search" entry to fill in found character
  for yi:=0; yi<h; yi++ {
    bitc=0
    bytec=0
    oct=0
    for xi:=0; xi<w; xi++ {
      if xi+xl<=xr && yi+yt<=yb {
        if (isPixOut(xi+xl,yi+yt))&&(isBlue(xi+xl,yi+yt)) {
          oct|=1<<(7-bitc)
          // color pixel green, preventing re-scan next character
          imgOut.Set(xi+xl,yi+yt,color.RGBA{0,255,0,255})
        }
      }
      bitc++
      if (bitc>7) {
        font["search"][yi][bytec]=oct
        bitc=0
        bytec++
        oct=0
      }
    }
  }

  //now iterate to whole font, to find best matching char
  for k, v := range font {
    if k != "search" {
      //fmt.Println("checking k",k)
      hitscore=0
      for r:=0; r<h; r++ {
        bitc=0
        bytec=0
        oct=0
        for c:=0; c<w; c++ {
          if v[r][bytec]&(1<<(7-bitc))>0 {
            if font["search"][r][bytec]&(1<<(7-bitc)) >0 {
              hitscore+=2
            } else {
              hitscore--
            }
          }
          bitc++
          if (bitc>7) {
            bitc=0
            bytec++
            oct=0
          }
        }
      }
      if hitscore>highscore {
        //fmt.Println("Got hitscore:",hitscore)
        highscore=hitscore
        found=k[:1]
      }
    }
  }
  return found
}




func getChar(xleft int,yleft int,font map[string][][]int) (string,int) {
  var cleanCol=false

  xl:=xleft
  xr:=0
  yt:=height
  yb:=0
  ch:="?"

  // Find all neighbouring pixels and color them yellow
  filled=0
  floodFill(xleft,yleft)
  if filled<5 {// just picked up noise, ignore
    return "",0
  }

  // Get borders of character, and color it blue
  for x := xl; x < width && !cleanCol; x++ {
    cleanCol=true
    for y := 0; y < height; y++ {
      if isPixOut(x,y)&&(isYellow(x,y)) {
        cleanCol=false
        imgOut.Set(x,y,color.RGBA{0,0,255,255})
        if xr<x {
          xr=x
        }
        if yb<y {
          yb=y
        }
        if yt>y {
          yt=y
        }
      }
    }
  }

  // Output image to includable format, based on array size of "search" in font 
  //dumpChar(xl,yt,xr,yb,font)

  // check character against font
  ch=searchChar(xl,yt,xr,yb,font)


  // Draw found borders (only to check in debug)
  for xi:=xl; xi<xr; xi++ {
    if (!isPixOut(xi,yt)) {
      imgOut.Set(xi,yt,color.RGBA{uint8(thresshold-1),uint8(thresshold-1),uint8(thresshold-1),255})
    }
    if (!isPixOut(xi,yb)) {
      imgOut.Set(xi,yb,color.RGBA{uint8(thresshold-1),uint8(thresshold-1),uint8(thresshold-1),255})
    }
  }
  for yi:=yt; yi<yb; yi++ {
    if (!isPixOut(xl,yi)) {
      imgOut.Set(xl,yi,color.RGBA{uint8(thresshold-1),uint8(thresshold-1),uint8(thresshold-1),255})
    }
    if (!isPixOut(xr,yi)) {
      imgOut.Set(xr,yi,color.RGBA{uint8(thresshold-1),uint8(thresshold-1),uint8(thresshold-1),255})
    }
  }
  return ch,xr-xl
}


func Img2Title(img *image.RGBA,res int) string {
  switch res {
    case 0:
      setImg(img,128,70,1000,105)
    case 1:
      setImg(img,59,35,500,53)
  }
  line:=ReadTitle(res)
  //SaveFile("test.png")
  return line
}


func Img2CurrentSR(img *image.RGBA,res int) int {
  switch res {
    case 0:
      setImg(img,3155,245,3283,305)
    case 1:
      setImg(img,1577,122,1641,153)
  }
  line:=ReadSR(res,0)
  //SaveFile("test.png")
  i,_:=strconv.Atoi(line)
  return i
}

func Img2HighSR(img *image.RGBA,res int) int {
  switch res {
    case 0:
      setImg(img,3155,405,3283,465)
    case 1:
      setImg(img,1577,202,1641,233)
  }
  line:=ReadSR(res,0)
  //SaveFile("test.png")
  i,_:=strconv.Atoi(line)
  return i
}

func Img2CompSR(img *image.RGBA,res int) int {
  switch res {
    case 0:
      setImg(img,1738,915,2030,1060)
    case 1:
      setImg(img,859,457,1018,530)
  }
  line:=ReadSR(res,1)
  //SaveFile("test.png")
  i,_:=strconv.Atoi(line)
  return i
}

// Try to decipher stats digits
// (black part on bottom, every time you press tab during game )
// 
// Resolution type
// 0=3840x2160 (4K)
// 1=1920x1080 (1080p)

func GetStats(img *image.RGBA, col int, row int, res int) string{
  var line=""
  var ccnt=1
  var font map[string][][]int
  var cols_l []int
  var cols_r []int
  var rows_u []int
  var rows_b []int


  switch res {
    case 0:
      cols_l= []int{255,755,1255,2060,2600,3140}
      cols_r= []int{500,1000,1500,2360,3000,3540}
      rows_u= []int{1780,1830}
      rows_b= []int{1910,1960}
      font=owfonts.FontStats4K
    case 1:
      cols_l= []int{128,378,627,1030,1300,1570}
      cols_r= []int{250,500,750,1180,1500,1770}
      rows_u= []int{890,915}
      rows_b= []int{955,980}
      font=owfonts.FontStats1080
  }

  setImg(img,cols_l[col],rows_u[row],cols_r[col],rows_b[row])
  blackOut()

  for xi := 0; xi < width; xi++ {
    for yi := 0; yi < height; yi++ {
      if isPixIn(xi,yi) && !isPixOut(xi,yi) {
        //fmt.Println("Reading char #: ",ccnt)

        //Intepret character from TitleFont
        ch,_:=getChar(xi,yi,font)
        line+=ch
        ccnt++
      }
    }
  }
  line=strings.Replace(line,",","",1) // : 9,999 => 9999 (easier to convert to number)
  line=strings.Replace(line,"%%%","%",1) // : is intepreted as three seperate characters
  line=strings.Replace(line,"..",":",1) // : is intepreted as two single dots
  return line
}


// Try to decipher SR font
// every time you press tab during game )
// 
// Resolution type
// 0=3840x2160 (4K)
// 1=1920x1080 (1080p)
// 
// size=0 => SR size for text in top right of Overview screen
// size=1 => SR size for text in middle for Gain/Loss screen


func ReadSR(res int,size int) string{
  var line=""
  var ccnt=1
  var font map[string][][]int

  switch res {
    case 0:
      if (size==0) {
        font=owfonts.FontSR4K
      } else {
        font=owfonts.FontBigSR4K
      }
    case 1:
      if (size==0) {
        font=owfonts.FontSR1080
      } else {
        font=owfonts.FontBigSR1080
      }
  }

  blackOut()
  for xi := 0; xi < width; xi++ {
    for yi := 0; yi < height; yi++ {
      if isPixIn(xi,yi) && !isPixOut(xi,yi) {
        //fmt.Println("Reading char #: ",ccnt)
        //Intepret character from SRFont
        ch,_:=getChar(xi,yi,font)
        //fmt.Println(">> ",ch)
        line+=ch
        ccnt++
      }
    }
  }
  line=strings.Replace(line,".","",-1) // get rid of artifacts
  line=strings.Replace(line,"?","",-1) // get rid of artifacts
  //fmt.Println("Got line:",line)
  return line
}

// Try to decipher titel and game type (displayed on left upper corner
// every time you press tab during game )
// 
// Resolution type
// 0=3840x2160 (4K)
// 1=1920x1080 (1080p)


func ReadTitle(res int) string{
  var line=""
  var ccnt=1
  var px=0
  var pw=0
  var space=0
  var divider=0
  var font map[string][][]int

  switch res {
    case 0:
      font=owfonts.FontTitle4K
      space=10
      divider=30
    case 1:
      font=owfonts.FontTitle1080
      space=5
      divider=15
  }


  blackOut()
  for xi := 0; xi < width; xi++ {
    for yi := 0; yi < height; yi++ {
      if isPixIn(xi,yi) && !isPixOut(xi,yi) {
        // Did we hit more then a space ? (division between map title and game type)
        if px>0 && xi-px-pw>divider{
          line+="|"
        } else {
          // Did we hit a space ?
          if px>0 && xi-px-pw>space {
            line+=" "
          }
        }

        //fmt.Println("Reading char #: ",ccnt)

        //Intepret character from TitleFont
        ch,w:=getChar(xi,yi,font)
        line+=ch
        ccnt++
        px=xi
        pw=w
      }
    }
  }
  line=strings.Replace(line,"..",":",1) // : is intepreted as two single dots
  return line
}


// ----------------------------------------------------------------------------
// check if pixel is white (R,G,B > 224)

func isWhite(x int,y int) bool {
  return isAbove(x,y,224)
}

// ----------------------------------------------------------------------------
// check if pixel all R,G and B do have value higher then given value

func isAbove(x int,y int,tr int) bool {
  if pixIn(x,y).R > tr && pixIn(x,y).G > tr && pixIn(x,y).B > tr {
    return true
  }
  return false
}
