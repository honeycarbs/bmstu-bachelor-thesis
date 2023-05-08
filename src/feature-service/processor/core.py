from pathlib import Path
import os
import json


def _resolve_absolute_path(wav_path):
    p = Path(wav_path)
    return str(p.resolve())


class DatasetProcessor:
    def __init__(self, dataset_path):
        self._dataset_path = dataset_path
        self._dataset_meta_file = ""

        self.wavs = []
        # self._get_wavs()

    def set_metafile(self, dataset_metafile):
        self._dataset_meta_file = dataset_metafile

    def get_wavs(self):
        os.chdir(self._dataset_path)
        with open(self._dataset_meta_file, 'r') as mf:
            raw_data = list(mf)
            for i, rec in enumerate(raw_data):
                json_string = json.loads(rec)
                json_string["audio_path"] = _resolve_absolute_path(json_string["audio_path"])
                self.wavs.append(json_string)
