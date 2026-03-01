import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from './services/auth.service';
import { ChatService } from './services/chat.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  title = 'Chat Application';
  isAuthenticated = false;
  currentUser: any = null;
  chats: any[] = [];
  selectedChat: any = null;

  constructor(
    private authService: AuthService,
    private chatService: ChatService,
    private router: Router
  ) {}

  ngOnInit() {
    this.checkAuthStatus();
  }

  checkAuthStatus() {
    const token = localStorage.getItem('token');
    const user = localStorage.getItem('user');
    
    if (token && user) {
      this.isAuthenticated = true;
      this.currentUser = JSON.parse(user);
      this.loadChats();
    }
  }

  loadChats() {
    this.chatService.getChats().subscribe(
      (chats: any[]) => {
        this.chats = chats;
      },
      (error) => {
        console.error('Error loading chats:', error);
      }
    );
  }

  selectChat(chat: any) {
    this.selectedChat = chat;
  }

  logout() {
    this.authService.logout();
    this.isAuthenticated = false;
    this.currentUser = null;
    this.chats = [];
    this.selectedChat = null;
    this.router.navigate(['/login']);
  }

  createNewChat() {
    // Navigate to create chat component
    this.router.navigate(['/chat/new']);
  }

  goToProfile() {
    this.router.navigate(['/profile']);
  }
}
