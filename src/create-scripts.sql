CREATE TABLE sample (
    uuid VARCHAR(37) PRIMARY KEY,
    batch VARCHAR(10),
    audio_path VARCHAR(500),
    emotion VARCHAR(10)
);

CREATE TABLE centroid (
    uuid VARCHAR(37) PRIMARY KEY
);

CREATE TABLE cluster (
     uuid VARCHAR(37)  PRIMARY KEY,
     index INT,
     centroid_uuid VARCHAR(37) REFERENCES centroid(uuid)
);

CREATE TABLE frame (
    uuid VARCHAR(37) PRIMARY KEY,
    sample_uuid VARCHAR(37) REFERENCES sample(uuid),
    index INT,
    cluster_uuid VARCHAR(37) REFERENCES cluster(uuid)
);

CREATE TABLE mfcc (
    uuid VARCHAR(37) PRIMARY KEY,
    frame_uuid VARCHAR(37) REFERENCES frame(uuid),
    index INT CHECK (index >= 1 AND index <= 13),
    value decimal
);

CREATE TABLE centroid_coords (
     uuid VARCHAR(37) PRIMARY KEY,
     centroid_uuid VARCHAR(37) REFERENCES centroid(uuid),
     index INT,
     value NUMERIC
);