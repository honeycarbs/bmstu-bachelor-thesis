import json
from fastapi import APIRouter, Response
from fastapi.responses import JSONResponse
from pydantic import BaseModel


class CreatorBody(BaseModel):
    dataset_metafile: str


class SampleHTTPRequestHandler:
    def __init__(self, sample_service, logger):
        self.router = APIRouter()
        self._service = sample_service
        self._logger = logger

        origins = [
            "http://localhost",
            "http://localhost:3000",
        ]

        @self.router.post("/api/v1/samples")
        async def create(response: Response, creator_body: CreatorBody):
            self._logger.info("got POST request at api/v1/samples")
            try:
                self._service.set_dataset_metafile(creator_body.dataset_metafile)
                logger.info(f"dataset metafile set to {creator_body.dataset_metafile}")

                samples = self._service.create()
                response.status_code = 201

                sample_dicts = [s.__dict__() for s in samples]
                return JSONResponse(content=sample_dicts)
            except Exception as e:
                self._logger.error(f"Error getting samples: {str(e)}")
                response.status_code = 500
                return JSONResponse(content={"error": f"Failed to get samples: {str(e)}"})

        @self.router.get("/api/v1/samples")
        async def get(response: Response):
            self._logger.info("got GET request at api/v1/samples")
            try:
                samples = self._service.get()
                response.status_code = 200

                sample_dicts = [s.__dict__() for s in samples]
                # response.headers["Access-Control-Allow-Origin"] = "http://localhost:3000"
                return JSONResponse(content=sample_dicts)
            except Exception as e:
                self._logger.error(f"Error getting samples: {str(e)}")
                response.status_code = 500
                return JSONResponse(content={"error": f"Failed to get samples: {str(e)}"})
