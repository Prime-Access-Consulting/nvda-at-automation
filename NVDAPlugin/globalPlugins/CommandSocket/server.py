from http.server import HTTPServer
from .handler import RequestHandler
import threading


class Server:
	def __init__(self, port):
		self._running = False
		self._port = port
		self._server = HTTPServer(('localhost', port), RequestHandler)

	def start(self):
		try:
			thread = threading.Thread(target=self._server.serve_forever)
			thread.daemon = True
			thread.start()
			print(f'ARIA-AT: Starting server on port {self._server.server_port}')
			self._running = True
			return True
		except Exception as e:
			print(f'ARIA-AT: Unable to start server on port {self._server.server_port}: {e}')
			return False

	def stop(self):
		try:
			if not self._running:
				return False
			self._server.shutdown()
			print(f'ARIA-AT: Stopping server on port {self._server.server_port}')
			return False
		except Exception as e:
			print(f'ARIA-AT: Unable to shutdown server: {e}')
			return True
