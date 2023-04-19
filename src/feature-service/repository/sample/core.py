from repository.core import Repository


class SampleRepository(Repository):
    def __init__(self, db_client):
        self.db_client = db_client

    def create(self, entity):
        cursor = self.db_client.cnx.cursor()
        query = f"INSERT INTO sample VALUES (%s, %s, %s);"
        cursor.execute(query, (entity.uuid, entity.audio_path, entity.emo))
        self.db_client.cnx.commit()
        cursor.close()
