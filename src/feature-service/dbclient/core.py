import psycopg2


class PostgresqlClient:
    def __init__(self, username, password, dbname, host="localhost", port=5432, ):
        dsn = {
            'user': username,
            'password': password,
            'host': host,
            'port': port,
            'database': dbname
        }
        self.cnx = psycopg2.connect(**dsn)
        print("connection to db is opened")

    def __del__(self):
        self.cnx.close()
        print("connection to db is closed")
