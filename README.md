# Terminal Video Viewer

Welcome to the Terminal Video Viewer! This project is a Go application designed to view videos directly in your terminal using ANSI colors and ffmpeg. Created as a learning exercise, this application aims to provide a hands-on experience with Go programming and terminal-based UIs.

## Preview

<video width="640" controls>
  <source src="assets/preview.mp4" type="video/mp4">
</video>

## Features

- **Download Videos**: Download videos from YouTube or load local video files.
- **Video Frame Extraction**: Extract video frames using ffmpeg.
- **Rendering**: Render video frames in the terminal.
- **Audio Playback**: Play audio with controls for play, pause, seek forward, seek backward, and mute.
- **Progress Bars**: Visual progress bars for downloading, processing tasks and watching the video.
- **User Input Handling**: Manage playback with keyboard inputs.

## Installation

To get started with Terminal Video Viewer, ensure you have Go installed on your system, along with ffmpeg and dependencies for audio playback.

1. **Clone the repository:**

   ```sh
   git clone https://github.com/victor247k/TerminalVideoViewer.git
   ```

2. **Navigate to the project directory:**

   ```sh
   cd TerminalVideoViewer
   ```

3. **Install dependencies:**

   ```sh
   go mod tidy
   ```

4. **Build the project:**

   ```sh
   make build
   ```

5. **Run the application:**

   ```sh
   ./bin/main.out
   ```

## Usage

1. **Launch the application:**

   ```sh
   ./bin/main.out
   ```

2. **Choose an option from the menu:**
   - **YouTube Option**: Paste a YouTube video link to download and view.
   - **Local Option**: Paste the path to a local video file to view it.

3. **Follow the on-screen prompts** to download the video or prepare the local file.

4. **Watch the video** rendered in your terminal with audio playback.

## Functionality

- **Download from YouTube**: The application downloads video and audio streams concurrently with a progress bar.
- **Local Video**: Copy a local video file to a temporary location for processing.
- **Extract Frames**: Extract video frames using ffmpeg with a progress bar.
- **Render Frames**: Render extracted frames in the terminal. Frames are preloaded and displayed using ANSI colors.
- **Audio Playback**: Manage audio playback with controls for pausing, seeking, and adjusting volume.
- **Keyboard Input**: Use keyboard shortcuts to control playback:
  - **Space**: Toggle play/pause.
  - **Z**: Seek backward.
  - **X**: Seek forward.
  - **M**: Mute/unmute volume.
  - **Q**: Quit the application.

## Dependencies

- [ffmpeg](https://ffmpeg.org) - For video frame extraction.
- [oto](https://pkg.go.dev/github.com/ebitengine/oto) - For audio playback.
- [bubbletea](https://github.com/charmbracelet/bubbletea) - For terminal UI components.
- [keyboard](https://pkg.go.dev/github.com/eiannone/keyboard) - For handling keyboard input.

## Code Overview

### `cmd/main/main.go`

The entry point of the application, handling the menu, video downloading, frame extraction, and rendering.

### `internal/audio/audio.go`

Handles audio playback with capabilities for play, pause, seek, and mute. It uses the `oto` library for audio processing.

### `internal/render/render.go`

Manages the rendering of video frames in the terminal. Frames are preloaded and processed into ANSI color codes for display. Includes a progress bar for rendering progress.

### `internal/input/input.go`

Handles keyboard input for controlling playback. Implements debounce logic to manage user inputs efficiently.

### `internal/download/download.go`

The download.go module is responsible for handling the download and processing of video and audio files from YouTube or a local file path.

### `internal/extractvideoframes/extractvideoframes.go`

The extractvideoframes package is responsible for extracting frames from a video file and managing the extracted frames. It uses ffmpeg for the extraction process and provides functions to count and clean up the frames.

### `internal/menu/menu.go`

The menu package provides a terminal-based menu for selecting the video source. It uses the bubbletea library for terminal-based user interfaces and allows the user to choose between watching a video from YouTube or a local file.

### `internal/message/message.go`

The message package provides a terminal-based message screen that appears before playing a video. It displays instructions for the user to set up their terminal and provides key control information.

### `internal/progressbar/progressbar.go`

The progressbar package provides a customizable progress bar for terminal applications, using the bubbletea and lipgloss libraries. It displays progress and updates based on the specified percentage.

## Contributing

If you’d like to contribute to this project, please follow these steps:

1. Fork the repository.
2. Create a new branch.
3. Make your changes.
4. Submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Charmbracelet Bubble Tea](https://github.com/charmbracelet/bubbletea) - For terminal UI components.
- [Oto](https://pkg.go.dev/github.com/ebitengine/oto) - For audio playback.
- [Keyboard](https://pkg.go.dev/github.com/eiannone/keyboard) - For handling keyboard input.
- [CSharpTeoMan911](https://github.com/CSharpTeoMan911) - For helping me solve some of the issues I have faced.