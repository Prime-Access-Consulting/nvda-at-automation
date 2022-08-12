import datetime
import platform
import threading
import time
from http.server import BaseHTTPRequestHandler
from http import HTTPStatus
import json
import uuid

import versionInfo
from synthDrivers.sapi5 import SynthDriver

VERSION = "0.1"
HOST = 'localhost'
PORT = 8765

_info = "NVDA Command Socket v%s on %s:%d, NVDA v%s" % (VERSION, HOST, PORT, versionInfo.version)


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
		if not self.path or self.path == '/':
			self._set_headers('text/plain')
			self.wfile.write(_info.encode('utf-8'))
		else:
			self._set_headers('text/plain', HTTPStatus.NOT_FOUND)
			self.wfile.write(f'{self.path} 404 not found'.encode('utf-8'))

	def do_POST(self):
		if not self.path or self.path != '/command':
			self._set_headers('text/plain', HTTPStatus.NOT_FOUND)
			return

		length = int(self.headers.get('content-length'))
		payload = json.loads(self.rfile.read(length))

		response = RequestHandler._parse_command(payload)
		self._set_headers()
		self.wfile.write(json.dumps(response).encode('utf-8'))

	def do_OPTIONS(self):
		self.send_response(HTTPStatus.NO_CONTENT.value)
		self.send_header('Access-Control-Allow-Origin', '*')
		self.send_header('Access-Control-Allow-Methods', 'GET, POST')
		self.send_header('Access-Control-Allow-Headers', 'content-type')
		self.end_headers()

	@staticmethod
	def _parse_command(command):
		commands = {
			'session.new': RequestHandler._handle_create_session_command,
			'settings.getSettings': RequestHandler._handle_get_settings_command
		}

		if 'method' not in command or command['method'] not in commands:
			return RequestHandler._handle_invalid_command(command)

		return commands[command['method']](command)

	@staticmethod
	def _handle_create_session_command(_command):
		return {
			'sessionId': str(uuid.uuid4()),
			'capabilities': {
				'atName': 'NVDA',
				'atVersion': versionInfo.version,
				'platformName': platform.system().lower()
			}
		}

	@staticmethod
	def _handle_invalid_command(command):
		return {
			'id': command['id'] if 'id' in command else None,
			'error': 'unknown command',
			'message': 'the command was not recognised'
		}

	@staticmethod
	def _handle_get_settings_command(command):
		import config
		c = config.conf

		return {
			'general.language': c['general']['language'],
			'general.saveConfigurationOnExit': c['general']['saveConfigurationOnExit'],
			'general.askToExit': c['general']['askToExit'],
			'general.playStartAndExitSounds': c['general']['playStartAndExitSounds'],
			'general.loggingLevel': c['general']['loggingLevel'],
			'general.showWelcomeDialogAtStartup': c['general']['showWelcomeDialogAtStartup'],

			'speech.synth': c['speech']['synth'],
			'speech.symbolLevel': c['speech']['symbolLevel'],
			'speech.trustVoiceLanguage': c['speech']['trustVoiceLanguage'],
			'speech.includeCLDR': c['speech']['includeCLDR'],
			'speech.beepSpeechModePitch': c['speech']['beepSpeechModePitch'],
			'speech.outputDevice': c['speech']['outputDevice'],
			'speech.autoLanguageSwitching': c['speech']['autoLanguageSwitching'],
			'speech.autoDialectSwitching': c['speech']['autoDialectSwitching'],
		}
