import { Component, OnInit } from '@angular/core';
import { AuthService } from '../../auth/services/auth.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  userName: string = '';
  
  // Estatísticas de exemplo
  statistics = {
    negocios: { total: 45, novos: 12, fechados: 8, valor: 'R$ 240.000,00' },
    projetos: { total: 18, emAndamento: 10, atrasados: 2, concluidos: 6 },
    tarefas: { total: 86, pendentes: 32, concluidas: 54, atrasadas: 8 },
    clientes: { total: 124, ativos: 98, novos: 15 }
  };
  
  // Dados de exemplo para gráficos
  recentActivities = [
    { type: 'task', title: 'Tarefa "Reunião com cliente" concluída', time: '10 minutos atrás', user: 'Você' },
    { type: 'project', title: 'Projeto "Implementação CRM" atualizado', time: '1 hora atrás', user: 'Ana Silva' },
    { type: 'business', title: 'Negócio "Consultoria ABC" movido para Fechado', time: '2 horas atrás', user: 'Carlos Mendes' },
    { type: 'client', title: 'Cliente "XYZ Ltda." adicionado', time: '3 horas atrás', user: 'Maria Souza' },
    { type: 'task', title: 'Comentário adicionado à tarefa "Desenvolver API"', time: '5 horas atrás', user: 'Paulo Santos' }
  ];
  
  constructor(private authService: AuthService) { }

  ngOnInit(): void {
    const user = this.authService.getCurrentUser();
    if (user) {
      this.userName = user.nome;
    }
  }
}