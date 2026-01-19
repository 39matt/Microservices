'use client'

import {useEffect, useRef, useState} from "react";
import mqtt, {MqttClient} from 'mqtt';

export default function Page() {
    const [readings, setReadings] = useState<string[]>([]);
    const [connectStatus, setConnectStatus] = useState('');
    const [topic, setTopic] = useState('/limit');
    const clientRef = useRef<MqttClient>(null);

    useEffect(() => {
        const client = mqtt.connect('ws://mosquitto:1884');
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
            setReadings(prev => [...prev, message.toString()]);
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

    return(
        <>
            <div className={"m-auto"}>
                <div>Status: {connectStatus}</div>
                <div>
                    {readings.map((reading, idx) => (
                        <div key={idx}>{reading}</div>
                    ))}
                </div>
            </div>
        </>
    )
}