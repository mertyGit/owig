# owig
Overwatch Information Grabber
This program runs alongside your Overwatch game and get statistics and other information by grabbing screen on certain intervals and deciphering it.
Information will be displayed on seperate window and written to disk to be analyzed later on with tool of your choice.

At this moment, will work only with two resolutions: 1080p (1920x1080) or 4K (3840x2160). 
Don't use programs that are using overlays like Discord, it will mess with the recognition

# Release information
WARNING: STILL BETA Usable to track on going progress and stats, but no garantees whatsover

Version 0.81: 
- Restructuring of files and rewrite of code for using simplified methods on image struct. 
- Rewrite and improvement of OCR code.  

# TODO:
* Write CSV output
* Write Excel output (?)
* Recognision of Gaming screen to gather gamestats (score, position payload, objectives captured)
* Card information on end of game
* Recognize and notify about (enemy-) composition changes


# Bugs:
* Wrong choice of statistic lines for own hero: When joining game as group, but leaving when in same game, your own hero position in TAB screen isn't on the left anymore..
* Assembly screen isn't recognized or recognized at wrong time
* Can't determine attack/defend on assembly screen: Algorithm uses color to find out, but in some cases "Attack" is written in blue....
* When viewing stats of someone who has same icon, it will assume these are your stats
* Victory or Defeat message with scores under it will not be recognized

# More about OWIG
OWIG is program that will gather game statics while Overwatch game is running or by feeding it screenshots
It consists out of a windows program that will run alongside the overwatch game and will do a screenshot every second (configurable). This screenshot will be intepreted, digits will be ocr-ed and information in the screen will be interpreted by looking at distinguish features like pixel colors and appereance of buttons on the screen. 

In short, it takes a screenshot like this:
![Example screenshot](https://raw.githubusercontent.com/mertyGit/owig/master/doc/screenshot_example.png)

And tries to decrypt the information writes it to disk and to the Windows GUI screen: 

<img src="https://raw.githubusercontent.com/mertyGit/owig/master/doc/example.png" width="400">


# Features

Right now, it is able to figure out:

In game statistics (what you get when pressing TAB):
* Team composition, heroes choosen during game
* Groups
* Medals
* All stats on lower end of screen
* Mapname
* Gametype (Quickplay, Mystery Heroes and so on)

Assemble screen (starting of match):
* Attacking or defending

End screen (end of match, voting for cards and experience progression) :
* Victory, Draw or Defeat

SR Gain/Loss screen (after end of match, when playing comp):
* Current SR

Statistics overview: (right click own icon, view career statistics)
* Current SR
* Highest SR
