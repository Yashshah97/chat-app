import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { BehaviorSubject, Observable } from 'rxjs';
import { AuthService } from './auth.service';

@Injectable({
  providedIn: 'root'
})
export class ChatService {
  private apiUrl = 'http://localhost:8080/api/chats';
  private messagesSubject = new BehaviorSubject<any[]>([]);
  public messages$ = this.messagesSubject.asObservable();

  constructor(
    private http: HttpClient,
    private authService: AuthService
  ) {}

  private getHeaders(): HttpHeaders {
    const token = this.authService.getToken();
    return new HttpHeaders({
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    });
  }

  getChats(): Observable<any[]> {
    return this.http.get<any[]>(this.apiUrl, {
      headers: this.getHeaders()
    });
  }

  getChatMessages(chatId: number): Observable<any[]> {
    return this.http.get<any[]>(
      `${this.apiUrl}/${chatId}/messages`,
      { headers: this.getHeaders() }
    );
  }

  createPrivateChat(userId: number): Observable<any> {
    return this.http.post(
      this.apiUrl,
      { userId },
      { headers: this.getHeaders() }
    );
  }

  createGroupChat(name: string, members: number[]): Observable<any> {
    return this.http.post(
      `${this.apiUrl}/group`,
      { name, members },
      { headers: this.getHeaders() }
    );
  }

  sendMessage(chatId: number, content: string): Observable<any> {
    return this.http.post(
      `${this.apiUrl}/${chatId}/messages`,
      { content },
      { headers: this.getHeaders() }
    );
  }

  updateMessages(messages: any[]) {
    this.messagesSubject.next(messages);
  }
}
