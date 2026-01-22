import os
import pickle; import numpy as np

from models.reading import Reading


def predict(reading: Reading) -> float:

    data = pickle.load(open('../models/temp_forecast.pk', 'rb'))
    knn, scaler, lags = data['model'], data['scaler'], data['lags']

    new_data = np.array(reading).reshape(1,-1)
    next_temp = knn.predict(scaler.transform(new_data))[0]

    return next_temp

# Usage: last 4 readings → next 3 temps
# next_temps = predict_multi(history, h=3)
# print([f"{t:.1f}°C" for t in next_temps])
def predict_multi(history: list[Reading], h: int = 3, lags: int = 3) -> list[float]:
    if len(history) < lags + 1:
        raise ValueError(f"Need {lags + 1} readings")

    model_path = os.path.join(os.path.dirname(__file__), '..', 'models', 'temp_forecast.pk')
    data = pickle.load(open(model_path, 'rb'))
    knn, scaler = data['model'], data['scaler']

    predictions = []
    current_history = history.copy()

    for _ in range(h):
        features = [
            *[r.temperature for r in current_history[-lags:]],
            current_history[-1].humidity,
            current_history[-1].co,
            current_history[-1].lpg,
            current_history[-1].smoke
        ]
        new_data = np.array([features]).reshape(1, -1)
        next_temp = knn.predict(scaler.transform(new_data))[0]
        predictions.append(next_temp)

        fake_reading = Reading(
            timestamp="", device_id="", co=current_history[-1].co,
            humidity=current_history[-1].humidity, light=False,
            lpg=current_history[-1].lpg, motion=False, smoke=current_history[-1].smoke,
            temperature=next_temp
        )
        current_history.append(fake_reading)

    return predictions


