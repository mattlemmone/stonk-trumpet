import logging

class StatusProcessor:
    def __init__(self, api_client, sentiment_analyzer, notifier, persistence):
        self.api_client = api_client
        self.sentiment_analyzer = sentiment_analyzer
        self.notifier = notifier
        self.persistence = persistence
        
    def process_statuses(self, last_known_id=None):
        """
        Fetch, analyze, and potentially notify for new statuses since the last known ID.
        Returns the ID of the latest status processed in this batch, or the last_known_id if none were new.
        """
        latest_id_in_batch = last_known_id
        new_statuses_processed = 0
        statuses_processed_this_run = [] # Keep track of IDs processed in this specific run
        
        try:
            # Get new statuses
            for status in self.api_client.get_new_statuses(last_known_id):
                status_id_str = str(status.get('id'))
                if not status_id_str:
                    logging.warning("Found status without ID, skipping.")
                    continue

                # Double-check against last_known_id (pull_statuses should handle this, but belt-and-suspenders)
                if last_known_id is not None and status_id_str <= last_known_id:
                    logging.warning(f"Got status ID {status_id_str} which is not newer than {last_known_id}. Skipping.")
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
                        self.notifier.send_notification(status, sentiment, reasoning)
                else:
                    logging.info(f"Status ID {status_id_str} has no text content.")

                # Keep track of the absolute latest ID seen in this batch
                if latest_id_in_batch is None or status_id_str > latest_id_in_batch:
                    latest_id_in_batch = status_id_str

            if new_statuses_processed == 0:
                logging.info("No new statuses found since last check (ID: {}).".format(last_known_id or 'None'))
            else:
                logging.info(f"Processed {new_statuses_processed} new statuses. Latest ID: {latest_id_in_batch}. Processed IDs: {statuses_processed_this_run}")
                # Only save if we actually processed new statuses and got a new latest ID
                if latest_id_in_batch != last_known_id:
                    self.persistence.save_last_processed_id(latest_id_in_batch)

        except Exception as e:
            logging.error(f"An error occurred processing statuses: {e}", exc_info=True)
        
        # Return the latest ID found in this batch, even if it's the same as the input last_known_id
        return latest_id_in_batch 