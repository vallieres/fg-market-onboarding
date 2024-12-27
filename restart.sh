#!/bin/bash

set -o pipefail
IFS=$'\n\t'

source $HOME/.bashrc
cd "$HOME/apps/$FGONBOARDING_APP_FOLDER"

fgonboardPID=$(cat $HOME/apps/$FGONBOARDING_APP_FOLDER/fgonboard.pid)
echo "→ Killing Running FG Onboarding with PID $fgonboardPID"
pkill -F $HOME/apps/$FGONBOARDING_APP_FOLDER/fgonboard.pid

sleep 1

echo "→ Relaunching FGONBOARDIND..."
nohup ./fgonboarding >>$HOME/logs/apps/$FGONBOARDING_APP_FOLDER/cron.log 2>&1 </dev/null &
disown
echo $! > $HOME/apps/FGONBOARDING_APP_FOLDER/fgonboard.pid

sleep 2

fgonboardPID=$(cat $HOME/apps/$FGONBOARDING_APP_FOLDER/fgonboard.pid)
echo -e "→ New FGONBOARDIND PID: $fgonboardPID \n"

exit 0
