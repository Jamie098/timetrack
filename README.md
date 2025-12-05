# Get started
1. Build the .exe using go build -o timetrack.exe .
2. Add the path to the directory containing the .exe to your PATH
3. Test it is working with timetrack help
4. Make it auto-start or run timetrack start-bg if you want notifications.

# Configure
Check out the Config section in the help to set up things like recurring meetings.

# Commands Help
```
Usage:
  timetrack                      Show today's status
  timetrack add <project> <%>    Add/update time to a project
  timetrack exclude <name> <%>   Exclude ceremony time (one-off)
  timetrack rm <project>         Remove a project entry
  timetrack rmex <name>          Remove an excluded meeting
  timetrack clear                Clear today's data
  timetrack history [days]       Show history (default: 7 days)

Config:
  timetrack config               Show current config
  timetrack config edit          Open config file in editor
  timetrack meeting add <name> <percent> <days>   Add recurring meeting
  timetrack meeting rm <name>    Remove recurring meeting
  timetrack reminder <times>     Set reminder times (e.g., "09:00,12:00,15:00")

Reminders:
  timetrack start                Start reminder service (foreground)
  timetrack start-bg             Start reminder service (background)
  timetrack stop                 Stop reminder service
  timetrack status               Check if reminder service is running

Examples:
  timetrack meeting add standup 6.25 weekdays
  timetrack meeting add retro 12.5 fri
  timetrack reminder 09:00,12:00,15:30,17:00
  timetrack add clinical-trials 25

Days: mon, tue, wed, thu, fri, sat, sun, daily, weekdays

Quick reference (8hr day):
  15min = 3.125%    30min = 6.25%    45min = 9.375%
  1hr   = 12.5%     1.5hr = 18.75%   2hr   = 25%
```
