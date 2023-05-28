import numpy as np
import librosa
import librosa.feature as feature


class AudioFeaturesExtractor:
    def __init__(self, file_path, frame_length=0.2, hop_length=0.1):
        sig, self.sampling_rate = librosa.load(file_path, res_type='kaiser_fast')

        self.frame_length = int(self.sampling_rate * frame_length)  # 100 мс
        self.hop_length = int(hop_length * self.sampling_rate)      # 50% перекрытие

        self.frames = librosa.util.frame(sig, frame_length=self.frame_length, hop_length=self.hop_length)

    def get_mfcc(self, n_mfcc=13):
        mfcc = []
        for frame in self.frames.T:
            mfccs = librosa.feature.mfcc(y=frame, sr=self.sampling_rate, n_mfcc=n_mfcc)
            mfcc_per_frame = np.mean(mfccs.T, axis=0)
            mfcc.append(mfcc_per_frame)

        return np.array(mfcc)