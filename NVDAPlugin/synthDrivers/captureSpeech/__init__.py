from synthDrivers import espeak as espeakSynth
from speech.types import SpeechSequence
from logHandler import log

class SynthDriver(espeakSynth.SynthDriver):
	name="captureSpeech"
	# Translators: Description for a speech synthesizer.
	description=_("Capture Speech")

	def speak(self, speechSequence: SpeechSequence):
		for speech in speechSequence:
			if not isinstance(speech, str):
				continue
			log.info("Speech: %s | pitch: %d, inflection: %d, rate: %d" % (speech, self.pitch, self.inflection, self.rate))

		super().speak(speechSequence)
