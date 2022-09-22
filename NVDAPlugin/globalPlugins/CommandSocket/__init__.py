# -*- coding: UTF-8 -*-
# Command Socket for NVDA
# Open up websocket and parse incoming commands

from .server import Server
import globalPluginHandler

COMMAND_PORT = 8765

class GlobalPlugin(globalPluginHandler.GlobalPlugin):
	def __init__(self):
		self._command_server = Server(COMMAND_PORT)
		self._command_server.start()
		print(f'ARIA-AT: Started http interface on port {COMMAND_PORT}')
		super().__init__()

	def terminate(self):
		self._command_server.stop()
		self._event_server.shutdown()