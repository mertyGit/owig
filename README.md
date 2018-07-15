# owig
Overwatch Information Grabber
This program runs alongside your Overwatch game and get statistics and other information by grabbing screen on certain intervals and deciphering it.
Information will be displayed on seperate window and written to disk to be analyzed later on with tool of your choice.

At this moment, will work only with three resolutions: 1080p (1920x1080), WQHD (2560x1440) or 4K (3840x2160). 
Don't use programs that are using overlays like Discord, it will mess with the recognition

Version 1.00 is now available to download to. Please check https://www.mertymade.com/owig/

# Release information

Version 1.00: 
- Updated to final version
- Minor cleanup in code

Version 0.94: 

New functionality:
- WQHD resolution (2560x1440) added
- Adjusted to latest changes in GUI of Overview Screen (current, high SR)
- Adjusted to latest changes in TAB Screen (position of icons)
- Added Hero "Wrecking Ball"
- Changed to latest changes in Hero Stats (like Symetra)

Bug Fixes:
- Torbjorn to Torbjörn
- Incidently not recognizing of TAB screen fixed

Version 0.93: 

New functionality:
- Added information about objectives (WAITING, A,B or PAYLOAD) and more info about payload; points taken, total amount of points, position of payload relative to track between two points and total distance start-end
- Added defend & attack score for competitive games
- Cooler simpler logo

Bug Fixes
- Any wrong intepretation of attack/defend during assembly screen, will be corrected via in-game information
- No crashes when loading with non-existing (test)file
- View cleanup & prevention settings of mismanagement of known mapnames
- Preventing recognition of map title , time , game type when chat icon is blocking it

Version 0.92: 
- Added XLSX output functionality; now you can specify ".xlsx" file in statsfile, it will create or add rows to the "owig" sheet in your excel file

Version 0.91: 
- Adjusting OCR mishaps, forcing to known mapnames and types
- You can use other ini files, given as first argument, so per ~~smurf~~account, you can have your own statistics (and settings)

Version 0.90: 
- Added CSV output functionality

Version 0.82: 
- Recognizion of choosen hero in TAB screen

Version 0.81: 
- Restructuring of files and rewrite of code for using simplified methods on image struct. 
- Rewrite and improvement of OCR code.  

# TODO:
* ~~Write CSV output~~
* ~~Using more then one .ini file (to track multiple profiles/accounts)~~
* ~~Write Excel output (?)~~
* ~~Recognition of Gaming screen to gather gamestats (score,position payload, objectives captured)~~
* Recognition of group SR at start of game

# MAYBE TODO:
* Recognize and notify about (enemy-) composition changes (tracked, but can easily determined or calculated in output)
* percentage of capture point taken (hard)
* Card information on end of game (hard, unknown who's card is yours)

# DONT:
* Scanning of username/battletag : to hard to write/learn ocr that understands ALL unicode characters for all fonts used


# Bugs:
* ~~Wrong choice of statistic lines for own hero: When joining game as group, but leaving when in same game, your own hero position in TAB screen isn't on the left anymore..~~
* ~~Assembly screen isn't recognized or recognized at wrong time~~
* ~~Can't determine attack/defend on assembly screen: Algorithm uses color to find out, but in some cases "Attack" is written in blue....~~
* When viewing stats of someone who has same icon, it will assume these are your stats
* Victory or Defeat message with scores under it will not be recognized
* ~~Crashes on missing or wrong files given at startup~~

# More about OWIG
OWIG is program that will gather game statics while Overwatch game is running or by feeding it screenshots
It consists out of a windows program that will run alongside the overwatch game and will do a screenshot every second (configurable). This screenshot will be intepreted, digits will be ocr-ed and information in the screen will be interpreted by looking at distinguish features like pixel colors and appereance of buttons on the screen. 

In short, it takes a screenshot like this:
![Example screenshot](https://raw.githubusercontent.com/mertyGit/owig/master/doc/screenshot_example.png)

And tries to decrypt the information writes it to disk and to the Windows GUI screen: 

<img src="https://raw.githubusercontent.com/mertyGit/owig/master/doc/example.png" width="400">


# Features

Right now, it is able to figure out:

Game screen (default game screen)
* Competitive scores defend/attack 
* Side defend, attack
* Objective (A,B,Payload)
* Payload progress % between points (track)
* Payload progress % between start and finish (total)
* Payload Points taken
* All Payload Points

In game statistics (what you get when pressing TAB):
* Team composition, heroes choosen during game
* Groups
* Medals
* All stats on lower end of screen
* Mapname
* Time
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
