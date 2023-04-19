import logging

from extractor.core import AudioFeaturesExtractor
from entity.frame import Frame
from entity.sample import DatasetEntity
from dbclient.core import PostgresqlClient
from processor.core import DatasetProcessor
from repository.sample.core import SampleRepository
from repository.frame.core import FrameRepository


def setup_logger(lg):
    lg.setLevel(logging.DEBUG)

    fh = logging.FileHandler("logs/all.log")
    fh.setLevel(logging.DEBUG)
    formatter = logging.Formatter("%(asctime)s - %(name)s - %(levelname)s - %(message)s")
    fh.setFormatter(formatter)

    lg.addHandler(fh)
    lg.addHandler(logging.StreamHandler())


if __name__ == '__main__':
    db_conn = PostgresqlClient("admin", "admin", "thesis")
    logger = logging.getLogger("main")
    setup_logger(logger)

    dataset_processor = DatasetProcessor(
        "/home/honeycarbs/BMSTU/bmstu-bachelor-thesis/src/DUSHA/processed_dataset_090/aggregated_dataset",
        "smaller_dataset.jsonl")

    sample_repository = SampleRepository(db_conn)
    frame_repository = FrameRepository(db_conn)

    wavs_length = len(dataset_processor.wavs)

    dataset_processor.wavs = dataset_processor.wavs

    for i, wav in enumerate(dataset_processor.wavs):
        sample = DatasetEntity(wav["uuid"], wav["audio_path"], wav["emotion"])
        sample_repository.create(sample)
        logger.debug(f"created {i + 1} out of {wavs_length} samples")

        audio_extractor = AudioFeaturesExtractor(file_path=sample.audio_path)
        mfcc_features = audio_extractor.get_mfcc(n_mfcc=13)
        frame_num = mfcc_features.shape[0]
        logger.debug(f"frame number is {frame_num}")
        for j in range(frame_num):
            mfcc = mfcc_features[j][0]
            frame = Frame(sample.uuid, j + 1, mfcc)
            frame_repository.create(frame)