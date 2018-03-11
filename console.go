package main

import (
  "fmt"
  "github.com/fatih/color"
)


// ----------------------------------------------------------------------------
// Print Medal information

func printMedal(m string) {
  switch m {
    case "G":
      color.Set(color.FgYellow)
      fmt.Print("(G)")
      color.Set(color.FgWhite)
    case "S":
      color.Set(color.FgHiWhite)
      fmt.Print("(S)")
      color.Set(color.FgWhite)
    case "B":
      color.Set(color.FgRed)
      fmt.Print("(B)")
      color.Set(color.FgWhite)
    default:
      fmt.Print("( )")
  }
}

// ----------------------------------------------------------------------------
// Dump statistics on console

func dumpTabStats() {
    fmt.Println()
    fmt.Print("             ")
    color.Set(color.FgWhite,color.Bold)
    fmt.Print(game.mapname)
    color.Set(color.FgWhite)
    fmt.Print(" | ")
    color.Set(color.FgYellow,color.Bold)
    fmt.Println(game.gametype)
    color.Unset()
    color.Set(color.FgHiBlue,color.Underline)
    fmt.Println("                                                                                ")
    color.Unset()
    fmt.Println()
    for x := 0; x < 6; x++ {
      if game.enemy.groupid[x]>0 {
        color.Set(color.FgHiWhite)
      }
      fmt.Printf(" %-10s",guessHero(x,0))
      if x<5 && game.enemy.groupid[x] >0 && game.enemy.groupid[x] == game.enemy.groupid[x+1] {
        fmt.Print(" - ");
      } else {
        fmt.Print("   ");
      }
      color.Unset()
    }
    fmt.Println()
    fmt.Println()
    color.HiBlue("-------------------------------------= V S =------------------------------------")
    fmt.Println()
    for x := 0; x < 6; x++ {
      if game.own.groupid[0]==1 && game.own.groupid[x]==1 {
        color.Set(color.FgHiGreen)
      } else {
        if game.own.groupid[x]>0 {
          color.Set(color.FgHiWhite)
        } else {
          color.Set(color.FgWhite)
        }
      }
      fmt.Printf(" %-10s",guessHero(x,1))
      if x<5 && game.own.groupid[x] >0 && game.own.groupid[x] == game.own.groupid[x+1] {
        fmt.Print(" - ")
      } else {
        fmt.Print("   ")
      }
      color.Unset()
    }
    fmt.Println()
    color.Set(color.FgHiBlue,color.Underline)
    fmt.Println("                                                                                ")
    color.Unset()
    fmt.Println()
    fmt.Print(" Eliminations     ")
    printMedal(getMedal(0))
    fmt.Printf(":%8s\n",getStats(0,0))
    fmt.Print(" Objective kills  ")
    printMedal(getMedal(1))
    fmt.Printf(":%8s\n",getStats(1,0))
    fmt.Print(" Objective time   ")
    printMedal(getMedal(2))
    fmt.Printf(":%8s\n",getStats(2,0))
    fmt.Print(" Hero Damage Done ")
    printMedal(getMedal(3))
    fmt.Printf(":%8s\n",getStats(0,1))
    fmt.Print(" Healing Done     ")
    printMedal(getMedal(4))
    fmt.Printf(":%8s\n",getStats(1,1))
    fmt.Print(" Deaths              ")
    fmt.Printf(":%8s\n",getStats(2,1))
    color.Set(color.FgHiBlue,color.Underline)
    fmt.Println("                                                                                ")
    fmt.Println()
    color.Unset()
    fmt.Println(" Stat 1: ",getStats(3,0))
    fmt.Println(" Stat 2: ",getStats(4,0))
    fmt.Println(" Stat 3: ",getStats(5,0))
    fmt.Println(" stat 4: ",getStats(3,1))
    fmt.Println(" stat 5: ",getStats(4,1))
    fmt.Println(" stat 6: ",getStats(5,1))
    color.Set(color.FgHiBlue,color.Underline)
    fmt.Println("                                                                                ")
    color.Unset()
    fmt.Println()
}
