import logging
import datetime
import time

import pytz

class Scheduler:
    def __init__(self, poll_start_hour, poll_end_hour, poll_interval_seconds, timezone):
        self.poll_start_hour = poll_start_hour
        self.poll_end_hour = poll_end_hour
        self.poll_interval_seconds = poll_interval_seconds
        self.timezone = pytz.timezone(timezone)
        
    def is_within_polling_hours(self):
        """Check if the current time is within the allowed polling hours in ET."""
        now_et = datetime.datetime.now(self.timezone)
        return self.poll_start_hour <= now_et.hour < self.poll_end_hour
        
    def sleep_until_next_run(self):
        """
        Sleep until the next run time.
        If within polling hours, sleep for the polling interval.
        If outside polling hours, sleep until the start of the next polling window.
        
        Returns:
            float: The sleep duration in seconds
        """
        if self.is_within_polling_hours():
            logging.debug(f"Check complete. Sleeping for {self.poll_interval_seconds} seconds.")
            time.sleep(self.poll_interval_seconds)
            return self.poll_interval_seconds
        else:
            now = datetime.datetime.now(self.timezone)
            start_time_today = now.replace(hour=self.poll_start_hour, minute=0, second=0, microsecond=0)
            start_time_tomorrow = start_time_today + datetime.timedelta(days=1)
            next_run_time = start_time_today if now < start_time_today else start_time_tomorrow
            sleep_duration = max((next_run_time - now).total_seconds(), 1)
            
            logging.info(f"Outside polling hours. Sleeping until ~{next_run_time.strftime('%Y-%m-%d %H:%M:%S %Z')} ({sleep_duration:.0f} seconds)...")
            time.sleep(sleep_duration)
            return sleep_duration 