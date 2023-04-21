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
