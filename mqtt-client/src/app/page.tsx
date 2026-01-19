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

export default function Page() {
  const [readings, setReadings] = useState<Reading[]>([]);
  const [connectStatus, setConnectStatus] = useState('');
  const [topic, setTopic] = useState('/limit');
  const clientRef = useRef<MqttClient>(null);

  useEffect(() => {
    const client = mqtt.connect('ws://localhost:9002');
    clientRef.current = client;


    client.on('connect', () => {
      setConnectStatus('Connected');
      client.subscribe(topic, (err)=>{
        if (err) {
          console.error(err);
        }
      })
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

    return () => {
      if (clientRef.current) {
        clientRef.current.end();
      }
    };
  }, []);

  if (!readings) {
    return (
        <div className={"m-auto"}>Status: {connectStatus}</div>
    )
  }

  return (
      <div className="min-h-screen bg-gray-100 p-8">
        <div className="max-w-7xl mx-auto">
          {/* Header */}
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

          {/* Readings Grid */}
          {readings.length === 0 ? (
              <div className="bg-white rounded-lg shadow-md p-12 text-center">
                <p className="text-gray-500 text-lg">Waiting for sensor data...</p>
              </div>
          ) : (
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {readings.map((reading, idx) => (
                    <div key={idx} className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
                      {/* Card Header */}
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

                      {/* Sensor Readings Grid */}
                      <div className="grid grid-cols-2 gap-4">
                        {/* Temperature */}
                        <div className="bg-blue-50 rounded-lg p-4">
                          <p className="text-sm text-blue-600 font-medium mb-1">Temperature</p>
                          <p className="text-2xl font-bold text-blue-900">
                            {reading.Temperature.toFixed(1)}°C
                          </p>
                        </div>

                        {/* Humidity */}
                        <div className="bg-cyan-50 rounded-lg p-4">
                          <p className="text-sm text-cyan-600 font-medium mb-1">Humidity</p>
                          <p className="text-2xl font-bold text-cyan-900">
                            {reading.Humidity.toFixed(1)}%
                          </p>
                        </div>

                        {/* CO */}
                        <div className="bg-orange-50 rounded-lg p-4">
                          <p className="text-sm text-orange-600 font-medium mb-1">CO Level</p>
                          <p className="text-2xl font-bold text-orange-900">
                            {reading.Co.toFixed(2)}
                          </p>
                        </div>

                        {/* LPG */}
                        <div className="bg-purple-50 rounded-lg p-4">
                          <p className="text-sm text-purple-600 font-medium mb-1">LPG</p>
                          <p className="text-2xl font-bold text-purple-900">
                            {reading.Lpg.toFixed(2)}
                          </p>
                        </div>

                        {/* Smoke */}
                        <div className="bg-red-50 rounded-lg p-4">
                          <p className="text-sm text-red-600 font-medium mb-1">Smoke</p>
                          <p className="text-2xl font-bold text-red-900">
                            {reading.Smoke.toFixed(2)}
                          </p>
                        </div>

                        {/* Light */}
                        <div className="bg-yellow-50 rounded-lg p-4">
                          <p className="text-sm text-yellow-600 font-medium mb-1">Light</p>
                          <p className="text-2xl font-bold text-yellow-900">
                            {reading.Light ? '💡 On' : '⚫ Off'}
                          </p>
                        </div>

                        {/* Motion */}
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
      </div>
  );
}