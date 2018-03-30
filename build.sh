rm -f rsrc.syso bin/owig.exe 
rsrc -manifest owig.manifest -ico "owig256.ico,owig48.ico,owig32.ico,owig16.ico" -o rsrc.syso
go build -ldflags="-H windowsgui" -o bin/owig.exe
