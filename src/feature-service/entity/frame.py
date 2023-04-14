import uuid


class Frame:
    def __init__(self, sample_id, index, mfcc):
        self.id = str(uuid.uuid4())
        self.sample_id = sample_id
        self.sample_num = index
        self.mfcc = mfcc
