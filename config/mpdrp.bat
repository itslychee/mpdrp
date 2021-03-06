@ECHO OFF

SetLocal EnableDelayedExpansion

IF NOT "%1"=="as_admin" (
    powershell -Command "Start-Process -FilePath %0 as_admin -Verb RunAs" 
    exit
)

:begin

echo [MPDRP Service]
echo 1) Add mpdrp to Windows Services
echo 2) Remove mpdrp from Windows Services

set /p option="Choose option: "
set Bpath=

IF "%option%" == "1" (
    set /p Bpath="[Path to executable]: "
    sc create mpdrp binPath= "!Bpath!" start= delayed-auto
    sc description mpdrp "A Discord Rich Presence for MPD (https://musicpd.org)"
    sc start mpdrp
) ELSE IF "%option%" == "2" (
    sc stop mpdrp
    sc delete mpdrp
) ELSE (
    echo invalid option, please try again
)
timeout /T 10
cls
goto begin