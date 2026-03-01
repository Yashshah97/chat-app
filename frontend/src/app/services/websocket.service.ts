import { Injectable } from '@angular/core';
import { Subject, Observable } from 'rxjs';
import { AuthService } from './auth.service';

@Injectable({
  providedIn: 'root'
})
export class WebSocketService {
  private socket: WebSocket | null = null;
  private messagesSubject = new Subject<any>();
  private connectionSubject = new Subject<boolean>();

  constructor(private authService: AuthService) {}

  connect(chatId: number): Observable<any> {
    const token = this.authService.getToken();
    const wsUrl = `ws://localhost:8080/ws/chat/${chatId}?token=${token}`;

    return new Observable((observer) => {
      this.socket = new WebSocket(wsUrl);

      this.socket.onopen = () => {
        observer.next({ type: 'connected' });
        this.connectionSubject.next(true);
      };

      this.socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        observer.next(message);
        this.messagesSubject.next(message);
      };

      this.socket.onerror = (error) => {
        observer.error(error);
        this.connectionSubject.next(false);
      };

      this.socket.onclose = () => {
        observer.complete();
        this.connectionSubject.next(false);
      };
    });
  }

  sendMessage(message: any) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(message));
    }
  }

  disconnect() {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
  }

  getMessages(): Observable<any> {
    return this.messagesSubject.asObservable();
  }

  getConnectionStatus(): Observable<boolean> {
    return this.connectionSubject.asObservable();
  }

  isConnected(): boolean {
    return this.socket !== null && this.socket.readyState === WebSocket.OPEN;
  }
}
