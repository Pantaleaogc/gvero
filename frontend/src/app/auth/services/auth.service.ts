import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable, of, throwError } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';
import { Router } from '@angular/router';
import { environment } from '../../../environments/environment';

export interface User {
  id: number;
  nome: string;
  email: string;
  tipo: string;
  empresa_id?: number;
}

export interface LoginResponse {
  token: string;
  user: User;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private apiUrl = `${environment.apiUrl}/auth`;
  private currentUserSubject = new BehaviorSubject<User | null>(null);
  public currentUser$ = this.currentUserSubject.asObservable();
  
  constructor(
    private http: HttpClient,
    private router: Router
  ) {
    // Verificar se há um usuário no localStorage ao iniciar
    const token = localStorage.getItem('auth_token');
    const storedUser = localStorage.getItem('current_user');
    
    if (token && storedUser) {
      this.currentUserSubject.next(JSON.parse(storedUser));
    }
  }

  login(email: string, password: string): Observable<LoginResponse> {
    return this.http.post<LoginResponse>(`${this.apiUrl}/login`, { email, password }).pipe(
      tap(response => {
        // Armazenar token e dados do usuário
        localStorage.setItem('auth_token', response.token);
        localStorage.setItem('current_user', JSON.stringify(response.user));
        this.currentUserSubject.next(response.user);
      }),
      catchError(error => {
        console.error('Erro ao fazer login:', error);
        return throwError(() => new Error('Credenciais inválidas. Por favor, verifique seu email e senha.'));
      })
    );
  }

  logout(): void {
    // Opcional: Chamar API para invalidar o token
    this.http.post(`${this.apiUrl}/logout`, {}).subscribe({
      next: () => console.log('Logout na API realizado com sucesso'),
      error: err => console.error('Erro ao fazer logout na API:', err),
      complete: () => {
        // Limpar dados locais
        localStorage.removeItem('auth_token');
        localStorage.removeItem('current_user');
        this.currentUserSubject.next(null);
        
        // Redirecionar para login
        this.router.navigate(['/login']);
      }
    });
  }

  isAuthenticated(): boolean {
    return !!localStorage.getItem('auth_token');
  }

  getCurrentUser(): User | null {
    return this.currentUserSubject.value;
  }

  // Verificar token e revalidar usuário
  verifyToken(): Observable<boolean> {
    const token = localStorage.getItem('auth_token');
    if (!token) {
      return of(false);
    }

    return this.http.get<User>(`${this.apiUrl}/verify`).pipe(
      tap(user => {
        localStorage.setItem('current_user', JSON.stringify(user));
        this.currentUserSubject.next(user);
      }),
      catchError(() => {
        this.logout();
        return of(false);
      })
    );
  }

  // Verificar se o usuário tem permissão específica
  hasPermission(permission: string): boolean {
    const user = this.currentUserSubject.value;
    if (!user) return false;
    
    // Aqui você implementaria a lógica baseada no seu modelo de permissões
    if (user.tipo === 'admin') return true;
    
    // Para outros tipos, verificar permissões específicas
    const permissionMap: Record<string, string[]> = {
      'gerente': ['view:all', 'edit:projects', 'view:reports'],
      'vendedor': ['view:sales', 'edit:leads'],
      'usuario': ['view:assigned']
    };
    
    return permissionMap[user.tipo]?.includes(permission) || false;
  }
}