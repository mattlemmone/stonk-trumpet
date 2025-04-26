import logging
from clients.notifier import Notifier
from clients.storage import write_file
from clients.truth_social import TruthSocialClient
from sentiment_analyzer import SentimentAnalyzer

class StatusProcessor:
    def __init__(self, api_client: TruthSocialClient, sentiment_analyzer:SentimentAnalyzer, notifier: Notifier, storage_path: str):
        self.api_client = api_client
        self.sentiment_analyzer = sentiment_analyzer
        self.notifier = notifier
        self.storage_path = storage_path
        
    def process_statuses(self, after_message_id=None):
        """
        Fetch, analyze, and potentially notify for new statuses since the last known ID.
        Returns the ID of the latest status processed in this batch, or the last_known_id if none were new.
        """
        latest_id_in_batch = after_message_id
        new_statuses_processed = 0
        statuses_processed_this_run = [] # Keep track of IDs processed in this specific run
        
        try:
            # Get new statuses
            for status in self.api_client.get_new_statuses(after_message_id):
                status_id_str = str(status.get('id'))
                if not status_id_str:
                    logging.warning("Found status without ID, skipping.")
                    continue

                # Double-check against last_known_id (pull_statuses should handle this, but belt-and-suspenders)
                if after_message_id is not None and status_id_str <= after_message_id:
                    logging.warning(f"Got status ID {status_id_str} which is not newer than {after_message_id}. Skipping.")
                    continue

                logging.info(f"Processing new status ID: {status_id_str}")
                statuses_processed_this_run.append(status_id_str) # Track processed ID for this run
                new_statuses_processed += 1

                status_text_html = status.get('content', '')

                if status_text_html:
                    # Analyze sentiment
                    sentiment, significant, reasoning = self.sentiment_analyzer.analyze(status_text_html)
                    
                    # Send notification if significant
                    if significant and (sentiment == "positive" or sentiment == "negative"):
                        title, message, tags = self._format_notification(status, sentiment, reasoning)
                        
                        # Use the notifier interface
                        self.notifier.send_notification(
                            title=title, 
                            message=message,
                            priority="high", 
                            tags=tags
                        )
                else:
                    logging.info(f"Status ID {status_id_str} has no text content.")

                # Keep track of the absolute latest ID seen in this batch
                if latest_id_in_batch is None or status_id_str > latest_id_in_batch:
                    latest_id_in_batch = status_id_str

            if new_statuses_processed == 0:
                logging.info("No new statuses found since last check (ID: {}).".format(after_message_id or 'None'))
            else:
                logging.info(f"Processed {new_statuses_processed} new statuses. Latest ID: {latest_id_in_batch}. Processed IDs: {statuses_processed_this_run}")
                # Only save if we actually processed new statuses and got a new latest ID
                if latest_id_in_batch != after_message_id:
                    write_file(self.storage_path, latest_id_in_batch)

        except Exception as e:
            logging.error(f"An error occurred processing statuses: {e}", exc_info=True)
        
        # Return the latest ID found in this batch, even if it's the same as the input last_known_id
        return latest_id_in_batch
        
    def _format_notification(self, status, sentiment, reasoning):
        """
        Format notification details for a status with market impact.
        
        Args:
            status (dict): The status data
            sentiment (str): The sentiment analysis result ("positive" or "negative")
            reasoning (str): The reasoning for the sentiment analysis
            
        Returns:
            tuple: (title, message, tags)
                title (str): The notification title
                message (str): The formatted message body
                tags (list): List of tag strings
        """
        username = status.get('account', {}).get('username', 'Unknown')
        raw_content = status.get('content', '')
        content_snippet = raw_content[:200] + ('...' if len(raw_content) > 200 else '')
        
        if sentiment == "positive":
            emoji = "📈"
            impact_type = "POSITIVE"
        else:
            emoji = "📉"
            impact_type = "NEGATIVE"
            
        title = f"🚨 {impact_type} IMPACT {emoji}"
        message = f"User: @{username}\nContent: {content_snippet}\n\nReasoning: {reasoning}"
        tags = ["market", sentiment]
        
        return title, message, tags 