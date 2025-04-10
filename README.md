# Stonk Trumpet - Truth Social Sentiment Analyzer

This Python application monitors a specified Truth Social account during configured hours, analyzes new posts ("statuses") for potential stock market impact using OpenAI, and sends notifications for statuses deemed to have significant market impact.

![CleanShot 2025-04-10 at 15 02 49@2x](https://github.com/user-attachments/assets/8d53092d-019e-419d-8778-0b6148193cf8)
![CleanShot 2025-04-10 at 15 03 14@2x](https://github.com/user-attachments/assets/8685c426-b41c-43bc-abcc-0c7927c69759)


## Features

- Polls a specific Truth Social user's statuses using the `truthbrush` library.
- Operates only within a configurable time window (defaults to 7 AM - 11 PM ET).
- Analyzes the sentiment of new statuses using OpenAI's API (specifically checking for positive/negative/neutral sentiment regarding stock market impact and whether the impact is significant).
- Sends notifications via ntfy.sh for statuses classified as having significant positive or negative market impact.
- Persistently tracks the last processed status ID to avoid reprocessing, ensuring continuity across restarts.
- Configurable via environment variables or a `.env` file.
- Modular architecture with clear separation of concerns.

## Project Structure

The application has been modularized into the following components:

- `config.py` - Configuration management
- `persistence.py` - Handling of status ID persistence
- `api_client.py` - Truth Social API interactions
- `sentiment_analyzer.py` - OpenAI-based sentiment analysis
- `notifier.py` - Notification service
- `scheduler.py` - Time-based scheduling logic
- `processor.py` - Core business logic for processing statuses
- `main.py` - Application entry point

## Configuration

Create a `.env` file in the project root directory with the following variables:

```dotenv
TRUTHSOCIAL_USERNAME=your_truthsocial_username
TRUTHSOCIAL_PASSWORD=your_truthsocial_password
OPENAI_API_KEY=your_openai_api_key # Required for sentiment analysis
NTFY_TOPIC=your_ntfy_topic         # Required for notifications

# Optional: Override defaults
# TARGET_HANDLE=some_other_handle  # Default: realDonaldTrump
# POLL_START_HOUR=8                # Default: 7 (7 AM ET)
# POLL_END_HOUR=22                 # Default: 23 (11 PM ET)
# POLL_INTERVAL_SECONDS=600        # Default: 300 (5 minutes)
```

- Get an OpenAI API key from [https://openai.com/](https://openai.com/).
- Ensure your Truth Social account is active.
- Choose a unique ntfy.sh topic for receiving notifications.

## Notifications with ntfy.sh

This application uses [ntfy.sh](https://ntfy.sh/) for sending push notifications when significant market impacts are detected. ntfy.sh is a simple HTTP-based pub-sub notification service that:

- Requires no registration
- Works with any device that can receive web push notifications
- Can be used with the ntfy mobile app or web browser

To receive notifications:

1. Set a unique `NTFY_TOPIC` in your `.env` file
2. Subscribe to your topic via:
   - Web browser: Visit `https://ntfy.sh/YOUR_TOPIC`
   - Mobile app: Download the ntfy app from [Google Play](https://play.google.com/store/apps/details?id=io.heckel.ntfy) or [App Store](https://apps.apple.com/us/app/ntfy/id1625396347) and subscribe to your topic

Notifications include:

- Market impact type (positive/negative)
- Username of the poster
- Timestamp
- Content snippet
- AI reasoning for the market impact assessment

## Running the Application

1. Install dependencies:

```bash
pip install -r requirements.txt
```

2. Run the application:

```bash
python main.py
```

The application will start logging its activity to the console and will send notifications via ntfy.sh when significant market impacts are detected.

To stop the application gracefully, press `Ctrl+C`. This will save the last processed ID before exiting.

## How Persistence Works

The application stores the ID of the most recent status it has processed in the `last_processed_id.txt` file. When starting, it reads this ID. When fetching new statuses, it only processes statuses with an ID greater than the stored one (assuming IDs are chronologically sortable).

**Important Note:** This method assumes status IDs are reliably sequential or sortable. It also means if the script is stopped for a long time, it might miss statuses if more than the API's page limit (typically 20-40) were posted during the downtime.

## Customization

- **Target User:** Change `TARGET_HANDLE` in your `.env` file.
- **Polling Window:** Adjust `POLL_START_HOUR` and `POLL_END_HOUR` (0-23) in `.env`.
- **Polling Frequency:** Change `POLL_INTERVAL_SECONDS` in `.env`.
- **Notification Topic:** Change `NTFY_TOPIC` in `.env` to use a different notification channel.
- **Analysis Model:** Change the `model` parameter in the `SentimentAnalyzer.analyze` method (e.g., to `"gpt-4"`) if desired.
- **Analysis Prompt:** Refine the `system_prompt` in `SentimentAnalyzer.analyze` to tailor the analysis criteria.
