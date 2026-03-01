import { Component, Input, OnInit, OnDestroy } from '@angular/core';
import { ChatService } from '../../services/chat.service';
import { WebSocketService } from '../../services/websocket.service';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-chat-window',
  templateUrl: './chat-window.component.html',
  styleUrls: ['./chat-window.component.css']
})
export class ChatWindowComponent implements OnInit, OnDestroy {
  @Input() chat: any;

  messages: any[] = [];
  newMessage = '';
  loading = false;
  private destroy$ = new Subject<void>();

  constructor(
    private chatService: ChatService,
    private webSocketService: WebSocketService
  ) {}

  ngOnInit() {
    if (this.chat) {
      this.loadMessages();
      this.connectWebSocket();
    }
  }

  ngOnDestroy() {
    this.webSocketService.disconnect();
    this.destroy$.next();
    this.destroy$.complete();
  }

  loadMessages() {
    this.chatService.getChatMessages(this.chat.id)
      .pipe(takeUntil(this.destroy$))
      .subscribe(
        (messages: any[]) => {
          this.messages = messages;
          this.scrollToBottom();
        },
        (error) => {
          console.error('Error loading messages:', error);
        }
      );
  }

  connectWebSocket() {
    this.webSocketService.connect(this.chat.id)
      .pipe(takeUntil(this.destroy$))
      .subscribe(
        (message) => {
          if (message.type !== 'connected') {
            this.messages.push(message);
            this.scrollToBottom();
          }
        },
        (error) => {
          console.error('WebSocket error:', error);
        }
      );
  }

  sendMessage() {
    if (!this.newMessage.trim()) {
      return;
    }

    this.loading = true;
    const messageText = this.newMessage;
    this.newMessage = '';

    this.webSocketService.sendMessage({
      content: messageText,
      type: 'text',
      chatId: this.chat.id
    });

    // Alternatively, send via HTTP
    this.chatService.sendMessage(this.chat.id, messageText)
      .pipe(takeUntil(this.destroy$))
      .subscribe(
        () => {
          this.loading = false;
        },
        (error) => {
          this.loading = false;
          console.error('Error sending message:', error);
        }
      );
  }

  private scrollToBottom() {
    setTimeout(() => {
      const messagesContainer = document.querySelector('.messages-container');
      if (messagesContainer) {
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
      }
    }, 0);
  }

  onKeyDown(event: KeyboardEvent) {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      this.sendMessage();
    }
  }
}
