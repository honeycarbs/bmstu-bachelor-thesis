from repository.core import Repository
import uuid


class FrameRepository(Repository):
    def __init__(self, db_client):
        self.db_client = db_client

    def create(self, entity):
        self._create_frame(entity)
        self._create_mfcc(entity.id, entity.mfcc)

    def get(self):
        pass

    def _create_frame(self, frame):
        cursor = self.db_client.cnx.cursor()
        query = "INSERT INTO frame (uuid, sample_uuid, index) VALUES (%s, %s, %s)"
        values = (frame.id, frame.sample_id, frame.sample_num)

        cursor.execute(query, values)
        self.db_client.cnx.commit()
        cursor.close()

    def _create_mfcc(self, frame_id, mfcc):
        cursor = self.db_client.cnx.cursor()

        for i, c in enumerate(mfcc):
            query = "INSERT INTO mfcc (uuid, frame_uuid, index, value)  VALUES (%s, %s, %s, %s)"
            values = (str(uuid.uuid4()), frame_id, i + 1, c.item())
            cursor.execute(query, values)
            self.db_client.cnx.commit()

        cursor.close()

    def create_many(self, frames):
        frame_ids = [frame.id for frame in frames]
        mfccs = [frame.mfcc for frame in frames]

        self._create_frames(frames)
        self._create_mfccs(frame_ids, mfccs)

    def _create_frames(self, frames):
        cursor = self.db_client.cnx.cursor()
        query = "INSERT INTO frame (uuid, sample_uuid, index) VALUES (%s, %s, %s)"
        values = [(frame.id, frame.sample_id, frame.sample_num) for frame in frames]

        cursor.executemany(query, values)
        self.db_client.cnx.commit()
        cursor.close()

    def _create_mfccs(self, frame_ids, mfccs):
        cursor = self.db_client.cnx.cursor()
        query = "INSERT INTO mfcc (uuid, frame_uuid, index, value) VALUES (%s, %s, %s, %s)"
        values = []

        for frame_id, mfcc in zip(frame_ids, mfccs):
            for i, c in enumerate(mfcc):
                mfcc_uuid = str(uuid.uuid4())
                value = c.item()
                values.append((mfcc_uuid, frame_id, i + 1, value))

        cursor.executemany(query, values)
        self.db_client.cnx.commit()
        cursor.close()
