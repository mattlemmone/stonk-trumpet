# Stonk Trumpet - Truth Social Sentiment Analyzer

A Python application that monitors a specified Truth Social account for new posts, analyzes them for potential stock market impact using OpenAI, and sends push notifications when significant market impact is detected.

![CleanShot 2025-04-10 at 15 02 49@2x](https://github.com/user-attachments/assets/8d53092d-019e-419d-8778-0b6148193cf8)
![CleanShot 2025-04-10 at 15 03 14@2x](https://github.com/user-attachments/assets/8685c426-b41c-43bc-abcc-0c7927c69759)

## Features

- Monitors a specified Truth Social account using the `truthbrush` library
- Operates only within configurable hours (default: 7 AM - 11 PM in your specified timezone)
- Analyzes post sentiment using OpenAI to identify market impact potential
- Sends push notifications via ntfy.sh for posts with significant market impact
- Maintains state across restarts by tracking the last processed post
- Fully configurable via environment variables or `.env` file

## Setup

1. Create a `.env` file with your configuration (see example below)
2. Install dependencies: `pip install -r requirements.txt`
3. Run the application: `python main.py`

## Notifications

This application uses [ntfy.sh](https://ntfy.sh/) for push notifications. Subscribe to your configured topic in the ntfy.sh mobile app or web browser to receive alerts.
