import platform
from http.server import BaseHTTPRequestHandler
from http import HTTPStatus
import json

import versionInfo

VERSION = "0.1"
HOST = 'localhost'
PORT = 8765

_info = json.dumps({
	'atName': 'NVDA',
	'atVersion': versionInfo.version,
	'platformName': platform.system().lower()
}).encode('utf-8')


class RequestHandler(BaseHTTPRequestHandler):
	def __init__(self, request, client_address, server):
		self.timeout = None
		BaseHTTPRequestHandler.__init__(self, request, client_address, server)

	def _set_headers(self, content_type='application/json', status=HTTPStatus.OK, extra=None):
		if extra is None:
			extra = {}

		self.send_response(status)
		self.send_header('Content-type', content_type)
		self.send_header('Access-Control-Allow-Origin', '*')

		for key, value in extra.items():
			self.send_header(key, value)

		self.end_headers()

	def do_GET(self):
		if not self.path or self.path == '/info':
			self._set_headers('text/plain')
			self.wfile.write(_info)
			return

	def do_POST(self):
		self._set_headers('text/plain', HTTPStatus.NOT_FOUND)
		return

	def do_OPTIONS(self):
		self.send_response(HTTPStatus.NO_CONTENT.value)
		self.send_header('Access-Control-Allow-Origin', '*')
		self.send_header('Access-Control-Allow-Methods', 'GET, POST')
		self.send_header('Access-Control-Allow-Headers', 'content-type')
		self.end_headers()
