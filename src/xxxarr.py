import logging
from enum import Enum
import requests
from prometheus_client import Gauge

from pimetrics.probe import APIProbe

GAUGES = {
    'calendar_count': Gauge('mediaserver_calendar_count', 'Number of upcoming episodes', ['server']),
    'queued_count': Gauge('mediaserver_queued_count', 'Number of queued torrents', ['server']),
    'monitored_count': Gauge('mediaserver_monitored_count', 'Number of monitored entries', ['server']),
    'unmonitored_count': Gauge('mediaserver_unmonitored_count', 'Number of unmonitored entries', ['server']),
}


class MonitorProbe(APIProbe):
    class App(Enum):
        sonarr = 0
        radarr = 1

    def __init__(self, host, app, api_key):
        super().__init__(f'http://{host}/')
        self.api_key = api_key
        self.app = app

    def report(self, output):
        if output:
            calendar = output['calendar']
            queue = output['queue']
            monitored = output['monitored'][0]
            unmonitored = output['monitored'][1]
            GAUGES['calendar_count'].labels(self.app.name).set(calendar)
            GAUGES['queued_count'].labels(self.app.name).set(queue)
            GAUGES['monitored_count'].labels(self.app.name).set(monitored)
            GAUGES['unmonitored_count'].labels(self.app.name).set(unmonitored)

    def call(self, endpoint):
        result = None
        try:
            headers = {'X-Api-Key': self.api_key}
            response = self.get(endpoint=endpoint, headers=headers)
            if response.status_code == 200:
                result = response.json()
            else:
                logging.error("%d - %s" % (response.status_code, response.reason))
        except requests.exceptions.RequestException as err:
            logging.warning(f'Failed to call "{self.url}": "{err}')
        return result

    def measure_calendar(self):
        calendar = self.call('api/calendar')
        if calendar:
            calendar = list(filter(lambda entry: not entry['hasFile'], calendar))
            return len(calendar)
        return 0

    def measure_queue(self):
        queue = self.call('api/queue')
        return len(queue) if queue else 0

    def measure_monitored(self):
        if self.app == self.App.sonarr:
            entries = self.call('api/series')
        elif self.app == self.App.radarr:
            entries = self.call('api/movie')
        else:
            entries = None
        if entries:
            monitored = list(filter(lambda entry: entry['monitored'], entries))
            unmonitored = list(filter(lambda entry: not entry['monitored'], entries))
            return len(monitored), len(unmonitored)
        return 0, 0

    def measure(self):
        return {
            'calendar': self.measure_calendar(),
            'queue': self.measure_queue(),
            'monitored': self.measure_monitored()
        }
