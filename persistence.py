import logging

class Persistence:
    def __init__(self, config):
        self.last_processed_id_file = config.last_processed_id_file
        
    def load_last_processed_id(self):
        """Loads the last processed status ID from a file."""
        last_id = None
        try:
            with open(self.last_processed_id_file, 'r') as f:
                last_id = f.readline().strip()
            if last_id:
                logging.info(f"Loaded last processed ID {last_id} from {self.last_processed_id_file}")
            else:
                logging.info(f"{self.last_processed_id_file} was empty. Starting fresh.")
        except FileNotFoundError:
            logging.info(f"{self.last_processed_id_file} not found. Starting fresh.")
        except IOError as e:
            logging.error(f"Error loading last processed ID from {self.last_processed_id_file}: {e}")
        return last_id

    def save_last_processed_id(self, last_id):
        """Saves the last processed status ID to a file."""
        if not last_id:
            return
        try:
            with open(self.last_processed_id_file, 'w') as f:
                f.write(f"{last_id}\n")
            logging.debug(f"Saved last processed ID {last_id} to {self.last_processed_id_file}")
        except IOError as e:
            logging.error(f"Error saving last processed ID {last_id} to {self.last_processed_id_file}: {e}") 