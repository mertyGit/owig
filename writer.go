package main

import (
  "fmt"
  "os"
  "time"
  "strings"
  "strconv"
  "github.com/360EntSecGroup-Skylar/excelize"
)

type Oww struct {
  init bool             // If structure has been initialized
  isxl bool             // If file is an xlsx file
  xlfile *excelize.File // File struct voor excel file
  xlidx int             // Index for owig sheet in excel file
  file *os.File         // If structure has been initialized
  ch chan GameInfo      // Channel to use 
  names []string
}

var wrt *Oww


func (w *Oww) Init() {
  w.names = []string{
    "logtime",
    "timestamp",
    "game time",
    "screen",
    "game state",
    "mapname",
    "gametype",
    "side",
    "objective",
    "pl point",
    "pl amount",
    "pl track%",
    "pl total%",
    "comp. defend score",
    "comp. attack score",
    "hero",
    "group id",
    "SR",
    "highest SR",
    "Eleminations",
    "Objective Kills",
    "Objective Time",
    "Hero Damage Done",
    "Healing Done",
    "Deaths",
    "Eleminations Medal",
    "Objective Kills Medal",
    "Objective Time Medal",
    "Hero Damage Done Medal",
    "Healing Done Medal",
    "Stat name 1",
    "Stat value 1",
    "Stat name 2",
    "Stat value 2",
    "Stat name 3",
    "Stat value 3",
    "Stat name 4",
    "Stat value 4",
    "Stat name 5",
    "Stat value 5",
    "Stat name 6",
    "Stat value 6",
    "enemy hero 1",
    "enemy hero 1 group id",
    "enemy hero 2",
    "enemy hero 2 group id",
    "enemy hero 3",
    "enemy hero 3 group id",
    "enemy hero 4",
    "enemy hero 4 group id",
    "enemy hero 5",
    "enemy hero 5 group id",
    "enemy hero 6",
    "enemy hero 6 group id",
    "own team hero 1",
    "own team hero 1 group id",
    "own team hero 2",
    "own team hero 2 group id",
    "own team hero 3",
    "own team hero 3 group id",
    "own team hero 4",
    "own team hero 4 group id",
    "own team hero 5",
    "own team hero 5 group id",
    "own team hero 6",
    "own team hero 6 group id",
  }
  w.ch=make(chan GameInfo)

  if strings.HasSuffix(config.stats,"xlsx") {
    w.isxl=true
  } else {
    w.isxl=false
  }

  if w.isxl {
    // Open Excel file
    file,err := excelize.OpenFile(config.stats)
    if err != nil {
      file = excelize.NewFile()
      if config.header {
        w.xlidx=file.NewSheet("owig")
        for i,v :=range w.names {
          file.SetCellValue("owig",excelize.ToAlphaString(i)+"1",v)
        }
      }
    }
    w.xlidx=file.GetSheetIndex("owig")
    if w.xlidx<1 {
      w.xlidx=file.NewSheet("owig")
    }
    w.xlfile=file
    w.xlfile.SaveAs(config.stats)
  } else {
    // Open csv file
    file,err := os.OpenFile(config.stats,os.O_APPEND|os.O_WRONLY,0666)
    if err != nil {
      file,err = os.OpenFile(config.stats,os.O_CREATE|os.O_WRONLY,0666)
      if err != nil {
        fmt.Println("WARNING: Statistics output file could not be opened or created")
        return
      }
      if config.header {
        for i,v :=range w.names {
          fmt.Fprintf(file,"\"%s\"",v)
          if i<len(w.names)-1 {
            fmt.Fprintf(file,"%s",config.divider)
          } else {
            fmt.Fprintf(file,"\n")
          }
        }
      }
    }
    w.file=file
  }
  w.init=true
}

func (w *Oww) Send() {
  if !w.init {
    return
  }
  w.ch<-game
}

func ChangedGI(a,b GameInfo) bool {
  if a.screen !=  b.screen {
    return true
  }
  if a.state !=  b.state {
    return true
  }
  if a.mapname !=  b.mapname {
    return true
  }
  if a.gametype !=  b.gametype {
    return true
  }
  if a.side !=  b.side {
    return true
  }
  if a.objective !=  b.objective {
    return true
  }
  if a.plpoint !=  b.plpoint {
    return true
  }
  if a.plamount !=  b.plamount {
    return true
  }
  if a.pltrack !=  b.pltrack {
    return true
  }
  if a.pltotal !=  b.pltotal {
    return true
  }
  if a.hero !=  b.hero {
    return true
  }
  if a.currentSR !=  b.currentSR {
    return true
  }
  if a.highestSR !=  b.highestSR {
    return true
  }
  if a.group !=  b.group {
    return true
  }
  if a.result !=  b.result {
    return true
  }
  r:=true
  for i:=0;i<6&&r;i++ {
    if a.medals[i] != b.medals[i] {
      r=false
    } else if a.lstats[i] != b.lstats[i] {
      r=false
    } else if a.rstats[i] != b.rstats[i] {
      r=false
    } else if a.own.hero[i] != b.own.hero[i] {
      r=false
    } else if a.own.groupid[i] != b.own.groupid[i] {
      r=false
    } else if a.enemy.hero[i] != b.enemy.hero[i] {
      r=false
    } else if a.enemy.groupid[i] != b.enemy.groupid[i] {
      r=false
    }
  }
  if !r {
    return true
  }
  return false
}

func SC2Name(id int) string {
  switch id {
    case SC_UNKNOWN:
      return "UNKNOWN"
    case SC_TAB:
      return "TAB"
    case SC_VICTORY:
      return "VICTORY"
    case SC_DEFEAT:
      return "DEFEAT"
    case SC_OVERVIEW:
      return "OVERVIEW"
    case SC_SRGAIN:
      return "SRGAIN"
    case SC_ASSEMBLE:
      return "ASSEMBLE"
    case SC_MAIN:
      return "MAIN"
    case SC_ENDING:
      return "ENDING"
    case SC_GAME:
      return "GAME"
    case SC_POTG:
      return "POTG"
    case SC_RESPAWN:
      return "RESPAWN"
  }
  return "UNASSIGNED"
}

func GS2Name (id int) string {
  switch id {
    case GS_NONE:
      return "START"
    case GS_START:
      return "START"
    case GS_RUN:
      return "RUN"
    case GS_END:
      return "END"
  }
  return "UNASSIGNED"
}



func (w *Oww) WriteCSV(g GameInfo) {
  t:=time.Now()
  // logtime
  fmt.Fprintf(w.file,"\"%s\"%s",t.Local(),config.divider)
  // timestamp
  fmt.Fprintf(w.file,"%d%s",g.ts,config.divider)
  // game time
  fmt.Fprintf(w.file,"\"%s\"%s",g.time,config.divider)
  // screen
  s:="UNKNOWN"
  fmt.Fprintf(w.file,"\"%s\"%s",SC2Name(g.screen),config.divider)
  // game state
  fmt.Fprintf(w.file,"\"%s\"%s",GS2Name(g.state),config.divider)
  // mapname
  fmt.Fprintf(w.file,"\"%s\"%s",g.mapname,config.divider)
  // gametype
  fmt.Fprintf(w.file,"\"%s\"%s",g.gametype,config.divider)
  // side
  fmt.Fprintf(w.file,"\"%s\"%s",g.side,config.divider)
  // objective
  fmt.Fprintf(w.file,"\"%s\"%s",g.objective,config.divider)
  // pl point
  fmt.Fprintf(w.file,"%d%s",g.plpoint,config.divider)
  // pl amount
  fmt.Fprintf(w.file,"%d%s",g.plamount,config.divider)
  // pl track%
  fmt.Fprintf(w.file,"%d%s",g.pltrack,config.divider)
  // pl total%
  fmt.Fprintf(w.file,"%d%s",g.pltotal,config.divider)
  // comp. defend score
  fmt.Fprintf(w.file,"%d%s",g.compdef,config.divider)
  // comp. attack score
  fmt.Fprintf(w.file,"%d%s",g.compatt,config.divider)
  // hero
  fmt.Fprintf(w.file,"\"%s\"%s",g.hero,config.divider)
  // Group id
  fmt.Fprintf(w.file,"%d%s",g.group,config.divider)
  // SR
  fmt.Fprintf(w.file,"%d%s",g.currentSR,config.divider)
  // highest SR
  fmt.Fprintf(w.file,"%d%s",g.highestSR,config.divider)
  // Eleminations
  fmt.Fprintf(w.file,"%s%s",g.lstats[0],config.divider)
  // Objective Kills
  fmt.Fprintf(w.file,"%s%s",g.lstats[1],config.divider)
  // Objective Time
  fmt.Fprintf(w.file,"\"%s\"%s",g.lstats[2],config.divider)
  // Hero Damage Done
  fmt.Fprintf(w.file,"%s%s",g.lstats[3],config.divider)
  // Healing Done
  fmt.Fprintf(w.file,"%s%s",g.lstats[4],config.divider)
  // Deaths
  fmt.Fprintf(w.file,"%s%s",g.lstats[5],config.divider)
  // Eleminations Medal
  fmt.Fprintf(w.file,"\"%s\"%s",g.medals[0],config.divider)
  // Objective Kills Medal
  fmt.Fprintf(w.file,"\"%s\"%s",g.medals[1],config.divider)
  // Objective Time Medal
  fmt.Fprintf(w.file,"\"%s\"%s",g.medals[2],config.divider)
  // Hero Damage Done Medal
  fmt.Fprintf(w.file,"\"%s\"%s",g.medals[3],config.divider)
  // Healing Done Medal
  fmt.Fprintf(w.file,"\"%s\"%s",g.medals[4],config.divider)
  // Stat name 1
  fmt.Fprintf(w.file,"\"%s\"%s",getStatsline(g.hero,0),config.divider)
  // Stat value 1
  s=strings.Replace(g.rstats[0],"%","",-1)
  fmt.Fprintf(w.file,"%s%s",s,config.divider)
  // Stat name 2
  fmt.Fprintf(w.file,"\"%s\"%s",getStatsline(g.hero,1),config.divider)
  // Stat value 2
  s=strings.Replace(g.rstats[1],"%","",-1)
  fmt.Fprintf(w.file,"%s%s",s,config.divider)
  // Stat name 3
  fmt.Fprintf(w.file,"\"%s\"%s",getStatsline(g.hero,2),config.divider)
  // Stat value 3
  s=strings.Replace(g.rstats[2],"%","",-1)
  fmt.Fprintf(w.file,"%s%s",s,config.divider)
  // Stat name 4
  fmt.Fprintf(w.file,"\"%s\"%s",getStatsline(g.hero,3),config.divider)
  // Stat value 4
  s=strings.Replace(g.rstats[3],"%","",-1)
  fmt.Fprintf(w.file,"%s%s",s,config.divider)
  // Stat name 5
  fmt.Fprintf(w.file,"\"%s\"%s",getStatsline(g.hero,4),config.divider)
  // Stat value 5
  s=strings.Replace(g.rstats[4],"%","",-1)
  fmt.Fprintf(w.file,"%s%s",s,config.divider)
  // Stat name 6
  fmt.Fprintf(w.file,"\"%s\"%s",getStatsline(g.hero,5),config.divider)
  // Stat value 6
  s=strings.Replace(g.rstats[5],"%","",-1)
  fmt.Fprintf(w.file,"%s%s",s,config.divider)
  // enemy hero 1
  fmt.Fprintf(w.file,"\"%s\"%s",g.enemy.hero[0],config.divider)
  // enemy hero 1 group id
  fmt.Fprintf(w.file,"%d%s",g.enemy.groupid[0],config.divider)
  // enemy hero 2
  fmt.Fprintf(w.file,"\"%s\"%s",g.enemy.hero[1],config.divider)
  // enemy hero 2 group id
  fmt.Fprintf(w.file,"%d%s",g.enemy.groupid[1],config.divider)
  // enemy hero 3
  fmt.Fprintf(w.file,"\"%s\"%s",g.enemy.hero[2],config.divider)
  // enemy hero 3 group id
  fmt.Fprintf(w.file,"%d%s",g.enemy.groupid[2],config.divider)
  // enemy hero 4
  fmt.Fprintf(w.file,"\"%s\"%s",g.enemy.hero[3],config.divider)
  // enemy hero 4 group id
  fmt.Fprintf(w.file,"%d%s",g.enemy.groupid[3],config.divider)
  // enemy hero 5
  fmt.Fprintf(w.file,"\"%s\"%s",g.enemy.hero[4],config.divider)
  // enemy hero 5 group id
  fmt.Fprintf(w.file,"%d%s",g.enemy.groupid[4],config.divider)
  // enemy hero 6
  fmt.Fprintf(w.file,"\"%s\"%s",g.enemy.hero[5],config.divider)
  // enemy hero 6 group id
  fmt.Fprintf(w.file,"%d%s",g.enemy.groupid[5],config.divider)
  // own team hero 1
  fmt.Fprintf(w.file,"\"%s\"%s",g.own.hero[0],config.divider)
  // own team hero 1 group id
  fmt.Fprintf(w.file,"%d%s",g.own.groupid[0],config.divider)
  // own team hero 2
  fmt.Fprintf(w.file,"\"%s\"%s",g.own.hero[1],config.divider)
  // own team hero 2 group id
  fmt.Fprintf(w.file,"%d%s",g.own.groupid[1],config.divider)
  // own team hero 3
  fmt.Fprintf(w.file,"\"%s\"%s",g.own.hero[2],config.divider)
  // own team hero 3 group id
  fmt.Fprintf(w.file,"%d%s",g.own.groupid[2],config.divider)
  // own team hero 4
  fmt.Fprintf(w.file,"\"%s\"%s",g.own.hero[3],config.divider)
  // own team hero 4 group id
  fmt.Fprintf(w.file,"%d%s",g.own.groupid[3],config.divider)
  // own team hero 5
  fmt.Fprintf(w.file,"\"%s\"%s",g.own.hero[4],config.divider)
  // own team hero 5 group id
  fmt.Fprintf(w.file,"%d%s",g.own.groupid[4],config.divider)
  // own team hero 6
  fmt.Fprintf(w.file,"\"%s\"%s",g.own.hero[5],config.divider)
  // own team hero 6 group id
  fmt.Fprintf(w.file,"%d",g.own.groupid[5])
  fmt.Fprintf(w.file,"\n")
}

func (w *Oww) WriteXLSX(g GameInfo) {
  var s string

  row:=fmt.Sprintf("%d",len(w.xlfile.GetRows("owig"))+1)
  style,_ := w.xlfile.NewStyle(`{"number_format":1}`)


  t:=time.Now()
  i:=0
  // logtime
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,t.Local())
  i++
  // timestamp
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.ts)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // game time
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.time)
  i++
  // screen
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,SC2Name(g.screen))
  i++
  // game state
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,GS2Name(g.state))
  i++
  // mapname
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.mapname)
  i++
  // gametype
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.gametype)
  i++
  // side
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.side)
  i++
  // objective
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.objective)
  i++
  // pl point
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.plpoint)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // pl amount
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.plamount)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // pl track%
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.pltrack)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // pl total%
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.pltotal)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // comp. defend score
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.compdef)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // comp. attack score
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.compatt)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // hero
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.hero)
  i++
  // Group id
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.group)
  i++
  // SR
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.currentSR)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // highest SR
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.highestSR)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // Eleminations
  v,_:=strconv.Atoi(g.lstats[0])
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,v)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // Objective Kills
  v,_=strconv.Atoi(g.lstats[1])
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,v)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // Objective Time
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.lstats[2])
  i++
  // Hero Damage Done
  v,_=strconv.Atoi(g.lstats[3])
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,v)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // Healing Done
  v,_=strconv.Atoi(g.lstats[4])
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,v)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // Deaths
  v,_=strconv.Atoi(g.lstats[5])
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,v)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // Eleminations Medal
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.medals[0])
  i++
  // Objective Kills Medal
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.medals[1])
  i++
  // Objective Time Medal
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.medals[2])
  i++
  // Hero Damage Done Medal
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.medals[3])
  i++
  // Healing Done Medal
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.medals[4])
  i++
  // Stat name 1
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,getStatsline(g.hero,0))
  i++
  // Stat value 1
  s=strings.Replace(g.rstats[0],"%","",-1)
  v,_=strconv.Atoi(s)
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,v)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // Stat name 2
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,getStatsline(g.hero,1))
  i++
  // Stat value 2
  s=strings.Replace(g.rstats[1],"%","",-1)
  v,_=strconv.Atoi(s)
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,v)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // Stat name 3
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,getStatsline(g.hero,2))
  i++
  // Stat value 3
  s=strings.Replace(g.rstats[2],"%","",-1)
  v,_=strconv.Atoi(s)
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,s)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // Stat name 4
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,getStatsline(g.hero,3))
  i++
  // Stat value 4
  s=strings.Replace(g.rstats[3],"%","",-1)
  v,_=strconv.Atoi(s)
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,v)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // Stat name 5
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,getStatsline(g.hero,4))
  i++
  // Stat value 5
  s=strings.Replace(g.rstats[4],"%","",-1)
  v,_=strconv.Atoi(s)
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,v)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // Stat name 6
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,getStatsline(g.hero,5))
  i++
  // Stat value 6
  s=strings.Replace(g.rstats[5],"%","",-1)
  v,_=strconv.Atoi(s)
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,v)
  w.xlfile.SetCellStyle("owig",excelize.ToAlphaString(i)+row,excelize.ToAlphaString(i)+row,style)
  i++
  // enemy hero 1
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.enemy.hero[0])
  i++
  // enemy hero 1 group id
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.enemy.groupid[0])
  i++
  // enemy hero 2
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.enemy.hero[1])
  i++
  // enemy hero 2 group id
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.enemy.groupid[1])
  i++
  // enemy hero 3
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.enemy.hero[2])
  i++
  // enemy hero 3 group id
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.enemy.groupid[2])
  i++
  // enemy hero 4
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.enemy.hero[3])
  i++
  // enemy hero 4 group id
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.enemy.groupid[3])
  i++
  // enemy hero 5
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.enemy.hero[4])
  i++
  // enemy hero 5 group id
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.enemy.groupid[4])
  i++
  // enemy hero 6
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.enemy.hero[5])
  i++
  // enemy hero 6 group id
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.enemy.groupid[5])
  i++
  // own team hero 1
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.own.hero[0])
  i++
  // own team hero 1 group id
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.own.groupid[0])
  i++
  // own team hero 2
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.own.hero[1])
  i++
  // own team hero 2 group id
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.own.groupid[1])
  i++
  // own team hero 3
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.own.hero[2])
  i++
  // own team hero 3 group id
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.own.groupid[2])
  i++
  // own team hero 4
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.own.hero[3])
  i++
  // own team hero 4 group id
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.own.groupid[3])
  i++
  // own team hero 5
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.own.hero[4])
  i++
  // own team hero 5 group id
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.own.groupid[4])
  i++
  // own team hero 6
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.own.hero[5])
  i++
  // own team hero 6 group id
  w.xlfile.SetCellValue("owig",excelize.ToAlphaString(i)+row,g.own.groupid[5])
  i++
  w.xlfile.SetActiveSheet(w.xlidx)
  w.xlfile.SaveAs(config.stats)
}

func (w *Oww) Run() {
  var old GameInfo
  if !w.init {
    return
  }
  for {
    // Get local copy of game statistics
    g:=<-w.ch
    skip:=false
    // prevent sending same information if information is the same
    if ChangedGI(g,old) && g.screen != SC_UNKNOWN {
      // prevent sending multiple SRGAIN intepretations
      if g.pscreen == g.screen && g.screen == SC_SRGAIN {
        skip=true
      } else if g.pscreen == SC_SRGAIN {
        if (w.isxl) {
          w.WriteXLSX(old)
        } else {
          w.WriteCSV(old)
        }
      }
      if g.screen != SC_UNKNOWN && !skip {
        if (w.isxl) {
          w.WriteXLSX(g)
        } else {
          w.WriteCSV(g)
        }
        old=g
      }
    }
  }
}
