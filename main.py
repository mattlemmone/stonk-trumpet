import os
import time
import datetime
import pytz
import json
import logging
import requests
from dotenv import load_dotenv
from truthbrush.api import Api
import openai

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

# Load environment variables
load_dotenv()

# --- Configuration ---
TRUTHSOCIAL_USERNAME = os.getenv("TRUTHSOCIAL_USERNAME")
TRUTHSOCIAL_PASSWORD = os.getenv("TRUTHSOCIAL_PASSWORD")
OPENAI_API_KEY = os.getenv("OPENAI_API_KEY")
TARGET_HANDLE = os.getenv("TARGET_HANDLE", "realDonaldTrump")
NTFY_TOPIC = os.getenv("NTFY_TOPIC")
LAST_PROCESSED_ID_FILE = "last_processed_id.txt"
EASTERN_TIMEZONE = pytz.timezone('US/Eastern')

# Configurable polling window with defaults
try:
    POLL_START_HOUR = int(os.getenv("POLL_START_HOUR", 7))
    POLL_END_HOUR = int(os.getenv("POLL_END_HOUR", 23))
    POLL_INTERVAL_SECONDS = int(os.getenv("POLL_INTERVAL_SECONDS", 300))
except ValueError:
    logging.error("Invalid non-integer value for POLL_START_HOUR, POLL_END_HOUR, or POLL_INTERVAL_SECONDS. Using defaults.")
    POLL_START_HOUR = 7
    POLL_END_HOUR = 23
    POLL_INTERVAL_SECONDS = 300

if not (0 <= POLL_START_HOUR < 24 and 0 <= POLL_END_HOUR <= 24 and POLL_START_HOUR < POLL_END_HOUR):
     logging.error("Invalid POLL_START_HOUR or POLL_END_HOUR configuration. Using defaults.")
     POLL_START_HOUR = 7
     POLL_END_HOUR = 23

# Initialize OpenAI client - required
if OPENAI_API_KEY:
    openai_client = openai.OpenAI(api_key=OPENAI_API_KEY)
else:
    logging.critical("OPENAI_API_KEY is not set. This is a required dependency.")
    raise ValueError("OPENAI_API_KEY environment variable must be set")
# --- End Configuration ---


# --- Persistence ---
def load_last_processed_id():
    """Loads the last processed status ID from a file."""
    last_id = None
    try:
        with open(LAST_PROCESSED_ID_FILE, 'r') as f:
            last_id = f.readline().strip()
        if last_id:
            logging.info(f"Loaded last processed ID {last_id} from {LAST_PROCESSED_ID_FILE}")
        else:
            logging.info(f"{LAST_PROCESSED_ID_FILE} was empty. Starting fresh.")
    except FileNotFoundError:
        logging.info(f"{LAST_PROCESSED_ID_FILE} not found. Starting fresh.")
    except IOError as e:
        logging.error(f"Error loading last processed ID from {LAST_PROCESSED_ID_FILE}: {e}")
    return last_id

def save_last_processed_id(last_id):
    """Saves the last processed status ID to a file."""
    if not last_id:
        return
    try:
        with open(LAST_PROCESSED_ID_FILE, 'w') as f:
            f.write(f"{last_id}\n")
        logging.debug(f"Saved last processed ID {last_id} to {LAST_PROCESSED_ID_FILE}")
    except IOError as e:
        logging.error(f"Error saving last processed ID {last_id} to {LAST_PROCESSED_ID_FILE}: {e}")
# --- End Persistence ---


def is_within_polling_hours():
    """Check if the current time is within the allowed polling hours in ET."""
    now_et = datetime.datetime.now(EASTERN_TIMEZONE)
    return POLL_START_HOUR <= now_et.hour < POLL_END_HOUR

def analyze_sentiment(status_text_html):
    """Analyzes sentiment using OpenAI for stock market impact."""
    log_text = status_text_html[:150] + ('...' if len(status_text_html) > 150 else '')
    logging.info(f"Analyzing HTML content: '{log_text}'")
    sentiment = "neutral"
    significant = False
    reasoning = ""

    system_prompt = """You are an AI assistant specialized in financial market sentiment analysis. 
    Analyze the following social media post (may contain HTML) regarding its potential impact on the stock market. 
    Classify the sentiment as 'positive' (number go up), 'negative' (number go down), or 'neutral' (no impact).
    Also, determine if the potential impact is 'significant' (true or false) and likely to have immediate impact on the market.
    Additionally, provide a brief reasoning (1-2 sentences) explaining why you believe this will influence the broader market.
    
    Respond ONLY with a JSON object containing three keys: 
    - "sentiment" (string): "positive", "negative", or "neutral"
    - "significant" (boolean): true or false
    - "reasoning" (string): brief explanation of market impact
    
    Example: {"sentiment": "positive", "significant": true, "reasoning": "The announcement of tax cuts for businesses could stimulate economic growth and boost investor confidence."}"""

    try:
        response = openai_client.chat.completions.create(
            model="gpt-4o-mini",
            messages=[
                {"role": "system", "content": system_prompt},
                {"role": "user", "content": status_text_html}
            ],
            temperature=1.0,
            max_tokens=150,  # Increased to accommodate reasoning
            response_format={"type": "json_object"}
        )
        content = response.choices[0].message.content
        logging.debug(f"OpenAI raw response: {content}")
        result = json.loads(content)
        sentiment = result.get("sentiment", "neutral").lower()
        significant = result.get("significant", False)
        reasoning = result.get("reasoning", "")

        if sentiment not in ["positive", "negative", "neutral"]:
            logging.warning(f"Received unexpected sentiment value: {sentiment}. Defaulting.")
            sentiment = "neutral"
        if not isinstance(significant, bool):
            logging.warning(f"Received non-boolean significance value: {significant}. Defaulting.")
            significant = False

        logging.info(f"Analysis result: Sentiment='{sentiment}', Significant={significant}, Reasoning='{reasoning}'")

    except json.JSONDecodeError as e:
         logging.error(f"Failed to decode JSON response from OpenAI: {e} - Response: {content}")
    except (openai.APIError, openai.APITimeoutError, openai.RateLimitError, Exception) as e:
        logging.error(f"OpenAI API call failed: {e}")

    return sentiment, significant, reasoning

def send_notification(status, sentiment, reasoning):
    """Sends a notification for a significant market impact using ntfy.sh."""
    raw_content = status.get('content', '')
    content_snippet = raw_content[:200] + ('...' if len(raw_content) > 200 else '')
    username = status.get('account', {}).get('username', 'Unknown')
    timestamp = status.get('created_at', 'Unknown time')
    status_id = status.get('id', 'Unknown ID')
    
    # Try to parse and format the timestamp nicely if possible
    formatted_timestamp = timestamp
    try:
        if timestamp and timestamp != 'Unknown time':
            # Truth Social timestamps are typically in ISO 8601 format
            dt = datetime.datetime.fromisoformat(timestamp.replace('Z', '+00:00'))
            # Convert to Eastern Time for consistency
            dt_eastern = dt.astimezone(EASTERN_TIMEZONE)
            formatted_timestamp = dt_eastern.strftime('%Y-%m-%d %H:%M:%S %Z')
    except Exception as e:
        logging.warning(f"Could not parse timestamp '{timestamp}': {e}")
    
    if sentiment == "positive":
        emoji = "📈"
        impact_type = "POSITIVE"
    else:
        emoji = "📉"
        impact_type = "NEGATIVE"
    
    message = f"🚨 SIGNIFICANT {impact_type} IMPACT {emoji}\nUser: @{username}\nTimestamp: {formatted_timestamp}\nContent: {content_snippet}\n\nReasoning: {reasoning}"
    
    logging.warning(f"--- Significant {impact_type} Impact Detected at {formatted_timestamp} ---")
    logging.warning(f"Status ID: {status_id}")
    logging.warning(f"User: {username}")
    logging.warning(f"Created at: {formatted_timestamp}")
    logging.warning(f"Raw Content: {content_snippet}")
    logging.warning(f"Reasoning: {reasoning}")
    logging.warning("--------------------------------------------")
    
    if NTFY_TOPIC:
        try:
            requests.post(f"https://ntfy.sh/{NTFY_TOPIC}",
                data=message.encode(encoding='utf-8'))
            logging.info(f"Notification sent to ntfy.sh/{NTFY_TOPIC}")
        except Exception as e:
            logging.error(f"Failed to send notification to ntfy.sh: {e}")
    else:
        logging.warning("NTFY_TOPIC not configured. Notification not sent.")


def process_statuses(api, last_known_id):
    """
    Fetch, analyze, and potentially notify for new statuses since the last known ID using pull_statuses.
    Returns the ID of the latest status processed in this batch, or the last_known_id if none were new.
    """
    logging.info(f"Fetching new statuses for @{TARGET_HANDLE} since ID: {last_known_id or 'None'}...")
    latest_id_in_batch = last_known_id
    new_statuses_processed = 0
    statuses_processed_this_run = [] # Keep track of IDs processed in this specific run

    try:
        # Use pull_statuses which handles fetching statuses newer than since_id
        # It returns a generator, so we iterate through it
        for status in api.pull_statuses(username=TARGET_HANDLE, since_id=last_known_id, replies=False, verbose=False):
            status_id_str = str(status.get('id'))
            if not status_id_str:
                 logging.warning("Found status without ID from pull_statuses, skipping.")
                 continue

            # Double-check against last_known_id (pull_statuses should handle this, but belt-and-suspenders)
            if last_known_id is not None and status_id_str <= last_known_id:
                logging.warning(f"pull_statuses returned status ID {status_id_str} which is not newer than {last_known_id}. Skipping.")
                continue

            logging.info(f"Processing new status ID: {status_id_str}")
            statuses_processed_this_run.append(status_id_str) # Track processed ID for this run
            new_statuses_processed += 1

            status_text_html = status.get('content', '')

            if status_text_html:
                sentiment, significant, reasoning = analyze_sentiment(status_text_html)
                if significant and (sentiment == "positive" or sentiment == "negative"):
                    send_notification(status, sentiment, reasoning)
            else:
                logging.info(f"Status ID {status_id_str} has no text content.")

            # Keep track of the absolute latest ID seen in this batch
            if latest_id_in_batch is None or status_id_str > latest_id_in_batch:
                 latest_id_in_batch = status_id_str


        if new_statuses_processed == 0:
            logging.info("No new statuses found since last check (ID: {}).".format(last_known_id or 'None'))
        else:
             logging.info(f"Processed {new_statuses_processed} new statuses. Latest ID processed in this run: {latest_id_in_batch}. Processed IDs: {statuses_processed_this_run}")
             # Only save if we actually processed new statuses and got a new latest ID
             if latest_id_in_batch != last_known_id:
                save_last_processed_id(latest_id_in_batch)

    except Exception as e:
        logging.error(f"An error occurred fetching or processing statuses: {e}", exc_info=True)
    
    # Return the latest ID found in this batch, even if it's the same as the input last_known_id
    return latest_id_in_batch


def main():
    """Main application loop."""
    if not TRUTHSOCIAL_USERNAME or not TRUTHSOCIAL_PASSWORD:
        logging.critical("TRUTHSOCIAL_USERNAME and TRUTHSOCIAL_PASSWORD must be set. Exiting.")
        return

    last_processed_id = load_last_processed_id()

    logging.info("Initializing Truth Social API client...")
    api = None
    try:
        api = Api(TRUTHSOCIAL_USERNAME, TRUTHSOCIAL_PASSWORD)
        logging.info("Truth Social API client initialized successfully.")
    except Exception as e:
        logging.critical(f"API initialization failed: {e}. Exiting.", exc_info=True)
        return

    logging.info(f"Starting polling loop for @{TARGET_HANDLE}.")
    logging.info(f"Polling window: {POLL_START_HOUR}:00 - {POLL_END_HOUR}:00 ET.")
    logging.info(f"Polling interval: {POLL_INTERVAL_SECONDS} seconds.")
    logging.info(f"Starting with last processed ID: {last_processed_id or 'None'}")
    logging.info(f"Notifications will be sent to: {'ntfy.sh/' + NTFY_TOPIC if NTFY_TOPIC else 'DISABLED - NTFY_TOPIC not set'}")

    try:
        while True:
            if is_within_polling_hours():
                new_last_id = process_statuses(api, last_processed_id)
                if new_last_id != last_processed_id:
                     logging.info(f"Updating last processed ID to: {new_last_id}")
                     last_processed_id = new_last_id
                logging.debug(f"Check complete. Sleeping for {POLL_INTERVAL_SECONDS} seconds.")
                time.sleep(POLL_INTERVAL_SECONDS)

            else:
                now = datetime.datetime.now(EASTERN_TIMEZONE)
                start_time_today = now.replace(hour=POLL_START_HOUR, minute=0, second=0, microsecond=0)
                start_time_tomorrow = start_time_today + datetime.timedelta(days=1)
                next_run_time = start_time_today if now < start_time_today else start_time_tomorrow
                sleep_duration = max((next_run_time - now).total_seconds(), 1)

                logging.info(f"Outside polling hours. Sleeping until ~{next_run_time.strftime('%Y-%m-%d %H:%M:%S %Z')} ({sleep_duration:.0f} seconds)...")
                time.sleep(sleep_duration)

    except KeyboardInterrupt:
        logging.info("Shutdown signal received. Saving final state...")
        save_last_processed_id(last_processed_id)
        logging.info("Exiting.")
    except Exception as e:
        logging.critical(f"An unexpected critical error occurred in the main loop: {e}", exc_info=True)
        save_last_processed_id(last_processed_id)
        logging.critical("Exiting due to critical error.")

if __name__ == "__main__":
    main() 