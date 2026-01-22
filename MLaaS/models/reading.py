from dataclasses import dataclass


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

def json_to_reading(r: dict) -> Reading:
    return Reading(
        timestamp=r['ts'],
        device_id=r['device'],
        temperature=float(r.get('temp', r.get('temperature', 0))),
        humidity=float(r.get('humidity', 0)),
        co=float(r.get('co', 0)),
        lpg=float(r.get('lpg', 0)),
        smoke=float(r.get('smoke', 0)),
        light=bool(r.get('light', False)),
        motion=bool(r.get('motion', False))
    )