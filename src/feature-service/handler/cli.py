import logging
import click
import os

from config.core import Config
from constants import CONFIG_FILE_PATH, LOG_PATH
import sys
from dbclient.core import PostgresqlClient
from service.sample import SampleService

from repository.sample.core import SampleRepository
from repository.frame.core import FrameRepository


def init_logger(config):
    logger = logging.getLogger("main")
    logger.setLevel(logging.INFO)

    if not os.path.exists(LOG_PATH):
        os.makedirs(LOG_PATH)

    fh = logging.FileHandler(f"{LOG_PATH}/all.log")
    fh.setLevel(logging.INFO)
    formatter = logging.Formatter("%(asctime)s - %(name)s - %(levelname)s - %(message)s")
    fh.setFormatter(formatter)
    logger.addHandler(fh)

    sh = logging.StreamHandler(sys.stdout)
    sh.setFormatter(formatter)
    logger.addHandler(sh)

    return logger


@click.group()
def cli():
    pass


@cli.command('fill', short_help='Prepare database for training and testing from dataset. For dataset arguments check config.yml.')
@click.option('-mf', '--dataset_metafile', type=str, required=True, default="merged.jsonl", show_default=True,
              help='Meta file path for your dataset relative to dataset root. '
                   'Check your dataset root in config.yml for feature service')
def fill(dataset_metafile):
    config = Config(CONFIG_FILE_PATH)
    logger = init_logger(config)
    db_conn = PostgresqlClient(config.PSQL_USER, config.PSQL_PASSWORD, config.PSQL_DATABASE)
    #
    sample_repository = SampleRepository(db_conn)
    frame_repository = FrameRepository(db_conn)

    sample_service = SampleService(sample_repository, frame_repository,
                                   config.DATASET_PATH,
                                   logger)

    logger.info("got POST request at api/v1/samples")
    try:
        sample_service.set_dataset_metafile(dataset_metafile)
        logger.info(f"dataset metafile set to {dataset_metafile}")

        sample_service.create_many()
        logger.info(f"all samples are created...")

    except Exception as e:
        logger.error(f"Error creating samples: {str(e)}")
        sys.exit(1)


@cli.command('add', short_help='Add to database your own audio file. For dataset arguments check config.yml.')
@click.option('-a', '--audio_path', type=click.Path(exists=True), required=True,
              help='Absolute audio file path recognition.')
def add(audio_path):
    config = Config(CONFIG_FILE_PATH)
    logger = init_logger(config)
    db_conn = PostgresqlClient(config.PSQL_USER, config.PSQL_PASSWORD, config.PSQL_DATABASE)
    #
    sample_repository = SampleRepository(db_conn)
    frame_repository = FrameRepository(db_conn)

    sample_service = SampleService(sample_repository, frame_repository,
                                   config.DATASET_PATH,
                                   logger)
    sample_service.create(audio_path)