import uuid

from entity.sample import Sample
from entity.frame import Frame
from processor.core import DatasetProcessor
from extractor.core import AudioFeaturesExtractor


# TODO: add logger
class SampleService:

    def __init__(self, sample_repository, frame_repository, dataset_absolute_path, logger):
        self._sample_repository = sample_repository
        self._frame_repository = frame_repository
        self._init_dataset_processor(dataset_absolute_path)
        self._logger = logger

    def _init_dataset_processor(self, _dataset_absolute_path):
        self.dataset_processor = DatasetProcessor(
            _dataset_absolute_path)

    def set_dataset_metafile(self, dataset_metafile):
        self.dataset_processor.set_metafile(dataset_metafile)

    def create_many(self):
        self.dataset_processor.get_wavs()

        samples = []
        frames = []

        self._logger.info(f"red {len(self.dataset_processor.wavs)} wavs from file")
        for i, wav in enumerate(self.dataset_processor.wavs):
            sample = Sample(wav["uuid"], wav["audio_path"], wav["emotion"], wav["batch"])
            samples.append(sample)

            audio_extractor = AudioFeaturesExtractor(file_path=sample.audio_path)
            mfcc_features = audio_extractor.get_mfcc(n_mfcc=13)
            sample_frames_num = mfcc_features.shape[0]

            for j in range(sample_frames_num):
                mfcc = mfcc_features[j]
                frame = Frame(sample.uuid, j + 1, mfcc)
                frames.append(frame)

            self._logger.info(f"created sample {i} out of {len(self.dataset_processor.wavs)} with {sample_frames_num} frames")

            if len(samples) >= 100:
                self._logger.info(f"batch limit of 100 reached, database successfully updated")
                self._sample_repository.create_many(samples)
                self._frame_repository.create_many(frames)

                samples = []
                frames = []

        return samples

    def create(self, audio_path):
        frames = []
        sample = Sample(str(uuid.uuid4()), audio_path, "", "")

        audio_extractor = AudioFeaturesExtractor(file_path=sample.audio_path)
        mfcc_features = audio_extractor.get_mfcc(n_mfcc=13)
        frames_mfccs = mfcc_features.shape[0]

        for j in range(frames_mfccs):
            mfcc = mfcc_features[j][0]
            frame = Frame(sample.uuid, j + 1, mfcc)
            frames.append(frame)

        self._sample_repository.create(sample)
        self._frame_repository.create_many(frames)

        self._logger.info(f"created sample with {len(frames)} frames")

    def get(self):
        return self._sample_repository.get()
