import logging
import requests

class Notifier:
    def __init__(self, ntfy_topic):
        self.ntfy_topic = ntfy_topic
        
    def send_notification(self, title, message, priority="default", tags=None):
        """
        Sends a notification using ntfy.sh.
        
        Args:
            title (str): The notification title
            message (str): The notification message
            priority (str, optional): Priority level (default, low, high, urgent). Defaults to "default".
            tags (list, optional): List of tag strings for the notification. Defaults to None.
        """
        logging.info(f"Preparing to send notification: {title}")
        
        if not self.ntfy_topic:
            logging.warning("NTFY_TOPIC not configured. Notification not sent.")
            return
            
        # Use query parameters for title and tags instead of headers to avoid encoding issues
        params = {}
        
        # Add title as a query parameter
        if title:
            params['title'] = title
            
        # Add tags as query parameters
        if tags:
            params['tags'] = ','.join(str(tag) for tag in tags)
            
        # Only use headers for priority which shouldn't contain Unicode
        headers = {
            "Priority": str(priority)
        }
            
        try:
            response = requests.post(
                f"https://ntfy.sh/{self.ntfy_topic}",
                data=message.encode(encoding='utf-8'),
                headers=headers,
                params=params
            )
            
            if response.status_code == 200:
                logging.info(f"Notification sent to ntfy.sh/{self.ntfy_topic}")
            else:
                logging.error(f"Failed to send notification. Status code: {response.status_code}")
                
        except Exception as e:
            logging.error(f"Failed to send notification to ntfy.sh: {e}") 