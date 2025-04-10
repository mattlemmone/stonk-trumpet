import logging
import time

from config import Config
from persistence import Persistence
from api_client import TruthSocialClient
from sentiment_analyzer import SentimentAnalyzer
from notifier import Notifier
from scheduler import Scheduler
from processor import StatusProcessor

def main():
    """Main application entry point."""
    # Initialize configuration
    try:
        config = Config()
        config.validate()
    except Exception as e:
        logging.critical(f"Configuration error: {e}")
        return

    # Initialize components
    persistence = Persistence(config)
    api_client = TruthSocialClient(config)
    sentiment_analyzer = SentimentAnalyzer(config)
    notifier = Notifier(config)
    scheduler = Scheduler(config)
    
    # Load last processed ID
    last_processed_id = persistence.load_last_processed_id()
    
    # Initialize API client
    if not api_client.initialize():
        return
        
    # Create status processor with all dependencies
    processor = StatusProcessor(api_client, sentiment_analyzer, notifier, persistence)

    # Log startup information
    logging.info(f"Starting polling loop for @{config.target_handle}.")
    logging.info(f"Polling window: {config.poll_start_hour}:00 - {config.poll_end_hour}:00 ET.")
    logging.info(f"Polling interval: {config.poll_interval_seconds} seconds.")
    logging.info(f"Starting with last processed ID: {last_processed_id or 'None'}")
    logging.info(f"Notifications will be sent to: {'ntfy.sh/' + config.ntfy_topic if config.ntfy_topic else 'DISABLED - NTFY_TOPIC not set'}")

    try:
        while True:
            if scheduler.is_within_polling_hours():
                new_last_id = processor.process_statuses(last_processed_id)
                if new_last_id != last_processed_id:
                    logging.info(f"Updating last processed ID to: {new_last_id}")
                    last_processed_id = new_last_id
                scheduler.sleep_until_next_run()
            else:
                scheduler.sleep_until_next_run()

    except KeyboardInterrupt:
        logging.info("Shutdown signal received. Saving final state...")
        persistence.save_last_processed_id(last_processed_id)
        logging.info("Exiting.")
    except Exception as e:
        logging.critical(f"An unexpected critical error occurred in the main loop: {e}", exc_info=True)
        persistence.save_last_processed_id(last_processed_id)
        logging.critical("Exiting due to critical error.")

if __name__ == "__main__":
    main() 