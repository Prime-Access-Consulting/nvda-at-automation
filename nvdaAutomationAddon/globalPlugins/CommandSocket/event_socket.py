import queue
import socket
import threading
import socketserver
import json


class ThreadedTCPRequestHandler(socketserver.StreamRequestHandler):

	def handle(self):
		print(f'connection opened.')
		while True:
			response = self.server.get_response()

			if not response:
				break

			print(f'sending {response}')

			self.wfile.write(bytes(response, 'utf-8'))
			self.wfile.flush()


class ThreadedTCPServer(socketserver.ThreadingMixIn, socketserver.TCPServer):
	def __init__(self, server_address, request_handler_class, bind_and_activate, q):
		super().__init__(server_address, request_handler_class, bind_and_activate)
		self._queue = q

	def push_event(self, event):
		self._queue.put(event)
		size = str(self._queue.qsize())
		print(f'pushed event {event}, thread {threading.current_thread()}, {size} events total')

	def get_response(self):
		e = list()

		while not self._queue.empty():
			e.append(self._queue.get())

		if not e:
			return False

		return json.dumps([{'event': event} for event in e]) + "\x1E"


class EventServer:
	def __init__(self, port):
		self._port = port
		self._queue = queue.Queue()
		self._server = ThreadedTCPServer(('127.0.0.1', port), ThreadedTCPRequestHandler, True, self._queue)
		print(f'started event server on port {port}')

	def push_event(self, event):
		print(f'event {event}')
		self._server.push_event(event)

	def start(self):
		server_thread = threading.Thread(target=self._server.serve_forever)
		server_thread.daemon = True
		server_thread.start()

	def shutdown(self):
		self._server.shutdown()
