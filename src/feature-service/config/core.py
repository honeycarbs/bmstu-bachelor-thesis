import os

import yaml


class Config:
    __slots__ = ["filepath",
                 "DATASET_PATH",
                 "DATASET_METAFILE",
                 "PSQL_PORT",
                 "PSQL_USER",
                 "PSQL_PASSWORD",
                 "PSQL_DATABASE",
                 "APP_HOST",
                 "APP_PORT"]

    def __init__(self, config_file_path):
        self.filepath = config_file_path
        self._read()

    def _read(self):
        if not os.path.exists(self.filepath):
            raise AttributeError(f"config yaml doesnt exist: {self.filepath}")

        with open(self.filepath) as config_file:
            config_content = config_file.read()
            config_yaml = yaml.safe_load(config_content)

        for k, v in config_yaml.items():
            k = k.upper()
            if k in self.__slots__:
                setattr(self, k, v)
