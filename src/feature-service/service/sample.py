from entity.sample import Sample
from entity.frame import Frame
from processor.core import DatasetProcessor
from extractor.core import AudioFeaturesExtractor


# TODO: add logger
class SampleService:

    def __init__(self, sample_repository, frame_repository, dataset_absolute_path, dataset_metafile, logger):
        self._sample_repository = sample_repository
        self._frame_repository = frame_repository
        self._init_dataset_processor(dataset_absolute_path, dataset_metafile)
        self._logger = logger

    def _init_dataset_processor(self, _dataset_absolute_path, _dataset_metafile):
        self.dataset_processor = DatasetProcessor(
            _dataset_absolute_path,
            _dataset_metafile)

    def create(self):
        for i, wav in enumerate(self.dataset_processor.wavs):
            sample = Sample(wav["uuid"], wav["audio_path"], wav["emotion"])
            self._sample_repository.create(sample)
            self._logger.debug(f"created frame {i} out of {len(self.dataset_processor.wavs)}")

            audio_extractor = AudioFeaturesExtractor(file_path=sample.audio_path)
            mfcc_features = audio_extractor.get_mfcc(n_mfcc=13)
            frame_num = mfcc_features.shape[0]

            for j in range(frame_num):
                mfcc = mfcc_features[j][0]
                frame = Frame(sample.uuid, j + 1, mfcc)
                self._frame_repository.create(frame)

    def get(self):
        return self._sample_repository.get()
