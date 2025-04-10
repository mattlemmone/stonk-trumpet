# Stonk Trumpet - Truth Social Sentiment Analyzer

This Python application monitors a specified Truth Social account during configured hours, analyzes new posts ("statuses") for potential stock market impact using OpenAI, and logs notifications for statuses deemed to have significant positive impact.

## Features

- Polls a specific Truth Social user's statuses using the `truthbrush` library.
- Operates only within a configurable time window (defaults to 7 AM - 11 PM ET).
- Analyzes the sentiment of new statuses using OpenAI's API (specifically checking for positive/negative/neutral sentiment regarding stock market impact and whether the impact is significant).
- Sends notifications (currently logged to the console) for statuses classified as having a significant positive impact.
- Persistently tracks the last processed status ID to avoid reprocessing, ensuring continuity across restarts.
- Configurable via environment variables or a `.env` file.

## Setup

1.  **Clone the repository (or create the files):**

    ```bash
    # If you have a repo:
    # git clone <your-repo-url>
    # cd <your-repo-dir>
    ```

2.  **Create a Python virtual environment (recommended):**

    ```bash
    python -m venv venv
    source venv/bin/activate # On Windows use `venv\Scripts\activate`
    ```

3.  **Install dependencies:**

    ```bash
    pip install -r requirements.txt
    ```

4.  **Create a `.env` file:**
    Create a file named `.env` in the project root directory and add your credentials and any desired configuration overrides:

    ```dotenv
    TRUTHSOCIAL_USERNAME=your_truthsocial_username
    TRUTHSOCIAL_PASSWORD=your_truthsocial_password
    OPENAI_API_KEY=your_openai_api_key # Required for sentiment analysis

    # Optional: Override defaults
    # TARGET_HANDLE=some_other_handle # Default: realDonaldTrump
    # POLL_START_HOUR=8                 # Default: 7 (7 AM ET)
    # POLL_END_HOUR=22                  # Default: 23 (11 PM ET)
    # POLL_INTERVAL_SECONDS=600         # Default: 300 (5 minutes)
    ```

    - Get an OpenAI API key from [https://openai.com/](https://openai.com/).
    - Ensure your Truth Social account is active.

## Running the Application

Simply run the `main.py` script:

```bash
python main.py
```

The application will start logging its activity to the console. It will check for new statuses from the `TARGET_HANDLE` every `POLL_INTERVAL_SECONDS` seconds, but only perform fetching and analysis between `POLL_START_HOUR` and `POLL_END_HOUR` (Eastern Time).

- New statuses are processed and analyzed.
- The ID of the last processed status is stored in `last_processed_id.txt`.
- Significant positive statuses trigger a notification log message.

To stop the application gracefully, press `Ctrl+C`. This will save the last processed ID before exiting.

## How Persistence Works

The application stores the ID of the most recent status it has processed in the `last_processed_id.txt` file. When starting, it reads this ID. When fetching new statuses, it only processes statuses with an ID greater than the stored one (assuming IDs are chronologically sortable).

**Important Note:** This method assumes status IDs are reliably sequential or sortable. It also means if the script is stopped for a long time, it might miss statuses if more than the API's page limit (typically 20-40) were posted during the downtime. A more robust solution would involve storing all processed IDs or using API features like `since_id` if available and reliable in `truthbrush`.

## Customization

- **Target User:** Change `TARGET_HANDLE` in your `.env` file.
- **Polling Window:** Adjust `POLL_START_HOUR` and `POLL_END_HOUR` (0-23) in `.env`.
- **Polling Frequency:** Change `POLL_INTERVAL_SECONDS` in `.env`.
- **Notification Method:** Modify the `send_notification` function in `main.py` to use services like Twilio for SMS, `smtplib` for email, or other notification platforms instead of just logging.
- **Analysis Model:** Change the `model` parameter in the `analyze_sentiment` function (e.g., to `"gpt-4"`) if desired.
- **Analysis Prompt:** Refine the `system_prompt` in `analyze_sentiment` to tailor the analysis criteria.
