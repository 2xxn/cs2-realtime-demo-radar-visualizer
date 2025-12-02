# cs2-realtime-demo-radar-visualizer

CS2 Demo-Based Real-Time Radar Visualizer

This application visualizes player positions in real-time on a radar by reading data from a CS2 demo file.

![Preview](https://github.com/2xxn/cs2-realtime-demo-radar-visualizer/blob/main/assets/preview.png)

## How it works
Demo buffer in CS2 is written to the file very often (once every few seconds). This application reads the demo file and extracts the player data from it in real-time, then visualizes it similarly to in-game radar.

This is only a PoC though, so it may not work perfectly in all cases.

## Requirements
- Go 1.24.5

## Installation

```bash
git clone https://github.com/2xxn/cs2-realtime-demo-radar-visualizer
cd cs2-realtime-demo-radar-visualizer
go build -o radar-visualizer main.go
```

## Usage
```bash
./radar-visualizer
```

Once it prompts you for the demo file path, provide the path to your CS2 demo file (in your CS2 dir, e.g., `C:/SteamLibrary/steamapps/common/Counter-Strike Global Offensive/game/csgo/demo.dem`).

When it prompts you for the map name, provide the map name (e.g., `de_dust2`).

Then, start a CS2 game and record a demo using the command:
```bash
record demo
```
Once you confirm by pressing Enter, the application will start visualizing the radar in real-time on [localhost:8080](http://localhost:8080).

## Showcase Preview
![Preview](https://github.com/2xxn/cs2-realtime-demo-radar-visualizer/blob/main/assets/preview.mp4)
<video src="assets/preview.mp4" controls></video>

## Disclaimer
Do not use this tool in competitive matches or in any way that violates the game's terms of service. While undetectable by the game, this tool is intended for educational and personal use only.