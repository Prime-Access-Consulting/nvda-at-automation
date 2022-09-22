import platform
from http.server import BaseHTTPRequestHandler
from http import HTTPStatus
import json
from urllib.parse import urlsplit, parse_qs

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
		elif self.path.startswith('/settings?'):
			self._set_headers('text/plain')
			query = urlsplit(self.path).query
			query_parts = parse_qs(query)
			settings = query_parts['q'][0].split(',') if 'q' in query_parts and len(query_parts['q']) == 1 else []
			self.wfile.write(json.dumps(RequestHandler._get_settings(settings)).encode('utf-8'))


	def do_POST(self):
		self._set_headers('text/plain', HTTPStatus.NOT_FOUND)
		return

	def do_OPTIONS(self):
		self.send_response(HTTPStatus.NO_CONTENT.value)
		self.send_header('Access-Control-Allow-Origin', '*')
		self.send_header('Access-Control-Allow-Methods', 'GET, POST')
		self.send_header('Access-Control-Allow-Headers', 'content-type')
		self.end_headers()

	@staticmethod
	def _get_settings(names):
		import config

		data = RequestHandler._dot_join_settings({}, config.conf.dict(), None)

		if not names:
			return data

		return {k: v for k, v in data.items() if k in names}

	@staticmethod
	def _dot_join_settings(output, settings_dict, parent):
		for k, v in settings_dict.items():
			if type(v) is dict:
				output = RequestHandler._dot_join_settings(
					output,
					v,
					'.'.join([x for x in [parent, k] if x is not None])
				)
			else:
				key = k
				if parent is not None:
					key = '.'.join([parent, k])
				output[key] = v

		return output
