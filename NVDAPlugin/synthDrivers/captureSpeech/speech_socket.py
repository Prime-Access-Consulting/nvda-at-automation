import threading
import socketserver
import re
from logHandler import log


class ThreadedTCPRequestHandler(socketserver.StreamRequestHandler):

	def handle(self):
		last_event_id_msg = self.request.recv(1024).decode('utf-8')

		match = re.match(r'^Last-Event-ID:(\d+)$', last_event_id_msg)
		if not match:
			raise Exception(f'ARIA-AT: Invalid Last-Event-ID')

		last_event_id = int(match.group(1))

		log.info(f'ARIA-AT: Connection opened, last event id {last_event_id}')
		while True:
			response = self.server.get_response(last_event_id)

			if not response:
				response = ""

			self.wfile.write(bytes(response, 'utf-8'))
			self.wfile.flush()


class ThreadedTCPServer(socketserver.ThreadingMixIn, socketserver.TCPServer):
	def __init__(self, server_address, request_handler_class, bind_and_activate, items):
		super().__init__(server_address, request_handler_class, bind_and_activate)
		self._items = items
		self._last_event_id = 0
		self._zero_index_event_id = 0

	def push_speech(self, speech):
		self._items.append((self._last_event_id, speech))
		self._last_event_id = self._last_event_id + 1
		size = len(self._items)

	def get_response(self, _last_event_id):
		e = self._items[:]
		self._items = []

		if not e:
			return None

		return "\n".join(["event: speech\ndata: %s\nid: %s\n" % (data, index) for (index, data) in e])


class SpeechServer:
	def __init__(self, port):
		self._port = port
		self._items = list()
		self._server = ThreadedTCPServer(('127.0.0.1', port), ThreadedTCPRequestHandler, True, self._items)
		log.info(f'ARIA-AT: Started speech capture server on port {port}')

	def push_speech(self, speech):
		self._server.push_speech(speech)

	def start(self):
		server_thread = threading.Thread(target=self._server.serve_forever)
		server_thread.daemon = True
		server_thread.start()
