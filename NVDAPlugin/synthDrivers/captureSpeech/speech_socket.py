import threading
import socketserver
import re
from logHandler import log


class ThreadedTCPRequestHandler(socketserver.StreamRequestHandler):

	def handle(self):
		# opening a new connection should clear the speech queue
		last_event_id_msg = self.request.recv(1024).decode('utf-8')

		match = re.match(r'^Last-Event-ID:(\d+)$', last_event_id_msg)
		if not match:
			raise Exception(f'ARIA-AT: Invalid Last-Event-ID')

		last_event_index = int(match.group(1))

		log.info(f'ARIA-AT: Connection opened, last event index {last_event_index}')
		while True:
			response_last_event_index, response = self.server.get_response(last_event_index)

			if response is None:
				response = "\n"
			else:
				last_event_index = response_last_event_index

			self.wfile.write(bytes(response, 'utf-8'))
			self.wfile.flush()


class ThreadedTCPServer(socketserver.ThreadingMixIn, socketserver.TCPServer):
	def __init__(self, server_address, request_handler_class, bind_and_activate, events):
		super().__init__(server_address, request_handler_class, bind_and_activate)
		self._last_event_index = 0
		self._events = events

	def push_speech(self, speech):
		self._events.append((self._last_event_index, speech))
		self._last_event_index = self._last_event_index + 1

	def get_response(self, last_event_index):
		if last_event_index < 0:
			last_event_index = 0

		e = self._events[last_event_index:]

		if not e:
			return None, None

		next_index = last_event_index + len(e)

		return next_index, "\n".join(["event: speech\ndata: %s\nid: %s\n" % (data, index) for (index, data) in e])


class SpeechServer:
	def __init__(self, port):
		self._port = port
		self._events = list()
		self._server = ThreadedTCPServer(('127.0.0.1', port), ThreadedTCPRequestHandler, True, self._events)
		log.info(f'ARIA-AT: Started speech capture server on port {port}')

	def push_speech(self, speech):
		self._server.push_speech(speech)

	def start(self):
		server_thread = threading.Thread(target=self._server.serve_forever)
		server_thread.daemon = True
		server_thread.start()
