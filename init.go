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
  game.enemy.isChanged=false
  game.own.isChanged=false
  game.stats.eleminations=0
  game.stats.objectiveKills=0
  game.stats.objectiveTime=""
  game.stats.objectiveSecs=0
  game.stats.damage=0
  game.stats.deaths=0
  game.time=""
  game.ts=time.Now().UnixNano()/1000000
  for i:=0;i<6;i++ {
    game.stats.medals[i]=""
    game.stats.stats[i]=""
    game.stats.statsText[i]=""
    game.enemy.hero[i]=""
    game.enemy.groupid[i]=0
    game.enemy.switches[i]=0
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
  heroStats=map[string][]string{
    "Ana":{
       "Unscoped Accuracy",
       "Scoped Accuracy",
       "Defensive Assists",
       "Nano Boost Assists",
       "Enemies Slept",
       "",
    },
    "Bastion":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Recon Kills",
       "Sentry Kills",
       "Tank Kills",
       "Self Healing",
    },
    "Brigitte":{
       "Offensive Assists",
       "Defensive Assists",
       "Damage Blocked",
       "Armor Provided",
       "",
       "",
       "",
    },
    "Doomfist":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Ability Damage Done",
       "Meteor Strike Kill",
       "Shields Created",
    },
    "D.Va":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Damage Blocked",
       "Self-Destruct Kills",
       "Mechs Called",
       "",
    },
    "Genji":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Damage Reflected",
       "Dragonblade Kills",
       "",
    },
    "Hanzo":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Damage Reflected",
       "Critical Hits",
       "Recon Assists",
       "Dragonstrike Kills",
    },
    "Junkrat":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Enemies Trapped",
       "Rip-Tire Kills",
       "",
    },
    "Lucio":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Sound Barriers Provided",
       "Offensive Assists",
       "Defensive Assists",
       "",
    },
    "McCree":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Critical Hits",
       "Deadeye Kills",
       "Fan The Hammer Kill",
    },
    "Mei":{
       "Damage Blocked",
       "Kill Streak - Best",
       "Enemies Frozen",
       "Blizzard Kills",
       "Self Healing",
       "",
    },
    "Mercy":{
       "Offensive Assists",
       "Defense Assists",
       "Player Resurrected",
       "Blaster Kills",
       "",
       "",
    },
    "Moira":{
       "Secondary Fire Accuracy",
       "Kill Streak - Best",
       "Defensive Assists",
       "Coalescence Kills",
       "Coalescence Healing",
       "Self Healing",
    },
    "Orisa":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Damage Blocked",
       "Offensive Assists",
       "Damage Amplified",
       "",
    },
    "Pharah":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Barrage Kill",
       "Rocket Direct Hits",
       "",
    },
    "Reaper":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Death Blossom Kills",
       "Self Healing",
       "",
    },
    "Reinhardt":{
       "Damage Blocked",
       "Kill Streak - Best",
       "Charge Kills",
       "Fire Strike Kills",
       "Earthshatter Kills",
       "",
    },
    "Roadhog":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Enemies Hooked",
       "Hook Accuracy",
       "Self Healing",
       "Whole Hog Kills",
    },
    "Soldier 76":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Helix Rockets Kills",
       "Tactical Visor Kills",
       "",
    },
    "Sombra":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Offensive Assists",
       "Enemies Hacked",
       "Enemies EMP'D",
       "",
    },
    "Symmetra":{
       "Sentry Turret Kills",
       "Kill Streak - Best",
       "Damage Blocked",
       "Players Teleported",
       "Teleporter Uptime - Average",
       "",
    },
    "Torbjorn":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Torbjorn Kills",
       "Turret Kills",
       "Molten Core Kills",
       "Armor Packs Created",
    },
    "Tracer":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Pulse Bomb kills",
       "Pulse Bombs Attached",
       "",
    },
    "Widowmaker":{
       "Scoped Accuracy",
       "Kill Streak - Best",
       "Final Blows",
       "Scoped Critical Hits",
       "Recon Assists",
       "",
    },
    "Winston":{
       "Damage Blocked",
       "Kill Streak - Best",
       "Melee Kill",
       "Players Knocked Back",
       "",
       "",
    },
    "Zarya":{
       "Damage Blocked",
       "Kill Streak - Best",
       "High Energy Kills",
       "Average Energy",
       "Graviton Surge Kills",
       "",
    },
    "Zenyatta":{
       "Weapon Accuracy",
       "Kill Streak - Best",
       "Offensive Assists",
       "Defense Assists",
       "Best Transcendence Heal",
       "",
    },
  }
}
