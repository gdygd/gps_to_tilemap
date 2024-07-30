# 'Get tilemap(z/x/y) corresponding to the gps(lat, lon) range.'

## requirement
 * You have to get entire tilemap like vworld in advance.

 * It only supports within the Korean gps range.

## usage 
1. clone and build
2. configure "conf.ini"
3. execute process

## description conf.ini
 * srcpath: your entire tilemap path 
 * destpath: Path to store tiles corresponding to the gps range
 * lat1 : start latitude
 * lon1 : start longitude
 * lat2 : end latitude
 * lon2 : end longitude
 * startlv : start zoom level
 * endlv : end zoom level

 ## conf.ini example
 ```
[MAPINFO]
srcpath = /home/user/vworldmap
destpath = /home/user/tilemap
lat1 = 37.20747576600081
lon1 = 126.97826010879602
lat2 = 37.11759097226632
lon2 = 127.10547191439518
startlv = 6
endlv = 11
 ```