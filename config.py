import configparser
import os

config_path = "./config.ini" 

senformubin

config = configparser.ConfigParser()

try:
    with open(config_path, 'r') as configfile:
        config.read_file(configfile)
except FileNotFoundError:
    print("config.init not found.")
    raise FileNotFoundError("config not found. at {config_path}")
except Exception as e:
    raise e

print("config.ini loaded successfully.")
print(config.sections())
print(config["Database"])
print(config.get("Database", "username"))