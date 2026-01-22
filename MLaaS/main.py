import pickle

from flask import Flask, request, jsonify

from ml.predict.predict import predict_multi
from models.reading import json_to_reading

app = Flask(__name__)

MODEL_PATH = 'ml/models/temp_forecast.pk'
data = pickle.load(open(MODEL_PATH, 'rb'))
knn, scaler = data['model'], data['scaler']
lags = data['lags']

@app.route('/predict', methods=['POST'])
def predict():
    try:
        data = request.get_json()
        history = [json_to_reading(r) for r in data.get('readings', [])]
        horizon = data.get('horizon', 1)
        result = predict_multi(history, horizon)
        return jsonify({'next_temps': result, 'status': 'success'})
    except Exception as e:
        return jsonify({'error': str(e)}), 400


if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=8100)