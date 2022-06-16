from http.server import BaseHTTPRequestHandler
from http import HTTPStatus
import json
import uuid

import versionInfo

VERSION = "0.1"
HOST = 'localhost'
PORT = 8765

_info = "NVDA Command Socket v%s on %s:%d, NVDA v%s" % (VERSION, HOST, PORT, versionInfo.version)


class RequestHandler(BaseHTTPRequestHandler):
    def _set_headers(self, content_type='application/json', status=HTTPStatus.OK):
        self.send_response(status)
        self.send_header('Content-type', content_type)
        self.send_header('Access-Control-Allow-Origin', '*')
        self.end_headers()

    def do_GET(self):
        if not self.path or self.path == '/':
            self._set_headers('text/plain')
            self.wfile.write(_info.encode('utf-8'))
        else:
            self._set_headers('text/plain', HTTPStatus.NOT_FOUND)
            self.wfile.write(f"{self.path} 404 not found".encode('utf-8'))

    def do_POST(self):
        if not self.path or self.path != '/command':
            self._set_headers('text/plain', HTTPStatus.NOT_FOUND)
            return

        length = int(self.headers.get('content-length'))
        payload = json.loads(self.rfile.read(length))

        try:
            response = RequestHandler._parse_command(payload)
            self._set_headers()
            self.wfile.write(json.dumps({'success': True, 'command': payload, 'response': response}).encode('utf-8'))
        except (LookupError, RuntimeError) as e:
            self._set_headers('application/json', HTTPStatus.BAD_REQUEST)
            self.wfile.write(json.dumps({'success': False, 'command': payload, 'error': str(e)}).encode('utf-8'))

    def do_OPTIONS(self):
        self.send_response(HTTPStatus.NO_CONTENT.value)
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST')
        self.send_header('Access-Control-Allow-Headers', 'content-type')
        self.end_headers()

    @staticmethod
    def _parse_command(command):
        commands = {
            'startSession': lambda: {'message': 'Session started', 'id': str(uuid.uuid4())}
        }

        if 'name' not in command:
            raise RuntimeError('Unexpected Payload')

        if command['name'] not in commands:
            name = command['name']
            raise LookupError(f'Unknown Command Name \'{name}\'')

        return commands[command['name']]()
