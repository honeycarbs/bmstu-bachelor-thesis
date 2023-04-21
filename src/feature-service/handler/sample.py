import json
from fastapi import APIRouter, Response
from fastapi.responses import JSONResponse


class SampleHTTPRequestHandler:
    def __init__(self, sample_service, logger):
        self.router = APIRouter()
        self._service = sample_service
        self._logger = logger

        @self.router.post("/samples")
        async def create(response: Response):
            self._logger.debug("got POST request at api/v1/samples")
            self._service.create()
            response.status_code = 201

        @self.router.get("/samples")
        async def get(response: Response):
            self._logger.debug("got GET request at api/v1/samples")
            samples = self._service.get()
            response.status_code = 200

            sample_dicts = [s.__dict__() for s in samples]
            return JSONResponse(content=sample_dicts)
