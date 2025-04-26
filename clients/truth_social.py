import logging
from truthbrush.api import Api

class TruthSocialClient:
    def __init__(self, username, password, target_handle):
        self.username = username
        self.password = password
        self.target_handle = target_handle
        self.api_client = None
        
    def initialize(self):
        """Initialize the Truth Social API client."""
        try:
            self.api_client = Api(self.username, self.password)
            logging.info("Truth Social API client initialized successfully.")
            return True
        except Exception as e:
            raise Exception(f"API initialization failed: {e}.")
            
    def get_new_statuses(self, last_known_id=None):
        """
        Fetch new statuses for the target user since the last known ID.
        Returns a generator of status objects.
        """
        if not self.api_client:
            logging.error("API client not initialized. Call initialize() first.")
            return []
            
        logging.info(f"Fetching new statuses for @{self.target_handle} since ID: {last_known_id or 'None'}...")
        
        try:
            # pull_statuses returns a generator, so we can iterate through it
            return self.api_client.pull_statuses(
                username=self.target_handle, 
                since_id=last_known_id, 
                replies=False, 
                verbose=False
            )
        except Exception as e:
            logging.error(f"Error fetching statuses: {e}", exc_info=True)
            return [] 