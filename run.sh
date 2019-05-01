#!/bin/bash
# Shell script to open terminals
# and execute a separate command in each

# Commands to run (one per terminal)
cmds=('go run main.go 6688', 'go run main.go 6689', 'go run main.go 6690')

# Loop through commands, open terminal, execute command
for i in "${cmds[@]}"
do
    xterm -e "$i && /bin/tcsh" &
done