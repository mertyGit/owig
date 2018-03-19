package main

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
