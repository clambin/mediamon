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
    def __init__(self, testfiles=None):
        self.testfiles = testfiles if testfiles is not None else dict()

    # FIXME: accept any arguments
    def call(self, endpoint, headers=None):
        if endpoint in self.testfiles:
            with open(self.testfiles[endpoint]['filename'], 'r') as f:
                content = f.read()
                try:
                    raw = self.testfiles[endpoint]['raw']
                except KeyError:
                    raw = False
                return json.loads(content) if raw is False else content
        return None
