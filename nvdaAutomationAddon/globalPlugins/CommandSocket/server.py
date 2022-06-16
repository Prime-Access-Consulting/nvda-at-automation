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
            print(f'starting server on port {self._server.server_port}')
            self._running = True
            return True
        except:
            print(f'unable to start server on port {self._server.server_port}')
            return False

    def stop(self):
        try:
            if not self._running:
                return False
            self._server.shutdown()
            print(f'stopping server on port {self._server.server_port}')
            return False
        except:
            print(f'unable to shutdown server')
            return True
