import click
import subprocess
import os
import webbrowser


@click.group()
def launcher():
    pass


@launcher.command('train', short_help='Launch HMM training procedure. Type train --help for arguments.')
@click.option('-mf', '--dataset_metafile', type=str, required=True, show_default=True,
              help='Meta file path for your dataset relative to dataset root. '
                   'Check your dataset root in config.yml for feature service')
def train(dataset_metafile):
    os.chdir("../feature-service")
    subprocess.run(['python3', 'main.py', 'fill', '--dataset_metafile', dataset_metafile])
    # os.chdir("../cluster-service")
    # subprocess.run(['go', 'run', 'cmd/main/main.go', '--mode', 'assign'])
    # os.chdir("../ml-service")
    # subprocess.run(['go', 'run', 'cmd/main/main.go', '--mode', 'train'])


@launcher.command('test', short_help='Launch HMM testing procedure. Type test --help for arguments.')
def test():
    os.chdir("../ml-service")
    subprocess.run(['go', 'run', 'cmd/main/main.go', '--mode', 'test'])
    webbrowser.open('file://' + '/home/honeycarbs/BMSTU/bmstu-bachelor-thesis/src/ml-service/etc/static/heatmap.html')


@launcher.command('recognize', short_help='Launch HMM recognizing procedure. Type recognize --help for arguments.')
@click.option('-a', '--audio_path', type=click.Path(exists=True), required=True,
              help='Absolute audio file path recognition.')
def recognize(audio_path):
    os.chdir("../feature-service")
    subprocess.run(['python3', 'main.py', 'add', '--audio_path', audio_path])
    os.chdir("../cluster-service")
    subprocess.run(['go', 'run', 'cmd/main/main.go', '--mode', 'add', '--path', audio_path])
    os.chdir("../ml-service")
    subprocess.run(['go', 'run', 'cmd/main/main.go', '--mode', 'recognize', '--path', audio_path])