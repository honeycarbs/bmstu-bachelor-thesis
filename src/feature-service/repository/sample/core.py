from repository.core import Repository
from entity.sample import Sample


class SampleRepository(Repository):
    def __init__(self, db_client):
        self.db_client = db_client

    def create(self, entity):
        cursor = self.db_client.cnx.cursor()
        query = f"INSERT INTO sample (uuid, audio_path, emotion, batch) VALUES (%s, %s, %s, %s)"
        cursor.execute(query, (entity.uuid, entity.audio_path, entity.emotion, entity.batch))
        self.db_client.cnx.commit()
        cursor.close()

    def create_many(self, entities):
        cursor = self.db_client.cnx.cursor()
        query = "INSERT INTO sample (uuid, audio_path, emotion, batch) VALUES (%s, %s, %s, %s)"
        values = [(entity.uuid, entity.audio_path, entity.emotion, entity.batch) for entity in entities]
        cursor.executemany(query, values)
        self.db_client.cnx.commit()
        cursor.close()

    def get(self):
        cursor = self.db_client.cnx.cursor()
        query = "SELECT uuid, audio_path, emotion, batch FROM sample"
        cursor.execute(query)

        samples = []
        for row in cursor.fetchall():
            uuid, audio_path, emotion, batch = row
            sample = Sample(uuid, audio_path, emotion, batch)
            samples.append(sample)

        cursor.execute(query)
        self.db_client.cnx.commit()
        cursor.close()

        return samples