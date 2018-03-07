# owig
Overwatch Information Grabber , get statistics and other information by grabbing screen and deciphering it
# About
Warning: Still very, very beta version

OWIG is program that will gather game statics will Overwatch game is running.
It consists out of a windows program that will run alongside the overwatch game and will do a screenshot every second. This screenshot will be intepreted, digits will be ocr-ed and information in the screen will be interpreted by looking at distinguish features like pixel colors and appereance of buttons on the screen. 

At this moment, will work only with two resolutions: 1080p (1920x1080) or 4K (3840x2160). 

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

