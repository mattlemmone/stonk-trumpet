import logging
import datetime
import requests

class Notifier:
    def __init__(self, config):
        self.ntfy_topic = config.ntfy_topic
        self.eastern_timezone = config.eastern_timezone
        
    def send_notification(self, status, sentiment, reasoning):
        """
        Sends a notification for a significant market impact using ntfy.sh.
        
        Args:
            status (dict): The status data from Truth Social
            sentiment (str): The sentiment analysis result (positive/negative)
            reasoning (str): The reasoning for the sentiment analysis
        """
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
                dt_eastern = dt.astimezone(self.eastern_timezone)
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
        
        if self.ntfy_topic:
            try:
                requests.post(f"https://ntfy.sh/{self.ntfy_topic}",
                    data=message.encode(encoding='utf-8'))
                logging.info(f"Notification sent to ntfy.sh/{self.ntfy_topic}")
            except Exception as e:
                logging.error(f"Failed to send notification to ntfy.sh: {e}")
        else:
            logging.warning("NTFY_TOPIC not configured. Notification not sent.") 