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

<img width="577" height="343" alt="Screenshot 2026-01-28 at 21 59 32" src="https://github.com/user-attachments/assets/6c2090d4-7d86-4fa7-bfc9-c05ebe6c65af" />

## Add Time manually
<img width="575" height="394" alt="Screenshot 2026-01-30 at 12 37 03" src="https://github.com/user-attachments/assets/5d759c6b-5185-45bf-96a1-694c4d676904" />
