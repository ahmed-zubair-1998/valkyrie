import time

from locust import HttpUser, task, constant


def get_current_unix_timestamp():
    """Returns the current Unix timestamp in milliseconds."""
    return time.time() * 1000


class EventGenerator(HttpUser):
    wait_time = constant(1)

    @task
    def send_http_requests(self):
        self.client.post('/events/broadcast', json={"id": 1, "msg": str(get_current_unix_timestamp())})
