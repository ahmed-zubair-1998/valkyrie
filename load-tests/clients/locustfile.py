import time

import websocket
from locust import task, between, User


def get_current_unix_timestamp():
    """Returns the current Unix timestamp in milliseconds."""
    return time.time() * 1000


class WebsocketClient(User):
    wait_time = between(0.2, 2)

    def on_start(self):
        self.connection = False

    def connect(self):
        self.ws = websocket.WebSocketApp(
            'ws://localhost:8080/topics/subscribe?topicId=1',
            on_message=self.on_message,
            on_error=self.on_error,
            on_close=self.on_close
        )
        self.ws.run_forever() 

    @task
    def create_client(self):
        if not self.connection:
            self.connect()
            self.connection = True
    
    def on_message(self, ws, message):
        broadcast_time = float(message.split(':')[1])
        time_diff = get_current_unix_timestamp() - broadcast_time
        self.environment.events.request.fire(
            request_type="WSR",
            name="EVENT",
            response_time=time_diff,
            response_length=len(message),
            exception=None,
            context=self.context(),
        )
    
    def on_error(self, ws, error):
        if type(error) != websocket.WebSocketConnectionClosedException:
            print("Error:", error)
            if str(error) != "Connection to remote host was lost":
                self.environment.events.request.fire(
                    request_type="WSR",
                    name="EVENT",
                    response_time=0,
                    response_length=0,
                    exception=Exception("Websocket Exception"),
                    context=self.context(),
                )
    
    def on_close(self, ws, close_status_code, close_msg):
        print("Closed:", close_status_code, close_msg)
        self.environment.events.request.fire(
            request_type="WSR",
            name="EVENT",
            response_time=0,
            response_length=0,
            exception=Exception("Websocket Closed"),
            context=self.context(),
        )
        self.connection = False
