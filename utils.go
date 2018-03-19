package main
import (
  "fmt"
  "time"
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
// Print a timestamp, for timing & debugging purposes

func ts(id string) {
  if config.dbg_time {
    fmt.Println("TIME:",time.Now().UnixNano()/1000000 - game.ts,id)
  }
}

// ----------------------------------------------------------------------------
// Get name of statistics on right bottom corner with TAB screen

func getStatsline(hero string, i int) string {
  lines,ok:=heroStats[hero]
  if ok {
    return lines[i]
  } else {
    return ""
  }
}
