package main

import (
  "strings"
)

// Functions to adjust any OCR fuckups
// Check on specific unique parts that "survived" most of the OCR runs
// and use that to identify right mapname


func forceName(n string, fl map[string]string ) string {
  for k,v := range fl {
    if strings.Contains(n,k) {
      n=v
    }
  }
  return n
}

func checkName(n string, fl map[string]string ) bool {
  for _,v := range fl {
    if n==v {
      return true
    }
  }
  return false
}


func cleanGametype(m string) string {
  fl := map[string]string{
    "COM" : "COMPETITIVE PLAY",
    "UIC" : "QUICK PLAY",
    "MYS" : "MYSTERY HEROES",
    "LIM" : "NO LIMITS",
    "GAME": "GAME BROWSER",
  }
  r:=forceName(m,fl)
  if r!=m || checkName(r,fl) {
    game.forceT=true
    return m
  }
  // Gametype not in list and found earlier one that was in list ?
  // dont use found name, but already found name
  if game.forceT {
    return game.gametype
  }

  return r
}


func cleanMapname(m string) string {
  fl := map[string]string{
    "WORLD" : "BLIZZARD WORLD",
    "RADO"  : "DORADO",
    "WAL"   : "EICHENWALDE",
    "NAMUR" : "HANAMURA",
    "WOOD"  : "HOLLYWOOD",
    "COLON" : "HORIZON LUNAR COLONY",
    "LIO"   : "ILIOS",
    "UNKER" : "JUNKERTOWN",
    "TOWN"  : "JUNKERTOWN",
    "ROW"   : "KING'S ROW",
    "TOWER" : "LIJIANG TOWER",
    "BANI"  : "NUMBANI",
    "EPAL"  : "NEPAL",
    "ASIS"  : "OASIS",
    "UTE"   : "ROUTE 66",
    "EMPLE" : "TEMPLE OF ANUBIS",
    "INDU"  : "VOLSKAYA INDUSTRIES",
    "OINT"  : "WATCHPOINT: GIBRALTAR",
  }
  r:=forceName(m,fl)
  if r!=m || checkName(r,fl) {
    game.forceM=true
    return m
  }
  // Mapname not in list and found earlier one that was in list ?
  // dont use found name, but already found name
  if game.forceM {
    return game.mapname
  }
  return r
}
