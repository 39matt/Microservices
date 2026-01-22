import pickle
import pandas as pd
from sklearn.metrics import mean_absolute_error
from sklearn.model_selection import train_test_split
from sklearn.neighbors import KNeighborsClassifier, KNeighborsRegressor
from sklearn.preprocessing import StandardScaler


def GenerateMLModel():
    df = pd.read_csv("../../data/iot_telemetry_data.csv")

    df["ts"] = pd.to_datetime(df['ts'].astype(float), unit='s')
    df = df.sort_values('ts').reset_index(drop=True)


    lags = 3
    for lag in range(1, lags + 1):
        df[f'temp_lag{lag}'] = df['temp'].shift(lag)
        df[f'hum_lag{lag}'] = df['humidity'].shift(lag)
    print(f"After lags: {df.shape[0]:,} rows")

    df = df.dropna()
    X =  df[['temp_lag1', 'temp_lag2', 'temp_lag3', 'hum_lag1', 'co', 'lpg', 'smoke']].values
    Y = df['temp'].shift(-1)[:-1].values

    X = X[:-1]

    print(f"X/Y shape: {X.shape}")
    X_train, X_test, Y_train, Y_test = train_test_split(X, Y, test_size=0.2, random_state=100, shuffle=False)

    scaler = StandardScaler()
    X_train = scaler.fit_transform(X_train)
    X_test = scaler.transform(X_test)

    knn = KNeighborsRegressor(n_neighbors=50)
    knn.fit(X_train, Y_train)

    train_mae = mean_absolute_error(Y_train, knn.predict(X_train))
    test_mae = mean_absolute_error(Y_test, knn.predict(X_test))
    print(f"Train MAE: {train_mae:.2f}°C | Test MAE: {test_mae:.2f}°C")

    with open('../models/temp_forecast.pk', 'wb') as f:
        pickle.dump({'model': knn, 'scaler': scaler, 'lags': lags}, f)
    print("Model saved!")