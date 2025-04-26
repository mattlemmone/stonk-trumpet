import logging
import time

from config import Config
from clients.storage import read_file, write_file
from clients.truth_social import TruthSocialClient
from sentiment_analyzer import SentimentAnalyzer
from clients.notifier import Notifier
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
    truth_social_client = TruthSocialClient(
        username=config.truthsocial_username,
        password=config.truthsocial_password,
        target_handle=config.target_handle
    )
    
    truth_social_client.initialize()
    
    sentiment_analyzer = SentimentAnalyzer(config.openai_api_key)
    notifier = Notifier(config.ntfy_topic)
    scheduler = Scheduler(
        poll_start_hour=config.poll_start_hour,
        poll_end_hour=config.poll_end_hour,
        poll_interval_seconds=config.poll_interval_seconds,
        timezone=config.timezone
    )
    
    # Load last processed ID
    last_processed_id = read_file(config.storage_file)
    
    # Create status processor with all dependencies
    processor = StatusProcessor(truth_social_client, sentiment_analyzer, notifier, config.storage_file)

    # Log startup information
    logging.info(f"Starting polling loop for @{config.target_handle}.")
    logging.info(f"Polling window: {config.poll_start_hour}:00 - {config.poll_end_hour}:00 {config.timezone}.")
    logging.info(f"Polling interval: {config.poll_interval_seconds} seconds.")
    logging.info(f"Starting with last processed ID: {last_processed_id or 'None'}")
    logging.info(f"Notifications will be sent to: {'ntfy.sh/' + config.ntfy_topic if config.ntfy_topic else 'DISABLED - NTFY_TOPIC not set'}")

    try:
        while True:
            if scheduler.is_within_polling_hours():
                newest_message_id = processor.process_statuses(last_processed_id)

                if newest_message_id != last_processed_id:
                    logging.info(f"Updating last processed ID to: {newest_message_id}")
                    last_processed_id = newest_message_id

            scheduler.sleep_until_next_run()

    except KeyboardInterrupt:
        logging.info("Shutdown signal received. Saving final state...")
        write_file(config.storage_file, last_processed_id)
        logging.info("Exiting.")
    except Exception as e:
        logging.critical(f"An unexpected critical error occurred in the main loop: {e}", exc_info=True)
        write_file(config.storage_file, last_processed_id)
        logging.critical("Exiting due to critical error.")

if __name__ == "__main__":
    main() 