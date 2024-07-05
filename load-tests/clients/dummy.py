from locust import User, task

class Dummy(User):
    @task
    def hello(self):
        pass
