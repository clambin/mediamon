from enum import Enum
import json


class FakeResponse:
    def __init__(self, status_code, headers, text):
        self.status_code = status_code
        self.headers = headers
        self.text = text

    def json(self):
        return self.text


class APIStub:
    class Mode(Enum):
        raw = 0
        json = 1

    def __init__(self):
        self.fake_response = None

    def set_response(self, filename, mode=Mode.raw):
        with open(filename, 'r') as f:
            content = f.read()
            self.fake_response = json.loads(content) if mode == APIStub.Mode.json else content

    # FIXME: accept any arguments
    def call(self, endpoint, headers=None):
        return self.fake_response
