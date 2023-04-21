import logging
import os

from fastapi import FastAPI
import uvicorn

from config.core import Config
from constants import CONFIG_FILE_PATH, LOG_PATH
from dbclient.core import PostgresqlClient
from handler.sample import SampleHTTPRequestHandler
from service.sample import SampleService

from repository.sample.core import SampleRepository
from repository.frame.core import FrameRepository


if __name__ == '__main__':
    config = Config(CONFIG_FILE_PATH)
    logger = logging.getLogger("main")
    logger.setLevel(logging.INFO)

    if not os.path.exists(LOG_PATH):
        os.makedirs(LOG_PATH)

    fh = logging.FileHandler(f"{LOG_PATH}/all.log")
    fh.setLevel(logging.INFO)
    formatter = logging.Formatter("%(asctime)s - %(name)s - %(levelname)s - %(message)s")
    fh.setFormatter(formatter)
    logger.addHandler(fh)

    db_conn = PostgresqlClient(config.PSQL_USER, config.PSQL_PASSWORD, config.PSQL_DATABASE)
    #
    sample_repository = SampleRepository(db_conn)
    frame_repository = FrameRepository(db_conn)

    sample_service = SampleService(sample_repository, frame_repository,
                                   config.DATASET_PATH,
                                   config.DATASET_METAFILE,
                                   logger)

    app = FastAPI()
    sample_handler = SampleHTTPRequestHandler(sample_service, logger)
    app.include_router(sample_handler.router)

    uvicorn.run(app, host=config.APP_HOST, port=config.APP_PORT)
