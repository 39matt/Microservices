'use client'

import {useEffect, useRef, useState} from "react";
import mqtt, {MqttClient} from 'mqtt';

interface Reading {
  ID: string;
  Timestamp: string;
  DeviceID: string;
  Co: number;
  Humidity: number;
  Light: boolean;
  Lpg: number;
  Motion: boolean;
  Smoke: number;
  Temperature: number;
}

interface Prediction {
  temps: number[];
  timestamp: string;
  count: number;
}

export default function Page() {
  const [readings, setReadings] = useState<Reading[]>([]);
  const [predictions, setPredictions] = useState<Prediction | null>(null);
  const [connectStatus, setConnectStatus] = useState('');
  const [topic, setTopic] = useState('/limit');
  const mqttClientRef = useRef<MqttClient>(null);

  function InitMqttClient() {
    const client = mqtt.connect('ws://localhost:9002');
    mqttClientRef.current = client;

    client.on('connect', () => {
      setConnectStatus('Connected');
      client.subscribe(topic, (err) => {
        if (err) console.error(err);
      });
    });
    client.on('error', (err) => {
      console.error('Connection error: ', err);
      setConnectStatus('Error');
      client.end();
    });
    client.on('reconnect', () => {
      setConnectStatus('Reconnecting');
    });
    client.on('message', (topic, message) => {
      const reading = JSON.parse(message.toString());
      console.log(reading);
      setReadings(prev => [...prev, reading]);
    });
  }

  useEffect(() => {
    if (typeof window === 'undefined') return;

    const eventSource = new EventSource('/api/predictions');
    eventSource.onopen = () => console.log('🔓 SSE opened');
    eventSource.onmessage = (event) => {
      try {
        console.log("aaaa", JSON.parse(event.data));
        const temps: number[] = JSON.parse(event.data);
        const pred: Prediction = {
          temps,
          timestamp: new Date().toLocaleTimeString(),
          count: temps.length,
        };
        setPredictions(pred);
        console.log('Predicted temps:', temps);
      } catch (err) {
        console.error('SSE parse error:', err);
      }
    };

    eventSource.onerror = (err) => {
      console.error('SSE connection lost:', err);
    };

    return () => {
      eventSource.close();
    };
  }, []);

  useEffect(() => {
    InitMqttClient();

    return () => {
      if (mqttClientRef.current) {
        mqttClientRef.current.end();
      }
    };
  }, []);

  const getTempEmoji = (temp: number) => {
    if (temp >= 30) return '🔥'
    if (temp <= 15) return '🧊';
    return '☀️';
  };

  return (
      <div className="min-h-screen bg-gray-100 p-8">
        <div className="max-w-7xl mx-auto">
          <div className="grid lg:grid-cols-4 gap-8">
            {/* Header */}
            <div className="lg:col-span-4">
              <div className="bg-white rounded-lg shadow-md p-6 mb-6">
                <h1 className="text-3xl font-bold text-gray-800 mb-2">
                  IoT Sensor Dashboard
                </h1>
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium text-gray-600">Status:</span>
                  <span className={`px-3 py-1 rounded-full text-sm font-semibold ${
                      connectStatus === 'Connected'
                          ? 'bg-green-100 text-green-800'
                          : connectStatus === 'Error'
                              ? 'bg-red-100 text-red-800'
                              : 'bg-yellow-100 text-yellow-800'
                  }`}>
                  {connectStatus}
                </span>
                </div>
              </div>
            </div>

            {/* Readings Grid */}
            <div className="lg:col-span-3">
              {readings.length === 0 ? (
                  <div className="bg-white rounded-lg shadow-md p-12 text-center">
                    <p className="text-gray-500 text-lg">Waiting for sensor data...</p>
                  </div>
              ) : (
                  <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    {readings.sort((a, b)=> +a.Timestamp - +b.Timestamp).map((reading, idx) => (
                        <div key={idx} className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
                          <div className="border-b pb-4 mb-4">
                            <h2 className="text-xl font-semibold text-gray-800">
                              Device: {reading.DeviceID}
                            </h2>
                            <p className="text-sm text-gray-500 mt-1">
                              {(() => {
                                const timestamp = Number(reading.Timestamp) * 1000;
                                const date = new Date(timestamp);
                                return isNaN(date.getTime()) ? reading.Timestamp : date.toLocaleString();
                              })()}
                            </p>
                          </div>
                          <div className="grid grid-cols-2 gap-4">
                            <div className="bg-blue-50 rounded-lg p-4">
                              <p className="text-sm text-blue-600 font-medium mb-1">Temperature</p>
                              <p className="text-2xl font-bold text-blue-900">{reading.Temperature.toFixed(1)}°C</p>
                            </div>
                            <div className="bg-cyan-50 rounded-lg p-4">
                              <p className="text-sm text-cyan-600 font-medium mb-1">Humidity</p>
                              <p className="text-2xl font-bold text-cyan-900">{reading.Humidity.toFixed(1)}%</p>
                            </div>
                            <div className="bg-orange-50 rounded-lg p-4">
                              <p className="text-sm text-orange-600 font-medium mb-1">CO Level</p>
                              <p className="text-2xl font-bold text-orange-900">{reading.Co.toFixed(2)}</p>
                            </div>
                            <div className="bg-purple-50 rounded-lg p-4">
                              <p className="text-sm text-purple-600 font-medium mb-1">LPG</p>
                              <p className="text-2xl font-bold text-purple-900">{reading.Lpg.toFixed(2)}</p>
                            </div>
                            <div className="bg-red-50 rounded-lg p-4">
                              <p className="text-sm text-red-600 font-medium mb-1">Smoke</p>
                              <p className="text-2xl font-bold text-red-900">{reading.Smoke.toFixed(2)}</p>
                            </div>
                            <div className="bg-yellow-50 rounded-lg p-4">
                              <p className="text-sm text-yellow-600 font-medium mb-1">Light</p>
                              <p className="text-2xl font-bold text-yellow-900">
                                {reading.Light ? '💡 On' : '⚫ Off'}
                              </p>
                            </div>
                            <div className="bg-green-50 rounded-lg p-4 col-span-2">
                              <p className="text-sm text-green-600 font-medium mb-1">Motion</p>
                              <p className="text-2xl font-bold text-green-900">
                                {reading.Motion ? '🏃 Detected' : '⭕ None'}
                              </p>
                            </div>
                          </div>
                        </div>
                    ))}
                  </div>
              )}
            </div>

            <div className="lg:col-span-1">
              <div className="bg-gradient-to-br from-sky-400 via-sky-500 to-sky-600 rounded-xl shadow-2xl p-6 text-white sticky top-8 h-fit">
                <h2 className="text-2xl font-bold mb-6 flex items-center gap-3">
                  🤖 ML Predictions
                </h2>
                {predictions === null ? (
                    <div className="text-sky-100 text-center py-12">
                      <p className="text-lg mb-4">Waiting for predictions...</p>
                      <div className="w-24 h-24 bg-white/10 rounded-2xl mx-auto animate-pulse"></div>
                    </div>
                ) : (
                    <div className="space-y-4 max-h-96 overflow-y-auto">
                      <div className="bg-white/10 backdrop-blur-sm rounded-2xl p-4 border border-white/20">
                        <div className="text-xs opacity-75 mb-3 font-mono">{predictions.timestamp}</div>
                        {/* Wider grid, smaller fonts */}
                        <div className="grid grid-cols-5 gap-1 text-xs mb-3">  {/* gap-1, text-xs */}
                          {Array.from({length: 5}, (_, h) => predictions.temps[h] ?? '?').map((temp, tIdx) => (
                              <div key={tIdx} className={`p-2 rounded-lg text-center transition-all ${
                                  typeof temp === 'number'
                                      ? temp >= 30 ? 'bg-red-500/30'
                                          : temp <= 15 ? 'bg-blue-500/30'
                                              : 'bg-green-500/30'
                                      : 'bg-gray-500/30'
                              }`}>
                                <div className="text-[10px] opacity-75 mb-1">H{(tIdx+1).toString().padStart(2, '0')}</div>
                                <div className="text-sm font-bold flex items-center justify-center gap-0.5">  {/* text-sm, gap-0.5 */}
                                  {getTempEmoji(temp as number)}
                                  {typeof temp === 'number' ? temp.toFixed(1) + '°' : '—'}
                                </div>
                              </div>
                          ))}
                        </div>
                        <div className="text-xs opacity-75 flex items-center gap-2">
                          📊 Batch: {readings.length} readings
                          <span className="ml-auto text-sky-200">5h forecast</span>
                        </div>
                      </div>
                    </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
  );
}
