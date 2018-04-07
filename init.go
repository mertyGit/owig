package main

import (
  "fmt"
  "time"
)

func initGameInfo() {

  if config.dbg_screen {
    fmt.Println("== initGameInfo ==")
  }
  game.mapname=""
  game.gametype=""
  game.hero=""
  game.group=0
  game.result=""
  game.forceM=false
  game.forceT=false
  game.enemy.isChanged=false
  game.own.isChanged=false
  game.time=""
  game.ts=time.Now().UnixNano()/1000000
  for i:=0;i<6;i++ {
    game.enemy.hero[i]=""
    game.enemy.groupid[i]=0
    game.enemy.switches[i]=0
    game.own.hero[i]=""
    game.own.groupid[i]=0
    game.own.switches[i]=0
    game.lstats[i]=""
    game.rstats[i]=""
    game.medals[i]=""
  }
}

func initOwig() {
  if config.dbg_screen {
    fmt.Println("== initOwig ==")
  }
  getIni()
  initGameInfo()
  loadIcons()
  initHStats()
  wrt.Init()
}
