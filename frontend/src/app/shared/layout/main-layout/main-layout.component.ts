import { Component, OnInit } from '@angular/core';
import { RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatMenuModule } from '@angular/material/menu';
import { MatDividerModule } from '@angular/material/divider';
import { MatBadgeModule } from '@angular/material/badge';

import { AuthService, User } from '../../../auth/services/auth.service';

@Component({
  selector: 'app-main-layout',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatSidenavModule,
    MatToolbarModule,
    MatListModule,
    MatIconModule,
    MatButtonModule,
    MatExpansionModule,
    MatMenuModule,
    MatDividerModule,
    MatBadgeModule
  ],
  templateUrl: './main-layout.component.html',
  styleUrls: ['./main-layout.component.scss']
})
export class MainLayoutComponent implements OnInit {
  currentUser: User | null = null;
  sidenavOpened = true;
  notifications = [
    { title: 'Novo cliente cadastrado', time: '10 minutos atrás' },
    { title: 'Tarefa atribuída a você', time: '1 hora atrás' },
    { title: 'Negócio movido para coluna "Fechado"', time: '2 horas atrás' }
  ];

  constructor(private authService: AuthService) { }

  ngOnInit(): void {
    // Inscrever para mudanças no usuário atual
    this.authService.currentUser$.subscribe(user => {
      this.currentUser = user;
    });
    
    // Verificar se o dispositivo é móvel
    if (window.innerWidth < 768) {
      this.sidenavOpened = false;
    }
  }

  logout(): void {
    this.authService.logout();
  }
  
  toggleSidenav(): void {
    this.sidenavOpened = !this.sidenavOpened;
  }
}