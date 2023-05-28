import os
import jsonlines
import random
import uuid

# Путь к папке с файлами JSONL
folder_path = "processed_dataset_090/aggregated_dataset"

# Количество записей для каждой эмоции
num_samples_per_emotion = 1500

# Диапазон длительности записи (в секундах)
duration_min = 1.5
duration_max = 4.0

# Путь и название для сохранения JSONL-файла
output_file_path = "processed_dataset_090/aggregated_dataset/dataset.jsonl"

# Словарь для хранения отфильтрованных записей
filtered_samples = {"sad": [], "positive": [], "neutral": [], "angry": []}

# Перебор файлов JSONL в папке
for file_name in os.listdir(folder_path):
    file_path = os.path.join(folder_path, file_name)
    with jsonlines.open(file_path) as reader:
        for sample in reader:
            duration = sample.get("duration", 0)
            emotion = sample.get("emotion")

            # Фильтрация записей по длительности и эмоции
            if duration > duration_min and duration < duration_max and emotion in filtered_samples:
                filtered_samples[emotion].append(sample)

# Сбалансированная выборка для каждой эмоции
balanced_samples = []
for emotion in filtered_samples:
    samples = filtered_samples[emotion]
    random.shuffle(samples)
    balanced_samples.extend(samples[:num_samples_per_emotion])

# Разделение на тестовую и обучающую выборку
random.shuffle(balanced_samples)
train_samples = balanced_samples[: int(len(balanced_samples) * 0.8)]
test_samples = balanced_samples[int(len(balanced_samples) * 0.8):]

# Добавление поля "batch" и генерация UUID для каждой записи
for sample in train_samples:
    sample["batch"] = "train"
    sample["uuid"] = str(uuid.uuid4())

for sample in test_samples:
    sample["batch"] = "test"
    sample["uuid"] = str(uuid.uuid4())

# Выбор выбранных полей для записи в финальный файл
selected_fields = ["audio_path", "duration", "emotion", "batch", "uuid"]

# Соединение обучающей и тестовой выборки
combined_samples = train_samples + test_samples

# Сохранение выбранных полей в JSONL-файл
with jsonlines.open(output_file_path, mode="w") as writer:
    for sample in combined_samples:
        selected_data = {field: sample[field] for field in selected_fields}
        writer.write(selected_data)
