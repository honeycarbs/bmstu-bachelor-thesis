import os
import sys

import subprocess
from subprocess import check_output

from pathlib import Path
import dash
from dash import dcc, html
from dash.dependencies import Input, Output, State
import dash_bootstrap_components as dbc
from dash.exceptions import PreventUpdate


app = dash.Dash(__name__, external_stylesheets=[dbc.themes.BOOTSTRAP])

app.layout = dbc.Container(
    [
        dbc.Container(
            html.H3("Метод распознавания эмоций из звучащей речи на основе скрытой марковской модели", style={"text-align": "center"}),
            className="align-items-center justify-content-center mt-4"
        ),
        dbc.Row(
            [
                dbc.Col(
                    [
                        dcc.Upload(
                            id='upload-audio',
                            children=dbc.Button('Загрузить', id='load-button', color="primary")
                        ),
                    ],
                    className="mt-2",
                ),
                dbc.Col(
                    dbc.Button("Распознать", id="recognize-button", color="primary", disabled=True),
                    className="mt-2",
                ),
            ],
            justify="center",
            className="mt-3",
        ),
        dbc.Container(id="emotion-container", className="mt-3"),
        html.Div(id="audio-player-container", className="text-center"),
        dcc.Loading(id="loading-spinner", type="circle", fullscreen=True),
    ],
    className="d-flex flex-column align-items-center justify-content-center mt-4",
)


@app.callback(
    Output("recognize-button", "disabled"),
    Input('upload-audio', 'contents')
)
def update_recognize_button(contents):
    if contents is not None:
        return False  # Enable the button if audio is loaded
    return True  # Disable the button if audio is not loaded


def parse_contents(contents):
    return html.Audio(src=contents, controls=True)


@app.callback(
    Output("audio-player-container", "children"),
    Output("loading-spinner", "children"),
    Input('upload-audio', 'contents'),
    State('upload-audio', 'filename'),
    Input("recognize-button", "n_clicks")
)
def update_audio_player(contents, filename, recognize_button_clicks):
    ctx = dash.callback_context
    if not ctx.triggered:
        raise PreventUpdate

    prop_id = ctx.triggered[0]['prop_id']
    loading_spinner = None

    if 'upload-audio.contents' in prop_id:
        if contents is not None:
            return parse_contents(contents), loading_spinner
    elif 'recognize-button.n_clicks' in prop_id:
        if recognize_button_clicks is not None:
            emotion = driver(filename)  # Call your long-running function here to get the emotion
            # loading_spinner = dbc.Alert(emotion, color=get_emotion_color(emotion), className="mt-3")
            color, text = get_emotion_color(emotion)
            loading_spinner = dbc.Alert(
                text,
                color=color,
                className="mt-3"
            )

        return dash.no_update, loading_spinner


def get_emotion_color(emotion):
    emotion_map = {
        'positive': ('success', 'позитив'),
        'neutral': ('light', 'нейтраль'),
        'angry': ('danger', 'агрессия'),
        'sad': ('primary', 'грусть')
    }
    return emotion_map.get(emotion, ('light', 'неизвестно?'))  # Default to 'light' if emotion is not in the map


def driver(filename):
    home = str(Path.home())
    audio_path = find(filename, home)
    print(audio_path)
    os.chdir("../feature-service")
    subprocess.run(['python3', 'main.py', 'add', '--audio_path', audio_path])
    os.chdir("../cluster-service")
    subprocess.run(['go', 'run', 'cmd/main/main.go', '--mode', 'add', '--path', audio_path])
    os.chdir("../ml-service")
    result = check_output(['go', 'run', 'cmd/main/main.go', '--mode', 'recognize', '--path', audio_path])
    result.decode(sys.stdout.encoding).strip()

    emotion = extract_emotion_from_log(str(result))

    return emotion


def extract_emotion_from_log(log_string):
    emotions = ['neutral', 'sad', 'positive', 'angry']
    for emotion in emotions:
        if emotion in log_string:
            return emotion
    return None


def find(filename, path):
    for root, dirs, files in os.walk(path):
        if filename in files:
            return os.path.join(root, filename)
