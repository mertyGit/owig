package main

import (
  "fmt"
  "os"
  "path/filepath"
  "github.com/go-ini/ini"
)
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
  dbg_time  bool
}

var config Ini



// ----------------------------------------------------------------------------
// Read "owig.ini" file

func getIni() {

  // Fill in default values
  config.sleep=1000
  config.stats="owig_stats.csv"
  config.divider=","
  config.header=true
  config.dbg_screen=false
  config.dbg_window=false
  config.dbg_ocr=false
  config.dbg_time=false
  config.dbg_pause=2000

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
    config.dbg_screen,_=cfg.Section("debug").Key("screen").Bool()
  }
  if cfg.Section("debug").HasKey("window") {
    config.dbg_window,_=cfg.Section("debug").Key("window").Bool()
  }
  if cfg.Section("debug").HasKey("time") {
    config.dbg_time,_=cfg.Section("debug").Key("time").Bool()
  }
  if cfg.Section("debug").HasKey("ocr") {
    config.dbg_ocr,_=cfg.Section("debug").Key("ocr").Bool()
  }
  if cfg.Section("debug").HasKey("pause") {
    config.dbg_pause,_=cfg.Section("debug").Key("pause").Int()
  }
  // Initialize vars program startup here
  game.dmsg[0]=""
  game.dmsg[1]=""
  game.dmsg[2]=""
  game.dmsg[3]="INI file loaded"

  game.state=GS_NONE
  game.side=""
}

