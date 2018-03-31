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

func cleanGametype(m string) string {
  fl := map[string]string{
    "COM": "COMPETITIVE PLAY",
    "UIC": "QUICK PLAY",
  }

  return forceName(m,fl)
}

func cleanMapname(m string) string {
  fl := map[string]string{
    "ROW": "KING's ROW",
    "EMPLE": "TEMPLE OF ANUBIS",
    "RADO": "DORADO",
    "LIO": "ILIOS",
    "WOOD": "HOLLYWOOD",
    "OINT": "WATCHPOINT: GIBRALTAR",
    "ASIS": "OASIS",
    "NAMUR": "HANAMURA",
    "TOWN": "JUNKERTOWN",
    "UNKER": "JUNKERTOWN",
  }

  return forceName(m,fl)
}
