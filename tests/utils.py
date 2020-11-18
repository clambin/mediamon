import json


class APIStub:
    def __init__(self, testfiles=None):
        self.testfiles = testfiles if testfiles is not None else dict()

    def call(self, endpoint, headers=None, body=None, params=None):
        if endpoint in self.testfiles:
            with open(self.testfiles[endpoint]['filename'], 'r') as f:
                content = f.read()
                try:
                    raw = self.testfiles[endpoint]['raw']
                except KeyError:
                    raw = False
                return json.loads(content) if raw is False else content
        return None
