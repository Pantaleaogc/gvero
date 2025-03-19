import { Routes } from '@angular/router';
import { LoginComponent } from './auth/login/login.component';
import { MainLayoutComponent } from './shared/layout/main-layout/main-layout.component';
import { AuthGuard } from './auth/guards/auth.guard';

export const routes: Routes = [
  // Rotas públicas
  { path: 'login', component: LoginComponent },
  
  // Rotas protegidas (com layout principal)
  {
    path: '',
    component: MainLayoutComponent,
    canActivate: [AuthGuard],
    canActivateChild: [AuthGuard],
    children: [
      {
        path: 'dashboard',
        loadChildren: () => import('./dashboard/dashboard.module').then(m => m.DashboardModule)
      },
      {
        path: 'kanban',
        loadChildren: () => import('./kanban/kanban.module').then(m => m.KanbanModule)
      },
      {
        path: 'clientes',
        loadChildren: () => import('./clientes/clientes.module').then(m => m.ClientesModule)
      },
      {
        path: 'usuarios',
        loadChildren: () => import('./usuarios/usuarios.module').then(m => m.UsuariosModule),
        data: { permission: 'admin' } // Requer permissão de administrador
      },
      {
        path: 'configuracoes',
        loadChildren: () => import('./configuracoes/configuracoes.module').then(m => m.ConfiguracoesModule)
      },
      // Página inicial (redireciona para dashboard)
      { path: '', redirectTo: 'dashboard', pathMatch: 'full' }
    ]
  },
  
  // Qualquer outra rota redireciona para login
  { path: '**', redirectTo: 'login' }
];