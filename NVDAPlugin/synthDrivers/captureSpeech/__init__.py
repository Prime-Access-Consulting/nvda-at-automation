from synthDrivers import espeak as espeakSynth
from speech.types import SpeechSequence
from logHandler import log
from .speech_socket import SpeechServer

class SynthDriver(espeakSynth.SynthDriver):
	name="captureSpeech"
	# Translators: Description for a speech synthesizer.
	description=_("Capture Speech")

	def __init__(self):
		self.server = SpeechServer(5678)
		self.server.start()
		super().__init__()

	def speak(self, speechSequence: SpeechSequence):
		for speech in speechSequence:
			if not isinstance(speech, str):
				continue
			self.server.push_speech(speech)
		super().speak(speechSequence)

	def terminate(self):
		self.server.stop()
		super().terminate()
