import logging
from truthbrush.api import Api

class TruthSocialClient:
    def __init__(self, config):
        self.username = config.truthsocial_username
        self.password = config.truthsocial_password
        self.target_handle = config.target_handle
        self.api = None
        
    def initialize(self):
        """Initialize the Truth Social API client."""
        try:
            self.api = Api(self.username, self.password)
            logging.info("Truth Social API client initialized successfully.")
            return True
        except Exception as e:
            logging.critical(f"API initialization failed: {e}.", exc_info=True)
            return False
            
    def get_new_statuses(self, last_known_id=None):
        """
        Fetch new statuses for the target user since the last known ID.
        Returns a generator of status objects.
        """
        if not self.api:
            logging.error("API client not initialized. Call initialize() first.")
            return []
            
        logging.info(f"Fetching new statuses for @{self.target_handle} since ID: {last_known_id or 'None'}...")
        
        try:
            # pull_statuses returns a generator, so we can iterate through it
            return self.api.pull_statuses(
                username=self.target_handle, 
                since_id=last_known_id, 
                replies=False, 
                verbose=False
            )
        except Exception as e:
            logging.error(f"Error fetching statuses: {e}", exc_info=True)
            return [] 