import time

from locust import HttpUser, task, constant


def get_current_unix_timestamp():
    """Returns the current Unix timestamp in milliseconds."""
    ms = str(time.time() * 1000)
    return ms.split('.')[0]


class EventGenerator(HttpUser):
    wait_time = constant(1)

    @task
    def send_http_requests(self):
        self.client.post('/events/broadcast', json={"topic_id": 1, "message": get_current_unix_timestamp()})
