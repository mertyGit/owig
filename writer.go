package main

import (
  "fmt"
  "os"
  "time"
  "strings"
)

type Oww struct {
  init bool            // If structure has been initialized
  file *os.File        // If structure has been initialized
  ch chan GameInfo     // Channel to use 
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
  file,err := os.OpenFile(config.stats,os.O_APPEND|os.O_WRONLY,0666)
  if (err != nil) {
    file,err = os.OpenFile(config.stats,os.O_CREATE|os.O_WRONLY,0666)
    if (err != nil) {
      fmt.Println("WARNING: Statistics output file could not be opened")
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
  w.init=true
}

func (w *Oww) Send() {
  if !w.init {
    return
  }
  w.ch<-game
}

func (w *Oww) Run() {
  if !w.init {
    return
  }
  for {
    // Get local copy of game statistics
    g:=<-w.ch
    if g.screen != SC_UNKNOWN {
      t:=time.Now()
      // "logtime"
      fmt.Fprintf(w.file,"\"%s\"%s",t.Local(),config.divider)
      // "timestamp"
      fmt.Fprintf(w.file,"%d%s",g.ts,config.divider)
      //"game time"
      fmt.Fprintf(w.file,"\"%s\"%s",g.time,config.divider)
      // "screen"
      s:="UNKNOWN"
      switch g.screen {
        case SC_TAB:
          s="TAB"
        case SC_VICTORY:
          s="VICTORY"
        case SC_DEFEAT:
          s="DEFEAT"
        case SC_OVERVIEW:
          s="OvERVIEW"
        case SC_SRGAIN:
          s="SRGAIN"
        case SC_ASSEMBLE:
          s="ASSEMBLE"
        case SC_MAIN:
          s="MAIN"
        case SC_ENDING:
          s="ENDING"
        case SC_GAME:
          s="GAME"
      }
      fmt.Fprintf(w.file,"\"%s\"%s",s,config.divider)
      // "game state"
      s="NONE"
      switch g.state {
        case GS_START:
          s="START"
        case GS_RUN:
          s="RUN"
        case GS_END:
          s="END"
      }
      fmt.Fprintf(w.file,"\"%s\"%s",s,config.divider)
      // mapname
      fmt.Fprintf(w.file,"\"%s\"%s",g.mapname,config.divider)
      // gametype
      fmt.Fprintf(w.file,"\"%s\"%s",g.gametype,config.divider)
      // side
      fmt.Fprintf(w.file,"\"%s\"%s",g.side,config.divider)
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
  }
}
