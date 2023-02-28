import platform
from http.server import BaseHTTPRequestHandler
from http import HTTPStatus
import json
from urllib.parse import urlsplit, parse_qs

import inputCore
import versionInfo

from . import keyboard_input

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
		valid_paths = ['/settings', '/presskeys']

		if not self.path or self.path not in valid_paths:
			self._set_headers('text/plain', HTTPStatus.NOT_FOUND)
			return

		length = int(self.headers.get('content-length'))
		payload = json.loads(self.rfile.read(length))

		if self.path == '/settings':
			RequestHandler._handle_set_settings_command(payload)
		elif self.path == '/presskeys':
			RequestHandler._handle_press_keys_command(payload)

		self._set_headers()
		self.wfile.write(json.dumps({}).encode('utf-8'))

	def do_OPTIONS(self):
		self.send_response(HTTPStatus.NO_CONTENT.value)
		self.send_header('Access-Control-Allow-Origin', '*')
		self.send_header('Access-Control-Allow-Methods', 'GET, POST')
		self.send_header('Access-Control-Allow-Headers', 'content-type')
		self.end_headers()

	@staticmethod
	def _handle_press_keys_command(keys):
		import keyboardHandler

		gesture_name = None

		try:
			gesture_name = keyboard_input.create_gesture_name(keys)
			print(f'executing gesture "{gesture_name}"')

			gesture = keyboardHandler.KeyboardInputGesture.fromName(gesture_name)
			inputCore.manager.executeGesture(gesture)
		except KeyError as e:
			print(f'invalid gesture "{gesture_name}')
		except Exception as e:
			print(f'error executing gesture {gesture_name}: {e}')

	@staticmethod
	def _handle_set_settings_command(settings):
		import config
		settings_dict = tuple((k, k.split('.'), v) for k, v in settings.items())
		RequestHandler._parse_set_settings_command_data(config.conf.dict(), settings_dict)

	@staticmethod
	def _parse_set_settings_command_data(conf_dict, data):
		import config
		for (key, key_parts, set_value) in data:
			if not key_parts:
				continue
			(root_key, value) = RequestHandler._set_setting(key_parts, set_value, conf_dict)
			try:
				config.conf[root_key] = value
			except ValueError as error:
				print(f'Error setting {key} to {value}: {error}')

	@staticmethod
	def _set_setting(keys, value, setting_part, root_key=None):
		if root_key is None:
			root_key = keys.pop(0)

			if root_key not in setting_part:
				return setting_part

			return root_key, RequestHandler._set_setting(keys, value, setting_part[root_key], root_key)

		key = keys.pop(0)

		if key not in setting_part:
			return setting_part

		# leaf key
		if not keys:
			setting_part[key] = value
			return setting_part

		setting_part[key] = RequestHandler._set_setting(keys, value, setting_part[key], root_key)
		return setting_part

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
