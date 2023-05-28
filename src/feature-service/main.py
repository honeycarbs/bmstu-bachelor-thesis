from handler.cli import cli
from extractor.core import AudioFeaturesExtractor


if __name__ == '__main__':
    cli()
    # ext = AudioFeaturesExtractor(file_path="/home/honeycarbs/BMSTU/bmstu-bachelor-thesis/src/DUSHA/crowd_train/wavs/2780f2347749fb1ef513a752897b49ca.wav")
    # mfcc = ext.get_mfcc()
    # # print(len(mfcc))
    # print(mfcc)