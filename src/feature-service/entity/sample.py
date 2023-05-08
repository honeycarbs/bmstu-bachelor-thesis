

class Sample:
    def __init__(self, uuid, audio_path, emotion, batch):
        self.uuid = uuid
        self.audio_path = audio_path
        self.emotion = emotion
        self.batch = batch

    def __dict__(self):
        return {
            "uuid": self.uuid,
            "audio_path": self.audio_path,
            "emotion": self.emotion,
            "batch": self.batch
        }
