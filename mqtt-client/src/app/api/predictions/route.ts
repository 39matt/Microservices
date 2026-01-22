// app/api/predictions/route.ts
import { connect, Subscription } from 'nats';

export async function GET(request: Request) {
    try {
        console.log('🔗 /api/predictions: Starting NATS connect...');
        const nc = await connect({
            servers: 'nats://nats:4222',
            timeout: 10_000,
        });
        console.log('✅ NATS connected to', nc.getServer());

        return new Response(new ReadableStream({
            async start(controller) {
                console.log('📡 Creating subscription to "predictions"...');
                const sub: Subscription = nc.subscribe('predictions', { queue: 'web' });
                console.log('✅ Subscription created:', sub.getSubject(), 'SID:', sub.getID());

                let msgCount = 0;
                try {
                    for await (const msg of sub) {
                        msgCount++;
                        const dataStr = new TextDecoder().decode(msg.data);
                        console.log(`📨 Msg #${msgCount}: "${dataStr}" from ${msg.subject}`);
                        controller.enqueue(`data: ${dataStr}\n\n`);
                    }
                } catch (subErr) {
                    console.error('❌ Sub iterator error:', subErr);
                } finally {
                    console.log(`🔚 Subscription drained after ${msgCount} msgs`);
                    sub.unsubscribe();
                    await nc.close();
                }
            },
            cancel() {
                console.log('🛑 SSE stream cancelled');
            }
        }), {
            headers: {
                'Content-Type': 'text/event-stream',
                'Cache-Control': 'no-cache',
                'Connection': 'keep-alive',
            },
        });
    } catch (error) {
        console.error('💥 NATS / SSE failed:', error);
        return Response.json({ error: (error as Error).message }, { status: 500 });
    }
}
