import logging
import os
import sys

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
import uvicorn

from config.core import Config
from constants import CONFIG_FILE_PATH, LOG_PATH
from dbclient.core import PostgresqlClient
from handler.sample import SampleHTTPRequestHandler
from service.sample import SampleService

from repository.sample.core import SampleRepository
from repository.frame.core import FrameRepository


origins = [
    "http://localhost",
    "http://localhost:3000",
]


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

    sh = logging.StreamHandler(sys.stdout)
    logger.addHandler(sh)

    db_conn = PostgresqlClient(config.PSQL_USER, config.PSQL_PASSWORD, config.PSQL_DATABASE)
    #
    sample_repository = SampleRepository(db_conn)
    frame_repository = FrameRepository(db_conn)

    sample_service = SampleService(sample_repository, frame_repository,
                                   config.DATASET_PATH,
                                   logger)

    logger.info("service initialized")
    app = FastAPI()
    app.add_middleware(
        CORSMiddleware,
        allow_origins=origins,
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    sample_handler = SampleHTTPRequestHandler(sample_service, logger)
    app.include_router(sample_handler.router)
    logger.info("handler registered")

    uvicorn.run(app, host=config.APP_HOST, port=config.APP_PORT)
