import os
import logging
import pytz
from dotenv import load_dotenv

class Config:
    def __init__(self):
        # Load environment variables
        load_dotenv()
        
        # Configure logging
        logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
        
        # API credentials
        self.truthsocial_username = os.getenv("TRUTHSOCIAL_USERNAME")
        self.truthsocial_password = os.getenv("TRUTHSOCIAL_PASSWORD")
        self.openai_api_key = os.getenv("OPENAI_API_KEY")
        
        # Application settings
        self.target_handle = os.getenv("TARGET_HANDLE", "realDonaldTrump")
        self.ntfy_topic = os.getenv("NTFY_TOPIC")
        self.last_processed_id_file = "last_processed_id.txt"
        self.eastern_timezone = pytz.timezone('US/Eastern')
        
        # Polling configuration
        try:
            self.poll_start_hour = int(os.getenv("POLL_START_HOUR", 7))
            self.poll_end_hour = int(os.getenv("POLL_END_HOUR", 23))
            self.poll_interval_seconds = int(os.getenv("POLL_INTERVAL_SECONDS", 300))
            
            if not (0 <= self.poll_start_hour < 24 and 0 <= self.poll_end_hour <= 24 and self.poll_start_hour < self.poll_end_hour):
                logging.error("Invalid POLL_START_HOUR or POLL_END_HOUR configuration. Using defaults.")
                self.poll_start_hour = 7
                self.poll_end_hour = 23
        except ValueError:
            logging.error("Invalid non-integer value for POLL_START_HOUR, POLL_END_HOUR, or POLL_INTERVAL_SECONDS. Using defaults.")
            self.poll_start_hour = 7
            self.poll_end_hour = 23
            self.poll_interval_seconds = 300
            
    def validate(self):
        """Validate critical configuration settings"""
        if not self.truthsocial_username or not self.truthsocial_password:
            raise ValueError("TRUTHSOCIAL_USERNAME and TRUTHSOCIAL_PASSWORD must be set")
            
        if not self.openai_api_key:
            raise ValueError("OPENAI_API_KEY environment variable must be set") 