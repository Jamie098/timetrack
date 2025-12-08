# TimeTrack

A powerful command-line time tracking tool built in Go. Track your workday in hours, visualize with colors, and export to CSV for reporting.

## Features

- **Interactive Mode** - User-friendly menu when running without arguments
- **Fuzzy Matching** - Smart project name matching (e.g., "ctg" matches "CT.GOV Automation")
- **Color-Coded Display** - Visual feedback with green/yellow/red indicators
- **Smart Validation** - Automatic warnings for over-allocation or unusual values
- **Multiple Export Formats** - CSV, JSON, weekly or full history
- **CSV Import** - Import existing timesheet data
- **Comprehensive Reports** - Weekly summaries, project-specific reports, and statistics
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

### 1. Configure Your Projects

Parse your Excel header to automatically set up projects and aliases:

```bash
timetrack projects parse "Date,CT.GOV Automation,Bugs & Issues,Feature Development,Total Time Spent"
```

This will:
- Extract project names from the header
- Auto-generate short aliases (e.g., "ctgo", "bugs", "feat")
- Configure your project columns for export

### 2. Set Up Recurring Meetings

Add ceremonies that happen regularly:

```bash
timetrack meeting add standup 0.5 weekdays
timetrack meeting add retro 1 fri
timetrack meeting add planning 2 mon
```

**Note:** All time input is in **hours** (based on 8-hour workday).

### 3. Track Your Time

**Interactive mode** (easiest for beginners):
```bash
timetrack
```

**Command-line mode**:
```bash
timetrack add ctgo 2        # Add 2 hours to CT.GOV Automation
timetrack add bugs 1.5      # Add 1.5 hours to Bugs & Issues
timetrack add feat 3        # Add 3 hours to Feature Development
```

Fuzzy matching works automatically:
```bash
timetrack add automation 2  # Matches "CT.GOV Automation"
timetrack add bug 1         # Matches "Bugs & Issues"
```

### 4. View Your Status

```bash
timetrack           # Interactive mode shows current status
timetrack summary   # Compact one-line view
```

### 5. Export for Reporting

```bash
timetrack week              # Print current week as CSV (stdout)
timetrack export csv        # Export week to timetrack-week.csv
timetrack export all        # Export all data to timetrack-all.csv
timetrack export json       # Export as JSON
```

## Commands

### Basic Tracking

```bash
timetrack                        # Interactive mode
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
timetrack summary              # Compact one-line status
timetrack history [days]       # Show history (default: 7 days)
```

### Reports & Analytics

```bash
timetrack report week             # Weekly summary
timetrack report project <name>   # Project-specific report
timetrack report stats            # Overall statistics
```

### Import/Export

```bash
timetrack import <file.csv>       # Import from CSV
timetrack export csv [file]       # Export current week to CSV
timetrack export all [file]       # Export all data to CSV
timetrack export json [file]      # Export as JSON
timetrack week                    # Print week to stdout (legacy)
```

### Project Configuration

```bash
timetrack projects parse "Date,P1,P2,...,Total"  # Parse Excel header
timetrack projects set "P1,P2,P3"                # Set manually
timetrack projects list                          # List projects
timetrack alias <short> <full>                   # Create alias
timetrack alias rm <short>                       # Remove alias
timetrack alias list                             # List all aliases
```

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
timetrack

# Throughout day - track work
timetrack add automation 2
timetrack add bugs 1.5
timetrack add features 3

# Accidentally added wrong time
timetrack undo
timetrack add features 2.5

# Check current status
timetrack summary

# End of day
timetrack                    # Verify all time tracked
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
# View weekly summary
timetrack report week

# Export to CSV for manager
timetrack export csv weekly-report.csv

# Check specific project
timetrack report project "CT.GOV Automation"
```

### Setup New Project List

```bash
# Copy header row from Excel: "Date,ProjectA,ProjectB,ProjectC,Total Time Spent"
timetrack projects parse "Date,ProjectA,ProjectB,ProjectC,Total Time Spent"

# Check generated aliases
timetrack alias list

# Customize if needed
timetrack alias a "Project Alpha"
timetrack alias b "Project Beta"
```

## Tips

1. **Use Interactive Mode**: When starting out, run `timetrack` without arguments for guided interface
2. **Fuzzy Match**: You don't need exact project names - "auto", "automation", "ctg" all work
3. **Track as You Go**: Add time throughout the day rather than at the end
4. **Use Fill for Efficiency**: If most of your day is on one project, use `timetrack fill <project>` instead of calculating hours
5. **Backdate When Needed**: Forgot yesterday? Use `--date` flag: `timetrack add project 5 --date 2024-12-07`
6. **Store Timesheet URL**: Set up your online timesheet URL once with `timetrack url set <url>` and open it anytime with `timetrack url open`
7. **Check History**: Use `timetrack history` to review past weeks and spot missing days
8. **Set Reminders**: Use `timetrack start-bg` to get notifications at set times
9. **Export Weekly**: Run `timetrack export csv` every Friday for weekly reports

## Troubleshooting

### "Project not found"
- Use `timetrack projects list` to see configured projects
- Use `timetrack alias list` to see available aliases
- Fuzzy matching may suggest alternatives

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
