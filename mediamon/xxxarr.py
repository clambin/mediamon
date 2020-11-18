import logging
from enum import Enum
import requests
from pimetrics.probe import APIProbe
from mediamon import metrics


class MonitorProbe(APIProbe):
    class App(Enum):
        sonarr = 0
        radarr = 1

    def __init__(self, host, app, api_key):
        super().__init__(f'http://{host}/')
        self.api_key = api_key
        self.app = app
        self.healthy = True

    @property
    def name(self):
        return 'sonarr' if self.app == MonitorProbe.App.sonarr else 'radarr'

    def report(self, output):
        logging.debug(f'{self.name}: {output}')
        metrics.report(output, self.name)

    def _call(self, endpoint):
        try:
            headers = {'X-Api-Key': self.api_key}
            response = self.get(endpoint=endpoint, headers=headers)
            if response.status_code == 200:
                return response.json()
            else:
                logging.error("%d - %s" % (response.status_code, response.reason))
        except requests.exceptions.RequestException as err:
            logging.warning(f'Failed to call "{self.url}": "{err}')
        return None

    def call(self, endpoint):
        if result := self._call(endpoint):
            if not self.healthy:
                logging.info(f'Connection with {self.name} re-established')
                self.healthy = True
        else:
            self.healthy = False
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
        entries = None
        if self.app == self.App.sonarr:
            entries = self.call('api/series')
        elif self.app == self.App.radarr:
            entries = self.call('api/movie')
        monitored = unmonitored = []
        if entries:
            monitored = list(filter(lambda entry: entry['monitored'], entries))
            unmonitored = list(filter(lambda entry: not entry['monitored'], entries))
        return len(monitored), len(unmonitored)

    def measure_version(self):
        entries = self.call('api/system/status')
        if entries and 'version' in entries:
            return entries['version']
        else:
            logging.debug('No version found')
        return None

    def measure(self):
        monitored, unmonitored = self.measure_monitored()
        return {
            'xxxarr_calendar': self.measure_calendar(),
            'xxxarr_queue': self.measure_queue(),
            'xxxarr_monitored': monitored,
            'xxxarr_unmonitored': unmonitored,
            'version': self.measure_version()
        }
