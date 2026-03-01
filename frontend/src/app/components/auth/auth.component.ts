import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-auth',
  templateUrl: './auth.component.html',
  styleUrls: ['./auth.component.css']
})
export class AuthComponent {
  isLoginMode = true;
  loading = false;
  error: string | null = null;
  
  loginForm = {
    username: '',
    password: ''
  };

  registerForm = {
    username: '',
    email: '',
    password: '',
    confirmPassword: ''
  };

  constructor(
    private authService: AuthService,
    private router: Router
  ) {}

  toggleMode() {
    this.isLoginMode = !this.isLoginMode;
    this.error = null;
  }

  login() {
    if (!this.loginForm.username || !this.loginForm.password) {
      this.error = 'Please fill in all fields';
      return;
    }

    this.loading = true;
    this.authService.login(this.loginForm.username, this.loginForm.password).subscribe(
      () => {
        this.loading = false;
        this.router.navigate(['/chat']);
      },
      (error) => {
        this.loading = false;
        this.error = error.error?.message || 'Login failed. Please try again.';
      }
    );
  }

  register() {
    if (!this.registerForm.username || !this.registerForm.email || !this.registerForm.password) {
      this.error = 'Please fill in all fields';
      return;
    }

    if (this.registerForm.password !== this.registerForm.confirmPassword) {
      this.error = 'Passwords do not match';
      return;
    }

    this.loading = true;
    this.authService.register(
      this.registerForm.username,
      this.registerForm.email,
      this.registerForm.password
    ).subscribe(
      () => {
        this.loading = false;
        this.router.navigate(['/chat']);
      },
      (error) => {
        this.loading = false;
        this.error = error.error?.message || 'Registration failed. Please try again.';
      }
    );
  }
}
