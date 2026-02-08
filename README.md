<img width="650" height="770" alt="Screenshot 2026-02-08 at 08 46 35" src="https://github.com/user-attachments/assets/941566cd-1318-4d6e-bf3a-bc71a0f5954e" />

## Single-task, single-user task manager

This tool is based on how I personally track time. Most existing time-tracking tools feel too opinionated and come with many extras I don’t need. I also don’t like Pomodoro-style workflows.

When I’m done with what I’m working on, I want to stop the task immediately. Fixed time frames don’t work well for me, especially with frequent context switching.

## Install 

### macOS (prebuild binary)

Download the macOS binary from the release page.

### Go Install

go install github.com/fmo/cmd/timer-cli@v1.x.y

### Other platforms

Requires Go installed.

```
git clone git@github.com:fmo/timer-cli.git
cd timer-cli
make build
./timer-cli
```

## Features

* Add, stop, show time for the running task.
* Add manual time for missing time blocks.
