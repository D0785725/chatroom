import json
from channels.generic.websocket import AsyncWebsocketConsumer
import websockets
import asyncio

class ChatConsumer(AsyncWebsocketConsumer):
    async def connect(self):
        await self.channel_layer.group_add(
            'chat_group',
            self.channel_name
        )
        await self.accept()

    async def disconnect(self, close_code):
        await self.channel_layer.group_discard(
            'chat_group',
            self.channel_name
        )

    async def receive(self, text_data):
        text_data_json = json.loads(text_data)
        username = text_data_json['username']
        message = text_data_json['message']

        await self.channel_layer.group_send(
            'chat_group',
            {
                'type': 'chat_message',
                'username': username,
                'message': message
            }
        )

        # 將消息發送到 Go API 服務器
        async with websockets.connect('ws://localhost:8080/ws/chat/') as ws:
            await ws.send(json.dumps({'username': username, 'message': message}))

    async def chat_message(self, event):
        username = event['username']
        message = event['message']

        await self.send(text_data=json.dumps({
            'username': username,
            'message': message
        }))
