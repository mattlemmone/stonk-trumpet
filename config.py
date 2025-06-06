import os
import logging
import pytz
from dotenv import load_dotenv

class Config:
    def __init__(self):
        # Load environment variables
        load_dotenv()
        
        # Configure logging
        # Create formatter with timestamp including milliseconds
        formatter = logging.Formatter(
            fmt='%(asctime)s.%(msecs)03d - %(levelname)s - %(message)s',
            datefmt='%Y-%m-%d %H:%M:%S'
        )
        
        # Configure root logger
        root_logger = logging.getLogger()
        root_logger.setLevel(logging.INFO)
        
        # Remove any existing handlers to avoid duplicates
        for handler in root_logger.handlers[:]:
            root_logger.removeHandler(handler)
            
        # Add console handler with our formatter
        console_handler = logging.StreamHandler()
        console_handler.setFormatter(formatter)
        root_logger.addHandler(console_handler)
        
        # API credentials
        self.truthsocial_username = os.getenv("TRUTHSOCIAL_USERNAME")
        self.truthsocial_password = os.getenv("TRUTHSOCIAL_PASSWORD")
        self.openai_api_key = os.getenv("OPENAI_API_KEY")
        
        # Application settings
        self.target_handle = os.getenv("TARGET_HANDLE")
        self.ntfy_topic = os.getenv("NTFY_TOPIC")
        self.storage_file = os.getenv("STORAGE_FILE")
        self.timezone = os.getenv("TIMEZONE")
        
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
        required_settings = {
            "TRUTHSOCIAL_USERNAME": self.truthsocial_username,
            "TRUTHSOCIAL_PASSWORD": self.truthsocial_password,
            "OPENAI_API_KEY": self.openai_api_key,
            "TARGET_HANDLE": self.target_handle,
            "NTFY_TOPIC": self.ntfy_topic,
            "TIMEZONE": self.timezone
        }
        
        for name, value in required_settings.items():
            if not value:
                print(f"{name} must be set")
                raise ValueError(f"{name} must be set")