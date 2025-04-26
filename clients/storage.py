import logging

def read_file(path):
    """
    Reads a value from a file.
    
    Args:
        path (str): Path to the file to read from
        
    Returns:
        str: The content of the file or None if file doesn't exist
    """
    value = None
    try:
        with open(path, 'r') as f:
            value = f.readline().strip()
        if value:
            logging.info(f"Read value '{value}' from {path}")
        else:
            logging.info(f"{path} was empty.")
    except FileNotFoundError:
        logging.info(f"{path} not found.")
    except IOError as e:
        logging.error(f"Error reading from {path}: {e}")
    return value

def write_file(path, value):
    """
    Writes a value to a file.
    
    Args:
        path (str): Path to the file to write to
        value (str): Value to write
    """
    if not value:
        return
    try:
        with open(path, 'w') as f:
            f.write(f"{value}\n")
        logging.debug(f"Wrote value '{value}' to {path}")
    except IOError as e:
        logging.error(f"Error writing value '{value}' to {path}: {e}") 