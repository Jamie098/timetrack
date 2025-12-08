# TimeTrack

A powerful command-line time tracking tool built in Go. Track your workday in hours, visualize with colors, and export to CSV for reporting.

## Features

- **Simple Default View** - Running `timetrack` shows today's status
- **Calendar View** - Visualize tracked time across multiple days in a clean table format
- **Auto-Discovery** - Projects automatically populate from your time entries (no manual setup needed)
- **Fuzzy Matching** - Smart project name matching (e.g., "ctg" matches "CT.GOV Automation")
- **Color-Coded Display** - Visual feedback with green/yellow/red indicators
- **Smart Validation** - Automatic warnings for over-allocation or unusual values
- **Multiple Export Formats** - CSV, JSON, weekly or full history with alphabetical project ordering
- **CSV Import** - Import existing timesheet data
- **Interactive Mode** - Optional user-friendly menu mode
- **Recurring Meetings** - Auto-exclude ceremony time on specific days
- **Desktop Notifications** - Optional reminder service
- **Project Aliases** - Short names for long project titles
- **Auto-Fill Remaining Time** - Quickly assign all untracked time to a project
- **Timesheet URL Management** - Store and open your online timesheet with one command
- **Backdate Entries** - Edit time for any previous day using the `--date` flag

## Installation

1. Build the executable:
   ```bash
   go build -o timetrack .
   ```

2. Add to your PATH:
   - **Linux/Mac**: Move to `/usr/local/bin/` or add current directory to PATH
   - **Windows**: Build with `go build -o timetrack.exe .` and add directory to PATH

3. Test installation:
   ```bash
   timetrack help
   ```

## Quick Start

### 1. Set Up Recurring Meetings (Optional)

Add ceremonies that happen regularly:

```bash
timetrack meeting add standup 0.5 weekdays
timetrack meeting add retro 1 fri
timetrack meeting add planning 2 mon
```

**Note:** All time input is in **hours** (based on 8-hour workday).

### 2. Track Your Time

**Command-line mode**:
```bash
timetrack add ctgo 2        # Add 2 hours to CT.GOV Automation
timetrack add bugs 1.5      # Add 1.5 hours to Bugs & Issues
timetrack add feat 3        # Add 3 hours to Feature Development
```

**Interactive mode** (optional menu interface):
```bash
timetrack interactive
```

Fuzzy matching works automatically:
```bash
timetrack add automation 2  # Matches "CT.GOV Automation"
timetrack add bug 1         # Matches "Bugs & Issues"
```

### 3. View Your Status

```bash
timetrack              # Show today's status
timetrack show         # Calendar view of last 7 days
timetrack show 14      # Calendar view of last 14 days
```

### 4. Export for Reporting

```bash
timetrack export csv        # Export week to timetrack-week.csv (alphabetical columns)
timetrack export all        # Export all data to timetrack-all.csv
timetrack export json       # Export as JSON
```

Projects are automatically discovered from your time entries and exported in alphabetical order.

## Commands

### Basic Tracking

```bash
timetrack                        # Show today's status
timetrack interactive            # Interactive menu mode
timetrack add <project> <hours>  # Add/update time
timetrack fill <project>         # Fill remaining time with project
timetrack edit <project> <hours> # Update existing entry
timetrack exclude <name> <hours> # Exclude ceremony time (one-off)
timetrack rm <project>           # Remove project
timetrack undo                   # Remove last entry
timetrack clear                  # Clear today's data

# Work with past dates (add, fill, edit, rm support --date flag)
timetrack add <project> <hours> --date 2024-12-05
timetrack fill <project> -d 12/05/2024
timetrack edit <project> <hours> --date 2024-12-05
timetrack rm <project> --date 2024-12-05
```

### Viewing Data

```bash
timetrack                  # Show today's status
timetrack show [days]      # Calendar view (default: 7 days)
timetrack show 14          # Show last 14 days
```

### Import/Export

```bash
timetrack import <file.csv>       # Import from CSV
timetrack export csv [file]       # Export current week to CSV (auto-discovered projects)
timetrack export all [file]       # Export all data to CSV
timetrack export json [file]      # Export as JSON
```

Projects are automatically discovered from your tracked time and exported in alphabetical order.

### Project Management

```bash
timetrack projects list            # List all projects (auto-discovered from time entries)
timetrack alias <short> <full>     # Create alias for quick entry
timetrack alias rm <short>         # Remove alias
timetrack alias list               # List all aliases
```

**Note:** Projects are auto-discovered - they automatically appear when you track time to them. No need to manually configure project lists!

### Recurring Meetings

```bash
timetrack meeting add <name> <hours> <days>   # Add recurring meeting
timetrack meeting rm <name>                   # Remove meeting
```

**Days**: `mon`, `tue`, `wed`, `thu`, `fri`, `sat`, `sun`, `daily`, `weekdays`

### Configuration

```bash
timetrack config          # Show current config
timetrack config edit     # Edit config file
```

### Timesheet URL

```bash
timetrack url set <url>   # Set your online timesheet URL
timetrack url open        # Open timesheet in browser
timetrack url             # Show current URL
timetrack url rm          # Clear stored URL
```

### Reminders (Optional)

```bash
timetrack reminder 09:00,12:00,15:00  # Set reminder times
timetrack start                       # Start service (foreground)
timetrack start-bg                    # Start service (background)
timetrack stop                        # Stop service
timetrack status                      # Check if running
```

## Visual Indicators

The tool uses color coding for quick status understanding:

- **Green** - Good allocation, on track
- **Yellow** - Getting close to limit (< 10% remaining)
- **Red** - Over-allocated
- **Blue** - Tracked time values
- **Cyan** - Available time
- **Gray** - Excluded/ceremony time

## Smart Features

### Fuzzy Matching

The tool automatically matches project names:

```bash
timetrack add ctg 2        # Matches "CT.GOV Automation"
timetrack add auto 1.5     # Matches "CT.GOV Automation"
timetrack add automation 2 # Matches "CT.GOV Automation"
```

### Automatic Validation

- **Over 8 hours warning**: "⚠️ Warning: 10.00 hours is 125.0% of an 8-hour day (>100%). Did you mean 1.00 hours?"
- **Over-allocation**: "⚠️ Warning: Over-allocated by 25.0%!"
- **Near completion**: "💡 Only 6.2% remaining - almost done!"

### Fill Remaining Time

The `fill` command automatically assigns all remaining time to a project:

```bash
timetrack fill "Project Name"    # Fill all remaining time
```

Perfect for days where most time goes to one project. If you have a 5% recurring meeting, `fill` will automatically assign the remaining 95% (7.6 hours) to your project.

### Backdating Entries

You can add, edit, fill, or remove time for any previous day using the `--date` or `-d` flag:

```bash
# Add time to a past date
timetrack add "Project" 3 --date 2024-12-05
timetrack add "Project" 2 -d 12/05/2024

# Fill remaining time for yesterday
timetrack fill "Main Project" --date 2024-12-07

# Edit existing entry from last week
timetrack edit "Bugs" 1.5 --date 2024-12-01

# Remove incorrect entry from past
timetrack rm "Wrong Project" --date 2024-12-05
```

Supported date formats:
- `YYYY-MM-DD` (2024-12-05)
- `YYYY/MM/DD` (2024/12/05)
- `MM/DD/YYYY` (12/05/2024)
- `MM-DD-YYYY` (12-05-2024)

### Edit vs Add

- `add` creates or overwrites entries
- `edit` only updates existing entries (safer)
- `fill` assigns all remaining time to a project
- `undo` removes the most recent entry

## Data Storage

All data is stored in JSON format:

- **Config**: `~/.config/timetrack/config.json` (Linux/Mac) or `%APPDATA%\timetrack\config.json` (Windows)
- **Data**: `~/.local/share/timetrack/data.json` (Linux/Mac) or `%APPDATA%\timetrack\data.json` (Windows)

## Import/Export Format

### CSV Format

The tool expects/exports CSV in this format:

```csv
Date,Project1,Project2,Project3,Total Time Spent
2-Jan,25.0%,12.5%,37.5%,
3-Jan,50.0%,25.0%,12.5%,
```

- First column: Date (formats: `2-Jan`, `2-Jan-06`, `2006-01-02`, `01/02/2006`)
- Middle columns: Project percentages
- Last column: Total (optional, auto-calculated)

### JSON Format

```json
{
  "2024-12-06": {
    "date": "2024-12-06",
    "excluded_percent": 6.25,
    "projects": {
      "CT.GOV Automation": 25.0,
      "Bugs & Issues": 12.5
    },
    "excluded_meetings": {
      "standup": 6.25
    }
  }
}
```

## Examples

### Daily Workflow

```bash
# Morning - check status
timetrack                    # Shows today's status

# Throughout day - track work
timetrack add automation 2
timetrack add bugs 1.5
timetrack add features 3

# Accidentally added wrong time
timetrack undo
timetrack add features 2.5

# End of day - check what you've tracked
timetrack                    # Verify all time tracked
timetrack show               # See the week at a glance
```

### Quick Fill Workflow

```bash
# Morning - spent some time in meetings and bug fixes
timetrack add bugs 1.5

# Rest of the day on main project - use fill!
timetrack fill "Main Project"    # Automatically fills remaining 6.1 hours

# Or if you have recurring meetings (e.g., 5% standup)
timetrack fill "Main Project"    # Fills 95% - 18.75% = 76.25% (6.1 hours)
```

### Timesheet URL Setup

```bash
# One-time setup - store your online timesheet URL
timetrack url set "https://company.sharepoint.com/timesheets/mysheet.xlsx"

# Anytime you need to open it
timetrack url open               # Opens in default browser

# Check what's stored
timetrack url                    # Shows current URL
```

### Backfilling Previous Week

```bash
# Forgot to track Monday? No problem!
timetrack add "Project A" 4 --date 2024-12-02
timetrack add "Project B" 3 --date 2024-12-02

# Or use fill to complete the day
timetrack fill "Main Project" --date 2024-12-02

# Update an incorrect entry from last week
timetrack edit "Meetings" 1 --date 2024-12-01
```

### Weekly Reporting

```bash
# View calendar for the week
timetrack show

# View last 14 days
timetrack show 14

# Export to CSV for manager
timetrack export csv weekly-report.csv
```

### Setting Up Aliases

```bash
# Create shortcuts for frequently-used projects
timetrack alias auto "CT.GOV Automation"
timetrack alias bugs "Bugs & Issues"
timetrack alias feat "Feature Development"

# Check your aliases
timetrack alias list

# Now you can use short aliases
timetrack add auto 3
timetrack add bugs 1.5
```

## Tips

1. **Quick Status Check**: Just run `timetrack` to see today's status
2. **Calendar View**: Use `timetrack show` to see the week at a glance
3. **Auto-Discovery**: No need to set up projects - they appear automatically when you track time
4. **Fuzzy Match**: You don't need exact project names - "auto", "automation", "ctg" all work
5. **Track as You Go**: Add time throughout the day rather than at the end
6. **Use Fill for Efficiency**: If most of your day is on one project, use `timetrack fill <project>` instead of calculating hours
7. **Backdate When Needed**: Forgot yesterday? Use `--date` flag: `timetrack add project 5 --date 2024-12-07`
8. **Store Timesheet URL**: Set up your online timesheet URL once with `timetrack url set <url>` and open it anytime with `timetrack url open`
9. **Interactive Mode Available**: Run `timetrack interactive` for a guided menu interface
10. **Set Reminders**: Use `timetrack start-bg` to get notifications at set times
11. **Export Weekly**: Run `timetrack export csv` every Friday for weekly reports

## Troubleshooting

### "Project not found"
- Use `timetrack projects list` to see all projects you've tracked
- Use `timetrack alias list` to see available aliases
- Fuzzy matching will suggest alternatives
- Don't worry - if no match is found, your input will be used as a new project name

### Over-allocated warnings
- You've tracked more than 8 hours (100%)
- Review your entries with `timetrack` (shows current day)
- Use `timetrack edit <project> <hours>` to fix

### Import fails
- Ensure CSV has header row
- Date format should be recognizable (2-Jan, 2006-01-02, etc.)
- Percentages can include or omit the % symbol

## License

MIT License - feel free to modify and distribute.

## Contributing

This is a personal project, but suggestions and improvements are welcome!
