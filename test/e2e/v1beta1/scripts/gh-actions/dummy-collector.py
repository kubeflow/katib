import argparse
import logging
import time

from kubernetes import client, config

# The default logging config.
logging.basicConfig(level=logging.INFO)

def collect_metrics(metric_name : str):
    config.load_incluster_config()
    v1 = client.CoreV1Api()

    while True:
        dummy_metric_value = 42 
        logging.info(f"Collected dummy metric: {metric_name}={dummy_metric_value}")
        
        time.sleep(10)  # Collect metrics every 10 seconds

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--metric-name", type=str, required=True, help="Name of the metric to collect")
    args = parser.parse_args()

    collect_metrics(args.metric_name)

   