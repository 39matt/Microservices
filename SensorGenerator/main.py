import random
from dataclasses import dataclass
from random import choice
from typing import List, Any

import pandas as pd
import requests
from tabulate import tabulate


@dataclass
class Reading:
    timestamp: str
    device_id: str
    co: float
    humidity: float
    light: bool
    lpg: float
    motion: bool
    smoke: float
    temperature: float

    def __post_init__(self):
        self.co = round(self.co, 4)
        self.humidity = round(self.humidity, 1)
        self.lpg = round(self.lpg, 4)
        self.smoke = round(self.smoke, 4)
        self.temperature = round(self.temperature, 1)

    def __str__(self) -> str:
        return f"[{self.timestamp[:16]}] {self.device_id} | T:{self.temperature}°C H:{self.humidity}% | Motion:{self.motion}"

    def as_api_dict(self) -> dict[str, Any]:
        return {
            "Timestamp": self.timestamp,
            "DeviceId": self.device_id,
            "Co": self.co,
            "Humidity": self.humidity,
            "Light": self.light,
            "Lpg": self.lpg,
            "Motion": self.motion,
            "Smoke": self.smoke,
            "Temperature": self.temperature
        }

    def table_row(self) -> List:
        return [self.timestamp[:19], self.device_id, self.temperature, self.humidity,
                self.co, self.smoke, self.lpg, self.light, self.motion]

    @classmethod
    def from_api_dict(cls, data: dict) -> 'Reading':
        return cls(
            timestamp=data['timestamp'],
            device_id=data['deviceId'],
            co=data['co'],
            humidity=data['humidity'],
            light=data['light'],
            lpg=data['lpg'],
            motion=data['motion'],
            smoke=data['smoke'],
            temperature=data['temperature']
        )

def pretty_print_readings(readings: List[Reading], limit: int = 100):
    print(f"\n📊 First {min(limit, len(readings))} Readings:")
    print(tabulate(
        [r.table_row() for r in readings[:limit]],
        headers=["Timestamp", "Device", "Temp°C", "Humidity%", "CO", "Smoke", "LPG", "Light", "Motion"],
        tablefmt="grid"
    ))

    print(f"\n🔥 Summary: {len(readings)} total | "
          f"Avg Temp: {sum(r.temperature for r in readings) / len(readings):.1f}°C | "
          f"Motion Events: {sum(r.motion for r in readings)}")

def read_from_csv(csvFile, limit: int = 10) -> List[Reading]:
    df = pd.read_csv(csvFile)

    bool_cols = ['light', 'motion']
    for col in bool_cols:
        df[col] = df[col].map({'true': True, 'false': False}).infer_objects(copy=False)

    readings = []
    for i, row in df.iterrows():
        if i >= limit:
            break
        reading = Reading(
            timestamp=str(row['ts']),
            device_id=str(row['device']),
            co=float(row['co']),
            humidity=float(row['humidity']),
            light=bool(row['light']),
            lpg=float(row['lpg']),
            motion=bool(row['motion']),
            smoke=float(row['smoke']),
            temperature=float(row['temp'])
        )
        readings.append(reading)

    return readings

def send_get_request(url: str = "http://localhost:5237", request:str = "/api/Readings/") -> List[Reading] | None:
    full_url = f"{url}{request}"
    try:
        response = requests.get(full_url)
        response.raise_for_status()

        readings_json = response.json()
        return [Reading.from_api_dict(r) for r in readings_json]

    except requests.exceptions.RequestException as e:
        print(e)
    return None

def send_post_request(reading: Reading, url: str = "http://localhost:5237", request:str = "/api/Readings/") -> requests.Response | None:
    url = f"{url}{request}"
    try:
        response = requests.post(url, json=reading.as_api_dict())
        if response.ok:
            print(f"Successfully sent {reading.timestamp} to {url}")
            return response.json()
        else:
            raise Exception(f"Request to {url} failed with status code {response.status_code}: {response.text}")


    except requests.exceptions.RequestException as e:
        print(e)
    return None

def send_delete_request(url: str = "http://localhost:5237", request:str = "/api/Readings/") -> None:
    full_url = f"{url}{request}"
    try:
        response = requests.delete(full_url)
        response.raise_for_status()

        return None

    except requests.exceptions.RequestException as e:
        print(e)
    return None

def send_batch(readings: List[Reading]) -> int | None:
    success_count = 0
    for reading in readings:
        if send_post_request(reading):
            success_count += 1

    print(f"✅ Sent {success_count}/{len(readings)} readings")
    return success_count

if __name__ == "__main__":
    readings = read_from_csv("./../Data/iot_telemetry_data.csv", 1000)
    print(f"Loaded {len(readings)} readings")

    choice = input("Input how many random readings do you want to send to the Gateway?\nEnter 'D' if you want to delete all rows.\n").strip()

    if choice.upper() == 'D':
        send_delete_request()
        print("Successfully deleted all rows")
    else:
        choice = int(choice)
        random_readings = random.sample(readings, choice)
        send_batch(random_readings)

    existing_readings = send_get_request()
    if existing_readings:
        pretty_print_readings(existing_readings)