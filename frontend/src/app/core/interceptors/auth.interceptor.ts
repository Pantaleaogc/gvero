import { HttpInterceptorFn, HttpRequest, HttpHandlerFn, HttpErrorResponse } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, throwError } from 'rxjs';
import { MatSnackBar } from '@angular/material/snack-bar';

export const authInterceptor: HttpInterceptorFn = (req: HttpRequest<unknown>, next: HttpHandlerFn) => {
  const router = inject(Router);
  const snackBar = inject(MatSnackBar);
  
  // Obter token do localStorage
  const token = localStorage.getItem('auth_token');
  
  // Adicionar token de autoriza«®«ªo se dispon«¿vel
  if (token) {
    req = req.clone({
      setHeaders: {
        Authorization: `Bearer ${token}`
      }
    });
  }
  
  // Continuar e interceptar erros
  return next(req).pipe(
    catchError((error: HttpErrorResponse) => {
      if (error.status === 401) {
        // Token expirou ou inv«¡lido
        localStorage.removeItem('auth_token');
        localStorage.removeItem('current_user');
        
        router.navigate(['/login']);
        snackBar.open('Sua sess«ªo expirou. Por favor, fa«®a login novamente.', 'Fechar', {
          duration: 5000
        });
      } 
      else if (error.status === 403) {
        snackBar.open('Voc«´ n«ªo tem permiss«ªo para acessar este recurso.', 'Fechar', {
          duration: 5000
        });
      }
      else if (error.status === 0) {
        // Erro de conex«ªo
        snackBar.open('N«ªo foi poss«¿vel conectar ao servidor. Verifique sua conex«ªo.', 'Fechar', {
          duration: 5000
        });
      }
      
      return throwError(() => error);
    })
  );
};