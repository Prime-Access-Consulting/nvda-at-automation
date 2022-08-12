# -*- coding: UTF-8 -*-
# Command Socket for NVDA
# Open up websocket and parse incoming commands

from .server import Server
from .event_socket import EventServer
import globalPluginHandler

# from scriptHandler import script
# import ui

COMMAND_PORT = 8765
EVENT_PORT = 5432


class GlobalPlugin(globalPluginHandler.GlobalPlugin):
	def __init__(self):
		self._command_server = Server(COMMAND_PORT)
		self._command_server.start()
		self._event_server = EventServer(EVENT_PORT)
		self._event_server.start()
		super().__init__()

	def terminate(self):
		self._command_server.stop()
		self._event_server.shutdown()

	def event_gainFocus(self, obj, nextHandler):
		try:
			self._event_server.push_event(f'event_gainFocus - obj {obj.name}')
		except Exception as e:
			print(e)
		nextHandler()
