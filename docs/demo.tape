Output docs/demo.gif

Set Shell "bash"
Set FontSize 16
Set Width 1000
Set Height 650
Set Theme "Dracula"
Set TypingSpeed 60ms

Env BOBA_FILE "/tmp/boba-demo.json"

# seed a block covering "now" (hidden from the recording)
Hide
Type `rm -f /tmp/boba-demo.json; H=$((10#$(date +%H))); boba add ${H}:00-${H}:59 "Deep work"`
Enter
Type "clear"
Enter
Show

Type "boba"
Enter
Sleep 4s

Type "f"
Sleep 2s

Type "m"
Sleep 3s

Escape
Sleep 1s

Type "a"
Sleep 500ms
Type "22-23 Read"
Enter
Sleep 2s

Type "q"
Sleep 500ms
