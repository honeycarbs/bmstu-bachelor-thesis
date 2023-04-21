

class Sample:
    def __init__(self, uuid, audio_path, emo):
        self.uuid = uuid
        self.audio_path = audio_path
        self.emo = emo

    def __dict__(self):
        return {
            "uuid": self.uuid,
            "audio_path": self.audio_path,
            "emo": self.emo
        }
