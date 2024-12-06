import configparser
import json
import os
import re
import subprocess
import time
from datetime import datetime

from confluent_kafka import Producer


producer = Producer({'bootstrap.servers': '10.220.9.91:9092'})

# Kafka消息推送函数
def send_to_kafka(job_id, job_info):
    message = json.dumps(job_info)
    producer.produce("test1", key=str(job_id), value=message.encode('utf-8'))
    producer.flush()


if __name__ == '__main__':
    send_to_kafka("job_id","{'name':'Jorny'}")
