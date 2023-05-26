import numpy as np
import librosa
import librosa.feature as feature


class AudioFeaturesExtractor:
    def __init__(self, file_path, frame_length=0.02, hop_length=0.01):
        sig, self.sampling_rate = librosa.load(file_path)

        self.frame_length = int(self.sampling_rate * frame_length)  # 20 мс
        self.hop_length = int(hop_length * self.sampling_rate)      # 50% перекрытие

        self.frames = librosa.util.frame(sig, frame_length=self.frame_length, hop_length=self.hop_length)

    def get_mfcc(self, n_mfcc=13):
        mfcc = []
        for frame in self.frames.T:
            mfcc_per_frame = librosa.feature.mfcc(y=frame, sr=self.sampling_rate, n_fft=self.frame_length, n_mfcc=n_mfcc)
            mfcc.append(mfcc_per_frame.T)

        return np.array(mfcc)