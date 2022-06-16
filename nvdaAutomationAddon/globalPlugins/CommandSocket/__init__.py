# -*- coding: UTF-8 -*-
# Command Socket for NVDA
# Allows RPC over HTTP

from .server import Server
import globalPluginHandler

# from scriptHandler import script
# import ui

PORT = 8765


class GlobalPlugin(globalPluginHandler.GlobalPlugin):
    def __init__(self):
        self._server = Server(PORT)
        self._server.start()
        super().__init__()

    def terminate(self):
        self._server.stop()
