# Spotify Loop

A minimal tool to loop any section of your currently playing Spotify track.

Select a start and end point, and the app will continuously loop that segment.

**Features:**
* View currently playing Spotify track with album artwork
* Visual range slider to select loop points
* Real-time loop playback control
* Clean, responsive interface

**Technologies:**
* Go
* JavaScript
* HTML
* CSS

**Prerequisites:**
* [Go](https://go.dev/dl/) (version 1.21 or higher)
* Spotify Premium account

## Installation

### 1. Create a Spotify App

1. Go to the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard) to create your app. For more detailed instructions, follow the official guide from Spotify on [how to create an app](https://developer.spotify.com/documentation/web-api/concepts/apps).
2. When creating the Spotify app, set the Redirect URI to `http://127.0.0.1:8080/api/spotify/callback` if you want this tool to work out of the box with no further configuration. If you decide to host this app on a different/production domain, you'll have to set the Redirect URI to match. If you want to change the development server port, you can do so in `/cmd/main.go`.
3. After you have created the Spotify app, take note of the Client ID and Redirect URI. You will need these in the next step.

### 2. Configure Environment Variables
Create a .env file in the project root and configure it as follows:
```env
CLIENT_ID=your_client_id_here
REDIRECT_URI=http://127.0.0.1:8080/api/spotify/callback
SCOPE=user-read-playback-state user-modify-playback-state
```
You can add any additional scopes if needed, however, the ones provided above are all that's needed for this tool to run. [Read more about scopes here](https://developer.spotify.com/documentation/web-api/concepts/scopes).

### 3. Install Dependencies
```bash
go mod download
```

### 4. Run the Application
```bash
go run cmd/*.go
```

The app will start on `http://127.0.0.1:8080` (or the port specified in `/cmd/main.go`).

## Usage

1. Open `http://127.0.0.1:8080` in your browser.
2. Click **Connect to Spotify** and authorize the app.
3. Start playing a song on any Spotify device (this tool communicates directly with the Spotify Web API, so it works with any device).
4. Refresh the page to see the dashboard displaying your currently playing song.
5. Use the range sliders to select your loop points.
6. Click **Start loop** to begin looping.
7. The status indicator in the top right corner displays your current state: unauthorized, idle (no active playback), playing, or looping.

## Notes
* You must have an active Spotify session (playing on any device) to use the loop feature.
* The app requires periodic authorization renewal as Spotify tokens expire.
* This is designed for self-hosting and personal use.
* This app uses Spotify's PKCE OAuth flow, which does not require a Client Secret. Only the Client ID is needed.
