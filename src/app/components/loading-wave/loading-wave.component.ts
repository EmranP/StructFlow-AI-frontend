import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-loading-wave',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './loading-wave.component.html',
  styleUrls: ['./loading-wave.component.scss']
})
export class LoadingWaveComponent {
  @Input() message = 'Generating your project structure...';
  dots = Array(12).fill(0);
}
