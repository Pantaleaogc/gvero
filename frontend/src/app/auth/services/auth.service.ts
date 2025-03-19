import { Injectable } from '@angular/core';
import { ApiService } from '../../core/services/api.service';
import { BehaviorSubject, Observable, of, throwError } from 'rxjs';
import { catchError, map, tap } from 'rxjs/operators';
import { Router } from '@angular/router';

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
  private currentUserSubject = new BehaviorSubject<User | null>(null);
  public currentUser$ = this.currentUserSubject.asObservable();
  
  constructor(
    private apiService: ApiService,
    private router: Router
  ) {
    // Verificar se há um usuário no localStorage ao iniciar
    const storedUser = localStorage.getItem('current_user');
    if (storedUser) {
      this.currentUserSubject.next(JSON.parse(storedUser));
    }
  }

  login(email: string, password: string): Observable<LoginResponse> {
    return this.apiService.post<LoginResponse>('auth/login', { email, password }).pipe(
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
    this.apiService.post('auth/logout', {}).subscribe({
      next: () => console.log('Logout na API realizado com sucesso'),
      error: err => console.error('Erro ao fazer logout na API:', err)
    });
    
    // Limpar dados locais
    localStorage.removeItem('auth_token');
    localStorage.removeItem('current_user');
    this.currentUserSubject.next(null);
    
    // Redirecionar para login
    this.router.navigate(['/auth/login']);
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

    return this.apiService.get<User>('auth/verify').pipe(
      map(user => {
        localStorage.setItem('current_user', JSON.stringify(user));
        this.currentUserSubject.next(user);
        return true;
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
    // Por exemplo, verificando o tipo de usuário ou uma lista de permissões
    
    if (user.tipo === 'admin') return true;
    
    // Para outros tipos, verificar permissões específicas
    // Este é apenas um exemplo - ajuste conforme seu modelo
    const permissionMap: Record<string, string[]> = {
      'gerente': ['view:all', 'edit:projects', 'view:reports'],
      'vendedor': ['view:sales', 'edit:leads'],
      'usuario': ['view:assigned']
    };
    
    return permissionMap[user.tipo]?.includes(permission) || false;
  }
}